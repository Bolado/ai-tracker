#!/usr/bin/env python3
"""
Demo script for crawl4ai LLM-powered structured data extraction with caching
"""
import asyncio
import json
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

from scraper import scrape_articles, get_cache_mode_from_env
from crawl4ai import CacheMode
from models import Website

def load_websites_from_json() -> list[Website]:
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

        print(f"✅ Loaded {len(websites_list)} websites from websites.json")
        return websites_list
    except Exception as e:
        print(f"❌ Failed to load websites from JSON: {e}")
        return []

async def demo_llm_extraction_with_caching():
    """Demo LLM-powered extraction with different cache modes"""
    print("🤖 Testing LLM-powered article extraction with crawl4ai")
    print("=" * 60)

    # Load websites from JSON file
    websites = load_websites_from_json()
    if not websites:
        print("❌ No websites loaded. Exiting.")
        return

    # Use the first website for demo (or you can loop through all)
    test_website = websites[0]  # Use the first website from the list

    try:
        print(f"📰 Extracting articles from: {test_website.url}")
        print("🧠 Using LLM-powered structured extraction...")

        articles_found = await scrape_articles(test_website, cache_mode=CacheMode.ENABLED)
        print(f"✅ Found {len(articles_found)} articles.")

        # print the articles
        for article in articles_found:
            print("="*60)
            print(f"Title: {article['title']}")
            print(f"URL: {article['url']}")
            print(f"Image: {article['image_url']}")
            print(f"Source: {article['source_website']}")
            print("="*60)
        print("="*60)

    except Exception as e:
        print(f"❌ Error: {e}")
        print("💡 Make sure OPENAI_API_KEY is set in your .env file")

async def main():
    """Run LLM-only demo with caching"""
    print("🚀 crawl4ai LLM Extraction Demo with Caching")
    print("🤖 Modern structured data extraction for AI article tracking")

    # Check if OpenAI API key is available
    api_key = os.getenv('OPENAI_API_KEY')
    if not api_key:
        print("⚠️  Warning: OPENAI_API_KEY not found in environment")
        print("   Demo will fail - please set your OpenAI API key")
        return
    else:
        print(f"✅ OpenAI API key found: {api_key[:10]}...")

    # Show current cache setting
    cache_mode = os.getenv('AI_TRACKER_CACHE_MODE', 'ENABLED')
    print(f"🗄️ Current cache mode setting: {cache_mode}")

    try:
        # Run LLM demo with caching
        await demo_llm_extraction_with_caching()

        print("\n" + "=" * 60)
        print("🎉 Demo completed!")

    except KeyboardInterrupt:
        print("\n🛑 Demo interrupted by user")
    except Exception as e:
        print(f"\n❌ Demo failed: {e}")

if __name__ == "__main__":
    asyncio.run(main())
