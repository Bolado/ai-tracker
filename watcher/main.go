package watcher

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
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

// StartWatcher starts the watcher process.
func StartWatcher() error {
	// Start the Rod browser instance
	browser, err := startRod()
	if err != nil {
		return err
	}

	// Read the list of websites from a JSON file
	websites, err := readWebsitesJSON()
	if err != nil {
		return err
	}

	// Read the list of words to monitor from a JSON file
	words, err = readWordsJSON()
	if err != nil {
		return err
	}

	log.Printf("Loaded %d websites and %d words\n", len(websites), len(words))

	// Iterate over each website and start watching it
	for _, website := range websites {
		log.Printf("Watching website %s\n", website.Name)
		err := watchWebsite(website, browser)
		if err != nil {
			log.Printf("Error eccured while watching website %s: %v\n", website.Name, err)
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
	log.Printf("Adding article %s\n", article.Title)

	if err := database.InsertArticle(article); err != nil {
		return err
	}

	Articles = append(Articles, article)

	return nil
}

func startRod() (*rod.Browser, error) {
	// Create a new launcher instance
	launcher := launcher.New()

	// set path to the browser executable (optional)
	launcher.Bin(getChromiumPath())

	// Set the launcher to run in headless mode
	launcher.Headless(false)

	// Launch the browser and get the URL
	url, err := launcher.Launch()
	if err != nil {
		return nil, err
	}

	// Create a new browser instance and set the control URL
	browser := rod.New().ControlURL(url)

	// Connect to the browser
	err = browser.Connect()
	if err != nil {
		return nil, err
	}

	// Return the browser instance
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
	log.Printf("Navigating to %s\n", website.Url)

	//wait two seconds for network activity to settle
	page.WaitRequestIdle(2*time.Second, nil, nil, nil)

	//scroll some of the page
	page.Mouse.Scroll(0, 2000, 50)

	//wait for 10 seconds for the articles elements to load on the page
	els, err := page.Timeout(10 * time.Second).ElementsX(website.MainElement)
	if err != nil {
		return err
	}
	log.Printf("Looking at %d articles on %s\n", len(els), website.Name)

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
		log.Printf("Checking article %s\n", articleTitle)

		//check for relevant words on the title itself
		for _, word := range words {
			// check if article title contains relevant word
			if strings.Contains(articleTitle, word) {
				// analyze the article further
				existant, added, err := analyzeArticle(e, page, website, articleTitle)
				if err != nil {
					return err
				}
				if existant || added {
					log.Printf("Article %s is already existant or added\n", articleTitle)
					// if the article is already existant or added, continue the loop
					continue articlesLoop
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
	log.Printf("Analyzing article %s, got the url: %s, checking if it's already saved.\n", title, url)

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

	log.Printf("Article %s is not existant, navigating to the article page.\n", title)

	//navigate to the article page
	err = page.Navigate(url)
	if err != nil {
		return false, false, err
	}

	//wait for 10 seconds for the article page to load
	contentElements, err := page.Timeout(10 * time.Second).ElementsX(website.ContentElement)
	if err != nil {
		return false, false, err
	}

	// iterate over the content elements and concatenate the text
	var content strings.Builder
	for _, element := range contentElements {
		contentText := element.MustText()
		content.WriteString(contentText)
		content.WriteString("\n")
	}

	// get the final content string
	article.Content = content.String()
	log.Printf("Got the content of the article %s\n", title)

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

	log.Printf("Adding the article %s to the database\n", title)

	//add article to the database and on the program struct array
	err = addArticle(article)
	if err != nil {
		return false, false, err
	}

	return false, true, nil
}

// try to find on the system where chromium is by using which command
func getChromiumPath() string {
	cmd := exec.Command("which", "chromium")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	log.Printf("Chromium path: %s\n", string(out))

	// remove the newline character from the output
	out = out[:len(out)-1]

	return string(out)
}
