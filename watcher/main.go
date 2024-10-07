package watcher

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	database "github.com/Bolado/ai-tracker/database"
	types "github.com/Bolado/ai-tracker/types"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

var (
	Articles []types.Article
	words    []string
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

	words, err = readWordsJSON()
	if err != nil {
		return err
	}

	for _, website := range websites {
		err := watchWebsite(website, browser)
		if err != nil {
			return err
		}
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

func isExistant(link string) bool {
	for _, a := range Articles {
		if a.Link == link {
			return true
		}
	}
	return false
}

func addArticle(article types.Article) error {

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

// read the words from ./words.json and put into a []string
func readWordsJSON() ([]string, error) {
	var words []string
	file, err := os.ReadFile("./words.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &words)
	if err != nil {
		return nil, err
	}
	return words, nil
}

func watchWebsite(website types.Website, browser *rod.Browser) error {
	page, err := browser.Page(proto.TargetCreateTarget{URL: website.Url})
	if err != nil {
		return err
	}

	el, err := page.Timeout(10 * time.Second).ElementsX(website.MainElement)
	if err != nil {
		return err
	}

articlesLoop:
	for _, e := range el {
		articleTitleElement, err := e.ElementX(website.TitleElement)
		if err != nil {
			return err
		}

		articleTitle := articleTitleElement.MustText()

		for _, word := range words {
			// check if article title contains relevant word
			if strings.Contains(articleTitle, word) {

				//get anchor element and then url of the article
				urlElement, err := e.ElementX(website.AnchorElement)
				if err != nil {
					return err
				}
				url := urlElement.MustProperty("href").String()

				// if its already in the array, continue the loop
				if isExistant(url) {
					continue articlesLoop
				}
			}
		}
	}

	return nil
}
