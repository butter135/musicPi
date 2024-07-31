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
	time := time.Now().Format("15:04:05")
	arguments := os.Args[1:]

	// コマンドライン引数があるかどうかをチェックします。
	if len(arguments) == 0 {
		fmt.Println("NO URL")
	} else {
		url := arguments[0]
		filename := time + " | " + replaceInvalidChars(gettitle(url))
		file, err := os.Create("/home/musicPi/.queue/" + filename)
		if err != nil {
			log.Fatal("Failed to create file", err)
		}
		_, err = file.WriteString(url)
		if err != nil {
			log.Fatal("Failed to write", err)
		}
		fmt.Println("[" + filename + "] added to queue")
		defer file.Close()
	}
}

func gettitle(url string) string {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(developerKey))
	if err != nil {
		log.Fatal("YouTube service creation failed: ", err)
	}

	search := youtube.NewSearchService(service).
		List([]string{"snippet"}).
		Q(*flag.String("query", url, "Search term")).
		MaxResults(*flag.Int64("max-results", 1, "Max YouTube results"))
	response, err := search.Do()
	if err != nil {
		log.Fatal("Search service execution failed: ", err)
	}
	if len(response.Items) == 0 {
		log.Fatal("No items found in response")
	}

	return response.Items[0].Snippet.Title
}

func replaceInvalidChars(filename string) string {
	invalidChars := "/\\0:*?\"<>|"
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, string(char), " ")
	}
	return filename
}
