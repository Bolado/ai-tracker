# AI Tracker - LLM-Powered News Monitoring System

A sophisticated AI news tracking system that uses **crawl4ai** for web scraping and **LLM-powered structured data extraction** to monitor AI-related news from multiple sources.

## 🚀 Features

- **LLM-Powered Scraping**: Uses GPT-4o-mini to intelligently extract articles from news websites
- **Smart Content Extraction**: Automatically extracts article titles, summaries, publication dates, and images
- **AI Classification**: Filters articles to ensure they're AI-related using LLM classification
- **Robust Error Handling**: Comprehensive logging and error recovery
- **Web Interface**: Beautiful web UI to browse and search articles
- **REST API**: Full API for programmatic access
- **Scheduled Monitoring**: Runs every hour to check for new articles
- **Duplicate Prevention**: Prevents saving duplicate articles

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Interface │    │   REST API      │    │   Watcher       │
│   (FastAPI)     │    │   (FastAPI)     │    │   (Scheduler)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Database      │
                    │   (MongoDB)     │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Scraper       │
                    │   (crawl4ai)    │
                    └─────────────────┘
```

## 📋 Requirements

- Python 3.8+
- MongoDB
- OpenAI API key
- crawl4ai

## 🛠️ Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd ai-tracker
   ```

2. **Install dependencies**:
   ```bash
   pip install -r requirements.txt
   ```

3. **Set up environment variables**:
   Create a `.env` file:
   ```env
   MONGODB_URI=mongodb://localhost:27017
   OPENAI_API_KEY=your_openai_api_key_here
   AI_TRACKER_CACHE_MODE=ENABLED
   WATCHER_INTERVAL=60
   ```

4. **Configure websites**:
   Edit `websites.json` to add your target news sources:
   ```json
   [
     {
       "name": "Wired",
       "url": "https://www.wired.com/tag/artificial-intelligence/"
     },
     {
       "name": "TechCrunch",
       "url": "https://techcrunch.com/tag/artificial-intelligence/"
     }
   ]
   ```

5. **Start MongoDB**:
   ```bash
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:latest

   # Or install MongoDB locally
   ```

## 🚀 Usage

### Start the Application

```bash
# Start with watcher (default)
python main.py

# Start without watcher (web server only)
python main.py --disable-watcher

# Custom port
python main.py --port 8080 --host 0.0.0.0
```

### Run Integration Tests

```bash
# Test the complete system
python test_integration.py

# Test single article extraction
python test_single_article.py

# Test watcher functionality
python test_watcher.py
```

### Access the Web Interface

- **Main page**: http://localhost:8080
- **API documentation**: http://localhost:8080/docs
- **Health check**: http://localhost:8080/health

## 📊 API Endpoints

### Articles
- `GET /api/articles` - Get all articles with pagination
- `GET /api/articles/{article_id}` - Get specific article
- `GET /` - Web interface with pagination

### Websites
- `GET /api/websites` - Get configured websites
- `POST /api/websites` - Add new website

### Health
- `GET /health` - Health check endpoint

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `OPENAI_API_KEY` | OpenAI API key | Required |
| `AI_TRACKER_CACHE_MODE` | Cache mode (ENABLED/DISABLED/BYPASS) | `ENABLED` |
| `WATCHER_INTERVAL` | Watcher interval in minutes | `60` |

### Configuration Files

#### `websites.json`
Configure news sources to monitor:
```json
[
  {
    "name": "Website Name",
    "url": "https://example.com/ai-news"
  }
]
```

#### `words.json`
AI-related keywords for classification:
```json
[
  "artificial intelligence",
  "machine learning",
  "AI",
  "ML",
  "automation"
]
```

## 📈 Monitoring & Logging

The system provides comprehensive logging:

- **File logging**: `ai-tracker.log`
- **Console logging**: Real-time status updates
- **Emoji indicators**: Easy-to-read status messages
- **Statistics**: Detailed performance metrics

### Log Levels
- `INFO`: Normal operation
- `WARNING`: Non-critical issues
- `ERROR`: Critical errors

## 🔍 How It Works

1. **Scheduled Monitoring**: Watcher runs every hour
2. **Website Scraping**: Uses crawl4ai to extract article lists
3. **Content Extraction**: LLM extracts detailed article information
4. **Classification**: Filters for AI-related content
5. **Storage**: Saves to MongoDB with deduplication
6. **Web Interface**: Displays articles with pagination

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Ensure MongoDB is running
   - Check `MONGODB_URI` in `.env`

2. **OpenAI API Errors**
   - Verify `OPENAI_API_KEY` is set
   - Check API key validity and quota

3. **No Articles Found**
   - Check `websites.json` configuration
   - Verify website URLs are accessible
   - Check logs for scraping errors

4. **Watcher Not Running**
   - Check logs for errors
   - Verify `WATCHER_INTERVAL` setting
   - Use `--disable-watcher` flag to test web interface only

### Debug Mode

Enable detailed logging:
```bash
export PYTHONPATH=.
python -c "import logging; logging.basicConfig(level=logging.DEBUG)"
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- **crawl4ai**: For powerful web scraping capabilities
- **OpenAI**: For LLM-powered content extraction
- **FastAPI**: For the web framework
- **MongoDB**: For data storage
