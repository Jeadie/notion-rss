package main

import (
	"os"
	"testing"
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

func TestGetOldUnstarredRSSItems(t *testing.T) {

}

func TestArchivePages(t *testing.T) {

}

func TestGetEnabledRssFeeds(t *testing.T) {

}

func TestGetRssFeedFromDatabaseObject(t *testing.T) {

}

func TestAddRssItem(t *testing.T) {

}
