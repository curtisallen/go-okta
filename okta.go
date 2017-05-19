package okta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const basePath = "/api/v1"

// Service Okta client interface
type Service interface {
	Group(ctx context.Context, id string) (Group, error)
	CreateGroup(ctx context.Context, group Group) (Group, error)
	UpdateGroup(ctx context.Context, group Group) (Group, error)
	DeleteGroup(ctx context.Context, groupID string) error
	CreateMembership(ctx context.Context, groupID, email string) error
	DeleteMembership(ctx context.Context, groupID, email string) error
	MembershipExists(ctx context.Context, groupID, email string) (bool, error)
}

// ErrNotFound returned when an Okta resource isn't found
var ErrNotFound = errors.New("Okta resource not found")

// Client okta client
type Client struct {
	token        string
	organization string
	host         string
	client       *http.Client
}

// NewClient creates a new Okta service interface
// token - Okta API token
// organization - okta org id e.g dev-1234
// preview - true if targeting okta preview false otherwise (defaults to true)
// client - http client to use
func NewClient(token, organization string, preview bool, client *http.Client) Service {

	service := &Client{
		token:        token,
		organization: organization,
		client:       client,
	}
	if preview {
		service.host = "oktapreview.com"
	} else {
		service.host = "okta.com"
	}

	if client == nil {
		service.client = http.DefaultClient
	}

	return service
}

// sendRequest sends a http request to the provided URL
func (c *Client) sendRequest(ctx context.Context,
	method, url string,
	body io.Reader,
	target interface{}) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", c.token))

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return err
	}

	err = checkCode(resp, target)
	return err
}

func checkCode(resp *http.Response, target interface{}) error {
	defer resp.Body.Close() // nolint: errcheck

	switch resp.StatusCode {
	case http.StatusOK, http.StatusAccepted, http.StatusCreated:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if len(body) != 0 && target != nil {
			err = json.Unmarshal(body, target)
			if err != nil {
				return err
			}
		}
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusNoContent:
		return nil
	default:
		return fmt.Errorf("Okta returned status code %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) getRootURL() string {
	if c.organization != "" {
		return fmt.Sprintf("https://%s.%s%s/", c.organization, c.host, basePath)
	}
	return fmt.Sprintf("%s%s/", c.host, basePath)
}
