package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

var (
	username *string
	password *string
	routerIP *string
)

func main() {
	setupFlags()

	success := checkSites()

	if !success {
		fmt.Println("Looks like the internet is down")
		rebootRouter()
	} else {
		fmt.Println("Everything is fine")
	}
}

func setupFlags() {
	username = flag.String("username", "", "Your router username")
	password = flag.String("password", "", "Your router password")
	routerIP = flag.String("routerip", "", "Your router IP")

	flag.Parse()
}

func checkSites() bool {
	sitesToCheck := []string{
		"https://www.google.co.nz",
		"http://www.stuff.co.nz",
	}

	requestsOK := 0

	for _, siteToCheck := range sitesToCheck {
		requestsOK += checkSite(siteToCheck)
	}

	if requestsOK == 0 {
		return false
	} else {
		return true
	}
}

func checkSite(url string) int {
	request, _ := http.NewRequest("GET", url, nil)

	timeout := time.Duration(5 * time.Second)

	client := &http.Client{
		Timeout: timeout,
	}

	response, error := client.Do(request)

	if error != nil {
		return 0
	}

	defer response.Body.Close()

	return 1
}

func rebootRouter() {
	loginCookie := loginToRouter()
	sendRebootRequest(loginCookie)
}

func loginToRouter() *http.Cookie {
	url := "http://" + *routerIP + "/log/in"
	query := "?un=" + *username + "&pw=" + *password + "&rd=%2Fuir%2Find.htm&rd2=%2Fuir%2Fbsc_login.htm&Nrd=1"

	request, _ := http.NewRequest("GET", url+query, nil)
	response := sendRequest(request)

	return response.Cookies()[0]
}

func sendRequest(request *http.Request) *http.Response {
	var defaultTransport http.RoundTripper = &http.Transport{}

	response, error := defaultTransport.RoundTrip(request)

	if error != nil {
		panic(error)
	}

	defer response.Body.Close()

	return response
}

func sendRebootRequest(loginCookie *http.Cookie) {
	url := "http://" + *routerIP + "/uir/rebo.htm?Nrd=0"
	request, _ := http.NewRequest("GET", url, nil)
	request.AddCookie(loginCookie)
	sendRequest(request)
}
