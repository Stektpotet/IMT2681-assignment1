package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const BASE_PATH = "/projectinfo/v1/"

// GITHUB API HOST URL for access to github api
const GITHUB_API_HOST_URL = "https://api.github.com/"

type GithubOwner struct {
	Name string `json:"login"`
	Id   uint   `json:"id"`
}

type GithubContributor struct {
	Name          string `json:"login"`
	Contributions uint   `json:"contributions"`
}

type GithubRepo struct {
	Name         string              `json:"name"`
	Owner        GithubOwner         `json:"owner"`
	Contributors []GithubContributor `json:"contributions"`
	Languages    map[string]uint     `json:"languages"`
}

func githubProjectInfoHandler(w http.ResponseWriter, r *http.Request) {

	repoPath := strings.TrimPrefix(r.URL.Path, BASE_PATH)

	pathVars := strings.Split(r.URL.Path, "/")

	fmt.Fprintf(w, "user: %s\n", pathVars[4])
	fmt.Fprintf(w, "repo: %s\n", pathVars[5])

	// repo := GithubRepo{}
	// body := getRequestBody(repoPath)
	//
	// if err := json.Unmarshal(body, &repo); err != nil {
	// 	log.Println(string(body))
	// 	log.Println(err)
	// }
	fmt.Fprintln(w, GITHUB_API_HOST_URL+repoPath)
	fmt.Fprintln(w, "Projectinfo: "+string(getRequestBody(repoPath)))

}

func getRequestBody(repoPath string) []byte {
	response, err := http.Get(GITHUB_API_HOST_URL + repoPath)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {

	}
	response.Body.Close()
	return body
}

func main() {
	// client := github.NewClient(nil)

	http.HandleFunc(BASE_PATH, githubProjectInfoHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
