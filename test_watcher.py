#!/usr/bin/env python3
"""
Test script for the updated watcher system
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
        logging.StreamHandler()  # Output to console
    ]
)

# Load environment variables
load_dotenv()

from database import connect_to_mongo, close_mongo_connection
from watcher import start_watcher

async def test_watcher():
    """Test the watcher system"""
    print("🧪 Testing AI Tracker Watcher System")
    print("=" * 50)

    # Check if OpenAI API key is available
    api_key = os.getenv('OPENAI_API_KEY')
    if not api_key:
        print("⚠️  Warning: OPENAI_API_KEY not found in environment")
        print("   The system will use fallback summarization and classification")
    else:
        print(f"✅ OpenAI API key found: {api_key[:10]}...")

    # Check if MongoDB URI is available
    mongo_uri = os.getenv('MONGO_URI')
    if not mongo_uri:
        print("❌ Error: MONGO_URI not found in environment")
        print("   Please set MONGO_URI in your .env file")
        return
    else:
        print(f"✅ MongoDB URI found: {mongo_uri[:20]}...")

    try:
        # Connect to database
        print("\n🔌 Connecting to database...")
        await connect_to_mongo()
        print("✅ Database connected")

        print("\n🚀 Starting watcher test...")
        await start_watcher()
        print("\n✅ Watcher test completed successfully!")

    except Exception as e:
        print(f"\n❌ Watcher test failed: {e}")
        raise
    finally:
        # Close database connection
        await close_mongo_connection()
        print("🔌 Database connection closed")

if __name__ == "__main__":
    asyncio.run(test_watcher())
