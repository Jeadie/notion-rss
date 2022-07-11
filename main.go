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

	newContentErr := AddNewContent(nDao)
	archiveErr := ArchiveOldUnstarredContent(nDao)

	PanicOnErrors([]error{newContentErr, archiveErr})
}

// PanicOnErrors prints all non-nil err in errors and panics if there is at least one non-nil
// error in errors. Otherwise, return normally.
func PanicOnErrors(errors []error) {
	// Only used if one error (for better error handling).
	var firstErr error
	errN := 0

	// Print all non-nil errors.
	for _, err := range errors {
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			errN++
			firstErr = err
		}
	}

	// Multiple errors, panic with generic message.
	if errN > 1 {
		panic(fmt.Errorf("Multiple errors occured. Check output for details"))
	}

	if errN == 1 {
		panic(firstErr)
	}
}

// ArchiveOldUnstarredContent from the content database that is older than 30 days and is not starred.
func ArchiveOldUnstarredContent(nDao *NotionDao) error {
	pageIds := nDao.GetOldUnstarredRSSItems(time.Now().Add(-30 * time.Hour * time.Duration(24)))
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
			fmt.Printf("Could not create page for %s, URL: %s. Error: %s\n", item.title, item.link, err.Error())
			failedCount++
		}
	}

	// Fail after all RSS items are processed to minimise impact.
	if failedCount > 0 {
		return fmt.Errorf("%d Rss item/s failed to be created in the notion database. See errors above", failedCount)
	}
	return nil
}
