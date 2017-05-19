package okta

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"id\": \"test\"}")
	}))
	defer ts.Close()
	client := &Client{
		token:  "foo",
		client: http.DefaultClient,
		host:   ts.URL,
	}
	ctx := context.Background()
	group, err := client.Group(ctx, "123")
	if err != nil {
		t.Error("Error getting group", err)
	}
	if group.ID != "test" {
		t.Error("received unexpected ID", group.ID)
	}
	fmt.Printf("%+v", group)
}

func TestCreateGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"id\":\"00gaix5bquOfAfk1b0h7\", \"objectClass\":[\"okta:user_group\"], \"type\":\"OKTA_GROUP\", \"profile\":{\"Name\":\"test group 1\",\"description\":\"This is test group 1\"}}")
	}))
	defer ts.Close()
	client := &Client{
		token:  "foo",
		client: http.DefaultClient,
		host:   ts.URL,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	group := Group{
		Profile: GroupProfile{
			Name:        "test group 1",
			Description: "This is test group 1",
		},
	}
	createdGroup, err := client.CreateGroup(ctx, group)
	if err != nil {
		t.Error("Error creating group", err)
	}
	if createdGroup.ID != "00gaix5bquOfAfk1b0h7" {
		t.Errorf("Group ID incorrect got %s", createdGroup.ID)
	}
}

func TestUpdateGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"id\":\"00gaix5bquOfAfk1b0h7\", \"objectClass\":[\"okta:user_group\"], \"type\":\"OKTA_GROUP\", \"profile\":{\"Name\":\"test\",\"description\":\"updated test group\"}}")
	}))
	defer ts.Close()
	client := &Client{
		token:  "foo",
		client: http.DefaultClient,
		host:   ts.URL,
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	group := Group{
		ID: "00gaix5bquOfAfk1b0h7",
		Profile: GroupProfile{
			Name:        "test",
			Description: "updated test group",
		},
	}
	createdGroup, err := client.UpdateGroup(ctx, group)
	if err != nil {
		t.Error("Error creating group", err)
	}
	if createdGroup.ID != "00gaix5bquOfAfk1b0h7" {
		t.Errorf("Group ID incorrect got %s", createdGroup.ID)
	}
	if !strings.EqualFold(createdGroup.Profile.Description, "updated test group") {
		t.Errorf("Group Description doesn't match")
	}
}

func TestDeleteGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	client := &Client{
		token:  "foo",
		client: http.DefaultClient,
		host:   ts.URL,
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	err := client.DeleteGroup(ctx, "00gai2qvimqQ8lIQN0h7")
	if err != nil {
		t.Errorf("Unable to delete group, error %s", err)
	}
}
