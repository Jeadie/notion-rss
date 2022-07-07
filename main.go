package main

import (
	"fmt"
	"time"
)

func main() {

	nDao, err := ConstructNotionDaoFromEnv()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	rssUrls := nDao.GetEnabledRssFeeds()

	timeSince := time.Now().Add(-1 * time.Hour * time.Duration(24))
	rssContent := GetRssContent(rssUrls, timeSince)
	for item := range rssContent {
		nDao.AddRssItem(item)
	}

}
