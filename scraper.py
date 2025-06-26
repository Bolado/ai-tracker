import logging
import os
import json
from typing import List, Optional
from pydantic import BaseModel, ValidationError, Field, field_validator
from datetime import datetime
from urllib.parse import urljoin

from crawl4ai import AsyncWebCrawler, CacheMode
from crawl4ai.extraction_strategy import LLMExtractionStrategy
from crawl4ai.async_configs import CrawlerRunConfig, LLMConfig

from models import Website
from database import get_articles_collection

logger = logging.getLogger(__name__)


# Pydantic models for structured data extraction
class ArticleInfo(BaseModel):
    """Pydantic model for individual article extraction"""
    title: str = Field(..., description="The main title/headline of the article")
    url: str = Field(..., description="The full URL link to the article")
    image_url: str = Field(..., description="Featured image URL if available")


class ArticleListExtraction(BaseModel):
    """Pydantic model for extracting lists of articles from a page"""
    articles: List[ArticleInfo] = Field(..., description="List of AI-related articles found on the page")


class DetailedArticleContent(BaseModel):
    title: str = Field(..., description="The main headline of the article")
    content: str = Field(..., description="A very concise summary of the article with around 200 characters, **no more** than 240 characters.", max_length=240)
    publication_date: str = Field(..., description="The date/time the article was published at.")

    @field_validator('content')
    @classmethod
    def validate_content(cls, v):
        if not v or len(v.strip()) == 0:
            raise ValueError("Content cannot be empty")

        # Remove any leading/trailing whitespace
        v = v.strip()

        # Ensure it's within character limit
        if len(v) > 240:
            v = v[:237] + "..."
            logger.info(f"Content truncated to 240 characters during validation")

        return v


def get_cache_mode_from_env() -> CacheMode:
    """Get cache mode from environment variable AI_TRACKER_CACHE_MODE"""
    cache_mode_str = os.getenv("AI_TRACKER_CACHE_MODE", "ENABLED").upper()

    # Clean up the string (remove any extra text after the mode)
    cache_mode_str = cache_mode_str.split()[0] if cache_mode_str else "ENABLED"

    try:
        return CacheMode[cache_mode_str]
    except KeyError:
        logger.warning(f"Invalid cache mode '{cache_mode_str}', using ENABLED as default")
        return CacheMode.ENABLED


def parse_date_to_timestamp(date_str: str) -> int:
    """Parse various date formats to timestamp"""
    if not date_str:
        return int(datetime.now().timestamp())

    date_formats = [
        "%Y-%m-%d", "%Y-%m-%d %H:%M:%S", "%Y/%m/%d", "%m/%d/%Y", "%d/%m/%Y",
        "%B %d, %Y", "%b %d, %Y", "%d %B %Y", "%d %b %Y",
        "%Y-%m-%dT%H:%M:%S", "%Y-%m-%dT%H:%M:%SZ", "%Y-%m-%dT%H:%M:%S.%fZ"
    ]

    date_str = date_str.strip()
    for date_format in date_formats:
        try:
            parsed_date = datetime.strptime(date_str, date_format)
            return int(parsed_date.timestamp())
        except ValueError:
            continue

    # Fallback to current timestamp
    logger.warning(f"Could not parse date: {date_str}, using current timestamp")
    return int(datetime.now().timestamp())


def get_llm_extraction_instruction() -> str:
    """Get the LLM instruction for article extraction"""
    return """
        You are an expert at finding AI and technology articles on news websites.

        Look carefully through the entire webpage for article cards, headlines, or links that are related to:
        - Artificial Intelligence (AI)
        - Machine Learning (ML)
        - Automation
        - Robotics
        - Technology news that mentions AI/ML

        For each AI-related article you find, extract:
        1. The article title/headline
        2. The full URL to the article
        3. The featured image URL (if available)

        IMPORTANT INSTRUCTIONS:
        - Look for article cards, headlines, or links throughout the page
        - Check both visible content and any "load more" sections
        - If you see article titles but no images, still include them with empty image_url
        - Be thorough and look at all content on the page
        - Return ALL AI-related articles you can find, not just the first few

        Return a JSON object with this structure:
        {
            "articles": [
                {
                    "title": "Article Title Here",
                    "url": "https://example.com/full-article-url",
                    "image_url": "https://example.com/image-url.jpg"
                }
            ]
        }

        IMPORTANT INSTRUCTIONS:
        - Do not return a list of objects. Return only a single object with the "articles" array.
        - If you find no AI-related articles, return: {"articles": []}
        """


