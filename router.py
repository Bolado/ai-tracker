from fastapi import FastAPI, Request, Query
from fastapi.responses import HTMLResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from typing import Optional
import math
from datetime import datetime
from fastapi import HTTPException, status
import json

from database import get_articles_collection
from models import Article, Website

app = FastAPI(title="AI Tracker", description="AI News Tracker API")

# Mount static files
app.mount("/static", StaticFiles(directory="website/static"), name="static")

# Setup templates
templates = Jinja2Templates(directory="website/templates")

# Add custom Jinja2 filters
def timestamp_to_date(timestamp):
    """Convert timestamp to readable date"""
    try:
        dt = datetime.fromtimestamp(timestamp)
        return dt.strftime("%B %d, %Y")
    except:
        return "Unknown date"

templates.env.filters["timestamp_to_date"] = timestamp_to_date

ARTICLES_PER_PAGE = 20


def load_websites_from_json() -> list[dict]:
    """Load website configurations from websites.json file"""
    try:
        with open('websites.json', 'r') as f:
            websites_data = json.load(f)
        return websites_data
    except Exception as e:
        print(f"Failed to load websites from JSON: {e}")
        return []


def save_websites_to_json(websites_data: list[dict]):
    """Save website configurations to websites.json file"""
    try:
        with open('websites.json', 'w') as f:
            json.dump(websites_data, f, indent=2)
        return True
    except Exception as e:
        print(f"Failed to save websites to JSON: {e}")
        return False


@app.get("/", response_class=HTMLResponse)
async def read_root(request: Request, page: Optional[int] = Query(0, ge=0)):
    """
    Main page showing AI articles with pagination
    """
    try:
        collection = await get_articles_collection()

        # Calculate skip value
        skip = page * ARTICLES_PER_PAGE

        # Get articles with pagination
        cursor = collection.find({}).sort("timestamp", -1).skip(skip).limit(ARTICLES_PER_PAGE)
        articles = []
        async for doc in cursor:
            article = Article(**doc)
            articles.append(article)

        # Get total count for pagination
        total_articles = await collection.count_documents({})
        total_pages = math.ceil(total_articles / ARTICLES_PER_PAGE)

        return templates.TemplateResponse("index.html", {
            "request": request,
            "articles": articles,
            "current_page": page,
            "total_pages": total_pages,
            "has_previous": page > 0,
            "has_next": page < total_pages - 1,
            "previous_page": page - 1 if page > 0 else 0,
            "next_page": page + 1 if page < total_pages - 1 else page
        })

    except Exception as e:
        print(f"Error in read_root: {e}")
        return templates.TemplateResponse("index.html", {
            "request": request,
            "articles": [],
            "current_page": 0,
            "total_pages": 0,
            "has_previous": False,
            "has_next": False,
            "previous_page": 0,
            "next_page": 0
        })


@app.get("/api/articles")
async def get_articles(page: Optional[int] = Query(0, ge=0), limit: Optional[int] = Query(20, ge=1, le=100)):
    """
    API endpoint to get articles
    """
    try:
        collection = await get_articles_collection()

        skip = page * limit
        cursor = collection.find({}).sort("timestamp", -1).skip(skip).limit(limit)

        articles = []
        async for doc in cursor:
            article = Article(**doc)
            articles.append(article.dict())

        total_articles = await collection.count_documents({})

        return {
            "articles": articles,
            "page": page,
            "limit": limit,
            "total": total_articles,
            "total_pages": math.ceil(total_articles / limit)
        }

    except Exception as e:
        return {"error": str(e), "articles": [], "page": 0, "limit": limit, "total": 0, "total_pages": 0}


@app.get("/api/articles/{article_id}")
async def get_article(article_id: str):
    """
    Get a specific article by ID
    """
    try:
        from bson import ObjectId
        collection = await get_articles_collection()

        doc = await collection.find_one({"_id": ObjectId(article_id)})
        if doc:
            article = Article(**doc)
            return article.dict()
        else:
            return {"error": "Article not found"}

    except Exception as e:
        return {"error": str(e)}


@app.get("/api/websites")
async def get_websites():
    """
    Get all websites from the JSON file.
    """
    try:
        websites = load_websites_from_json()
        return {"websites": websites, "count": len(websites)}
    except Exception as e:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=f"Failed to load websites: {e}")


@app.post("/api/websites", status_code=status.HTTP_201_CREATED)
async def add_website(website: Website):
    """
    Add a new website to the JSON file.
    """
    try:
        # Load current websites
        websites = load_websites_from_json()

        # Check if website with same name or URL already exists
        for existing_website in websites:
            if existing_website.get('name') == website.name or existing_website.get('url') == website.url:
                raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="Website with this name or URL already exists.")

        # Create new website entry
        new_website = {
            "name": website.name,
            "url": website.url
        }

        # Add to list and save
        websites.append(new_website)
        if save_websites_to_json(websites):
            return {"message": "Website added successfully", "website": new_website}
        else:
            raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="Failed to save website to JSON file")

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=f"Failed to add website: {e}")


@app.get("/health")
async def health_check():
    """
    Health check endpoint
    """
    return {"status": "healthy", "service": "ai-tracker"}
