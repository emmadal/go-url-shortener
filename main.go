package main

import (
	"fmt"
	"log"
)

func main() {
	var shortener URLShortener

	shortener.url = "https://twitter.com/emmanuel_dal"

	url, err := shortener.GetShortenedURL()
	if err != nil {
		log.Fatalln(err)
	}

	shortener.short = url
	fmt.Printf("Opening %s\n", url["short"])

	if err = shortener.OpenBrowser(); err != nil {
		log.Fatalln(err)
	}
}
