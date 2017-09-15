package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const BASE_PATH = "/projectinfo/v1/"

// GITHUB API HOST URL for access to github api
const GITHUB_API_HOST_URL = "https://api.github.com/"

type ProjectInfo struct {
	Name  string     `json:"name"`
	Owner GithubUser `json:"owner"`
}
type GithubUser struct {
	Username string `json:"login"`
	ID       uint   `json:"id"`
	Commits  uint16 `json:"contributions"`
}

type ResponseData struct {
	Project   string   `json:"project"`
	Owner     string   `json:"owner"`
	Committer string   `json:"committer"`
	Commits   uint16   `json:"commits"`
	Languages []string `json:"language"`
}

func githubProjectInfoHandler(writer http.ResponseWriter, r *http.Request) {

	// fmt.Fprintf(writer "user: %s\n", pathVars[4])
	// fmt.Fprintf(writer "repo: %s\n", pathVars[5])
	repoPath := getRepoPath(r.URL.Path)
	contributorsPath := repoPath + "/contributors"
	languagesPath := repoPath + "/languages"

	var project ProjectInfo
	jsonBody := getRequestBody(repoPath)

	if err := json.Unmarshal(jsonBody, &project); err != nil {
		log.Println(string(jsonBody))
		log.Println(err)
	}
	fmt.Fprintln(writer, "=============================\tProject info\t=============================")
	fmt.Fprintln(writer, "project name:\t"+project.Name)
	fmt.Fprintln(writer, "owner name:\t"+project.Owner.Username)

	topContributor := GithubUser{Username: "", ID: 0, Commits: 0}
	fmt.Fprintln(writer, "=============================\tContributors\t=============================")
	jsonBody = getRequestBody(contributorsPath)
	contributors := make([]GithubUser, 0) //https://coderwall.com/p/4c2zig/decode-top-level-json-array-into-a-slice-of-structs-in-golang
	if err := json.Unmarshal(jsonBody, &contributors); err != nil {
		log.Println(string(jsonBody))
		log.Println(err)
	}
	for i := 0; i < len(contributors); i++ {
		fmt.Fprintf(writer, "Contributor:\t%s, commits: %v\n", contributors[i].Username, contributors[i].Commits)
		if contributors[i].Commits > topContributor.Commits {
			topContributor = contributors[i]
		}
	}
	fmt.Fprintf(writer, "\nTop Contributor:\t%s, who made %v commits\n", topContributor.Username, topContributor.Commits)
	fmt.Fprintln(writer, "=============================\tLanguages\t=============================")
	jsonBody = getRequestBody(languagesPath)

	var topLanguage string
	topLanguageBytes := 0
	languages := map[string]int{} //https://coderwall.com/p/4c2zig/decode-top-level-json-array-into-a-slice-of-structs-in-golang
	if err := json.Unmarshal(jsonBody, &languages); err != nil {
		log.Println(string(jsonBody))
		log.Println(err)
	}
	for key, value := range languages {
		fmt.Fprintf(writer, "Language:\t%s,\t bytes of code: %v\n", key, value)
		if topLanguageBytes < value {
			topLanguage = key
			topLanguageBytes = value
		}
	}
	fmt.Fprintf(writer, "\nMost used Language:\t%s, with %v bytes of code\n", topLanguage, topLanguageBytes)

	// fmt.Fprintln(writer GITHUB_API_HOST_URL+repoPath)
	// fmt.Fprintln(writer "Projectinfo: "+string(getRequestBody(repoPath)))

}

func getRepoPath(originalPath string) string {
	URLPath := originalPath

	pathVars := strings.Split(URLPath, "/")
	if len(pathVars) < 5 {
		URLPath = "repos/stektpotet/Amazeking"
	}

	return strings.TrimPrefix(URLPath, BASE_PATH)
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
