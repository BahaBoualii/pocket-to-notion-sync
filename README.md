
<h1 align="center">
  <br>
  <p align="center"> <img src="https://upload.wikimedia.org/wikipedia/commons/2/2e/Pocket_App_Logo.png" alt="Pocket Logo" width="350"/> &nbsp;&nbsp;&nbsp; <img src="https://upload.wikimedia.org/wikipedia/commons/e/e9/Notion-logo.svg" alt="Notion Logo" width="100"/> </p>
  Pocket to Notion Sync Tool
  <br>
</h1>

<h4 align="center">A minimal CLI tool to synchronize all your saved articles in Pocket with your Notion database.</h4>

![screenshot-sync](https://github.com/user-attachments/assets/07e028d9-d344-490f-bdfd-234e970e7153)

## Key Features

* Sync Pocket articles to a Notion database. 
* Saves article metadata (title, URL, tags). 

* Simple CLI tool.

## How To Use

Prerequisites:
- [Go](https://golang.org/) installed. 
- Pocket API and Notion API tokens.
```bash 
git clone https://github.com/BahaBoualii/pocket-to-notion-sync.git 
cd pocket-to-notion-sync 
go mod tidy 
go build -o pocket-to-notion-sync
```

> **Note**
> Set environment variables:
    
>  ```bash    
>   export POCKET_API_TOKEN=your_pocket_api_token
>   export NOTION_API_TOKEN=your_notion_api_token
 
> Or create a `.env` file with the above variables.

Then just run this command:
`./pocket-to-notion-sync`

## Roadmap

 - [ ] Add automatic sync scheduling.
 - [ ] Filter Pocket articles by tags/date.
 - [ ] Enhance logging and error handling.
 - [ ] Build more customization options for the Notion page properties.

## Contributions

Feel free to contribute by opening issues or submitting PRs.

