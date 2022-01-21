package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
)

var output = flag.String("out", "", "path to output event file")
var verify = flag.String("verify", "", "path to output verification file")
var count = flag.Int("count", 100, "number of customers to create")
var extraattrs = flag.Int("attrs", 5, "number of extra random attributes to create")
var events = flag.Int("events", 500, "number of events to create")
var maxevents = flag.Int("maxevents", 5, "max number of event types to create per customer")
var dupes = flag.Int("dupes", 20, "1 / N events will be duplicated (N defaults to 20)")
var anon = flag.Int("anon", 500, "1 / N events will be not be assigned a user_id")
var seed = flag.Int("seed", int(time.Now().Unix()), "timestamp to use as a random seed")

func main() {
	flag.Parse()

	// gofakeit just uses the stock rand internally
	rand.Seed(int64(*seed))

	out, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	customers := makeCustomers(*count, *extraattrs)

	for i := 0; i < *events; i++ {
		customer := customers[strconv.Itoa(rand.Intn(len(customers))+1)]

		attrs, timestamp := completeAttributes(customer, 3, *events / *count / 2)

		if len(attrs) > 0 {
			attributes := map[string]interface{}{
				"id":        gofakeit.UUID(),
				"type":      "attributes",
				"user_id":   customer.id,
				"data":      attrs,
				"timestamp": timestamp,
			}

			js, _ := json.Marshal(attributes)
			out.Write(append(js, '\n'))
		}

		name := strings.ToLower(strings.Replace(gofakeit.HackerVerb()+gofakeit.Color(), " ", "", -1))

		if len(customer.eventSummary) >= *maxevents {
			name = sampleEvent(customer.eventSummary)
		}

		event := map[string]interface{}{
			"id":   gofakeit.UUID(),
			"type": "event",
			"name": name,
			"data": map[string]interface{}{
				strings.ToLower(strings.Replace(gofakeit.HipsterWord()+gofakeit.HackerAdjective(), " ", "", -1)): gofakeit.HackerVerb() + gofakeit.BuzzWord(),
				strings.ToLower(gofakeit.HackerAdjective() + gofakeit.Color()):                                   strings.Replace(gofakeit.BeerName(), " ", "", -1),
			},
			"timestamp": timestamp,
		}

		if i%*anon != *anon-1 {
			event["user_id"] = customer.id
			customer.eventSummary[name] += 1
		}

		js, _ := json.Marshal(event)

		out.Write(append(js, '\n'))

		if *dupes > 0 && rand.Intn(*dupes) == 0 {
			out.Write(append(js, '\n'))
		}
	}

	for i := 1; i < len(customers)+1; i++ {
		c := customers[strconv.Itoa(i)]

		attrs, timestamp := completeAttributes(c, 3, 0)

		if len(attrs) > 0 {
			attributes := map[string]interface{}{
				"id":        gofakeit.UUID(),
				"type":      "attributes",
				"user_id":   c.id,
				"data":      attrs,
				"timestamp": timestamp,
			}

			js, _ := json.Marshal(attributes)
			out.Write(append(js, '\n'))
		}
	}

	v, err := os.Create(*verify)
	if err != nil {
		log.Fatal(err)
	}
	defer v.Close()

	for _, id := range sortedIds(customers) {
		c := customers[id]
		v.Write([]byte(c.id + "," + sortedAttributes(c.attributes) + "," + sortedEvents(c.eventSummary) + "\n"))
	}
}

type customer struct {
	id             string
	attributes     map[string]string
	attrsCompleted map[string]bool
	eventSummary   map[string]int
}

func makeCustomers(count, maxExtraAttrs int) map[string]customer {
	customers := make(map[string]customer)

	for i := 0; i < count; i++ {
		id := strconv.Itoa(i + 1)

		customers[id] = customer{
			id:             id,
			attributes:     makeAttrs(maxExtraAttrs),
			attrsCompleted: make(map[string]bool),
			eventSummary:   make(map[string]int),
		}
	}

	return customers
}

func makeAttrs(maxExtra int) map[string]string {
	attrs := map[string]string{
		"first_name": randomValueFor("first_name"),
		"last_name":  randomValueFor("last_name"),
		"email":      randomValueFor("email"),
		"city":       randomValueFor("city"),
		"ip":         randomValueFor("ip"),
		"created_at": randomValueFor("created_at"),
	}

	if maxExtra > 0 {
		extras := rand.Intn(maxExtra) + 1

		for i := 0; i < extras; i++ {
			attrs[strings.ToLower(gofakeit.CarMaker()+gofakeit.FirstName())] = gofakeit.JobTitle() + gofakeit.LastName()
		}
	}

	return attrs
}

func completeAttributes(c customer, num int, rate int) (attrs map[string]string, ts int) {
	attrs = map[string]string{}

	if rate == 0 && len(c.attributes) <= len(c.attrsCompleted) {
		return
	}

	if (rate == 0 || rand.Intn(rate) == 0) && len(c.attributes) > len(c.attrsCompleted) {
		ts = *seed

		count := 0

		keys := make([]string, 0, len(c.attributes))
		for k := range c.attributes {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			if !c.attrsCompleted[k] {
				count += 1

				if rate == 0 || count <= num {
					attrs[k] = c.attributes[k]
					c.attrsCompleted[k] = true
				}
			}
		}

	} else {
		ts = int(time.Unix(int64(*seed), 0).Add(-time.Duration(rand.Intn(365*60)+10) * time.Minute).Unix())
		attrs = sampleAttributes(c.attributes, num)
	}

	return
}

func sampleAttributes(attrs map[string]string, count int) map[string]string {
	keys := make([]string, 0, len(attrs))
	for k, _ := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	randAttrs := make(map[string]string)

	for i := 0; i < count; i++ {
		key := keys[rand.Intn(len(keys))]
		randAttrs[key] = randomValueFor(key)
	}

	return randAttrs
}

func sampleEvent(summary map[string]int) string {

	keys := make([]string, 0, len(summary))
	for name, _ := range summary {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	return keys[rand.Intn(len(keys))]
}

func sortedIds(customers map[string]customer) []string {
	ids := make([]string, 0, len(customers))

	for id, _ := range customers {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	return ids
}

func sortedAttributes(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	list := make([]string, 0, len(keys))

	for _, key := range keys {
		list = append(list, key+"="+m[key])
	}

	return strings.Join(list, ",")
}

func sortedEvents(m map[string]int) string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	list := make([]string, 0, len(keys))

	for _, key := range keys {
		list = append(list, key+"="+strconv.Itoa(m[key]))
	}

	return strings.Join(list, ",")
}

func randomValueFor(attr string) string {
	switch attr {
	case "first_name":
		return gofakeit.FirstName()
	case "last_name":
		return gofakeit.LastName()
	case "email":
		return gofakeit.Email()
	case "city":
		return gofakeit.City()
	case "ip":
		return gofakeit.IPv4Address()
	case "created_at":
		return strconv.Itoa(int(time.Unix(int64(*seed), 0).Add(-time.Duration(rand.Intn(365)*24) * time.Hour).Unix()))
	default:
		return strings.Replace(gofakeit.HipsterWord()+gofakeit.HackerNoun(), " ", "", -1)
	}
}
