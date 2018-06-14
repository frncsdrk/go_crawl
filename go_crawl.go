package main

import (
    "fmt"
    "github.com/anaskhan96/soup"
    "os"
    "strings"
    "flag"
)

var (
    page string
)
var visitedUrlsMap = make(map[string]int)

type Url struct {
    href string
}

func (u *Url) fixUrl() string {
    var fixUrl string = u.href
    // add protocol https protocol
    if strings.Index(u.href, "//") == 0 {
        fixUrl = "https:" + u.href
    }
    // add base url to relative urls, e.g. /blog
    if strings.Index(u.href, "/") == 0 {
        fixUrl = page + u.href
    }

    return fixUrl
}

func crawlPage(crawlUrl string, baseUrl string, file *os.File) {
    // _, err := file.WriteString("CRAWLING" + crawlUrl + "\n")
    // logError(err)
    logToFile("CRAWLING" + crawlUrl + "\n", file)
    fmt.Println("CRAWLING ", crawlUrl)
    url := Url{crawlUrl}
    crawlUrl = url.fixUrl()
    res, err := soup.Get(crawlUrl)
    logError(err)
    baseUrlWOProtocol := strings.Split(baseUrl, "://")[1]
    baseUrlWOPath := strings.Split(baseUrlWOProtocol, "/")[0]
    // fmt.Println(baseUrlWOPath)
    doc := soup.HTMLParse(res)
    anchors := doc.FindAll("a")
    for _, anchor := range anchors {
        href := anchor.Attrs()["href"]
        if len(href) > 0 {
            _, err := file.WriteString("Found anchor with href: " + href + "\n")
            logError(err)
            fmt.Println("Found anchor with href:", href)
            if _, ok := visitedUrlsMap[href]; !ok &&
                strings.Index(href, baseUrlWOPath) != -1 &&
                strings.Index(href, "mailto:") == -1 {
                visitedUrlsMap[href] = 1
                crawlPage(href, baseUrl, file)
            } else {
                _, err := file.WriteString("SKIPPING " + href + "\n")
                logError(err)
                fmt.Println("SKIPPING ", href)
            }
        }
    }
}

func logError(err error) {
    if err != nil {
        fmt.Println(err)
    }
}

func logToFile(s string, f *os.File) {
    _, err := f.WriteString(s)
    logError(err)
}

func main() {
   flag.StringVar(&page, "page", "http://example.com", "the base page e.g. example.com")
   flag.Parse()

   f, err := os.Create("go_crawl.log")
   logError(err)
   defer f.Close()

   crawlPage(page, page, f)
   f.Sync()
}
