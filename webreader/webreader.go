package webreader

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	goose "github.com/advancedlogic/GoOse"
)

// Link is a title + href
type Link struct {
	Title string
	Href  string
}

// GetWebPage returns the full html of a web page
func GetWebPage(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// Read everything from Body.
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	// Convert body bytes to string.
	bodyText := string(body)

	resp.Body.Close()

	return bodyText
}

// GetBaseURL returns the base url of an url (http[s]://domain.tld)
func GetBaseURL(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	return u.Scheme + "://" + u.Host
}

// GetLinksFromQuery returns an array of links found in the html body
func GetLinksFromQuery(body string, query string, baseURL string) []Link {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	links := []Link{}

	// Find the requested items
	doc.Find(query).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title and absolute link
		title := s.Text()
		href, exists := s.Attr("href")
		if exists {
			links = append(
				links,
				Link{
					Title: title,
					Href:  baseURL + href,
				})
		}
	})

	return links
}

// GetTextFromQuery returns the (first) compatible text for a goquery within body
func GetTextFromQuery(body string, query string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	textFromQuery := ""
	// Find the requested item
	doc.Find(query).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title and absolute link
		textFromQuery = s.Text()
	})

	return textFromQuery

}

// GetArticleContentAndImage returns the main features of an article
func GetArticleContentAndImage(url string, body string) (string, string) {
	g := goose.New()
	article, _ := g.ExtractFromRawHTML(url, body)
	// println("title", article.Title)
	// println("description", article.MetaDescription)
	// println("keywords", article.MetaKeywords)
	cleanedText := article.CleanedText
	// remove the phrases without a dot at the end
	regex, _ := regexp.Compile(`(?mi)^[^.]+$`)
	cleanedText = regex.ReplaceAllString(cleanedText, "")
	// remove the phrases that are too short
	regex, _ = regexp.Compile(`(?mi)^.{0,80}$`)
	cleanedText = regex.ReplaceAllString(cleanedText, "")

	return cleanedText, article.TopImage
}
