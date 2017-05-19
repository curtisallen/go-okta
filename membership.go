package okta

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// users list of Users
type users []*User

// User Okta user
type User struct {
	ID              string      `json:"id"`
	Status          string      `json:"status"`
	Created         *time.Time  `json:"created,omitempty"`
	Activated       *time.Time  `json:"activated,omitempty"`
	StatusChanged   *time.Time  `json:"statusChanged,omitempty"`
	LastLogin       *time.Time  `json:"lastLogin,omitempty"`
	LastUpdated     *time.Time  `json:"lastUpdated,omitempty"`
	PasswordChanged *time.Time  `json:"passwordChanged,omitempty"`
	Profile         UserProfile `json:"profile,omitempty"`
}

// UserProfile Okta User profile sans credentials
type UserProfile struct {
	Login             string `json:"login,omitempty"`
	FirstName         string `json:"firstName,omitempty"`
	LastName          string `json:"lastName,omitempty"`
	NickName          string `json:"nickName,omitempty"`
	DisplayName       string `json:"displayName,omitempty"`
	Email             string `json:"email,omitempty"`
	SecondEmail       string `json:"secondEmail,omitempty"`
	ProfileURL        string `json:"profileUrl,omitempty"`
	PreferredLanguage string `json:"preferredLanguage,omitempty"`
	UserType          string `json:"userType,omitempty"`
	Organization      string `json:"organization,omitempty"`
	Title             string `json:"title,omitempty"`
	Division          string `json:"division,omitempty"`
	Department        string `json:"department,omitempty"`
	CostCenter        string `json:"costCenter,omitempty"`
	EmployeeNumber    string `json:"employeeNumber,omitempty"`
	MobilePhone       string `json:"mobilePhone,omitempty"`
	PrimaryPhone      string `json:"primaryPhone,omitempty"`
	StreetAddress     string `json:"streetAddress,omitempty"`
	City              string `json:"city,omitempty"`
	State             string `json:"state,omitempty"`
	ZipCode           string `json:"zipCode,omitempty"`
	CountryCode       string `json:"countryCode,omitempty"`
}

// CreateMembership adds the given user to the given group
func (c *Client) CreateMembership(ctx context.Context, groupID, email string) error {
	userID, err := c.userIDbyEmail(ctx, email)
	if err != nil {
		return err
	}

	return c.sendRequest(
		ctx,
		"PUT",
		fmt.Sprintf("%s%s/%s/%s/%s", c.getRootURL(), groupSlug, groupID, userSlug, userID),
		nil,
		nil)
}

// DeleteMembership adds the given user to the given group
func (c *Client) DeleteMembership(ctx context.Context, groupID, email string) error {
	userID, err := c.userIDbyEmail(ctx, email)
	if err != nil {
		return err
	}

	return c.sendRequest(
		ctx,
		"DELETE",
		fmt.Sprintf("%s%s/%s/%s/%s", c.getRootURL(), groupSlug, groupID, userSlug, userID),
		nil,
		nil)
}

// MembershipExists checks if the user is a member of the given group
func (c *Client) MembershipExists(ctx context.Context, groupID, email string) (bool, error) {
	var result users
	err := c.sendRequest(
		ctx,
		"GET",
		fmt.Sprintf("%s%s/%s/%s", c.getRootURL(), groupSlug, groupID, userSlug),
		nil,
		&result)

	return result.contiansEmail(email), err
}

// Get a user's okta id by email address
func (c *Client) userIDbyEmail(ctx context.Context, email string) (string, error) {
	result := User{}
	err := c.sendRequest(
		ctx,
		"GET",
		fmt.Sprintf("%s%s/%s", c.getRootURL(), userSlug, email),
		nil,
		&result)

	return result.ID, err
}

func (up *UserProfile) contiansEmail(email string) bool {
	if strings.EqualFold(up.Email, email) {
		return true
	}
	return false
}

func (u users) contiansEmail(email string) bool {
	for _, v := range u {
		if v.Profile.contiansEmail(email) {
			return true
		}
	}
	return false
}
