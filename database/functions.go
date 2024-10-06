package database

import (
	"context"

	"github.com/Bolado/ai-tracker/types"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertArticle(article types.Article) error {
	_, err := collection.InsertOne(context.TODO(), article)
	if err != nil {
		return err
	}
	return nil
}

func GetArticles() ([]types.Article, error) {
	var articles []types.Article

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &articles); err != nil {
		return nil, err
	}

	return articles, nil
}
