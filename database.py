import os
import logging
from motor.motor_asyncio import AsyncIOMotorClient
from typing import Optional

logger = logging.getLogger(__name__)

class Database:
    client: Optional[AsyncIOMotorClient] = None
    database = None
    collection = None

db = Database()

async def connect_to_mongo():
    """Create database connection"""
    try:
        mongo_uri = os.getenv("MONGO_URI")
        if not mongo_uri:
            raise ValueError("MONGO_URI is not set")

        db.client = AsyncIOMotorClient(mongo_uri)
        db.database = db.client.get_database("ai_tracker")
        db.articles_collection = db.database.articles
        db.websites_collection = db.database.websites # New collection for websites

        # Test the connection
        await db.client.admin.command('ismaster')
        logger.info("Connected to MongoDB")

    except Exception as e:
        logger.error(f"Failed to connect to MongoDB: {e}")
        raise


async def close_mongo_connection():
    """Close database connection"""
    if db.client:
        db.client.close()
        logger.info("Disconnected from MongoDB")


async def get_database():
    return db.database


async def get_articles_collection():
    return db.articles_collection


async def get_websites_collection():
    return db.websites_collection
