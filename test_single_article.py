#!/usr/bin/env python3
"""
Test script for single article details extraction
"""
import asyncio
import os
import logging
from dotenv import load_dotenv

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler()
    ]
)

# Load environment variables
load_dotenv()

from database import connect_to_mongo, close_mongo_connection
from scraper import extract_article_details

async def test_single_article():
    """Test article details extraction for a single article"""
    print("🧪 Testing Single Article Details Extraction")
    print("=" * 50)

    # Test URL
    test_url = "https://www.wired.com/story/uncanny-valley-podcast-superintelligence/"

    try:
        # Connect to database
        print("\n🔌 Connecting to database...")
        await connect_to_mongo()
        print("✅ Database connected")

        print(f"\n📰 Testing article details extraction for: {test_url}")
        details = await extract_article_details(test_url)

        if details:
            print("\n✅ Article details extracted successfully!")
            print(f"Title: {details.title}")
            print(f"Content: {details.content}")
            print(f"Publication Date: {details.publication_date}")
        else:
            print("\n❌ Failed to extract article details")

    except Exception as e:
        print(f"\n❌ Test failed: {e}")
        raise
    finally:
        # Close database connection
        await close_mongo_connection()
        print("🔌 Database connection closed")

if __name__ == "__main__":
    asyncio.run(test_single_article())
