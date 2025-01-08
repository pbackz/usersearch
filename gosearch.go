package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ibnaleem/gobreach"
	"github.com/inancgumus/screen"
)

// Color output constants.
const (
	Red    = "\033[31m"
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)

// GoSearch ASCII logo.
const ASCII = `
 ________  ________  ________  _______   ________  ________  ________  ___  ___     
|\   ____\|\   __  \|\   ____\|\  ___ \ |\   __  \|\   __  \|\   ____\|\  \|\  \    
\ \  \___|\ \  \|\  \ \  \___|\ \   __/|\ \  \|\  \ \  \|\  \ \  \___|\ \  \\\  \   
 \ \  \  __\ \  \\\  \ \_____  \ \  \_|/_\ \   __  \ \   _  _\ \  \    \ \   __  \  
  \ \  \|\  \ \  \\\  \|____|\  \ \  \_|\ \ \  \ \  \ \  \\  \\ \  \____\ \  \ \  \ 
   \ \_______\ \_______\____\_\  \ \_______\ \__\ \__\ \__\\ _\\ \_______\ \__\ \__\
    \|_______|\|_______|\_________\|_______|\|__|\|__|\|__|\|__|\|_______|\|__|\|__|
                       \|_________|

`

// User-Agent header used in requests.
const UserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0"

// GoSearch version.
const VERSION = "v1.0.0"

var count uint16 = 0 // Maximum value for count is 65,535

type Website struct {
	Name            string   `json:"name"`
	BaseURL         string   `json:"base_url"`
	URLProbe        string   `json:"url_probe,omitempty"`
	FollowRedirects bool     `json:"follow_redirects,omitempty"`
	ErrorType       string   `json:"errorType"`
	ErrorMsg        string   `json:"errorMsg,omitempty"`
	ErrorCode       int      `json:"errorCode,omitempty"`
	Cookies         []Cookie `json:"cookies,omitempty"`
}

type Data struct {
	Websites []Website `json:"websites"`
}

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Stealer struct {
	TotalCorporateServices int         `json:"total_corporate_services"`
	TotalUserServices      int         `json:"total_user_services"`
	DateCompromised        string      `json:"date_compromised"`
	StealerFamily          string      `json:"stealer_family"`
	ComputerName           string      `json:"computer_name"`
	OperatingSystem        string      `json:"operating_system"`
	MalwarePath            string      `json:"malware_path"`
	Antiviruses            interface{} `json:"antiviruses"`
	IP                     string      `json:"ip"`
	TopPasswords           []string    `json:"top_passwords"`
	TopLogins              []string    `json:"top_logins"`
}

type HudsonRockResponse struct {
	Message  string    `json:"message"`
	Stealers []Stealer `json:"stealers"`
}

func UnmarshalJSON() (Data, error) {
	// GoSearch relies on data.json to determine the websites to search for.
	// Instead of forcing users to manually download the data.json file, we will fetch the latest version from the repository.
	// Therefore, we will do the following:
	// 1. Delete the existing data.json file if it exists as it will be outdated in the future
	// 2. Read the latest data.json file from the repository
	// Bonus: it does not download the data.json file, it just reads it from the repository.

	err := os.Remove("data.json")
	if err != nil && !os.IsNotExist(err) {
		return Data{}, fmt.Errorf("error deleting old data.json: %w", err)
	}

	url := "https://raw.githubusercontent.com/ibnaleem/gosearch/refs/heads/main/data.json"
	resp, err := http.Get(url)
	if err != nil {
		return Data{}, fmt.Errorf("error downloading data.json: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Data{}, fmt.Errorf("failed to download data.json, status code: %d", resp.StatusCode)
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return Data{}, fmt.Errorf("error reading downloaded content: %w", err)
	}

	var data Data
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return Data{}, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return data, nil
}

func WriteToFile(username string, content string) {
	filename := fmt.Sprintf("%s.txt", username)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		log.Fatal(err)
	}
}

func BuildURL(baseURL, username string) string {
	return strings.Replace(baseURL, "{}", username, 1)
}

