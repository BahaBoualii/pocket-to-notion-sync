package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/BahaBoualii/potion/internal/cli"
	"github.com/BahaBoualii/potion/internal/notion"
	"github.com/BahaBoualii/potion/internal/pocket"
	"github.com/BahaBoualii/potion/internal/sync"
)

var (
	notionDatabaseID string
)

func init() {
	flag.StringVar(&notionDatabaseID, "notion-db", "", "Notion Database ID")
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pocketConsumerKey := os.Getenv("POCKET_CONSUMER_KEY")
	notionToken := os.Getenv("NOTION_KEY")

	cli.DisplayTitle()
	cli.DisplayDescription()

	flag.Parse()
	cli.PromptForMissingFlags(&pocketConsumerKey, &notionToken, &notionDatabaseID)

	if !cli.ConfirmSync() {
		fmt.Println("Sync cancelled. Goodbye!")
		os.Exit(0)
	}

	fmt.Println("\nStarting the sync process...")

	pocketClient, err := pocket.NewClient(pocketConsumerKey)
	if err != nil {
		log.Fatalf("Error creating Pocket client: %v", err)
	}

	notionClient, err := notion.NewClient(notionToken, notionDatabaseID)
	if err != nil {
		log.Fatalf("Error creating Notion client: %v", err)
	}

	syncer := sync.NewSyncer(pocketClient, notionClient)
	err = syncer.Sync()
	if err != nil {
		log.Fatalf("Error during sync: %v", err)
	}

	fmt.Println("\nThank you for using Pocket to Notion Sync Tool!")
}
