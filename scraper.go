package main

import (
	"log"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// define a custom data type for the scraped data
type Product struct {
	Name, Price string
}

func main() {
	// where to store the scraped data
	// var products []Product

	// initialize a Chrome browser instance on port 4444
	service, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		log.Fatal("Error:", err)
	}

	defer service.Stop()

	// configure the browser options
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Path: "./chrome-linux64/chrome", Args: []string{
		"--headless", // comment out this line for testing
	}})

	// create a new remote client with the specified options
	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Error:", err)
	}

	// maximize the current window to avoid responsive rendering
	err = driver.MaximizeWindow("")
	if err != nil {
		log.Fatal("Error:", err)
	}

	// visit the target page
	// err = driver.Get("https://scrapingclub.com/exercise/list_infinite_scroll/")
	err = driver.Get("https://www.threads.net/")
	// err = driver.Get("https://www.threads.net/@ceoelonreevemusk.727252/post/DHm2GPKoj9l")
	if err != nil {
		log.Fatal("Error:", err)
	}

	// select the product elements
	// productElements, err := driver.FindElements(selenium.ByCSSSelector, ".post")
	productElements, err := driver.FindElements(selenium.ByCSSSelector, "span[dir='auto']")
	if err != nil {
		log.Fatal("Error:", err)
	}

	for _, product := range productElements {
		log.Println(product.Text())
	}

	// // iterate over the product elements
	// // and extract data from them
	// for _, productElement := range productElements {
	// 	// select the name and price nodes
	// 	nameElement, err := productElement.FindElement(selenium.ByCSSSelector, "span")
	// 	if err != nil {
	// 		log.Fatal("Error:", err)
	// 	}
	// 	priceElement, err := productElement.FindElement(selenium.ByCSSSelector, "h5")
	// 	if err != nil {
	// 		log.Fatal("Error:", err)
	// 	}

	// 	// extract the data of interest
	// 	name, err := nameElement.Text()
	// 	if err != nil {
	// 		log.Fatal("Error:", err)
	// 		name = "Empty"
	// 	}
	// 	price, err := priceElement.Text()
	// 	if err != nil {
	// 		log.Fatal("Error:", err)
	// 		price = "Empty"
	// 	}

	// 	// add the scraped data to the list
	// 	product := Product{}
	// 	product.Name = name
	// 	product.Price = price
	// 	product.Price = ""
	// 	products = append(products, product)
	// }

	// // export the scraped data to JSON
	// file, err := os.Create("products.json")
	// if err != nil {
	// 	log.Fatal("Error:", err)
	// }

	// defer file.Close()

	// // convert products slice to JSON
	// jsonData, err := json.MarshalIndent(products, "", "  ")
	// if err != nil {
	// 	log.Fatal("Error:", err)
	// }

	// // write JSON data to file
	// _, err = file.Write(jsonData)
	// if err != nil {
	// 	log.Fatal("Error:", err)
	// }
}
