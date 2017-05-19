package okta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const groupSlug = "groups"
const userSlug = "users"

// Group Okta group
type Group struct {
	ID                    string       `json:"id,omitempty"`
	Created               *time.Time   `json:"created,omitempty"`
	LastUpdated           *time.Time   `json:"lastUpdated,omitempty"`
	LastMembershipUpdated *time.Time   `json:"lastMembershipUpdated,omitempty"`
	ObjectClass           []string     `json:"objectClass,omitempty"`
	Type                  string       `json:"type,omitempty"`
	Profile               GroupProfile `json:"profile"`
}

// GroupProfile for okta
type GroupProfile struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Group gets the group for the given id
func (c *Client) Group(ctx context.Context, id string) (Group, error) {
	result := Group{}
	err := c.sendRequest(
		ctx,
		"GET",
		fmt.Sprintf("%s%s/%s", c.getRootURL(), groupSlug, id),
		nil,
		&result)
	return result, err
}

// CreateGroup create the given group in okta
func (c *Client) CreateGroup(ctx context.Context, group Group) (Group, error) {
	result := Group{}
	body, err := json.Marshal(group)
	if err != nil {
		return Group{}, err
	}
	err = c.sendRequest(
		ctx,
		"POST",
		fmt.Sprintf("%s%s", c.getRootURL(), groupSlug),
		bytes.NewReader(body),
		&result)
	if err != nil {
		return Group{}, err
	}

	return result, nil
}

// UpdateGroup updates the given group
func (c *Client) UpdateGroup(ctx context.Context, group Group) (Group, error) {
	result := Group{}
	body, err := json.Marshal(group)
	if err != nil {
		return Group{}, err
	}
	err = c.sendRequest(
		ctx,
		"PUT",
		fmt.Sprintf("%s%s/%s", c.getRootURL(), groupSlug, group.ID),
		bytes.NewReader(body),
		&result)
	if err != nil {
		return Group{}, err
	}

	return result, nil
}

// DeleteGroup removes the given group
func (c *Client) DeleteGroup(ctx context.Context, groupID string) error {
	return c.sendRequest(
		ctx,
		"DELETE",
		fmt.Sprintf("%s%s/%s", c.getRootURL(), groupSlug, groupID),
		nil,
		nil)
}
