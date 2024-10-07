package watcher

import (
	"encoding/json"
	"os"

	database "github.com/Bolado/ai-tracker/database"
	types "github.com/Bolado/ai-tracker/types"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

var (
	Articles []types.Article
)

func StartWatcher() error {
	//read ./websites to have

	browser, err := startRod()
	if err != nil {
		return err
	}

	websites, err := readWebsitesJSON()
	if err != nil {
		return err
	}

	for _, website := range websites {
		page, err := browser.Page(proto.TargetCreateTarget{URL: website.Url})
		if err != nil {
			return err
		}
		page.WaitLoad()
	}

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

func startRod() (*rod.Browser, error) {
	launcher := launcher.New()
	launcher.Headless(true)

	url, err := launcher.Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(url)

	err = browser.Connect()
	if err != nil {
		return nil, err
	}

	return browser, nil
}

// read websites from ./websites.json and put into a []types.Website
func readWebsitesJSON() ([]types.Website, error) {
	var websites []types.Website
	file, err := os.ReadFile("./websites.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &websites)
	if err != nil {
		return nil, err
	}

	return websites, nil
}
