package main

import "testing"

func TestAccessToken(t *testing.T) {
	iam := NewIAM("example-api-key")
	if iam.GetAccessToken() == "" {
		t.Error("Access token was empty.")
	}
}
