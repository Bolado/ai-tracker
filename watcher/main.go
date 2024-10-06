package watcher

import (
	database "github.com/Bolado/ai-tracker/database"
	types "github.com/Bolado/ai-tracker/types"
)

var (
	Articles []types.Article
)

func StartWatcher() error {
	return nil
}

func LoadArticles() error {
	var err error
	Articles, err = database.GetArticles()
	if err != nil {
		return err
	}
	return nil
}

func isExistant(article types.Article) bool {
	for _, a := range Articles {
		if a.Link == article.Link {
			return true
		}
	}
	return false
}

func addArticle(article types.Article) error {

	if isExistant(article) {
		return nil
	}

	if err := database.InsertArticle(article); err != nil {
		return err
	}

	Articles = append(Articles, article)

	return nil
}
