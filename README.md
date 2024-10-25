# AI Tracker

This repository contains the code for the "AI Tracker" project.

The goal of this project is to provide users with a quick way to get the gist of a news article without having to read the entire thing while also providing a link to the original article for those who want to read more.

The application will use a mix of tools to achieve this goal, examples are Go, HTMX, TailwindCSS, Docker and OpenAI's API.

## Delivery

The project is delivered as a Docker container, which the image is made through GitHub Actions. The image is then pushed to GitHub Container Registry and can be pulled from there.

## Goals

- [x] Golang Backend
- [x] Scraping
- [x] TailwindCSS
- [x] AI Summarization
- [ ] Content Pagination
- [ ] Docker
- [ ] GitHub Actions
- [ ] RSS Feed
- [ ] Webhooks
