package main

import (
	"fmt"
	"github.com/jomei/notionapi"
	"net/url"
	"os"
	"testing"
	"time"
)

func runConstructNotionDaoFromEnvWith(rsskey string, content_database_id string, content_feeds_id string) (*NotionDao, error) {
	if len(rsskey) > 0 {
		os.Setenv("NOTION_RSS_KEY", rsskey)
	}
	if len(content_database_id) > 0 {
		os.Setenv("NOTION_RSS_CONTENT_DATABASE_ID", content_database_id)
	}
	if len(content_feeds_id) > 0 {
		os.Setenv("NOTION_RSS_FEEDS_DATABASE_ID", content_feeds_id)
	}

	nDao, err := ConstructNotionDaoFromEnv()

	if len(rsskey) > 0 {
		os.Unsetenv("NOTION_RSS_KEY")
	}
	if len(content_database_id) > 0 {
		os.Unsetenv("NOTION_RSS_CONTENT_DATABASE_ID")
	}
	if len(content_feeds_id) > 0 {
		os.Unsetenv("NOTION_RSS_FEEDS_DATABASE_ID")
	}

	return nDao, err
}

// getEnvWithDefault returns the environment variable (string), whether it existed (boolean), and
// the value to use (either the defaultValue, or the environment variable).
func getEnvWithDefault(key string, defaultValue string) (string, bool, string) {
	value, exists := os.LookupEnv(key)
	if exists {
		return value, exists, value
	} else {
		return value, exists, defaultValue
	}
}

func TestConstructNotionDaoFromEnv(t *testing.T) {
	// Store environment variables, and calculate values (possibly with defaults).
	PRIOR_NOTION_RSS_KEY, PNRK_exists, NOTION_RSS_KEY := getEnvWithDefault("NOTION_RSS_KEY", "NOTION_RSS_KEY")
	PRIOR_NOTION_RSS_CONTENT_DATABASE_ID, PNRCDI_exists, NOTION_RSS_CONTENT_DATABASE_ID := getEnvWithDefault("NOTION_RSS_CONTENT_DATABASE_ID", "NOTION_RSS_CONTENT_DATABASE_ID")
	PRIOR_NOTION_RSS_FEEDS_DATABASE_ID, PNRFDI_exists, NOTION_RSS_FEEDS_DATABASE_ID := getEnvWithDefault("NOTION_RSS_FEEDS_DATABASE_ID", "NOTION_RSS_FEEDS_DATABASE_ID")

	_, err := runConstructNotionDaoFromEnvWith("", NOTION_RSS_CONTENT_DATABASE_ID, NOTION_RSS_FEEDS_DATABASE_ID)
	if err == nil {
		t.Errorf("ConstructNotionDaoFromEnvWith should return error if `NOTION_RSS_KEY` is not set")
	}

	_, err = runConstructNotionDaoFromEnvWith(NOTION_RSS_KEY, "", NOTION_RSS_FEEDS_DATABASE_ID)
	if err == nil {
		t.Errorf("ConstructNotionDaoFromEnvWith should return error if `NOTION_RSS_KEY` is not set")
	}

	_, err = runConstructNotionDaoFromEnvWith(NOTION_RSS_KEY, NOTION_RSS_CONTENT_DATABASE_ID, "")
	if err == nil {
		t.Errorf("ConstructNotionDaoFromEnvWith should return error if `NOTION_RSS_KEY` is not set")
	}

	nDao, err := runConstructNotionDaoFromEnvWith(NOTION_RSS_KEY, NOTION_RSS_CONTENT_DATABASE_ID, NOTION_RSS_FEEDS_DATABASE_ID)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if nDao.client == nil {
		t.Errorf("notion client was not constructed")
	}

	// Reset environment variables
	if PNRK_exists {
		os.Setenv("NOTION_RSS_KEY", PRIOR_NOTION_RSS_KEY)
	}
	if PNRCDI_exists {
		os.Setenv("NOTION_RSS_CONTENT_DATABASE_ID", PRIOR_NOTION_RSS_CONTENT_DATABASE_ID)
	}
	if PNRFDI_exists {
		os.Setenv("NOTION_RSS_FEEDS_DATABASE_ID", PRIOR_NOTION_RSS_FEEDS_DATABASE_ID)
	}
}

