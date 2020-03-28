package main

import (
  "flag"
  "fmt"
  "log"
  "net"
  "net/http"
  "strings"

  "github.com/PuerkitoBio/goquery"
)

// GetHistoryPage retrieves the history page URL from the wikipedia URL
func GetHistoryPage(url string) string {
  split := strings.SplitN(url, "/", 5)
  historyURL := "https://en.wikipedia.org/w/index.php?title=" + split[4] + "&offset=&limit=50000&action=history"
  return historyURL
}

// SortUnique returns a sorted array with unique records
func SortUnique(list []string) []string {
  uniqueIps := make(map[string]struct{})
  IPList := []string{}

  for _, ip := range list {
    if _, ok := uniqueIps[ip]; !ok {
      uniqueIps[ip] = struct{}{}
    }
  }

  for key := range uniqueIps {
    IPList = append(IPList, key)
  }

  return IPList
}

// GetIps does some stuff
func GetIps(historyPage string) []string {
  // Request the HTML page.
  IPList := []string{}

  res, err := http.Get(historyPage)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  doc.Find(".history-user").Each(func(i int, s *goquery.Selection) {
    name := s.Find("bdi").Text()
    if ip := net.ParseIP(string(name)); ip != nil {
      IPList = append(IPList, ip.String())
    }
  })

  return IPList
}

func main() {
  var url = flag.String("u", "", "Wikipedia URL to query")
  var sortUnique = flag.Bool("sort", false, "Sort IPs unique")

  flag.Parse()

  if *url != "" {
    historyPageURL := GetHistoryPage(*url)
    ips := GetIps(historyPageURL)

    if *sortUnique {
      ips = SortUnique(ips)
    }

    for _, element := range ips {
      fmt.Println(element)
    }
  } else {
    flag.PrintDefaults()
  }

}
