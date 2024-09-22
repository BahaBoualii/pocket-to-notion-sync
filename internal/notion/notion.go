package notion

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
)

type Client struct {
	client     *notionapi.Client
	databaseID string
}

func NewClient(token, databaseID string) (*Client, error) {
	client := notionapi.NewClient(notionapi.Token(token))
	return &Client{client: client, databaseID: databaseID}, nil
}

func (c *Client) ArticleExists(url string) (bool, error) {
	query := &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "URL",
			RichText: &notionapi.TextFilterCondition{
				Equals: url,
			},
		},
	}

	resp, err := c.client.Database.Query(context.Background(), notionapi.DatabaseID(c.databaseID), query)
	if err != nil {
		return false, fmt.Errorf("failed to query Notion database: %v", err)
	}

	return len(resp.Results) > 0, nil
}

func (c *Client) CreatePage(title, url, excerpt string, tags []string) error {
	pageCreateRequest := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(c.databaseID),
		},
		Properties: notionapi.Properties{
			"Title": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Type: notionapi.ObjectTypeText,
						Text: &notionapi.Text{
							Content: title,
						},
					},
				},
			},
			"URL": notionapi.URLProperty{
				URL: url,
			},
		},
	}

	if excerpt != "" {
		pageCreateRequest.Properties["Excerpt"] = notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: &notionapi.Text{
						Content: excerpt,
					},
				},
			},
		}
	}

	if len(tags) > 0 {
		notionTags := make([]notionapi.Option, len(tags))
		for i, tag := range tags {
			notionTags[i] = notionapi.Option{Name: tag}
		}
		pageCreateRequest.Properties["Tags"] = notionapi.MultiSelectProperty{
			MultiSelect: notionTags,
		}
	}

	_, err := c.client.Page.Create(context.Background(), pageCreateRequest)
	if err != nil {
		return fmt.Errorf("failed to create Notion page: %v", err)
	}

	return nil
}
