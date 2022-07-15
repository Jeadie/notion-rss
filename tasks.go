package main

import (
	"fmt"
	"time"
)

// NotionTask represents an independent unit of work to perform in Notion.so
type NotionTask struct {
	Run func(*NotionDao) error
}

// GetAllTasks that should be run.
func GetAllTasks() []NotionTask {
	return []NotionTask{
		{Run: ArchiveOldUnstarredContent},
		{Run: AddNewContent},
	}
}

// ArchiveOldUnstarredContent from the content database that is older than 30 days and is not starred.
func ArchiveOldUnstarredContent(nDao *NotionDao) error {
	pageIds := nDao.GetOldUnstarredRSSItemIds(time.Now().Add(-30 * time.Hour * time.Duration(24)))
	return nDao.ArchivePages(pageIds)
}

// AddNewContent from all enabled RSS Feeds that have been published within the last 24 hours.
func AddNewContent(nDao *NotionDao) error {
	rssFeeds := nDao.GetEnabledRssFeeds()
	last24Hours := time.Now().Add(-1 * time.Hour * time.Duration(24))
	rssContent := GetRssContent(rssFeeds, last24Hours)

	failedCount := 0
	for item := range rssContent {
		err := nDao.AddRssItem(item)
		if err != nil {
			fmt.Printf("Could not create page for %s, URL: %s. Error: %s\n", item.title, item.link.String(), err.Error())
			failedCount++
		}
	}

	// Fail after all RSS items are processed to minimise impact.
	if failedCount > 0 {
		return fmt.Errorf("%d Rss item/s failed to be created in the notion database. See errors above", failedCount)
	}
	return nil
}
