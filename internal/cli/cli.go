package cli

import (
	"fmt"
	"strings"
)

func DisplayTitle() {
	title := `
 ____            _        _     _        _   _       _   _             
|  _ \ ___   ___| | _____| |_  | |_ ___ | \ | | ___ | |_(_) ___  _ __  
| |_) / _ \ / __| |/ / _ \ __| | __/ _ \|  \| |/ _ \| __| |/ _ \| '_ \ 
|  __/ (_) | (__|   <  __/ |_  | || (_) | |\  | (_) | |_| | (_) | | | |
|_|   \___/ \___|_|\_\___|\__|  \__\___/|_| \_|\___/ \__|_|\___/|_| |_|
                                                                       
                          Sync Tool
`
	fmt.Println(title)
}

func DisplayDescription() {
	description := `
This tool syncs your saved Pocket articles to a Notion database.
It allows you to easily transfer your reading list from Pocket to Notion,
helping you organize and manage your content more effectively.

Before we begin, you'll need the following:
1. Pocket Consumer Key
2. Notion Integration Token
3. Notion Database ID

If you don't have these, here's how to get them:

- Pocket Consumer Key: Create an app at https://getpocket.com/developer/
- Notion Integration Token: Create an integration at https://www.notion.so/my-integrations
- Notion Database ID: Create a database in Notion and copy its ID from the URL

Let's get started!
`
	fmt.Println(description)
}

func PromptForMissingFlags(pocketKey, notionToken, notionDB *string) {
	if *pocketKey == "" {
		fmt.Print("Enter Pocket Consumer Key: ")
		fmt.Scanln(pocketKey)
	}
	if *notionToken == "" {
		fmt.Print("Enter Notion Integration Token: ")
		fmt.Scanln(notionToken)
	}
	if *notionDB == "" {
		fmt.Print("Enter Notion Database ID: ")
		fmt.Scanln(notionDB)
	}
}

func ConfirmSync() bool {
	var response string
	fmt.Print("\nDo you want to start syncing Pocket articles to Notion? (y/n): ")
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}
