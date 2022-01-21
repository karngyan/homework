package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/customerio/homework/serve"
)

var (
	verifyFile = flag.String("verify-file", "", "file to verify against")
	serverAddr = flag.String("server-addr", "http://localhost:1323", "server to verify")
)

func main() {
	flag.Parse()

	if *verifyFile == "" {
		log.Fatal("must specify --verify-file")
	}

	f, err := os.Open(*verifyFile)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)

	var verifyCustomers = make(map[int]*serve.Customer)

	var line int
	for scanner.Scan() {
		line++
		var customer = serve.Customer{
			Attributes: make(map[string]string),
			Events:     make(map[string]int),
		}
		for i, element := range strings.Split(scanner.Text(), ",") {
			if i == 0 {
				customer.ID, err = strconv.Atoi(element)
				if err != nil {
					log.Fatalf("error on line %d of verify file: %s", line, err)
				}
				continue
			}

			if element == "" {
				continue
			}

			s := strings.Split(element, "=")
			if len(s) != 2 {
				log.Fatalf("error on line %d verify file: %s", line, element)
			}

			if s[0] == "created_at" {
				customer.Attributes[s[0]] = s[1]
				continue
			}

			if count, err := strconv.Atoi(s[1]); err == nil {
				customer.Events[s[0]] = count
				continue
			}

			customer.Attributes[s[0]] = s[1]
		}
		verifyCustomers[customer.ID] = &customer
	}

	var serverCustomers = make(map[int]*serve.Customer)

	for i := 0; i < len(verifyCustomers); i++ { // should take no more than len(verifyCustomers) iterations, if pagination were set to 1 per page

		resp, err := http.Get(fmt.Sprintf("%s/customers?page=%d&per_page=10000", *serverAddr, i+1))
		if err != nil {
			log.Fatalf("error requesting customers: %s", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("error reading response body: %s", err)
		}

		var reply struct {
			Customers []*serve.Customer `json:"customers"`
			Meta      struct {
				Page    int `json:"page"`
				PerPage int `json:"per_page"`
				Total   int `json:"total"`
			} `json:"meta"`
		}
		if err := json.Unmarshal(body, &reply); err != nil {
			log.Fatalf("error unmarshalling response body: %s", err)
		}

		for _, cust := range reply.Customers {
			cust.LastUpdated = 0 // not comparing these in our summaries
			serverCustomers[cust.ID] = cust
		}

		if len(reply.Customers) == 0 {
			break
		}
	}

	for id, vrfy := range verifyCustomers {
		cust, ok := serverCustomers[id]
		if !ok {
			log.Fatalf("cust %d missing from server responses\n", id)
			continue
		}

		if !reflect.DeepEqual(cust, vrfy) {
			log.Fatalf("cust %d didn't match!\ngot:%#v\nwant:%#v\n", cust.ID, cust, vrfy)
		}
	}

	for id := range serverCustomers {
		if _, ok := verifyCustomers[id]; !ok {
			log.Fatalf("extra customer found %d in server responses %#v\n", id, serverCustomers[id])
		}
	}

	fmt.Println("server responses match verification file!")
}
