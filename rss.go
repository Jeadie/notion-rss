package main

import (
	"net/url"
	"time"
)

type RssItem struct {
	title   string
	link    url.URL
	content []string
}

func GetRssContent(urls chan *url.URL, lastNHours int) chan RssItem {
	result := make(chan RssItem)

	go func(urls chan *url.URL, lastNHours int) {
		defer close(urls)
		timeSince := time.Now().Add(-1 * time.Hour * lastNHours)

		rssContent := make(chan *RssItem)
		for u := range urls {
			rssContent <- GetRssContentFrom(u, timeSince)
		}
	}(urls, lastNHours)

	return result
}

func GetRssContentFrom(url *url.URL, since time.Time) *RssItem {

}
