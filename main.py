import asyncio
import logging
import os
import argparse
from contextlib import asynccontextmanager

import uvicorn
from dotenv import load_dotenv
from fastapi import FastAPI

from database import connect_to_mongo, close_mongo_connection, get_articles_collection
from watcher import start_watcher, load_articles
import router

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('ai-tracker.log'),
        logging.StreamHandler()
    ]
)

logger = logging.getLogger(__name__)

# Global variables
disable_watcher = False
watcher_interval = 60  # minutes (1 hour)


@asynccontextmanager
async def lifespan(app):
    """Application lifespan management"""
    # Startup
    logger.info("Starting AI Tracker application...")

    # Load environment variables
    load_dotenv()

    # Get watcher interval from environment
    global watcher_interval
    try:
        watcher_interval = int(os.getenv("WATCHER_INTERVAL", "60"))
    except ValueError:
        watcher_interval = 60

    # Connect to database
    await connect_to_mongo()
    logger.info("Database connected")

    # Load articles from database
    try:
        # Ensure the articles collection is initialized by calling it once
        await get_articles_collection()
        await load_articles()
        logger.info("Articles loaded")
    except Exception as e:
        logger.error(f"Failed to load articles: {e}")

    # Start watcher if not disabled
    if not disable_watcher:
        asyncio.create_task(watcher_loop())
        logger.info("Watcher started")

    yield

    # Shutdown
    logger.info("Shutting down AI Tracker...")
    await close_mongo_connection()


async def watcher_loop():
    """Run the watcher in a loop with specified interval"""
    first_run = True

    while True:
        try:
            if first_run:
                logger.info("🚀 Starting initial watcher run...")
                first_run = False
            else:
                logger.info(f"⏰ Running scheduled watcher (every {watcher_interval} minutes)...")

            await start_watcher()

            if first_run:
                logger.info("✅ Initial watcher run completed successfully")
            else:
                logger.info(f"✅ Watcher completed. Next run in {watcher_interval} minutes...")

            # Sleep until next run
            await asyncio.sleep(watcher_interval * 60)

        except Exception as e:
            logger.error(f"❌ Watcher failed: {e}")
            logger.info(f"🔄 Retrying in {watcher_interval} minutes...")
            await asyncio.sleep(watcher_interval * 60)


# Create FastAPI app with lifespan
app = FastAPI(title="AI Tracker", description="AI News Tracker API", lifespan=lifespan)

# Include routes from router.py
from router import app as router_app

# Copy all routes from router app
for route in router_app.routes:
    app.routes.append(route)


def main():
    """Main application entry point"""
    global disable_watcher

    # Parse command line arguments
    parser = argparse.ArgumentParser(description="AI Tracker - Monitor AI-related news")
    parser.add_argument("--disable-watcher", action="store_true",
                       help="Disable the watcher (web server only)")
    parser.add_argument("--port", type=int, default=8080,
                       help="Port to run the web server on (default: 8080)")
    parser.add_argument("--host", type=str, default="0.0.0.0",
                       help="Host to bind the web server to (default: 0.0.0.0)")

    args = parser.parse_args()
    disable_watcher = args.disable_watcher

    if disable_watcher:
        logger.info("Watcher disabled - running web server only")

    # Run the application
    uvicorn.run(
        "main:app",
        host=args.host,
        port=args.port,
        reload=False,
        log_level="info"
    )


if __name__ == "__main__":
    main()
