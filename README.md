# AI Tracker

This project aims to give users a quick, easy way to get the gist of news articles without needing to read the full text. For those interested in diving deeper, there’s also an option to navigate to the original article.

The app brings together a range of tools—like Go, HTMX, TailwindCSS, Docker, and OpenAI’s API—to make it all happen.

**The project can be found running [in here](https://aitracker.news).**

## Delivery

The project is delivered as a Docker container, which the image is made through GitHub Actions. The image is then pushed to GitHub Container Registry and can be pulled from there.

## Goals

| Feature             | Status |
|---------------------|--------|
| Golang Backend      | ✅     |
| Scraping            | ✅     |
| TailwindCSS         | ✅     |
| AI Summarization    | ✅     |
| Docker              | ✅     |
| GitHub Actions      | ✅    |
| Content Pagination  | ✅     |
| RSS Feed            | ⏳     |
| Webhooks            | ⏳     |

## Want to contribute?

Searching and adding new sources for articles is a constant task, if you are interested in contributing to that, sources can be added to the ```words.json``` file. ❤

## Links

[Source Code](https://github.com/Bolado/ai-tracker)