func HudsonRock(username string, wg *sync.WaitGroup) {
	defer wg.Done()

	url := "https://cavalier.hudsonrock.com/api/json/v2/osint-tools/search-by-username?username=" + username

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data for "+username+" in HudsonRock function:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response in HudsonRock function:", err)
		return
	}

	var response HudsonRockResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing JSON in HudsonRock function:", err)
		return
	}

	if response.Message == "This username is not associated with a computer infected by an info-stealer. Visit https://www.hudsonrock.com/free-tools to discover additional free tools and Infostealers related data." {
		fmt.Println(Green + ":: This username is not associated with a computer infected by an info-stealer." + Reset)
		WriteToFile(username, ":: This username is not associated with a computer infected by an info-stealer.")
		return
	}

	fmt.Println(Red + ":: This username is associated with a computer that was infected by an info-stealer, all the credentials saved on this computer are at risk of being accessed by cybercriminals." + Reset)

	for i, stealer := range response.Stealers {
		fmt.Println(Red + fmt.Sprintf("[-] Stealer #%d", i+1) + Reset)
		fmt.Println(Red + fmt.Sprintf("::    Stealer Family: %s", stealer.StealerFamily) + Reset)
		fmt.Println(Red + fmt.Sprintf("::    Date Compromised: %s", stealer.DateCompromised) + Reset)
		fmt.Println(Red + fmt.Sprintf("::    Computer Name: %s", stealer.ComputerName) + Reset)
		fmt.Println(Red + fmt.Sprintf("::    Operating System: %s", stealer.OperatingSystem) + Reset)
		fmt.Println(Red + fmt.Sprintf("::    Malware Path: %s", stealer.MalwarePath) + Reset)

		switch v := stealer.Antiviruses.(type) {
		case string:
			WriteToFile(username, fmt.Sprintf("::    Antiviruses: %s\n", v))
		case []interface{}:
			antiviruses := make([]string, len(v))

			for i, av := range v {
				antiviruses[i] = fmt.Sprint(av)
			}

			avs := strings.Join(antiviruses, ", ")
			WriteToFile(username, fmt.Sprintf("::    Antiviruses: %s\n", avs))
		}

		fmt.Println(Red + fmt.Sprintf("::    IP: %s", stealer.IP) + Reset)

		fmt.Println(Red + "[-] Top Passwords:" + Reset)
		for _, password := range stealer.TopPasswords {
			fmt.Println(Red + fmt.Sprintf("::    %s", password) + Reset)
		}

		fmt.Println(Red + "[-] Top Logins:" + Reset)
		for _, login := range stealer.TopLogins {
			fmt.Println(Red + fmt.Sprintf("::    %s", login) + Reset)
		}
	}

	// For performance reasons, we should not print and write to the file at the same time during a single for-loop iteration.
	// Therefore, there will be 2 for-loop iterations: one for printing, and one for writing to the file.
	// This ensures that GoSearch can print as quickly as possible since the terminal output is most important.

	for i, stealer := range response.Stealers {
		WriteToFile(username, fmt.Sprintf("[-] Stealer #%d\n", i+1))
		WriteToFile(username, fmt.Sprintf("::    Stealer Family: %s\n", stealer.StealerFamily))
		WriteToFile(username, fmt.Sprintf("::    Date Compromised: %s\n", stealer.DateCompromised))
		WriteToFile(username, fmt.Sprintf("::    Computer Name: %s\n", stealer.ComputerName))
		WriteToFile(username, fmt.Sprintf("::    Operating System: %s\n", stealer.OperatingSystem))
		WriteToFile(username, fmt.Sprintf("::    Malware Path: %s\n", stealer.MalwarePath))

		switch v := stealer.Antiviruses.(type) {
		case string:
			WriteToFile(username, fmt.Sprintf("::    Antiviruses: %s\n", v))
		case []interface{}:
			antiviruses := make([]string, len(v))

			for i, av := range v {
				antiviruses[i] = fmt.Sprint(av)
			}

			avs := strings.Join(antiviruses, ", ")
			WriteToFile(username, fmt.Sprintf("::    Antiviruses: %s\n", avs))
		}

		WriteToFile(username, fmt.Sprintf("::    IP: %s\n", stealer.IP))

		WriteToFile(username, "[-] Top Passwords:\n")
		for _, password := range stealer.TopPasswords {
			WriteToFile(username, fmt.Sprintf("::    %s\n", password))
		}

		WriteToFile(username, "[-] Top Logins:\n")
		for _, login := range stealer.TopLogins {
			WriteToFile(username, fmt.Sprintf("::    %s\n", login))
		}
	}
}

