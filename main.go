package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/motemen/go-pocket/api"
	"github.com/motemen/go-pocket/auth"
)

const (
	pocketConsumerKey = "112219-9d8455cb537b41a455004a4"
	notionToken       = "secret_KTHi7TS8uYzhZiwcBrvGQH6qNbMGYX4mVjstGD7WbHy"
	notionDatabaseID  = "102e3d5afe0c8098b51fd7b3f0968fa7"
)

func authenticatePocket() (string, error) {
	redirectURL := "http://localhost" // Can be any valid URL for this example

	requestToken, err := auth.ObtainRequestToken(pocketConsumerKey, redirectURL)
	if err != nil {
		return "", fmt.Errorf("failed to obtain request token: %v", err)
	}

	authorizationURL := auth.GenerateAuthorizationURL(requestToken, redirectURL)
	fmt.Printf("Please open the following URL in your browser and authorize the app:\n%s\n", authorizationURL)
	fmt.Print("Press Enter when you've authorized the app...")
	fmt.Scanln() // Wait for user to authorize

	authorizationResponse, err := auth.ObtainAccessToken(pocketConsumerKey, requestToken)
	if err != nil {
		return "", fmt.Errorf("failed to obtain access token: %v", err)
	}

	return authorizationResponse.AccessToken, nil
}

func getPocketArticles(accessToken string) ([]api.Item, error) {
	client := api.NewClient(pocketConsumerKey, accessToken)
	options := &api.RetrieveOption{
		State: api.StateAll,
	}

	result, err := client.Retrieve(options)
	if err != nil {
		return nil, err
	}

	var articles []api.Item
	for _, item := range result.List {
		articles = append(articles, item)
	}

	return articles, nil
}

func createNotionPage(client *notionapi.Client, article api.Item) error {

	exists, err := articleExists(client, article.ResolvedURL)
	if err != nil {
		return fmt.Errorf("error checking if article exists: %v", err)
	}
	if exists {
		fmt.Printf("Article already exists in Notion: %s\n", article.ResolvedURL)
		return nil
	}

	// Determine the title
	title := article.ResolvedTitle
	if strings.TrimSpace(title) == "" {
		if article.ResolvedURL != "" {
			title = "Article from " + article.ResolvedURL
		} else {
			title = "Untitled Pocket Article"
		}
	}
	// Create the page
	pageCreateRequest := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(notionDatabaseID),
		},
		Properties: notionapi.Properties{
			"Title": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Type: notionapi.ObjectTypeText,
						Text: &notionapi.Text{
							Content: article.ResolvedTitle,
						},
					},
				},
			},
			"URL": notionapi.URLProperty{
				URL: article.ResolvedURL,
			},
			"Excerpt": notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.ObjectTypeText,
						Text: &notionapi.Text{
							Content: article.Excerpt,
						},
					},
				},
			},
		},
	}

	// Add tags if available
	if len(article.Tags) > 0 {
		tags := make([]notionapi.Option, 0, len(article.Tags))
		for tag := range article.Tags {
			tags = append(tags, notionapi.Option{Name: tag})
		}
		pageCreateRequest.Properties["Tags"] = notionapi.MultiSelectProperty{
			MultiSelect: tags,
		}
	}

	// Create the page in Notion
	_, err = client.Page.Create(context.Background(), pageCreateRequest)
	if err != nil {
		return fmt.Errorf("failed to create Notion page: %v", err)
	}

	return nil
}

func articleExists(client *notionapi.Client, url string) (bool, error) {
	query := &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "URL",
			RichText: &notionapi.TextFilterCondition{
				Equals: url,
			},
		},
	}

	resp, err := client.Database.Query(context.Background(), notionapi.DatabaseID(notionDatabaseID), query)
	if err != nil {
		return false, fmt.Errorf("failed to query Notion database: %v", err)
	}

	return len(resp.Results) > 0, nil
}

func main() {
	accessToken, err := authenticatePocket()
	if err != nil {
		log.Fatalf("Error authenticating with Pocket: %v", err)
	}

	articles, err := getPocketArticles(accessToken)
	if err != nil {
		log.Fatalf("Error getting Pocket articles: %v", err)
	}

	notionClient := notionapi.NewClient(notionapi.Token(notionToken))

	for _, article := range articles {
		err := createNotionPage(notionClient, article)
		if err != nil {
			log.Printf("Error processing article %s: %v", article.ResolvedURL, err)
		}
	}
}
