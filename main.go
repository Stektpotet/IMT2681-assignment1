package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//ServiceBasePath -Base path for this service
const ServiceBasePath = "/projectinfo/v1/"

//DefaultRepoPath - Default path to handle if no/too few parameters are given
const DefaultRepoPath = "repos/apache/kafka"

//GithubAPIHostURL - for access to github api
const GithubAPIHostURL = "https://api.github.com/"

//ProjectInfo - The project name and its owner
type ProjectInfo struct {
	Name  string     `json:"full_name"`
	Owner GithubUser `json:"owner"`
}

//GithubUser - Holder for username and commits given
type GithubUser struct {
	Username string `json:"login"`
	Commits  uint16 `json:"contributions"`
}

//ResponseData - the data returned by this service
type ResponseData struct {
	Project   string   `json:"project"`
	Owner     string   `json:"owner"`
	Committer string   `json:"committer"`
	Commits   uint16   `json:"commits"`
	Languages []string `json:"language"`
}

//RequestBodies - local json bodies for the three requests that would be done when deployed
type RequestBodies struct {
	Project      []byte
	Contributors []byte
	Languages    []byte
}

func getProjectInfo(jsonBody []byte, r *ResponseData) {
	var project ProjectInfo
	if err := json.Unmarshal(jsonBody, &project); err != nil {
		log.Fatalln(string(jsonBody))
		log.Fatalln(err)
	}
	r.Project = "github.com/" + project.Name
	r.Owner = project.Owner.Username
}

func getContributorInfo(jsonBody []byte, r *ResponseData) {
	contributors := make([]GithubUser, 0)
	if err := json.Unmarshal(jsonBody, &contributors); err != nil {
		log.Fatalln(string(jsonBody))
		log.Fatalln(err)
	}
	r.Committer = contributors[0].Username
	r.Commits = contributors[0].Commits
}

func getLanguageInfo(jsonBody []byte, r *ResponseData) {

	languageMap := make(map[string]int) //https://coderwall.com/p/4c2zig/decode-top-level-json-array-into-a-slice-of-structs-in-golang
	if err := json.Unmarshal(jsonBody, &languageMap); err != nil {
		log.Fatalln(string(jsonBody))
		log.Fatalln(err)
	}
	languages := make([]string, 0, len(languageMap))
	for key := range languageMap {
		languages = append(languages, key)
	}
	r.Languages = languages
}

func runService(path string, devEnv bool) (ResponseData, error) {
	var serviceResponse ResponseData
	var jsonBodies RequestBodies
	if devEnv {
		//this next line is only to test if a given path is valid
		_, serviceError := getRepoPath(path)
		if serviceError != nil {
			return ResponseData{}, serviceError
		}
		jsonBodies = setupLocal() //local service only uses DefaultRepoPath
	} else {
		log.Println("runing remotely!")
		repoPath, serviceError := getRepoPath(path)
		if serviceError != nil {
			return ResponseData{}, serviceError
		}
		jsonBodies = setupRemote(repoPath)
	}
	getProjectInfo(jsonBodies.Project, &serviceResponse)
	getContributorInfo(jsonBodies.Contributors, &serviceResponse)
	getLanguageInfo(jsonBodies.Languages, &serviceResponse)
	return serviceResponse, nil
}

func serviceHandler(writer http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		devEnv, _ := strconv.ParseBool(os.Getenv("DEVENV"))
		writer.Header().Add("content-type", "application/json")
		serviceResponse, err := runService(r.URL.Path, devEnv)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest)+"\n"+err.Error(), http.StatusBadRequest)
		}
		json.NewEncoder(writer).Encode(serviceResponse)
	} else {
		http.Error(writer, "Only GET requests allowed", http.StatusMethodNotAllowed)
	}
}

func getRepoPath(originalPath string) (string, error) {
	URLPath := originalPath
	pathVars := strings.Split(URLPath, "/")
	if len(pathVars) != 6 {
		return DefaultRepoPath, errors.New("Invalid URL PATH: " + URLPath)
	}
	//@TODO regex filtering username

	URLPath = strings.TrimPrefix(URLPath, ServiceBasePath)
	return strings.Replace(URLPath, "github.com", "repos", 1), nil
}

func getRequestBody(repoPath string) []byte {
	response, err := http.Get(GithubAPIHostURL + repoPath)
	if err != nil {

	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {

	}
	response.Body.Close()
	return body
}

func setupLocal() RequestBodies {
	var requestBodies RequestBodies

	requestBodies.Project = readLocalWorkingFile("project")
	requestBodies.Languages = readLocalWorkingFile("languages")
	requestBodies.Contributors = readLocalWorkingFile("contributors")

	return requestBodies
}

func setupRemote(repoPath string) RequestBodies {
	return RequestBodies{
		Project:      getRequestBody(repoPath),
		Languages:    getRequestBody(repoPath + "/languages"),
		Contributors: getRequestBody(repoPath + "/contributors"),
	}
}

func readLocalWorkingFile(filename string) []byte {
	devPath := os.Getenv("LOCALPROJECT")
	data, err := ioutil.ReadFile(devPath + filename + ".json")
	if err != nil {
		log.Fatalf("Could not open local working file:%s", devPath+filename+".json")
	}
	return data
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.HandleFunc(ServiceBasePath, serviceHandler)
	http.ListenAndServe(":"+port, nil)
}
