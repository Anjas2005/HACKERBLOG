package main

import (
	// "encoding/csv"
	"fmt"
	"log"
	"net/http"

	// "os"
	"github.com/PuerkitoBio/goquery"
)

// type News struct{
// 	Rank,score,hours_ago,comments_number int
// 	title,url,uploaded_by string
// }

func call_scrape(x string)  {
	var NextPage string
	// We Will Now get all the HTML
	res, err := http.Get(x)
	// Error Handling so that there must not be any problems
	if err != nil {
		log.Fatal("Failed to coneect to the target page", err)
	}
	defer res.Body.Close() // This Saves Resources
	if res.StatusCode != 200 {
		log.Fatalf("HTTP Error %d: %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Failed to parse the HTML document", err)
	}

	Hacker_News := doc.Find("table#hnmain") // # is a css selector for id so this means table with ID 'hnmain'
	tables := Hacker_News.Find("table")
	Hacker_News_Table := tables.Eq(1)
	var NextRow *goquery.Selection // This will be a NextRow variable that will hold a goquery.Selction value in future but is currently empty
	Hacker_News_Table.Find("tr").Each(func(i int, row *goquery.Selection) {
	
		if row.HasClass("athing submission") {
			href,_ := row.Find(".titleline a").Attr("href")
			fmt.Printf(" Rank : %s, Title : %s ,LinkToPage: %s\n", row.Find(".rank").Text()[:len(row.Find(".rank").Text())-1], row.Find(".titleline").Text(),href)
			NextRow = row.Next()
			spans := NextRow.Find("a")
			as := spans.Last()
			fmt.Printf("Score: %s , Posted By: %s , %s , %s\n\n", NextRow.Find(".score").Text(), NextRow.Find(".hnuser").Text(), NextRow.Find(".age a").Text(), as.Text())
		}
		if row.HasClass("morelink"){
			href,_:=row.Attr("href")
			NextPage= "https://news.ycombinator.com/"+href
		}
		
	})
	call_scrape(NextPage)
}

func main(){
	call_scrape("https://news.ycombinator.com/")
}