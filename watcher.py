import json
import os
import asyncio
import logging
from typing import List, Dict, Any
from datetime import datetime

from models import Article, Website, ArticlesListItem
from database import get_articles_collection
from ai_service import summarize, classify_article
from scraper import scrape_articles, extract_article_details

logger = logging.getLogger(__name__)

# Global variables
articles: List[Article] = []
words: List[str] = []
websites: List[Website] = []


async def load_articles():
    """Load articles from database"""
    global articles
    try:
        collection = await get_articles_collection()
        cursor = collection.find({})
        articles = []
        async for doc in cursor:
            article = Article(**doc)
            articles.append(article)
        logger.info(f"Loaded {len(articles)} articles from database")
    except Exception as e:
        logger.error(f"Failed to load articles: {e}")
        raise


def load_words() -> List[str]:
    """Load AI-related keywords from words.json"""
    global words
    try:
        with open('words.json', 'r') as f:
            words = json.load(f)
        logger.info(f"Loaded {len(words)} keywords")
        return words
    except Exception as e:
        logger.error(f"Failed to load words: {e}")
        return []


def load_websites_from_json() -> List[Website]:
    """Load website configurations from websites.json file"""
    try:
        with open('websites.json', 'r') as f:
            websites_data = json.load(f)

        websites_list = []
        for website_data in websites_data:
            # Create Website object with only name and url
            website = Website(
                name=website_data['name'],
                url=website_data['url']
            )
            websites_list.append(website)

        logger.info(f"Loaded {len(websites_list)} websites from websites.json")
        return websites_list
    except Exception as e:
        logger.error(f"Failed to load websites from JSON: {e}")
        return []


async def load_websites() -> List[Website]:
    """Load website configurations from websites.json file"""
    global websites
    try:
        websites = load_websites_from_json()
        logger.info(f"Loaded {len(websites)} websites from websites.json")
        return websites
    except Exception as e:
        logger.error(f"Failed to load websites: {e}")
        return []


async def article_exists_in_db(article_url: str) -> bool:
    """Check if an article already exists in the database"""
    try:
        collection = await get_articles_collection()
        existing = await collection.find_one({"link": article_url})
        return existing is not None
    except Exception as e:
        logger.error(f"Error checking if article exists: {e}")
        return False


async def process_article(article_data: dict, website: Website) -> bool:
    """Process a single article: extract details, classify, and save"""
    try:
        # Check if article already exists
        if await article_exists_in_db(article_data['url']):
            logger.info(f"Article already exists: {article_data['title']}")
            return False

        # Extract details using crawl4ai LLM extraction
        logger.info(f"Extracting details for: {article_data['title']}")
        details = await extract_article_details(article_data['url'])
        if not details or not details.content:
            logger.warning(f"No content found for: {article_data['title']}")
            return False

        # Classify article (optional, can use details.content)
        if not await classify_article(details.content):
            logger.info(f"Article not AI-related: {article_data['title']}")
            return False

        # Create Article object
        article = Article(
            title=details.title or article_data['title'],
            summary=details.content,  # Already summarized by LLM
            link=article_data['url'],
            timestamp=int(datetime.now().timestamp()),
            source=article_data['source_website'],
            image=article_data.get('image_url', ''),
            content=details.content
        )
        await save_article(article)
        articles.append(article)
        logger.info(f"Successfully processed article: {article_data['title']}")
        return True

    except Exception as e:
        logger.error(f"Error processing article {article_data['title']}: {e}")
        return False


async def save_article(article: Article):
    """Save article to database"""
    try:
        collection = await get_articles_collection()

        # Convert to dict and insert
        article_dict = article.model_dump(by_alias=True)
        if '_id' in article_dict and article_dict['_id'] is None:
            del article_dict['_id']

        result = await collection.insert_one(article_dict)
        logger.info(f"Saved article: {article.title}")

    except Exception as e:
        logger.error(f"Error saving article: {e}")


async def start_watcher():
    """Start the website watcher process"""
    logger.info("Starting watcher...")

    try:
        # Load configuration
        load_words()
        await load_websites()

        # Process each website
        for website in websites:
            try:
                logger.info(f"Processing website: {website.name} ({website.url})")

                # Scrape articles using LLM-powered extraction
                scraped_articles = await scrape_articles(website)
                logger.info(f"Found {len(scraped_articles)} articles from {website.name}")

                # Process each scraped article
                processed_count = 0
                for article_data in scraped_articles:
                    try:
                        if await process_article(article_data, website):
                            processed_count += 1

                        # Rate limiting between articles
                        await asyncio.sleep(1)

                    except Exception as e:
                        logger.error(f"Error processing article {article_data.get('title', 'Unknown')}: {e}")
                        continue

                logger.info(f"Processed {processed_count} new articles from {website.name}")

                # Rate limiting between websites
                await asyncio.sleep(2)

            except Exception as e:
                logger.error(f"Error processing website {website.name}: {e}")
                continue

        logger.info("Watcher completed")

    except Exception as e:
        logger.error(f"Watcher failed: {e}")
        raise
