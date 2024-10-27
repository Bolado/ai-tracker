package templates

import (
	"github.com/Bolado/ai-tracker/types"
	"github.com/Bolado/ai-tracker/watcher"
)

var (
	articlesPerPage = 20
)

func GetNumberOfPages() int {
	totalArticles := len(watcher.Articles)
	totalPages := (totalArticles + articlesPerPage - 1) / articlesPerPage // Calculate total pages

	return totalPages
}

func GetPagedArticles(page int) []types.Article {
	totalArticles := len(watcher.Articles)
	totalPages := GetNumberOfPages()

	// If requested page is greater than total pages, set page to the last page
	if page > totalPages-1 {
		page = totalPages - 1
	}

	start := page * articlesPerPage
	end := start + articlesPerPage

	// Ensure end does not exceed the length of the articles array
	if end > totalArticles {
		end = totalArticles
	}

	// Slice the articles array to get only the articles for the current page
	pagedArticles := watcher.Articles[start:end]

	return pagedArticles
}
