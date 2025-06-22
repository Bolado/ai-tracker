import os
import asyncio
from typing import List
from openai import OpenAI
import logging

logger = logging.getLogger(__name__)

# Placeholder for AI functionality
# TODO: Implement AI summarization using OpenAI API

async def summarize(text: str) -> str:
    """
    Summarize text using OpenAI API

    Args:
        text: The text content to summarize

    Returns:
        str: Summarized text
    """
    try:
        api_key = os.getenv("OPENAI_API_KEY")
        if not api_key:
            logger.warning("OPENAI_API_KEY not set, using fallback summarization")
            return text[:200] + "..." if len(text) > 200 else text

        client = OpenAI(api_key=api_key)

        # Truncate text if too long for API
        if len(text) > 4000:
            text = text[:4000] + "..."

        response = client.chat.completions.create(
            model="gpt-4o-mini",
            messages=[
                {
                    "role": "system",
                    "content": "You are a helpful assistant that summarizes news articles about AI and technology. Provide concise, informative summaries in 2-3 sentences."
                },
                {
                    "role": "user",
                    "content": f"Please summarize this article: {text}"
                }
            ],
            max_tokens=150,
            temperature=0.7
        )

        summary = response.choices[0].message.content.strip()
        logger.info(f"Generated summary: {summary[:100]}...")
        return summary

    except Exception as e:
        logger.error(f"Error summarizing with OpenAI: {e}")
        # Fallback to simple truncation
        return text[:200] + "..." if len(text) > 200 else text


async def classify_article(text: str) -> bool:
    """
    Classify if an article is AI-related using OpenAI API

    Args:
        text: The article text to classify

    Returns:
        bool: True if article is AI-related, False otherwise
    """
    try:
        api_key = os.getenv("OPENAI_API_KEY")
        if not api_key:
            logger.warning("OPENAI_API_KEY not set, using keyword-based classification")
            return _keyword_classification(text)

        client = OpenAI(api_key=api_key)

        # Truncate text if too long for API
        if len(text) > 4000:
            text = text[:4000] + "..."

        response = client.chat.completions.create(
            model="gpt-4o-mini",
            messages=[
                {
                    "role": "system",
                    "content": "You are an expert at identifying AI and technology-related content. Respond with 'YES' if the article is about AI, machine learning, automation, robotics, or related technology topics. Respond with 'NO' otherwise."
                },
                {
                    "role": "user",
                    "content": f"Is this article about AI or technology? Article: {text}"
                }
            ],
            max_tokens=10,
            temperature=0.1
        )

        result = response.choices[0].message.content.strip().upper()
        is_ai_related = "YES" in result

        logger.info(f"AI classification result: {result} -> {is_ai_related}")
        return is_ai_related

    except Exception as e:
        logger.error(f"Error classifying with OpenAI: {e}")
        # Fallback to keyword-based classification
        return _keyword_classification(text)


def _keyword_classification(text: str) -> bool:
    """Fallback keyword-based classification"""
    ai_keywords = [
        "artificial intelligence", "ai", "machine learning", "neural network",
        "openai", "chatgpt", "llm", "deep learning", "computer vision",
        "automation", "robotics", "algorithm", "data science", "big data",
        "natural language processing", "nlp", "computer vision", "autonomous",
        "smart", "intelligent", "automated", "digital transformation"
    ]

    text_lower = text.lower()
    return any(keyword in text_lower for keyword in ai_keywords)


# Example of how to implement the real OpenAI integration:
"""
async def summarize_with_openai(text: str) -> str:
    try:
        token = os.getenv("OPENAI_API_TOKEN")
        if not token:
            raise ValueError("OPENAI_API_TOKEN not set")

        client = OpenAI(api_key=token)

        response = client.chat.completions.create(
            model="gpt-3.5-turbo",
            messages=[
                {
                    "role": "system",
                    "content": "You are a helpful assistant that summarizes news articles about AI and technology."
                },
                {
                    "role": "user",
                    "content": f"Please summarize this article in 2-3 sentences: {text}"
                }
            ],
            max_tokens=150,
            temperature=0.7
        )

        return response.choices[0].message.content.strip()

    except Exception as e:
        print(f"Error summarizing with OpenAI: {e}")
        return text[:200] + "..." if len(text) > 200 else text
"""