func BuildEmail(username string) []string {
	emailDomains := []string{
		"@gmail.com",
		"@yahoo.com",
		"@outlook.com",
		"@hotmail.com",
		"@icloud.com",
		"@aol.com",
		"@live.com",
		"@protonmail.com",
		"@zoho.com",
		"@msn.com",
		"@proton.me",
		"@onionmail.org",
		"@gmx.de",
		"@mail2world.com",
	}

	var emails []string

	for _, domain := range emailDomains {
		emails = append(emails, username+domain)
	}

	return emails
}

func BuildDomains(username string) []string {
	tlds := []string{
		".com",
		".net",
		".org",
		".biz",
		".info",
		".name",
		".pro",
		".cat",
		".co",
		".me",
		".io",
		".tech",
		".dev",
		".app",
		".shop",
		".fail",
		".xyz",
		".blog",
		".portfolio",
		".store",
		".online",
		".about",
		".space",
		".lol",
		".fun",
		".social",
	}

	var domains []string

	for _, tld := range tlds {
		domains = append(domains, username+tld)
	}

	return domains
}

func SearchDomains(username string, domains []string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	fmt.Println(Yellow+"[*] Searching", len(domains), "domains with the username", username, "..."+Reset)

	domaincount := 0

	for _, domain := range domains {
		url := "http://" + domain

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf("Error creating request for %s: %v\n", domain, err)
			continue
		}
		req.Header.Set("User-Agent", UserAgent)

		resp, err := client.Do(req)
		if err != nil {
			netErr, ok := err.(net.Error)

			// The following errors mean that the domain does not exist.
			noSuchHostError := strings.Contains(err.Error(), "no such host")
			networkTimeoutError := ok && netErr.Timeout()

			if !noSuchHostError && !networkTimeoutError {
				fmt.Printf("Error sending request for %s: %v\n", domain, err)
			}

			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println(Green+"[+] 200 OK:", domain+Reset)
			domaincount++
		}
	}

	if domaincount > 0 {
		fmt.Println(Green+"[+] Found", domaincount, "domains with the username", username+Reset)
	} else {
		fmt.Println(Red+"[-] No domains found with the username", username+Reset)
	}
}

func SearchBreachDirectory(emails []string, apikey string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get an API key (10 lookups for free) @ https://rapidapi.com/rohan-patra/api/breachdirectory
	client, err := gobreach.NewBreachDirectoryClient(apikey)
	if err != nil {
		log.Fatal(err)
	}

	for _, email := range emails {
		fmt.Println(Yellow + "[*] Searching " + email + " on Breach Directory for any compromised passwords..." + Reset)

		response, err := client.SearchEmail(email)
		if err != nil {
			log.Fatal(err)
		}

		if response.Found == 0 {
			fmt.Printf(Red+"[-] No breaches found for %s. Moving on...\n", email+Reset)
			continue
		}

		fmt.Printf(Green+"[+] Found %d breaches for %s:\n", response.Found, email+Reset)
		for _, entry := range response.Result {
			fmt.Println(Green+"[+] Password:", entry.Password+Reset)
			fmt.Println(Green+"[+] SHA1:", entry.Sha1+Reset)
			fmt.Println(Green+"[+] Source:", entry.Sources+Reset)
		}
	}
}

func MakeRequestWithErrorCode(website Website, url string, username string) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	client := &http.Client{
		Timeout:   85 * time.Second,
		Transport: transport,
	}

	if !website.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error creating request in function MakeRequestWithErrorCode: %v\n", err)
		return
	}

	req.Header.Set("User-Agent", UserAgent)

	if website.Cookies != nil {
		for _, cookie := range website.Cookies {
			cookieObj := &http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			}
			req.AddCookie(cookieObj)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making GET request to %s: %v\n", url, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != website.ErrorCode {
		fmt.Println(Green+"[+]", website.Name+":", url+Reset)
		WriteToFile(username, url+"\n")
		count++
	}
}

