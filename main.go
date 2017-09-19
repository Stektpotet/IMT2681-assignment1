package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const BASE_PATH = "/projectinfo/v1/"
const DEFAULT_REPO_PATH = "repos/apache/kafka"

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

	//TODO - Only handle r.Method == "GET"
	writer.Header().Add("content-type", "application/json")
	var serviceResponse ResponseData

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
	serviceResponse.Project = project.Name
	serviceResponse.Owner = project.Owner.Username

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

	serviceResponse.Committer = topContributor.Username
	serviceResponse.Commits = topContributor.Commits

	fmt.Fprintln(writer, "=============================\tLanguages\t=============================")
	jsonBody = getRequestBody(languagesPath)

	var topLanguage string
	topLanguageBytes := 0
	languageMap := make(map[string]int) //https://coderwall.com/p/4c2zig/decode-top-level-json-array-into-a-slice-of-structs-in-golang
	if err := json.Unmarshal(jsonBody, &languageMap); err != nil {
		log.Println(string(jsonBody))
		log.Println(err)
	}
	languages := make([]string, 0, len(languageMap))
	for key, value := range languageMap {
		fmt.Fprintf(writer, "Language:\t%s,\t bytes of code: %v\n", key, value)
		languages = append(languages, key)
		if topLanguageBytes < value {
			topLanguage = key
			topLanguageBytes = value
		}
	}
	fmt.Fprintf(writer, "\nMost used Language:\t%s, with %v bytes of code\n", topLanguage, topLanguageBytes)

	serviceResponse.Languages = languages
	json.NewEncoder(writer).Encode(serviceResponse)

	// fmt.Fprintln(writer "Projectinfo: "+string(getRequestBody(repoPath)))

}

func getRepoPath(originalPath string) string {
	//projectinfo/v1/repos/stektpotet/Amazeking
	URLPath := originalPath

	pathVars := strings.Split(URLPath, "/")
	if len(pathVars) < 5 {
		URLPath = DEFAULT_REPO_PATH
	}

	URLPath = strings.TrimPrefix(URLPath, BASE_PATH)
	return strings.Replace(URLPath, "github.com", "repos", 1)
}

func getRequestBody(repoPath string) []byte {
	response, err := http.Get(GITHUB_API_HOST_URL + repoPath)
	if err != nil {

	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {

	}
	response.Body.Close()
	return body
}

func main() {
	// client := github.NewClient(nil)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.HandleFunc(BASE_PATH, githubProjectInfoHandler)
	http.ListenAndServe(":"+port, nil)
}
