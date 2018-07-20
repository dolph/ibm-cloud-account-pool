package main

// This package contains a bunch of code copied from
// https://github.com/lancecarlson/couchgo/blob/master/couch.go (which does not
// currently have a LICENSE) for proof of concept purposes. I'm intending to
// supersede it here, possibly with https://github.com/go-kivik/kivik

import "bytes"
import "encoding/json"
import "fmt"
import "io"
import "io/ioutil"
import "net/http"
import "net/url"

type Client struct {
	IAM *iam
	URL *url.URL
}

type Response struct {
	Ok     bool
	ID     string
	Rev    string
	Error  string
	Reason string
}

func NewClient(IAM *iam, endpoint string) *Client {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	return &Client{
		IAM: IAM,
		URL: parsedURL,
	}
}

func (c *Client) ListAllDatabases() ([]string, error) {
	headers := c.IAM.GetHeaders()
	resp := []string{}
	_, err := c.execJSON("GET", "/_all_dbs", &resp, nil, nil, headers)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) execJSON(method string, path string, result interface{}, doc interface{}, values *url.Values, headers *http.Header) (int, error) {
	resBytes, code, err := c.execRead(method, path, doc, values, headers)
	if err != nil {
		return 0, err
	}
	if err = c.HandleResponseError(code, resBytes); err != nil {
		return code, err
	}
	if err = json.Unmarshal(resBytes, result); err != nil {
		return 0, err
	}
	return code, nil
}

func (c *Client) execRead(method string, path string, doc interface{}, values *url.Values, headers *http.Header) ([]byte, int, error) {
	r, code, err := c.exec(method, path, doc, values, headers)
	if err != nil {
		return nil, 0, err
	}
	defer r.Close()
	resBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, 0, err
	}
	return resBytes, code, nil
}

func (c *Client) exec(method string, path string, doc interface{}, values *url.Values, headers *http.Header) (io.ReadCloser, int, error) {
	reqReader, err := docReader(doc)
	if err != nil {
		return nil, 0, err
	}

	req, err := c.NewRequest(method, c.UrlString(path, values), reqReader, headers)
	if err != nil {
		return nil, 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, resp.StatusCode, nil
}

func docReader(doc interface{}) (io.Reader, error) {
	if doc == nil {
		return nil, nil
	}

	docJson, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	r := bytes.NewBuffer(docJson)
	return r, nil
}

func (c *Client) HandleResponse(resp *http.Response, result interface{}) (code int, err error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	code = resp.StatusCode
	if err = c.HandleResponseError(code, body); err != nil {
		return code, err
	}
	if err = json.Unmarshal(body, result); err != nil {
		return 0, err
	}
	return
}

func (c *Client) HandleResponseError(code int, resBytes []byte) error {
	if code < 200 || code >= 300 {
		res := Response{}
		if err := json.Unmarshal(resBytes, &res); err != nil {
			return err
		}
		return fmt.Errorf(fmt.Sprintf("Code: %d, Error: %s, Reason: %s", code, res.Error, res.Reason))
	}
	return nil
}

func (c *Client) NewRequest(method, url string, body io.Reader, headers *http.Header) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if headers != nil {
		req.Header = *headers
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return
}

func (c *Client) UrlString(path string, values *url.Values) string {
	u := c.URL
	u.Path = path
	if values != nil {
		u.RawQuery = values.Encode()
	}
	return u.String()
}
