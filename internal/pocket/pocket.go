package pocket

import (
	"fmt"

	"github.com/motemen/go-pocket/api"
	"github.com/motemen/go-pocket/auth"
)

type Client struct {
	consumerKey string
	accessToken string
}

func NewClient(consumerKey string) (*Client, error) {
	client := &Client{consumerKey: consumerKey}
	err := client.authenticate()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) authenticate() error {
	redirectURL := "http://localhost"
	requestToken, err := auth.ObtainRequestToken(c.consumerKey, redirectURL)
	if err != nil {
		return fmt.Errorf("failed to obtain request token: %v", err)
	}

	authorizationURL := auth.GenerateAuthorizationURL(requestToken, redirectURL)
	fmt.Printf("Please open the following URL in your browser and authorize the app:\n%s\n", authorizationURL)
	fmt.Print("Press Enter when you've authorized the app...")
	fmt.Scanln()

	accessToken, err := auth.ObtainAccessToken(c.consumerKey, requestToken)
	if err != nil {
		return fmt.Errorf("failed to obtain access token: %v", err)
	}

	c.accessToken = accessToken.AccessToken
	return nil
}

func (c *Client) GetArticles() ([]api.Item, error) {
	client := api.NewClient(c.consumerKey, c.accessToken)
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
