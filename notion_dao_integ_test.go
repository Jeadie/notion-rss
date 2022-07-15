//go:build integration
// +build integration

package main

import (
	"context"
	"fmt"
	"github.com/jomei/notionapi"
	"net/url"
	"testing"
	"time"
)

func TestGetOldUnstarredRSSItems(t *testing.T) {
	nDao, err := ConstructNotionDaoFromEnv()
	if err != nil {
		t.Fatalf(err.Error())
	}

	itemUrl, _ := url.Parse("https://github.com/Jeadie/notion-rss")
	published := time.Now().Add(-1 * time.Hour)

	t.Run("Should use created time, not published time", func(t *testing.T) {

		// Setup
		err := nDao.AddRssItem(RssItem{
			title:      "notion-rss integration test",
			link:       *itemUrl,
			content:    []string{},
			categories: []string{"Integration testing", "notion", "rss"},
			feedName:   "Integration Test",
			published:  &published,
		})
		if err != nil {
			t.Fatalf(err.Error())
		}

		pages := nDao.GetOldUnstarredRSSItems(time.Now().Add(-1 * time.Minute))
		if len(pages) > 0 {
			fmt.Println(pages)
			t.Errorf("No pages are expected to exist that are this old")
		}
	})
	t.Run("Should return items olderThan", func(t *testing.T) {

		// Setup
		err := nDao.AddRssItem(RssItem{
			title:      "notion-rss integration test",
			link:       *itemUrl,
			content:    []string{},
			categories: []string{"Integration testing", "notion", "rss"},
			feedName:   "Integration Test",
			published:  &published,
		})
		if err != nil {
			t.Fatalf(err.Error())
		}

		pageIds := nDao.GetOldUnstarredRSSItems(time.Now().Add(1 * time.Hour))
		if len(pageIds) == 0 {
			t.Errorf("Expected GetOldUnstarredRSSItems to return an item")
		}
	})
	cleanupContentDatabase(nDao)
}

func cleanupContentDatabase(nDao *NotionDao) {
	resp, _ := nDao.client.Database.Query(context.Background(), nDao.contentDatabaseId, &notionapi.DatabaseQueryRequest{})
	for _, p := range resp.Results {
		_, err := nDao.client.Page.Update(
			context.TODO(),
			notionapi.PageID(p.ID),
			&notionapi.PageUpdateRequest{
				Archived:   true,
				Properties: notionapi.Properties{},
			})
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}

//
//func TestArchivePages(t *testing.T) {
//
//}
//
//func TestGetEnabledRssFeeds(t *testing.T) {
//
//}
//
//func TestAddRssItem(t *testing.T) {
//
//}
