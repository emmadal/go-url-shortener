package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
)

// protocol used. only https will be serve for url
var protocol = "https://"

// URLShortener define the url struct
type URLShortener struct {
	url   string
	short map[string]string
}

// IsValidURL verify if the url user submited is valid
func IsValidURL(url string) bool {
	// Regular expression to match URLs
	regex := regexp.MustCompile(`^(http|https):\/\/[\w\-]+(\.[\w\-]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?$`)
	return regex.MatchString(url)
}

// RandomURL allow to generate a shortened url from original URL
func (u URLShortener) GetShortenedURL() (map[string]string, error) {

	shortUrl, err := u.ShortURL()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"original": u.url,
		"short":    shortUrl,
	}, nil
}

func (u URLShortener) ShortURL() (string, error) {

	if ok := IsValidURL(u.url); !ok {
		return "", errors.New("Cannot shortened the url. use a valid url")
	}

	charset := "abcdefghijklmnopqrstuvwxyz"
	randomCharset := make([]byte, 3)

	for i := range 3 {
		randomCharset[i] = byte(charset[rand.Intn(len(charset))])
	}

	shortURL := fmt.Sprintf("%s%s.dl", protocol, string(randomCharset))
	return shortURL, nil
}

// OpenBrowser open the shortened link in the default browser depending on user system
func (u URLShortener) OpenBrowser() error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", u.short["original"]).Start()
		getURLResponse(u.short["original"])
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", u.short["original"]).Start()
		getURLResponse(u.short["original"])
	case "darwin":
		err = exec.Command("open", u.short["original"]).Start()
		getURLResponse(u.short["original"])
	default:
		err = errors.New("Browser not found")
		return err
	}
	return err
}

// retrive url response
func getURLResponse(url string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Date: %s\n", response.Header.Get("Date"))
	fmt.Printf("Redirect to: %s\n", response.Request.URL)
	fmt.Printf("HTTP request sent, awaiting response... %s\n", response.Status)
	fmt.Printf("Length: %d [application/octet-stream]\n", response.ContentLength)
}