func MakeRequestWithErrorMsg(website Website, url string, username string) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	client := &http.Client{
		Timeout:   85 * time.Second,
		Transport: transport,
	}

	if !website.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error creating request in function MakeRequestWithErrorMsg: %v\n", err)
		return
	}

	req.Header.Set("User-Agent", UserAgent)

	if website.Cookies != nil {
		for _, cookie := range website.Cookies {
			cookieObj := &http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			}
			req.AddCookie(cookieObj)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making GET request to %s: %v\n", url, err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	if website.URLProbe != "" {
		url = BuildURL(website.BaseURL, username)
	}

	bodyStr := string(body)
	// if the error message is not found in the response body, then the profile exists
	if !strings.Contains(bodyStr, website.ErrorMsg) {
		fmt.Println(Green+"[+]", website.Name+":", url+Reset)
		WriteToFile(username, url+"\n")
		count++
	}
}

func MakeRequestWithProfilePresence(website Website, url string, username string) {
	// Some websites have an indicator that a profile exists
	// but do not have an indicator when a profile does not exist.
	// If a profile indicator is not found, we can assume that the profile does not exist.

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	client := &http.Client{
		Timeout:   85 * time.Second,
		Transport: transport,
	}

	if !website.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error creating request in function MakeRequestWithErrorMsg: %v\n", err)
		return
	}

	req.Header.Set("User-Agent", UserAgent)

	if website.Cookies != nil {
		for _, cookie := range website.Cookies {
			cookieObj := &http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			}
			req.AddCookie(cookieObj)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making GET request to %s: %v\n", url, err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	if website.URLProbe != "" {
		url = BuildURL(website.BaseURL, username)
	}

	bodyStr := string(body)
	// if the profile indicator is found in the response body, the profile exists
	if strings.Contains(bodyStr, website.ErrorMsg) {
		fmt.Println(Green+"[+]", website.Name+":", url+Reset)
		WriteToFile(username, url+"\n")
		count++
	}
}

func Search(data Data, username string, wg *sync.WaitGroup) {
	var url string

	for _, website := range data.Websites {
		go func(website Website) {
			defer wg.Done()

			if website.URLProbe != "" {
				url = BuildURL(website.URLProbe, username)
			} else {
				url = BuildURL(website.BaseURL, username)
			}

			switch website.ErrorType {
			case "status_code":
				MakeRequestWithErrorCode(website, url, username)
			case "errorMsg":
				MakeRequestWithErrorMsg(website, url, username)
			case "profilePresence":
				MakeRequestWithProfilePresence(website, url, username)
			default:
				fmt.Println(Yellow+"[?]", website.Name+":", url+Reset)
				WriteToFile(username, "[?] "+url+"\n")
				count++
			}
		}(website)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gosearch <username>\nIssues: https://github.com/ibnaleem/gosearch/issues")
		os.Exit(1)
	}

	var username = os.Args[1]
	var wg sync.WaitGroup

	data, err := UnmarshalJSON()
	if err != nil {
		fmt.Printf("Error unmarshalling json: %v\n", err)
		os.Exit(1)
	}

	screen.Clear()
	fmt.Print(ASCII)
	fmt.Println(VERSION)
	fmt.Println(strings.Repeat("⎯", 85))
	fmt.Println(":: Username                              : ", username)
	fmt.Println(":: Websites                              : ", len(data.Websites))
	fmt.Println(strings.Repeat("⎯", 85))
	fmt.Println("[!] A yellow link indicates that I was unable to verify whether the username exists on the platform.")

	start := time.Now()

	wg.Add(len(data.Websites))
	go Search(data, username, &wg)
	wg.Wait()

	wg.Add(1)
	fmt.Println(strings.Repeat("⎯", 85))
	fmt.Println(Yellow + "[*] Searching HudsonRock's Cybercrime Intelligence Database..." + Reset)
	go HudsonRock(username, &wg)
	wg.Wait()

	if len(os.Args) == 3 {
		apikey := os.Args[2]
		fmt.Println(strings.Repeat("⎯", 85))
		emails := BuildEmail(username)
		wg.Add(1)
		go SearchBreachDirectory(emails, apikey, &wg)
		wg.Wait()
	}

	domains := BuildDomains(username)
	fmt.Println(strings.Repeat("⎯", 85))
	wg.Add(1)
	go SearchDomains(username, domains, &wg)
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println(strings.Repeat("⎯", 85))
	fmt.Println(":: Number of profiles found              : ", count)
	fmt.Println(":: Total time taken                      : ", elapsed)
}
