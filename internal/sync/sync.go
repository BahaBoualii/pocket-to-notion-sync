package sync

import (
	"fmt"
	"log"
	"strings"

	"github.com/BahaBoualii/potion/internal/notion"
	"github.com/BahaBoualii/potion/internal/pocket"
)

type Syncer struct {
	pocketClient *pocket.Client
	notionClient *notion.Client
}

func NewSyncer(pocketClient *pocket.Client, notionClient *notion.Client) *Syncer {
	return &Syncer{
		pocketClient: pocketClient,
		notionClient: notionClient,
	}
}

func (s *Syncer) Sync() error {
	articles, err := s.pocketClient.GetArticles()
	if err != nil {
		return fmt.Errorf("error getting Pocket articles: %v", err)
	}

	fmt.Printf("Found %d articles in Pocket. Starting sync to Notion...\n", len(articles))

	synced := 0
	skipped := 0
	for i, article := range articles {
		fmt.Printf("\rProcessing article %d of %d...", i+1, len(articles))

		exists, err := s.notionClient.ArticleExists(article.ResolvedURL)
		if err != nil {
			log.Printf("\nError checking if article exists: %v", err)
			continue
		}
		if exists {
			skipped++
			continue
		}

		title := article.ResolvedTitle
		if strings.TrimSpace(title) == "" {
			if article.ResolvedURL != "" {
				title = "Article from " + article.ResolvedURL
			} else {
				title = "Untitled Pocket Article"
			}
		}

		tags := make([]string, 0, len(article.Tags))
		for tag := range article.Tags {
			tags = append(tags, tag)
		}

		err = s.notionClient.CreatePage(title, article.ResolvedURL, article.Excerpt, tags)
		if err != nil {
			log.Printf("\nError creating Notion page for article %s: %v", title, err)
		} else {
			synced++
		}
	}

	fmt.Printf("\n\nSync completed successfully!\n")
	fmt.Printf("Articles synced: %d\n", synced)
	fmt.Printf("Articles skipped (already existed): %d\n", skipped)

	return nil
}
