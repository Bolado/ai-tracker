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

// watch the websites provided for new articles, summarizing and adding only the non existant ones to the database and the struct array
func watchWebsite(website types.Website, browser *rod.Browser) error {
	//navigate to the website url
	page, err := browser.Page(proto.TargetCreateTarget{URL: website.Url})
	if err != nil {
		return err
	}
	//wait for 10 seconds for the articles elements to load on the page
	els, err := page.Timeout(10 * time.Second).ElementsX(website.MainElement)
	if err != nil {
		return err
	}

	// loop through the elements found, and navigate to each relevant article that is not existant, then do the necessary tasks
articlesLoop:
	for _, e := range els {

		//get article title element
		articleTitleElement, err := e.ElementX(website.TitleElement)
		if err != nil {
			return err
		}

		//get the article title text
		articleTitle := articleTitleElement.MustText()

		//check for relevant words on the title itself
		for _, word := range words {
			// check if article title contains relevant word
			if strings.Contains(articleTitle, word) {
				// analyze the article further
				if existant, added, err := analyzeArticle(e, page, website, articleTitle); err != nil {
					return err
				} else {
					if existant || added {
						continue articlesLoop
					}
				}

				break
			}
		}
	}

	return nil
}

// analyze article
func analyzeArticle(e *rod.Element, page *rod.Page, website types.Website, title string) (existent bool, added bool, err error) {
	var article types.Article
	//get anchor element and then url of the article
	urlElement, err := e.ElementX(website.AnchorElement)
	if err != nil {
		return false, false, err
	}
	url := urlElement.MustProperty("href").String()

	// if its already in the array, continue the loop
	if isExistant(url) {
		return true, false, nil
	}

	//get the article image url if provided
	if website.ImageElement != "" {
		imageElement, err := e.ElementX(website.ImageElement)
		if err != nil {
			return false, false, err
		}
		imgURL := imageElement.MustProperty("src").String()
		article.Image = imgURL
	}

	article.Link = url
	article.Title = title

	//navigate to the article page
	err = page.Navigate(url)
	if err != nil {
		return false, false, err
	}

	//wait for the article content element to load and get it
	contentElement, err := page.Timeout(10 * time.Second).ElementX(website.ContentElement)
	if err != nil {
		return false, false, err
	}

	//get the content from the article page
	article.Content = contentElement.MustText()

	//make summary of the article now
	//
	//
	//

	//get subtitle if there is
	if website.SubtitleElement != "" {
		subtitleElement, err := page.ElementX(website.SubtitleElement)
		if err != nil {
			return false, false, err
		}
		article.Content = subtitleElement.MustText() + "\n" + article.Content
	}

	//add article to the database and on the program struct array
	err = addArticle(article)
	if err != nil {
		return false, false, err
	}

	return false, true, nil
}
