package okta

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateMembership(t *testing.T) {
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
	err := client.CreateMembership(ctx, "00g9ptlnxbRxcvRK50h7", "jim.halpert@mailinator.com")
	if err != nil {
		t.Errorf("Error creating membership %s", err)
	}
}

func TestDeleteMembership(t *testing.T) {
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
	err := client.DeleteMembership(ctx, "00g9ptlnxbRxcvRK50h7", "jim.halpert@mailinator.com")
	if err != nil {
		t.Errorf("Error creating membership %s", err)
	}
}

func TestMembershipExists(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		absPath, _ := filepath.Abs("./testdata/group_membership.json")
		byteResponse, _ := ioutil.ReadFile(absPath)
		_, err := w.Write(byteResponse)
		if err != nil {
			t.Fatal(err)
		}
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
	exists, err := client.MembershipExists(ctx, "00g9ptlnxbRxcvRK50h7", "nonexistant")
	if err != nil {
		t.Errorf("Error checking membership %s", err)
	}
	if exists {
		t.Errorf("Membership shouldn't exist %v", exists)
	}
	exists, err = client.MembershipExists(ctx, "00g9ptlnxbRxcvRK50h7", "dr.dre@example.com")
	if err != nil {
		t.Errorf("Error checking membership %s", err)
	}
	if !exists {
		t.Errorf("Membership doesn't exist %v", exists)
	}
}
