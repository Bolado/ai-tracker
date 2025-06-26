#!/usr/bin/env python3
"""
Integration test script for AI Tracker
Tests the complete system: database, scraper, watcher, and API
"""

import asyncio
import logging
import json
from datetime import datetime

from database import connect_to_mongo, close_mongo_connection, get_articles_collection
from watcher import load_words, load_websites_from_json
from scraper import scrape_articles, extract_article_details
from models import Website

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


async def test_database_connection():
    """Test database connection"""
    logger.info("🔌 Testing database connection...")
    try:
        await connect_to_mongo()
        collection = await get_articles_collection()
        count = await collection.count_documents({})
        logger.info(f"✅ Database connected. Found {count} existing articles")
        return True
    except Exception as e:
        logger.error(f"❌ Database connection failed: {e}")
        return False


def test_configuration_loading():
    """Test configuration files loading"""
    logger.info("📚 Testing configuration loading...")

    # Test words.json
    try:
        words = load_words()
        logger.info(f"✅ Loaded {len(words)} keywords from words.json")
    except Exception as e:
        logger.error(f"❌ Failed to load words.json: {e}")
        return False

    # Test websites.json
    try:
        websites = load_websites_from_json()
        logger.info(f"✅ Loaded {len(websites)} websites from websites.json")
        for website in websites:
            logger.info(f"   🌐 {website.name}: {website.url}")
    except Exception as e:
        logger.error(f"❌ Failed to load websites.json: {e}")
        return False

    return True


async def test_scraper():
    """Test scraper functionality"""
    logger.info("🔍 Testing scraper functionality...")

    try:
        websites = load_websites_from_json()
        if not websites:
            logger.warning("⚠️ No websites configured, skipping scraper test")
            return True

        # Test with first website
        website = websites[0]  # This is already a Website object

        logger.info(f"📰 Testing scraper with: {website.name} ({website.url})")
        articles = await scrape_articles(website)

        logger.info(f"✅ Scraper found {len(articles)} articles")
        for i, article in enumerate(articles[:3], 1):  # Show first 3
            logger.info(f"   📄 {i}. {article['title'][:60]}...")

        return True

    except Exception as e:
        logger.error(f"❌ Scraper test failed: {e}")
        return False


async def test_article_details_extraction():
    """Test article details extraction"""
    logger.info("🔍 Testing article details extraction...")

    try:
        # Use a test URL (you can change this to a real AI article)
        test_url = "https://www.wired.com/story/uncanny-valley-podcast-superintelligence/"

        logger.info(f"📝 Testing details extraction for: {test_url}")
        details = await extract_article_details(test_url)

        if details:
            logger.info(f"✅ Article details extracted successfully:")
            logger.info(f"   📰 Title: {details.title}")
            logger.info(f"   📝 Content: {details.content[:100]}...")
            logger.info(f"   📅 Date: {details.publication_date}")
            logger.info(f"   📏 Content length: {len(details.content)} characters")
        else:
            logger.warning("⚠️ No article details extracted")

        return True

    except Exception as e:
        logger.error(f"❌ Article details extraction failed: {e}")
        return False


async def main():
    """Run all integration tests"""
    logger.info("🧪 Starting AI Tracker Integration Tests")
    logger.info("=" * 50)

    tests = [
        ("Database Connection", test_database_connection),
        ("Configuration Loading", lambda: test_configuration_loading()),
        ("Scraper Functionality", test_scraper),
        ("Article Details Extraction", test_article_details_extraction),
    ]

    results = {}

    for test_name, test_func in tests:
        logger.info(f"\n🔬 Running: {test_name}")
        logger.info("-" * 30)

        try:
            if asyncio.iscoroutinefunction(test_func):
                result = await test_func()
            else:
                result = test_func()

            results[test_name] = result
            status = "✅ PASSED" if result else "❌ FAILED"
            logger.info(f"{status}: {test_name}")

        except Exception as e:
            logger.error(f"❌ FAILED: {test_name} - {e}")
            results[test_name] = False

    # Summary
    logger.info("\n" + "=" * 50)
    logger.info("📊 INTEGRATION TEST SUMMARY")
    logger.info("=" * 50)

    passed = sum(1 for result in results.values() if result)
    total = len(results)

    for test_name, result in results.items():
        status = "✅ PASSED" if result else "❌ FAILED"
        logger.info(f"{status}: {test_name}")

    logger.info(f"\n🎯 Overall: {passed}/{total} tests passed")

    if passed == total:
        logger.info("🎉 All tests passed! System is ready to run.")
    else:
        logger.warning("⚠️ Some tests failed. Please check the issues above.")

    # Cleanup
    try:
        await close_mongo_connection()
        logger.info("🔌 Database connection closed")
    except:
        pass


if __name__ == "__main__":
    asyncio.run(main())
