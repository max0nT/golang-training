package githubapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreateRepoData struct {
	name        string
	description string
	homepage    string
	is_private  bool
	is_template bool
}

func (repoData *CreateRepoData) ToMap() map[string]any {
	return map[string]any{
		"name":        repoData.name,
		"description": repoData.description,
		"homepage":    repoData.homepage,
		"is_private":  repoData.is_private,
		"is_template": repoData.is_template,
	}
}

func (repoData *CreateRepoData) FromMap(data map[string]any) CreateRepoData {
	name := data["name"].(string)
	description := data["description"].(string)
	homepage := data["homepage"].(string)
	is_private := data["is_private"].(bool)
	is_template := data["is_template"].(bool)
	return CreateRepoData{
		name,
		description,
		homepage,
		is_private,
		is_template,
	}
}

func makeRequest(
	token string,
	url string,
	method string,
	data io.Reader,
) (*http.Response, error) {
	request, _ := http.NewRequest(
		method,
		url,
		data,
	)
	request.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %v", token),
	)
	client := &http.Client{
		CheckRedirect: nil,
	}
	return client.Do(request)
}

func GetMeInfo(
	token string,
) (*http.Response, error) {
	return makeRequest(
		token,
		"https://api.github.com/user",
		"GET",
		nil,
	)
}

func TokenIsValid(
	token string,
) bool {
	response, _ := GetMeInfo(token)
	return response.StatusCode == 200
}

func GetRepoList(
	token string,
	kind string,
	name string,
) (*http.Response, error) {
	return makeRequest(
		token,
		fmt.Sprintf(
			"https://api.github.com/%v/%v/repos",
			kind,
			name,
		),
		"GET",
		nil,
	)
}

func GetRepoDetailed(
	token string,
	owner string,
	repoName string,
) (*http.Response, error) {
	return makeRequest(
		token,
		fmt.Sprintf(
			"https://api.github.com/repos/%v/%v/",
			owner,
			repoName,
		),
		"GET",
		nil,
	)
}

func CreateRepo(
	token string,
	data CreateRepoData,
) (*http.Response, error) {
	jsonData, _ := json.Marshal(data.ToMap())
	return makeRequest(
		token,
		"https://api.github.com/user/repos",
		"POST",
		bytes.NewReader(jsonData),
	)
}