def parse_llm_response(llm_output) -> dict:
    """Parse and validate LLM response"""
    # Handle string responses
    if isinstance(llm_output, str):
        try:
            llm_output = json.loads(llm_output)
        except json.JSONDecodeError:
            logger.error(f"Failed to parse LLM response string: {llm_output}")
            return {"articles": []}

    # Handle list responses (merge articles from multiple objects)
    if isinstance(llm_output, list):
        if len(llm_output) > 0 and all(isinstance(item, dict) for item in llm_output):
            logger.warning("LLM returned a list of objects. Attempting to merge articles.")
            all_articles = []
            for item in llm_output:
                if isinstance(item, dict) and 'articles' in item and isinstance(item['articles'], list):
                    all_articles.extend(item['articles'])

            if all_articles:
                llm_output = {"articles": all_articles}
                logger.info(f"Merged {len(all_articles)} articles from list response")
            else:
                logger.warning("No articles found in list response")
                return {"articles": []}
        else:
            logger.error("Unexpected list response format")
            return {"articles": []}

    return llm_output


async def scrape_articles(website: Website, cache_mode: Optional[CacheMode] = None) -> List[dict]:
    """Extract articles from a website using LLM-powered structured data extraction"""
    if cache_mode is None:
        cache_mode = get_cache_mode_from_env()

    logger.info(f"Scraping {website.url} with cache mode: {cache_mode.value}")

    try:
        async with AsyncWebCrawler(verbose=True) as crawler:
            result = await crawler.arun(
                url=website.url,
                config=CrawlerRunConfig(
                    cache_mode=CacheMode.BYPASS,
                    scan_full_page=True,
                    wait_for_images=True,
                    wait_until="networkidle",
                    extraction_strategy=LLMExtractionStrategy(
                        llm_config=LLMConfig(
                            provider="openai/gpt-4o-mini",
                            api_token=os.getenv("OPENAI_API_KEY")
                        ),
                        schema=ArticleListExtraction.model_json_schema(),
                        extraction_type="schema",
                        instruction=get_llm_extraction_instruction()
                    ),
                )
            )

            if not result.success:
                logger.error(f"Crawling failed: {result.error_message}")
                return []

            if not result.extracted_content:
                logger.warning("No extracted content returned from LLM")
                return []

            logger.info(f"Raw LLM response: {result.extracted_content}")

            # Parse and validate the LLM response
            try:
                llm_output = parse_llm_response(result.extracted_content)
                extraction_result = ArticleListExtraction.model_validate(llm_output)

                # Convert to article dictionaries
                articles = []
                for article_info in extraction_result.articles:
                    absolute_url = urljoin(website.url, article_info.url) if article_info.url else None
                    absolute_image_url = urljoin(website.url, article_info.image_url) if article_info.image_url else None

                    article_dict = {
                        'title': article_info.title,
                        'url': absolute_url,
                        'image_url': absolute_image_url,
                        'source_website': website.name,
                    }
                    articles.append(article_dict)

                logger.info(f"Successfully extracted {len(articles)} articles using LLM")
                return articles

            except ValidationError as e:
                logger.error(f"Pydantic validation error: {e}")
                logger.error(f"Raw LLM response that failed validation: {result.extracted_content}")
                return []

    except Exception as e:
        logger.error(f"LLM extraction failed: {e}")
        return []


async def extract_article_details(url: str) -> Optional[DetailedArticleContent]:
    from crawl4ai import AsyncWebCrawler
    from crawl4ai.async_configs import CrawlerRunConfig, LLMConfig
    from crawl4ai.extraction_strategy import LLMExtractionStrategy

    async with AsyncWebCrawler(verbose=True) as crawler:
        result = await crawler.arun(
            url=url,
            config=CrawlerRunConfig(
                cache_mode=CacheMode.BYPASS,
                scan_full_page=True,
                wait_for_images=False,
                wait_until="domcontentloaded",
                extraction_strategy=LLMExtractionStrategy(
                    llm_config=LLMConfig(
                        provider="openai/gpt-4o-mini",
                        api_token=os.getenv("OPENAI_API_KEY")
                    ),
                    schema=DetailedArticleContent.model_json_schema(),
                    extraction_type="schema",
                    instruction="""
                        From the crawl content, extract the title, publication date and content summary.
                        Do not miss any information.
                        Return **ONE** article JSON format should look like this:
                        {
                            "title": "Article Title Here",
                            "content": "A very concise summary of the article with around 200 characters, **no more** than 240 characters.",
                            "publication_date": "The date and possible time the article was published at."
                        }
                        """
                ),
            )
        )
        if result.success and result.extracted_content:
            try:
                data = result.extracted_content

                # Log the raw LLM response for debugging
                logger.info(f"Raw LLM response for article details: {data}")

                if isinstance(data, str):
                    import json
                    data = json.loads(data)

                # Handle list responses (similar to main scraper)
                if isinstance(data, list):
                    logger.warning("LLM returned a list for article details. Using the first object only.")
                    data = data[0] if data else None
                    if not data:
                        logger.error("No valid objects in list response")
                        return None

                # Validate content length for single objects too
                if isinstance(data, dict) and data.get('content'):
                    if len(data['content']) > 240:
                        data['content'] = data['content'][:237] + "..."
                        logger.info(f"Truncated single object content from {len(data['content'])} to 240 characters")

                return DetailedArticleContent.model_validate(data)
            except Exception as e:
                logger.error(f"Failed to parse article details: {e}")
        else:
            logger.warning(f"Failed to extract details from {url}")
        return None
