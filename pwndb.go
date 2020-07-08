package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// User is the struct that contains leaked credentials
type User struct {
	email    string
	password string
}

// CheckDump checks pwndb for leaked credentials
func CheckDump(user string, domain string) string {
	postParam := url.Values{
		"luser":      {user},
		"domain":     {domain},
		"luseropr":   {"0"},
		"domainopr":  {"0"},
		"submitform": {"em"},
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.PostForm("https://pwndb2am4tzkvold.onion.ws/", postParam)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body)
}

// ParseDump parses the data returned from CheckDump
func ParseDump(rawData string) (users []User) {
	leaks := strings.Split(rawData, "Array")[2:]
	for _, leak := range leaks {
		username := strings.Split(strings.Split(leak, "[luser] => ")[1], "\n")[0]
		domain := strings.Split(strings.Split(leak, "[domain] => ")[1], "\n")[0]
		email := username + "@" + domain
		password := strings.Split(strings.Split(leak, "[password] => ")[1], "\n")[0]
		users = append(users, User{email, password})
	}
	return
}

func GetDump(username string, domain string) (users []User) {
	body := CheckDump(username, domain)

	// Error if data not returned correctly after a retry
	if !strings.Contains(body, "<pre>") {
		body := CheckDump(username, domain)
		if !strings.Contains(body, "<pre>") {
			log.Fatalln("Error contacting pwndb")
		}
	}
	rawData := strings.Split(strings.Split(body, "<pre>\n")[1], "</pre>")[0]
	users = ParseDump(rawData)
	return
}

func GetDumps(usernames []string, domains []string) (users []User) {
	//If domains and usernames are specified, only find matching permutations, otherwise find all
	if len(usernames) > 0 && len(domains) > 0 {
		for i := 0; i < len(domains); i++ {
			for j := 0; j < len(usernames); j++ {
				usersDump := GetDump(usernames[j], domains[i])
				users = append(users, usersDump...)
			}
		}
	} else if len(usernames) > 0 {
		for i := 0; i < len(usernames); i++ {
			usersDump := GetDump(usernames[i], "")
			users = append(users, usersDump...)
		}
	} else if len(domains) > 0 {
		for i := 0; i < len(domains); i++ {
			usersDump := GetDump("", domains[i])
			users = append(users, usersDump...)
		}
	}

	return
}

func init() {
	log.Println()
	fmt.Println(`                                         /$$ /$$      
                                        | $$| $$      
  /$$$$$$  /$$  /$$  /$$ /$$$$$$$   /$$$$$$$| $$$$$$$ 
 /$$__  $$| $$ | $$ | $$| $$__  $$ /$$__  $$| $$__  $$
| $$  \ $$| $$ | $$ | $$| $$  \ $$| $$  | $$| $$  \ $$
| $$  | $$| $$ | $$ | $$| $$  | $$| $$  | $$| $$  | $$
| $$$$$$$/|  $$$$$/$$$$/| $$  | $$|  $$$$$$$| $$$$$$$/
| $$____/  \_____/\___/ |__/  |__/ \_______/|_______/ 
| $$                                                  
| $$                                                  
|__/                                                  ` + "\n")
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var users, domains arrayFlags
	var totalString string
	flag.Var(&users, "user", "Username to check")
	flag.Var(&domains, "domain", "Domain to check")
	flag.Parse()

	if len(domains) == 0 && len(users) == 0 {
		flag.Usage()
		fmt.Println("  Example: pwndb -user foo -user bar -domain baz.com")
		fmt.Println()
		log.Fatalln("Please enter domains or users to check.")
	}

	dump := GetDumps(users, domains)
	if len(dump) < 1 {
		log.Fatalln("No data found")
	} else if len(dump) == 1 {
		totalString = "[1 User]\n"
	} else {
		totalString = fmt.Sprintf("[%d Users]\n", len(dump))
	}
	fmt.Printf(totalString)
	for _, user := range dump {
		fmt.Println(user.email + ":" + user.password)
	}
}
