# notion-rss
Get RSS feeds in notion.so daily

## Database Interface
This project uses two notion databases: to store RSS links (to subscribe to), to store the RSS content. 

### Feeds Database

| Property Name | Property Type |
| --- | :-- |
| Title | title |
| Link | url |

### Content Database

| Property Name | Property Type |
| --- | :-- |
| Title | title |
| Link | url |
| Enabled | boolean |
| From | Select |
| Categories | MultiSelect |
| Pubilshed | Date |

## Github Action Secrets
Github Secrets needed in the repository for the workflow actions to work:
- `NOTION_API_TOKEN`: notion.so api token for a specific integration. Integration must have access to `NOTION_RSS_CONTENT_DATABASE_ID` and `NOTION_RSS_FEEDS_DATABASE_ID`.     
- `NOTION_RSS_CONTENT_DATABASE_ID`: notion.so database id that stores RSS content (see Database Interface / Content Database).
- `NOTION_RSS_FEEDS_DATABASE_ID`: notion.so database id that stores RSS feed details (see Database Interface / Feeds Database).

## Nice to haves
1. Add all items from RSS feed when created date of notion database entry is within time window.
1. Use rss.Item.Content into notion blocks so the content can be viewed in Notion.
2. Add Categories to RSS items, and notion content tables.
3. Remove old, unread/starred items
4. Use release binary in `.github/workflows/release.yml`.
5. Convert combined RSS feeds into single feed: https://github.com/gorilla/feeds
