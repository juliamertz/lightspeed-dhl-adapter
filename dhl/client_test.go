package dhl_test

import (
	"lightspeed-dhl/config"
	"lightspeed-dhl/dhl"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	conf, err := config.LoadSecrets("../config.toml")
	check(err, t)

	client := dhl.New(conf, dhl.DefaultCluster)
	err = client.Authenticate()

	if err != nil {
		t.Fatalf("Failed to authenticate with dhl api, error: %v", err)
	}

	if client.Session == nil {
		t.Fatalf("Expected session to be set after authentication")
	}

	session := client.GetSession()
	if session == nil {
		t.Fatalf("Expected GetSession to return newly created auth session")
	}
}

func TestReauth(t *testing.T) {
	conf, err := config.LoadSecrets("../config.toml")
	check(err, t)

	client := dhl.New(conf, dhl.DefaultCluster)
	err = client.Authenticate()
	if err != nil {
		t.Fatalf("%v", err)
	}

	now := int(time.Now().Local().Unix())
	client.Session.AccessTokenExpiration = now - 1

	originalToken := client.Session.AccessToken
	session := client.GetSession()
	if session == nil {
		t.Fatalf("Session is nil")
	}

	if session.AccessToken == originalToken {
		t.Fatalf("Access token did not change, old: %v\nnew: %v\n", session.AccessToken, originalToken)
	}
}

func TestAuthFailsafe(t *testing.T) {
	conf, err := config.LoadSecrets("../config.toml")
	check(err, t)

	client := dhl.New(conf, dhl.DefaultCluster)
	err = client.Authenticate()
	if err != nil {
		t.Fatalf("%v", err)
	}

	now := int(time.Now().Local().Unix())
	client.Session.AccessTokenExpiration = now - 1
	client.Session.RefreshTokenExpiration = now - 1

	originalToken := client.Session.AccessToken
	session := client.GetSession()
	if session == nil {
		t.Fatalf("Session is nil")
	}

	if session.AccessToken == originalToken {
		t.Fatalf("Access token did not change, old: %v\nnew: %v\n", session.AccessToken, originalToken)
	}

	client.Session = nil
	client.GetSession()

	if session == nil {
		t.Fatalf("Session is nil")
	}
}
