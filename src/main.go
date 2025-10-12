package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type RawRepoData map[string]any

const insertQueryTemplate string = `
	INSERT INTO repodata
	(id, name, user_kind, user_name, is_private, is_fork)
	VALUES
	(%f, '%v', '%v', '%v', %v, %v)
`

type RepoData struct {
	id         float64
	name       string
	user_kind  string // organization or just user
	user_name  string
	is_private bool
	is_fork    bool
}

func (rd *RepoData) getInsertQuery() string {
	return fmt.Sprintf(
		insertQueryTemplate,
		float64(rd.id),
		rd.name,
		rd.user_kind,
		rd.user_name,
		rd.is_private,
		rd.is_fork,
	)
}

func main() {
	var token string
	var kind string
	var name string

	flag.StringVar(&token, "token", "", "Github token")
	flag.StringVar(&kind, "kind", "orgs", "Client kind")
	flag.StringVar(&name, "name", "", "Client name")
	flag.Parse()

	response, err := getRepos(token, kind, name)
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}

	log.Printf(
		"Status code: %d",
		response.StatusCode,
	)

	if response.StatusCode >= 400 {
		os.Exit(1)
	}

	body, _ := io.ReadAll(response.Body)
	var response_data []RawRepoData
	pars_err := json.Unmarshal(body, &response_data)

	if pars_err != nil {
		log.Fatalf("%v", pars_err)
	}

	var processed_data []RepoData = processRawData(response_data)
	fmt.Println(processed_data)
	saveRepoDataInDb(processed_data)
}

func getRepos(token string, kind string, name string) (*http.Response, error) {
	// Return string list
	request, _ := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://api.github.com/%v/%v/repos",
			kind,
			name,
		),
		nil,
	)
	request.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %v", token),
	)

	client := &http.Client{
		CheckRedirect: nil,
	}
	response, err := client.Do(request)
	return response, err
}

func processRawData(response_data []RawRepoData) []RepoData {
	var result []RepoData
	for _, row := range response_data[1:] {
		var repo_id float64 = row["id"].(float64)
		var repo_name string = row["name"].(string)
		var priv bool = row["private"].(bool)
		var fork bool = row["fork"].(bool)

		var owner_data map[string]interface{} = row["owner"].(map[string]interface{})
		var user_name string = owner_data["login"].(string)
		var user_kind string = owner_data["type"].(string)

		result = append(
			result,
			RepoData{
				repo_id,
				repo_name,
				user_kind,
				user_name,
				priv,
				fork,
			},
		)
	}
	return result

}

func getPostgresClient() (*sql.DB, error) {
	client, err := sql.Open(
		"postgres",
		"postgresql://postgres:postgres@0.0.0.0:5430/postgres?sslmode=disable",
	)
	return client, err
}

func saveRepoDataInDb(records []RepoData) error {
	client, err := getPostgresClient()

	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}

	transaction, tx_err := client.Begin()
	if tx_err != nil {
		fmt.Println(tx_err)
		return tx_err
	}

	for _, row := range records {
		_, query_error := transaction.Exec(row.getInsertQuery())
		if query_error != nil {
			log.Fatalln(query_error)
			return transaction.Rollback()
		}
	}

	cmtError := transaction.Commit()
	if cmtError != nil {
		fmt.Println(cmtError)
		return cmtError
	}
	return nil
}
