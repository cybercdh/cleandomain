/*
cleandomain
cleans a list of domains by removing special characters and converting to lowercase
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var domains = make(chan string, 200)
var domainRegexp = regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`)

func main() {

	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "set the concurrency level")

	flag.Parse()

	// spin the work up
	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			for domain := range domains {

				d := cleanAndValidateDomain(domain)
				if d != "" {
					fmt.Println(d)
				}

			}
			wg.Done()
		}()
	}

	// get user input in the channel
	_, err := GetUserInput()
	if err != nil {
		log.Fatalln(err)
	}

	// tidy up
	close(domains)
	wg.Wait()
}

/*
get a list of domains from the user input
*/
func GetUserInput() (bool, error) {

	seen := make(map[string]bool)

	// read from stdin or from arg
	var input io.Reader
	input = os.Stdin

	arg_domain := flag.Arg(0)
	if arg_domain != "" {
		input = strings.NewReader(arg_domain)
	}

	sc := bufio.NewScanner(input)

	for sc.Scan() {

		domain := strings.ToLower(sc.Text())

		// ignore domains we've seen
		if _, ok := seen[domain]; ok {
			continue
		}

		seen[domain] = true

		// send to channel
		domains <- domain

	}

	// check there were no errors reading stdin
	if err := sc.Err(); err != nil {
		return false, err
	}

	return true, nil
}

func cleanAndValidateDomain(d string) string {

	// regex to match invalid chars
	invalidStartCharsRegexp := regexp.MustCompile(`^[^a-z0-9]+`)

	// replace and clean the domain
	d = invalidStartCharsRegexp.ReplaceAllString(d, "")
	d = strings.TrimLeft(d, ".")
	d = strings.TrimRight(d, ".")
	d = strings.ToLower(d)

	if IsValidDomain(d) {
		return d
	}
	return ""

}

func IsValidDomain(domain string) bool {
	return domainRegexp.MatchString(domain)
}
