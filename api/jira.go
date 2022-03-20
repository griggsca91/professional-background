package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	rootURL = "/rest/api"
)

type JiraClient struct {
	baseURL string
	user    string
	token   string
}

func NewJiraClient(baseURL, user, token string) JiraClient {
	return JiraClient{
		baseURL: baseURL,
		user:    user,
		token:   token,
	}
}

func (c JiraClient) getAccountID() (string, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	encoded := base64.StdEncoding.EncodeToString([]byte(c.user + ":" + c.token))

	req, err := http.NewRequest("GET", c.baseURL+rootURL+"/3/myself", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("authorization", "Basic "+encoded)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsedPayload struct {
		AccountID string `json:"accountId"`
	}
	if err := json.Unmarshal(body, &parsedPayload); err != nil {
		return "", err
	}

	fmt.Printf("parsedPayload %+v\n", parsedPayload)
	fmt.Println(string(body))

	return parsedPayload.AccountID, nil
}

// assignee in (61ba89ec599f18006a4e216e) AND status in (BLOCKED, In-progress)

const jqlInprogressQuery = "assignee in (%s) AND status IN (In-progress)"
const jqlBlockedQuery = "assignee in (%s) AND status IN (BLOCKED)"

func (c JiraClient) getLatestTicket(query string) (string, error) {
	accountID, err := c.getAccountID()
	if err != nil {
		return "", err
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	encoded := base64.StdEncoding.EncodeToString([]byte(c.user + ":" + c.token))
	params := url.Values{}
	params.Add("jql", fmt.Sprintf(query, accountID))

	req, err := http.NewRequest("GET", c.baseURL+rootURL+"/2/search?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("authorization", "Basic "+encoded)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println(string(body))

	var result struct {
		Issues []struct {
			Key    string `json:"key"`
			Fields struct {
				Updated string `json:"updated"`
			} `json:"fields"`
		} `json:"issues"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Err unmarshalling", err)
		return "", err
	}

	if len(result.Issues) == 0 {
		return "", nil
	}

	fmt.Printf("result %+v\n", result)

	return result.Issues[0].Key, nil

}

func (c JiraClient) GetLatestInProgressTicket() (string, error) {
	return c.getLatestTicket(jqlInprogressQuery)
}

func (c JiraClient) GetLatestBlockedTicket() (string, error) {
	return c.getLatestTicket(jqlBlockedQuery)
}
