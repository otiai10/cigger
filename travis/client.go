package travis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/otiai10/jsonindent"
	tn "github.com/otiai10/ternary"
)

// Client represents an API client for Travis-CI.
// See https://docs.travis-ci.com/api for more information.
type Client struct {
	// HTTP client if you want to use customized one
	HTTPClient *http.Client
	// API Token available on
	Token string
	// API Version
	Version string
	// API Host, either [api.travis-ci.org api.travis-ci.com]
	Host string
	// Output for logging
	Output io.Writer
}

// NewClient constructs *Client.
func NewClient(token string) *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		Token:      token,
		Version:    "3",
		Host:       "https://api.travis-ci.org",
		Output:     os.Stdout,
	}
}

// TriggerPayload represents a payload for request trigger.
type TriggerPayload struct {
	Request struct {
		Branch string `json:"branch"`
	} `json:"request"`
}

// TriggerResponse represents a response of trigger request.
type TriggerResponse struct {
	Type              string `json:"@type"`
	RemainingRequests int    `json:"remaining_requests"`
	Repository        struct {
		Type           string `json:"@type"`
		Href           string `json:"@href"`
		Representation string `json:"@representation"`
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Slug           string `json:"slug"`
	} `json:"repository"`
	Request struct {
		Repository struct {
			ID        int    `json:"id"`
			OwnerName string `json:"owner_name"`
			Name      string `json:"name"`
		} `json:"repository"`
		User struct {
			ID int `json:"id"`
		} `json:"user"`
		ID      int         `json:"id"`
		Message interface{} `json:"message"`
		Branch  string      `json:"branch"`
		Config  struct {
		} `json:"config"`
	} `json:"request"`
	ResourceType string `json:"resource_type"`
}

// Trigger triggers build for specified project.
func (c *Client) Trigger(slug string, options ...string) error {

	if len(options) == 0 {
		options = append(options, "master")
	}
	payload := new(TriggerPayload)
	payload.Request.Branch = tn.String(options[0])("master")

	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/repo/%s/requests", c.Host, url.PathEscape(slug))
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Travis-API-Version", c.Version)
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))

	if c.HTTPClient == nil {
		return fmt.Errorf("no http client set on this API client")
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resp := new(TriggerResponse)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return err
	}

	switch resp.Type {
	case "error":
		b, _ := json.MarshalIndent(resp, "", "  ")
		return fmt.Errorf("api response with error: %d %s", res.StatusCode, string(b))
	case "pending":
		fallthrough
	default:
		return jsonindent.NewEncoder(c.Output, "", "  ").Encode(resp)
	}
}
