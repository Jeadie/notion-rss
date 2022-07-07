package main

import (
	"fmt"
	"time"
)

func main() {

	nDao, err := ConstructNotionDaoFromEnv()
	if err != nil {
		panic(fmt.Errorf("configuration error: %w", err))
	}

	rssUrls := nDao.GetEnabledRssFeeds()
	last24Hours := time.Now().Add(-1 * time.Hour * time.Duration(24))
	rssContent := GetRssContent(rssUrls, last24Hours)

	failedCount := 0
	for item := range rssContent {
		err := nDao.AddRssItem(item)
		if err != nil {
			fmt.Printf("Could not create page for %s, URL: %s. Error: %s\n", item.title, item.link, err.Error())
			failedCount++
		}
	}

	// Fail after all RSS items are processed to minimise impact.
	if failedCount > 0 {
		panic(fmt.Errorf("%d Rss item/s failed to be created in the notion database. See errors above", failedCount))
	}
}
