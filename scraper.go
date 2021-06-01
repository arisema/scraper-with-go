package scraper-with-go

import (
	"net/http"
	"io"
	"log"
	"golang.org/x/net/html"
	"strings"
)


func ScrapLinksInSection(url, sectionTitle string) [] string{

	log.Println("Starting Scrapping...")

	_, response := getDom(url)
	defer response.Body.Close()

	element := findDOMElement(response.Body, sectionTitle)

	links := getLinks(element)

	for _, link := range links {
		if !strings.HasPrefix(link, url) {
			link = url+link
		}
		log.Println(link)
	}

	log.Println("Completed")

	return links
}

/*
* Get DOM of website
*/
func getDom(url string) (dom *html.Tokenizer, response *http.Response){
	response, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	if response.StatusCode == 200 {
		dom = html.NewTokenizer(response.Body)
		return 
	}

	return
} 

/*
* Traverse DOM for section with sectionTitle
* takes DOM, element type (elementName), element name (eg. LatestBusiness)
*/
func findDOMElement(dom io.Reader, sectionTitle string) (node *html.Node){

	domElement, _ := html.Parse(dom)
	var parse func(*html.Node) 
	parse = func(currentNode *html.Node) {
		if currentNode.Type == html.TextNode &&  currentNode.Data == sectionTitle {
			node = currentNode.Parent.Parent
			return
		}
		for child := currentNode.FirstChild; child != nil; child = child.NextSibling {
			parse(child)
		}
	}
	parse(domElement)

	return	
}

/*
* Get all links in a given Node
*/
func getLinks(node *html.Node) ([] string){
	var links [] string

	var parse func(*html.Node)
	parse = func(currentNode *html.Node) {
		if currentNode.Type == html.ElementNode &&  currentNode.Data == "a" {
			for _, attribute := range currentNode.Attr {
				if attribute.Key == "href" {
					links = append(links, attribute.Val)
				}
			}
		}
		for child := currentNode.FirstChild; child != nil; child = child.NextSibling {
			parse(child)
		}
	}
	parse(node)

	return links
}

