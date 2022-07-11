package main

import (
	"context"
	"fmt"
	"github.com/jomei/notionapi"
	"net/url"
	"os"
)

type NotionDao struct {
	feedDatabaseId    notionapi.DatabaseID
	contentDatabaseId notionapi.DatabaseID
	client            *notionapi.Client
}

// ConstructNotionDaoFromEnv given environment variables: NOTION_RSS_KEY,
// NOTION_RSS_CONTENT_DATABASE_ID, NOTION_RSS_FEEDS_DATABASE_ID
func ConstructNotionDaoFromEnv() (*NotionDao, error) {
	integrationKey, exists := os.LookupEnv("NOTION_RSS_KEY")
	if !exists {
		return &NotionDao{}, fmt.Errorf("`NOTION_RSS_KEY` not set")
	}

	contentDatabaseId, exists := os.LookupEnv("NOTION_RSS_CONTENT_DATABASE_ID")
	if !exists {
		return &NotionDao{}, fmt.Errorf("`NOTION_RSS_CONTENT_DATABASE_ID` not set")
	}

	feedDatabaseId, exists := os.LookupEnv("NOTION_RSS_FEEDS_DATABASE_ID")
	if !exists {
		return &NotionDao{}, fmt.Errorf("`NOTION_RSS_FEEDS_DATABASE_ID` not set")
	}

	return ConstructNotionDao(feedDatabaseId, contentDatabaseId, integrationKey), nil
}

func ConstructNotionDao(feedDatabaseId string, contentDatabaseId string, integrationKey string) *NotionDao {
	return &NotionDao{
		feedDatabaseId:    notionapi.DatabaseID(feedDatabaseId),
		contentDatabaseId: notionapi.DatabaseID(contentDatabaseId),
		client:            notionapi.NewClient(notionapi.Token(integrationKey)),
	}
}

// GetEnabledRssFeeds from the Feed Database. Results filtered on property "Enabled"=true
func (dao *NotionDao) GetEnabledRssFeeds() chan *FeedDatabaseItem {
	rssFeeds := make(chan *FeedDatabaseItem)

	go func(dao *NotionDao, output chan *FeedDatabaseItem) {
		defer close(output)

		req := &notionapi.DatabaseQueryRequest{
			Filter: notionapi.PropertyFilter{
				Property: "Enabled",
				Checkbox: &notionapi.CheckboxFilterCondition{
					Equals: true,
				},
			},
		}

		//TODO: Get multi-page pagination results from resp.HasMore
		resp, err := dao.client.Database.Query(context.Background(), dao.feedDatabaseId, req)
		if err != nil {
			return
		}
		for _, r := range resp.Results {
			feed, err := GetRssFeedFromDatabaseObject(&r)
			if err == nil {
				rssFeeds <- feed
			}
		}
	}(dao, rssFeeds)
	return rssFeeds
}

func GetRssFeedFromDatabaseObject(p *notionapi.Page) (*FeedDatabaseItem, error) {
	urlProperty := p.Properties["Link"].(*notionapi.URLProperty).URL
	rssUrl, err := url.Parse(urlProperty)
	if err != nil {
		return &FeedDatabaseItem{}, err
	}

	nameRichTexts := p.Properties["Title"].(*notionapi.TitleProperty).Title
	if len(nameRichTexts) == 0 {
		return &FeedDatabaseItem{}, fmt.Errorf("RSS Feed database entry does not have any Title in 'Title' field")
	}

	return &FeedDatabaseItem{
		FeedLink:     rssUrl,
		Created:      p.CreatedTime,
		LastModified: p.LastEditedTime,
		Name:         nameRichTexts[0].PlainText,
	}, nil
}

// AddRssItem to Notion database as a single new page with Block content. On failure, no retry is attempted.
func (dao NotionDao) AddRssItem(item RssItem) error {
	categories := make([]notionapi.Option, len(item.categories))
	for i, c := range item.categories {
		categories[i] = notionapi.Option{
			Name: c,
		}
	}

	_, err := dao.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: dao.contentDatabaseId,
		},
		Properties: map[string]notionapi.Property{
			"Title": notionapi.TitleProperty{
				Type: "title",
				Title: []notionapi.RichText{{
					Type: "text",
					Text: notionapi.Text{
						Content: item.title,
					},
				}},
			},
			"Link": notionapi.URLProperty{
				Type: "url",
				URL:  item.link.String(),
			},
			"Categories": notionapi.MultiSelectProperty{
				MultiSelect: categories,
			},
			"From":      notionapi.SelectProperty{Select: notionapi.Option{Name: item.feedName}},
			"Published": notionapi.DateProperty{Date: &notionapi.DateObject{Start: (*notionapi.Date)(item.published)}},
		},
		Children: RssContentToBlocks(item),
	})
	return err
}

func RssContentToBlocks(item RssItem) []notionapi.Block {
	// TODO: implement when we know RssItem struct better
	return []notionapi.Block{}
}
