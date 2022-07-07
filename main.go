package main

import (
	"fmt"
)

func main() {

	nDao, err := ConstructNotionDaoFromEnv()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	rssUrls := nDao.GetEnabledRssFeeds()
	rssContent := GetRssContent(rssUrls, 24)
	for item := range rssContent {
		nDao.AddRssItem(item)
	}

}
