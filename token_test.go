package main

import "testing"
import "fmt"

func TestAccessToken(t *testing.T) {
	iam := New("cc1LbGzZ18enVKxnIJcyCilQO7JtKQQLDB7LiMbtFGXi")
	fmt.Println(iam.GetAccessToken())
}
