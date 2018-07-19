package main

import "testing"
import "fmt"

func TestAccessToken(t *testing.T) {
	iam := New("example-api-key")
	if iam.GetAccessToken() == "" {
		t.Fail("Access token was empty.")
	}
}
