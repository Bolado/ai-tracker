package watcher

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/Bolado/ai-tracker/ai"
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
	sort.Slice(Articles, func(i, j int) bool {
		return Articles[i].Timestamp > Articles[j].Timestamp
	})

	return nil
}

func startRod() (*rod.Browser, error) {
	// Create a new launcher instance
	launcher := launcher.New()

	// Get the path to the Chromium binary if running on NixOS
	if isNixOS() {
		launcher.Bin(getChromiumPath())
	}

	// Set the launcher to run in headless mode
	launcher.Headless(true)

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

func isNixOS() bool {
	cmd := exec.Command("nixos-version")
	err := cmd.Run()
	return err == nil
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
	//load the website url
	page, err := browser.Page(proto.TargetCreateTarget{URL: website.Url})
	if err != nil {
		return err
	}
	log.Printf("Navigating to %s\n", website.Url)
	defer page.Close()

	//wait for the page to load
	page.WaitStable(2 * time.Second)

	//zoom out the page to 1% to make sure all the elements are loaded
	page.Eval("document.body.style.zoom = '1%'")

	//wait for the page to load the rest of the elements
	page.WaitStable(2 * time.Second)

	//create context and cancel function to avoid the context ending before you get elements data
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//get the articles available on the website
	els, err := page.Context(ctx).ElementsX(website.MainElement)
	if err != nil {
		return err
	}

	var articlesListItems []types.ArticlesListItem

	//populate the articles array with the articles found on the website, which contains the title, link and image
	articlesListItems, err = populateArticleListItems(els, website)
	if err != nil {
		return err
	}

	log.Printf("Looking at %d articles on %s\n", len(articlesListItems), website.Name)

	//iterate over the articles and analyze if it contains relevant words
	for _, articleListItem := range articlesListItems {

		//iterate over the words and check if the article contains them
		for _, word := range words {
			if strings.Contains(articleListItem.Title, word) {
				// if it does, log it
				log.Printf("Article %s contains relevant word %s\n", articleListItem.Title, word)

				// if the article is already existant, continue the loop
				if isExistant(articleListItem.Link) {
					log.Printf("Article %s is already on database.\n", articleListItem.Title)
					break
				}

				// analyze article further
				added, err := analyzeArticle(articleListItem, browser, website)
				if err != nil {
					return err
				}

				// if the article was added, continue the loop
				if added {
					log.Printf("Article %s is already existant or was added successfully.\n", articleListItem.Title)
					break
				}
			}
		}
	}

	return nil
}

// analyze article
func analyzeArticle(articleListItem types.ArticlesListItem, browser *rod.Browser, website types.Website) (added bool, err error) {
	var article types.Article

	log.Printf("Analyzing article %s, got the url: %s, checking if it's already saved.\n", articleListItem.Title, articleListItem.Link)

	article.Image = articleListItem.Image
	article.Link = articleListItem.Link
	article.Title = articleListItem.Title

	log.Printf("Article %s is not existant, navigating to the article page.\n", article.Title)

	//navigate to the article page on a new tab
	page, err := browser.Page(proto.TargetCreateTarget{URL: article.Link})
	if err != nil {
		return false, err
	}
	defer page.Close()

	//scroll through the article a little bit
	page.Mouse.Scroll(0, 10000, 100)

	//get the content of the article
	contentElements, err := page.Timeout(10 * time.Second).ElementsX(website.ContentElement)
	if err != nil {
		return false, err
	}

	// iterate over the content elements and concatenate the text
	var content strings.Builder
	for _, element := range contentElements {
		contentText := element.MustText()
		content.WriteString(contentText)
		content.WriteString("\n")
	}

	// get the compiled content string
	article.Content = content.String()
	log.Printf("Got the content of the article %s\n", article.Title)

	//summarize the article
	summary, err := ai.Summarize(article.Content)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	article.Summary = summary

	//get subtitle if there is
	if website.SubtitleElement != "" {
		subtitleElement, err := page.ElementX(website.SubtitleElement)
		if err != nil {
			log.Println(err.Error())
			return false, err
		}
		article.Content = subtitleElement.MustText() + "\n" + article.Content
	}

	//get date if there is
	if website.DateElement != "" {
		dateElement, err := page.ElementX(website.DateElement)
		if err != nil {
			log.Println(err.Error())
			return false, err
		}
		datetime := *dateElement.MustAttribute("datetime")
		t, err := time.Parse(time.RFC3339, datetime)
		if err != nil {
			log.Println(err.Error())
			return false, err
		}
		article.Timestamp = t.Unix()
	}

	log.Printf("Adding the article %s to the database\n", article.Title)

	article.Source = website.Name

	//add article to the database and on the program struct array
	err = addArticle(article)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	return true, nil
}

// try to find on the system where chromium is by using which command
func getChromiumPath() string {
	cmd := exec.Command("which", "chromium")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	// remove the newline character from the output
	out = out[:len(out)-1]
	return string(out)
}

func populateArticleListItems(els rod.Elements, website types.Website) ([]types.ArticlesListItem, error) {
	var articlesListItems []types.ArticlesListItem

	// loop through the elements found, and navigate to each relevant article that is not existant, then do the necessary tasks
	for _, e := range els {
		var articleListItem types.ArticlesListItem

		//get the title of the article
		articleTitleElement, err := e.Timeout(2 * time.Second).ElementX(website.TitleElement)
		if err != nil {
			return nil, err
		}
		articleListItem.Title = articleTitleElement.MustText()

		//get the url of the article
		urlElement, err := e.Timeout(2 * time.Second).ElementX(website.AnchorElement)
		if err != nil {
			return nil, err
		}
		articleListItem.Link = urlElement.Timeout(2 * time.Second).MustProperty("href").String()

		//get the image of the article
		if website.ImageElement != "" {
			imageElement, err := e.Timeout(10 * time.Second).ElementX(website.ImageElement)
			if err != nil {
				articleListItem.Image = "https://images.unsplash.com/photo-1674027444485-cec3da58eef4?ixlib=rb-4.0.3&q=85&fm=jpg&crop=entropy&cs=srgb&dl=growtika-nGoCBxiaRO0-unsplash.jpg&w=640"
			} else {
				articleListItem.Image = imageElement.Timeout(2 * time.Second).MustProperty("src").String()
			}
		}

		articlesListItems = append(articlesListItems, articleListItem)

	}
	return articlesListItems, nil
}
