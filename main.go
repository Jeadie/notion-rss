package main

import (
	"context"
	"fmt"
	"os"
)
import "github.com/jomei/notionapi"

func main() {
	key, exists := os.LookupEnv("NOTION_KEY")
	if !exists {
		fmt.Printf("`NOTION_KEY` not set.")
		return
	}

	client := notionapi.NewClient(notionapi.Token(key))
	u, err := client.User.Get(context.Background(), "me")
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	fmt.Printf("Hello %s!\n", u.Name)
}
