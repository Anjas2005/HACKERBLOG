package main

import (
	// "encoding/csv"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type News struct {
	Scraped_Page,Rank, Title, Link, Score, Posted_by, Time_ago, Comments string
}

func call_scrape(x string) {
	// We Will Now get all the HTML
	log.Println("[INFO] SCRAPING URL ",x)
	res, err := http.Get(x)
	// Error Handling so that there must not be any problems
	if err != nil {
		log.Printf("[ERROR] Failed to connect to %s: %v\n", x, err)
	}
	defer res.Body.Close() // This Saves Resources
	if res.StatusCode != 200 {
		log.Printf("[ERROR] HTTP status %d from %s\n", res.StatusCode, x)
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Failed to parse the HTML document", err)
		return
	}
	// Variable to send all the news on one page together
	var allNews []News
	Hacker_News := doc.Find("table#hnmain") // # is a css selector for id so this means table with ID 'hnmain'
	tables := Hacker_News.Find("table")
	if tables.Length() < 2 {
		log.Println("[WARNING] Unexpected HTML structure: not enough tables")
		return
	}
	Hacker_News_Table := tables.Eq(1)

	var NextRow *goquery.Selection // This will be a NextRow variable that will hold a goquery.Selction value in future but is currently empty

	Hacker_News_Table.Find("tr").Each(func(i int, row *goquery.Selection) {
		if row.HasClass("athing submission") {

			href, exists := row.Find(".titleline a").Attr("href")

			if !exists{
				log.Println("[WARNING] Link is missing in the story row")
				return
			}

			if !strings.HasPrefix(href, "http") {
    	href = "https://news.ycombinator.com/" + href
			}

			NextRow = row.Next()

			if NextRow == nil {
				log.Println("[WARNING] Missing metadata row after title")
				return
			}
			// Scraping Rank
			rankText := row.Find(".rank").Text()
			// Checking Rank
			if len(rankText) < 2 {
			log.Println("[WARNING] Rank text too short:", rankText)
			return
			}

			score := strings.TrimSpace(NextRow.Find(".score").Text())
			if score == "" {
				score = "N/A (missing)"
			}

			postedBy := strings.TrimSpace(NextRow.Find(".hnuser").Text())
			if postedBy == "" {
				postedBy = "unknown_user"
			}
			// Creating  a news object with all the values
			news := News{
			Scraped_Page: x,
			Rank:      strings.TrimSuffix(rankText, "."),
			Title:     row.Find(".titleline").Text(),
			Link:      href,
			Score:     score,
			Posted_by: postedBy,
			Time_ago:  NextRow.Find(".age a").Text(),
			Comments:  NextRow.Find("a").Last().Text(),
			}
			// Appending to the list
			allNews = append(allNews, news)
		}

	})




	if len(allNews)==0{
		log.Println("[WARNING]NO NEWS ITEMS FOUND ON PAGE:",x)
		return
	}

	jsonData,err :=	json.Marshal(allNews)
	if err != nil{
		log.Printf("[ERROR] Failed to marshal JSON: %v \n",err)
		return
	}

	fmt.Println("\n [DEBUG] Sending JSON data :\n", string(jsonData))

	req, err := http.NewRequest("POST","http://127.0.0.1:8000/news/top-news/", bytes.NewBuffer(jsonData))
	if err != nil{
		log.Printf("[ERROR] Failed to create HTTP request: %v \n",err)
		return
	}
	req.Header.Set("Content-Type","application/json;charset=UTF-8")


	client := &http.Client{} // This is the actual object which will send the client
	resp,err := client.Do(req)
	if err != nil{
		log.Printf("[ERROR] Failed to send POST request: %v \n",err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("[ERROR] POST failed with status %d \n	",resp.StatusCode)
		//Save failing payload file
		filename := fmt.Sprintf("failed_payload_%d.json",time.Now().UnixNano())
		_=os.WriteFile(filename, jsonData,0644)
		log.Printf("[INFO] Payload saved to %s for debugging\n", filename)	
	}else {
		log.Println("[SUCCESS] POST successful")
	}
	

	// morelink := doc.Find("a.morelink")
	// if morelink.Length() > 0 {
	// 	href, exists := morelink.Attr("href")
	// 	if exists {
	// 		nextURL := "https://news.ycombinator.com/" + href
	// 		time.Sleep(2 * time.Second) // wait to be nice to the server
	// 		call_scrape(nextURL)
	// 	}
	// }
}

func main(){

// This configures the Go standard logger to include more information in its log output.

// Details:
// log.LstdFlags
// This adds the date and time to each log message.
// Example: 2009/01/23 01:23:23

// log.Lshortfile
// This adds the file name and line number where the log call was made.
// Example: main.go:15

// | (Bitwise OR)
// This combines both options, so both date/time and file:line info are included.

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	call_scrape("https://news.ycombinator.com/")
}



// 	if len(allNews) > 0{
// 			jsonData, _ := json.Marshal(allNews)
// 	request, _ := http.NewRequest("POST", "http://127.0.0.1:8000/news/top-news/", bytes.NewBuffer(jsonData))
// 	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
// 	client := &http.Client{}
// 	fmt.Println("\n Sending JSON data: ")
// 	fmt.Println(string(jsonData)) 
// 	response, error := client.Do(request)
// 	if error != nil {
// 		panic(error)
// 	}
// 	defer response.Body.Close()
// 	time.Sleep(5 * time.Millisecond)
// 	fmt.Println(string(jsonData))
// 	morelink := doc.Find("a.morelink") // more link class is an a tag
// 	if morelink.Length() > 0 {
// 		href, _ := morelink.Attr("href") // the _ takes the value and ignores so go is ok if the value is not used
// 		NextPage = "https://news.ycombinator.com/" + href
// 		call_scrape(NextPage)
// 	}
// 	}

// }

// func main() {
// 	call_scrape("https://news.ycombinator.com/")
// }
