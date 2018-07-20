package main

import "encoding/json"
import "fmt"
import "time"
import "io/ioutil"
import "net/http"
import "net/url"
import "strings"

var ZERO_TIME = time.Time{}

type iam struct {
	APIKey       string
	URL          string
	refreshToken string
	accessToken  string
	reauthAt     time.Time
	refreshAt    time.Time
}

func NewIAM(apiKey string) *iam {
	return &iam{
		APIKey: apiKey,
		URL:    "https://iam.bluemix.net/identity/token",
	}
}

func (t *iam) GetAccessToken() string {
	if t.accessToken == "" {
		t.authenticate()
	} else if t.needsRefresh() {
		t.refresh()
	}
	return t.accessToken
}

func (t *iam) GetHeaders() *http.Header {
	headers := &http.Header{}
	t.AddAuthorization(headers)
	return headers
}

func (t *iam) AddAuthorization(headers *http.Header) {
	headers.Add("Authorization", fmt.Sprintf("Bearer %v", t.GetAccessToken()))
}

func (t *iam) needsRefresh() bool {
	return time.Now().UTC().After(t.refreshAt)
}

func (t *iam) needsReauthenticate() bool {
	return time.Now().UTC().After(t.reauthAt)
}

func (t *iam) request(values url.Values) {
	requestBody := strings.NewReader(values.Encode())

	req, err := http.NewRequest("POST", t.URL, requestBody)
	if err != nil {
		panic(fmt.Sprintf("Unable to create new HTTP request: %v", err))
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic Yng6Yng=") // b64("bx:bx")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Unable to send HTTP request: %v", err))
	}
	defer resp.Body.Close()

	type AccessTokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		Expiration   int64  `json:"expiration"`
		Scope        string `json:"scope"`
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	accessTokenResp := AccessTokenResponse{}
	err = json.Unmarshal(responseBody, &accessTokenResp)
	if err != nil {
		panic(err)
	}
	t.refreshToken = accessTokenResp.RefreshToken
	t.accessToken = accessTokenResp.AccessToken

	// Refresh tokens only last about a month without any warning (the
	// expiration is never specified via the API), so let's totally
	// re-authenticate after two weeks to be safe.
	if t.reauthAt == ZERO_TIME {
		t.reauthAt = time.Unix(accessTokenResp.Expiration, 0).Add(14 * 24 * time.Hour)
	}

	// This is subjective, but it calculates a desired refresh deadline after
	// an arbitrary portion of the actual TTL as passed. If we waited until the
	// TTL was actually expired, then there's a good chance that requests would
	// be issued with valid access tokens which expire in flight.
	var bufferedTTL float64 = (-0.5) * float64(accessTokenResp.ExpiresIn)
	t.refreshAt = time.Unix(accessTokenResp.Expiration, 0).Add(time.Duration(int64(bufferedTTL)) * time.Second)
}
func (t *iam) authenticate() {
	// Clear the reauthentication time so that it will be set properly again.
	t.reauthAt = ZERO_TIME

	values := url.Values{}
	values.Add("grant_type", "urn:ibm:params:oauth:grant-type:apikey")
	values.Add("apikey", t.APIKey)
	values.Add("response_type", "cloud_iam")
	t.request(values)
}

func (t *iam) refresh() {
	values := url.Values{}
	values.Add("grant_type", "refresh_token")
	values.Add("refresh_token", t.refreshToken)
	t.request(values)
}
