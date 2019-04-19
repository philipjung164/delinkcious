package main

import (
	"encoding/json"
	"errors"
	"fmt"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	. "github.com/the-gigi/delinkcious/pkg/test_util"
	"io/ioutil"
	"log"
	"net/http"
	net_url "net/url"
	"os"
	"os/exec"
	"time"
)

var (
	delinkciousUrl   string
	delinkciousToken = os.Getenv("DELINKCIOUS_TOKEN")
	httpClient       = http.Client{}
)

type getLinksResponse struct {
	Err   string
	Links []om.Link
}

func getLinks() []om.Link {
	req, err := http.NewRequest("GET", string(delinkciousUrl)+"/links", nil)
	Check(err)

	req.Header.Add("Access-Token", delinkciousToken)
	r, err := httpClient.Do(req)
	Check(err)

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		Check(errors.New(r.Status))
	}

	var glr getLinksResponse
	body, err := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &glr)
	Check(err)
	if glr.Err != "" {
		Check(errors.New(glr.Err))
	}

	log.Println("=============")
	for _, link := range glr.Links {
		log.Println("title:", link.Title, "url:", link.Url, "status:", link.Status)
	}

	return glr.Links
}

func addLink(url string, title string) {
	params := net_url.Values{}
	params.Add("url", url)
	params.Add("title", title)
	qs := params.Encode()

	url = fmt.Sprintf("%s/links?%s", delinkciousUrl, qs)
	req, err := http.NewRequest("POST", url, nil)
	Check(err)

	req.Header.Add("Access-Token", delinkciousToken)
	r, err := httpClient.Do(req)
	Check(err)
	if r.StatusCode != http.StatusOK {
		defer r.Body.Close()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		Check(err)
		message := r.Status + " " + string(bodyBytes)
		Check(errors.New(message))
	}
}

func deleteLink(url string) {
	params := net_url.Values{}
	params.Add("url", url)
	qs := params.Encode()

	url = fmt.Sprintf("%s/links?%s", delinkciousUrl, qs)
	req, err := http.NewRequest("DELETE", url, nil)
	Check(err)

	req.Header.Add("Access-Token", delinkciousToken)
	r, err := httpClient.Do(req)
	Check(err)
	if r.StatusCode != http.StatusOK {
		defer r.Body.Close()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		Check(err)
		message := r.Status + " " + string(bodyBytes)
		Check(errors.New(message))
	}
}

func main() {
	tempUrl, err := exec.Command("minikube", "service", "api-gateway", "--url").CombinedOutput()
	delinkciousUrl = string(tempUrl[:len(tempUrl)-1]) + "/v1.0"
	Check(err)

	// Delete link
	deleteLink("https://github.com/the-gigi")

	// Get links
	getLinks()

	// Add a new link
	addLink("https://github.com/the-gigi", "Gigi on Github")

	// Get links again
	getLinks()

	// Wait a little and get links again
	time.Sleep(time.Second * 3)
	getLinks()
}