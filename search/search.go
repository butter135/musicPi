package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var developerKey string

func init() {
	developerKey = os.Getenv("YTKEY")
}

func main() {
	arguments := os.Args[1:]

	// コマンドライン引数があるかどうかをチェックします。
	if len(arguments) == 0 {
		fmt.Println("NO Keywords")
	} else {
		keywords := arguments[0]
		responses := getitems(keywords)
		for i := 0; i < 5; i++ {
			fmt.Printf("%d : %s \n", i, responses.Items[i].Snippet.Title)
		}
		fmt.Println("Select index number. Input 9 to cancel.")
		var rowidx int
		fmt.Scan(&rowidx)
		if rowidx == 9 {
			return
		}
		time := time.Now().Format("15:04:05")
		selected := responses.Items[rowidx]
		filename := time + " | " + replaceInvalidChars(selected.Snippet.Title)
		file, err := os.Create("/home/musicPi/.queue/" + filename)
		if err != nil {
			log.Fatal("Failed to create file", err)
		}
		_, err = file.WriteString("https://www.youtube.com/watch?v=" + string(selected.Id.VideoId))
		if err != nil {
			log.Fatal("Failed to write", err)
		}
		fmt.Println("[" + filename + "] added to queue")
		defer file.Close()

	}
}

func getitems(keywords string) *youtube.SearchListResponse {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(developerKey))
	if err != nil {
		log.Fatal("YouTube service creation failed: ", err)
	}

	search := youtube.NewSearchService(service).
		List([]string{"snippet"}).
		Q(*flag.String("query", keywords, "Search term")).
		MaxResults(*flag.Int64("max-results", 5, "Max YouTube results"))
	responses, err := search.Do()
	if err != nil {
		log.Fatal("Search service execution failed: ", err)
	}
	if len(responses.Items) == 0 {
		log.Fatal("No items found in response")
	}
	return responses
}

func replaceInvalidChars(filename string) string {
	invalidChars := "/\\0:*?\"<>|"
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, string(char), " ")
	}
	return filename
}