//GetRssFeedFromDatabaseObject(p *notionapi.Page) (*FeedDatabaseItem, error) {
//urlProperty := p.Properties["Link"].(*notionapi.URLProperty).URL
//rssUrl, err := url.Parse(urlProperty)
//if err != nil {
//return &FeedDatabaseItem{}, err
//}
//
//nameRichTexts := p.Properties["Title"].(*notionapi.TitleProperty).Title
//if len(nameRichTexts) == 0 {
//return &FeedDatabaseItem{}, fmt.Errorf("RSS Feed database entry does not have any Title in 'Title' field")
//}
//
//return &FeedDatabaseItem{
//FeedLink:     rssUrl,
//Created:      p.CreatedTime,
//LastModified: p.LastEditedTime,
//Name:         nameRichTexts[0].PlainText,
//}, nil
//}
func TestGetRssFeedFromDatabaseObject(t *testing.T) {
	type TestCase struct {
		page           *notionapi.Page
		expectedDbItem *FeedDatabaseItem
		expectedErr    error
		subTestName    string
	}

	editedTime := time.Now()
	repoUrl, _ := url.Parse("https://github.com/Jeadie/notion-rss")
	tests := []TestCase{
		{
			page: &notionapi.Page{
				LastEditedTime: editedTime,
				Properties: map[string]notionapi.Property{
					"Title":   &notionapi.TitleProperty{Title: []notionapi.RichText{{PlainText: "TestTitle"}}},
					"Link":    &notionapi.URLProperty{URL: repoUrl.String()},
					"Enabled": &notionapi.CheckboxProperty{Checkbox: true},
				},
			},
			expectedDbItem: &FeedDatabaseItem{
				FeedLink:     repoUrl,
				Name:         "TestTitle",
				LastModified: editedTime,
			},
			expectedErr: nil,
			subTestName: "valid parsing",
		},
		{
			page: &notionapi.Page{
				LastEditedTime: editedTime,
				Properties: map[string]notionapi.Property{
					"Title":   &notionapi.TitleProperty{},
					"Link":    &notionapi.URLProperty{URL: repoUrl.String()},
					"Enabled": &notionapi.CheckboxProperty{Checkbox: true},
				},
			},
			expectedDbItem: &FeedDatabaseItem{},
			expectedErr:    fmt.Errorf("failed"),
			subTestName:    "no Title element in TitleProperty",
		},
		{
			page: &notionapi.Page{
				LastEditedTime: editedTime,
				Properties: map[string]notionapi.Property{
					"Link":    &notionapi.URLProperty{URL: repoUrl.String()},
					"Enabled": &notionapi.CheckboxProperty{Checkbox: true},
				},
			},
			expectedDbItem: &FeedDatabaseItem{},
			expectedErr:    fmt.Errorf("failed"),
			subTestName:    "Missing TitleProperty",
		},
	}

	for _, test := range tests {
		t.Run(test.subTestName, func(t *testing.T) {
			item, err := GetRssFeedFromDatabaseObject(test.page)
			if (err != nil) != (test.expectedErr != nil) {
				if err != nil {
					t.Errorf("Unexpected error occurred. Error: %s \n", err.Error())
				} else {
					t.Errorf("Error was expected, but none returned. Expected error: %s \n", test.expectedErr.Error())
				}
			}

			if test.expectedErr == nil {
				expectedItem := test.expectedDbItem
				if item.Name != expectedItem.Name {
					t.Errorf("Incorrect name of item. Expected %s, returned %s", expectedItem.Name, item.Name)
				}
				if item.FeedLink.String() != expectedItem.FeedLink.String() {
					t.Errorf("Incorrect RSS feed url. Expected %s, returned %s", expectedItem.FeedLink, item.FeedLink)
				}
				if item.Created != expectedItem.Created {
					t.Errorf("Incorrect created timestamp. Expected %s, returned %s", expectedItem.Created, item.Created)
				}
				if item.LastModified != expectedItem.LastModified {
					t.Errorf("Incorrect last modified timestamp. Expected %s, returned %s", expectedItem.LastModified, item.LastModified)
				}
			}

		})
	}

}
