package infra

import (
	"bytes"
	"deni1688/gogie/internal/issues"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type gitlab struct {
	token  string
	host   string
	query  string
	client HttpClient
}

type gitlabIssue struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"description"`
	WebUrl string `json:"web_url"`
}

func NewGitlab(token, host, query string, client HttpClient) issues.GitProvider {
	return &gitlab{token, host, query, client}
}

func (r gitlab) GetRepos() (*[]issues.Repo, error) {
	req, err := r.request("GET", "projects")
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = r.query
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	var repos []issues.Repo

	return &repos, json.NewDecoder(resp.Body).Decode(&repos)
}

// Todo: Implement the CreateIssue method for Gitlab -> https://github.com/deni1688/gogie/issues/27
func (r gitlab) CreateIssue(repo *issues.Repo, issue *issues.Issue) error {
	req, err := r.request("POST", fmt.Sprintf("projects/%d/issues", repo.ID))
	if err != nil {
		return err
	}

	body, err := json.Marshal(gitlabIssue{Title: issue.Title, Desc: issue.Desc})
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(body))

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdIssue gitlabIssue
	if err = json.NewDecoder(resp.Body).Decode(&createdIssue); err != nil {
		return err
	}

	issue.ID = createdIssue.ID
	issue.Url = createdIssue.WebUrl

	return nil
}

func (r gitlab) request(method, resource string) (*http.Request, error) {
	req, err := http.NewRequest(method, r.endpoint(resource), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("PRIVATE-TOKEN", r.token)

	return req, err
}

func (r gitlab) endpoint(resource string) string {
	return r.host + "/api/v4/" + resource
}
