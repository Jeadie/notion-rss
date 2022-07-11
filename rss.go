package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"net/url"
	"time"
)

type RssItem struct {
	title      string
	link       url.URL
	content    []string
	categories []string
	feedName   string
	published  *time.Time
}

type FeedDatabaseItem struct {
	FeedLink     *url.URL
	Name         string
	Created      time.Time
	LastModified time.Time
}

// GetRssContent from a channel of RSS urls, parses new RSS items (that are from the lastNHours),
// and sends them to an output channel.
func GetRssContent(feedDatabaseItems chan *FeedDatabaseItem, since time.Time) chan RssItem {
	result := make(chan RssItem)

	go func(feeds chan *FeedDatabaseItem, since time.Time, rssContent chan RssItem) {
		defer close(result)

		for f := range feeds {
			for _, item := range GetRssContentFrom(f, since) {
				rssContent <- *item
			}
		}
	}(feedDatabaseItems, since, result)

	return result
}

// GetRssContentFrom since afterTime from the RSS feed found at url.
func GetRssContentFrom(feed *FeedDatabaseItem, afterTime time.Time) []*RssItem {
	fp := gofeed.NewParser()
	feedUrl := feed.FeedLink
	feedContent, err := fp.ParseURL(feedUrl.String())
	if err != nil {
		fmt.Println(fmt.Errorf("could not get content from rss url: %s. Error occurred %w", feedUrl, err).Error())
		return []*RssItem{}
	}

	// If Feed entry is new, publish all the content from it.
	publishAllItems := feed.Created.After(afterTime)
	result := make([]*RssItem, len(feedContent.Items))
	count := 0
	for _, item := range feedContent.Items {
		if publishAllItems || item.PublishedParsed.After(afterTime) {
			result[count] = convert(item, feed.Name)
			count++
		}
	}
	fmt.Printf("Feed %s has %d items. %d are eligible to be uploaded\n", feedUrl.String(), len(feedContent.Items), count)
	return result[:count]
}

// convert gofeed.Item into an internal RSSItem model.
func convert(item *gofeed.Item, itemFeedName string) *RssItem {
	link, _ := url.Parse(item.Link)
	return &RssItem{
		title:      item.Title,
		link:       *link,
		content:    []string{item.Content},
		categories: item.Categories,
		feedName:   itemFeedName,
		published:  item.PublishedParsed,
	}
}
