package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"net/url"
	"time"
)

type RssItem struct {
	title   string
	link    url.URL
	content []string
}

// GetRssContent from a channel of RSS urls, parses new RSS items (that are from the lastNHours),
// and sends them to an output channel.
func GetRssContent(urls chan *url.URL, lastNHours int) chan RssItem {
	result := make(chan RssItem)

	go func(urls chan *url.URL, lastNHours int, rssContent chan RssItem) {
		defer close(result)
		timeSince := time.Now().Add(-1 * time.Hour * time.Duration(lastNHours))

		for u := range urls {
			for _, item := range GetRssContentFrom(u, timeSince) {
				rssContent <- *item
			}
		}
	}(urls, lastNHours, result)

	return result
}

// GetRssContentFrom since afterTime from the RSS feed found at url.
func GetRssContentFrom(url *url.URL, afterTime time.Time) []*RssItem {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url.String())
	if err != nil {
		fmt.Println(fmt.Errorf("could not get content from rss url: %s. Error occurred %w", url, err).Error())
		return []*RssItem{}
	}

	result := make([]*RssItem, len(feed.Items))
	count := 0
	for _, item := range feed.Items {
		if item.PublishedParsed.After(afterTime) {
			result[count] = convert(item)
			count++
		}
	}
	fmt.Printf("Feed %s has %d items. %d are within timerange and will be uploaded\n", url.String(), len(feed.Items), count)
	return result[:count]
}

// convert gofeed.Item into an internal RSSItem model.
func convert(item *gofeed.Item) *RssItem {
	link, _ := url.Parse(item.Link)
	return &RssItem{
		title:   item.Title,
		link:    *link,
		content: []string{item.Content},
	}
}
