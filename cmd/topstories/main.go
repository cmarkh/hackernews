package main

import (
	"fmt"
	"log"

	"github.com/cmarkh/hackernews"
)

func main() {
	stories, err := hackernews.TopStories(10)
	if err != nil {
		log.Fatal(err)
	}

	for _, story := range stories {
		fmt.Println(story.Title)
		fmt.Println(story.TimeF)
		fmt.Println()
	}
}
