package main

import (
	"encoding/json"
	"log"
	"testing"
)

func Test_JSONUnmarshalProjectInfo(t *testing.T) {
	const jsonBody = `{
        "name": "Project1",
        "owner": {
            "login": "Bob",
            "id": 12743
        }
    }`
	// own := GithubUser{Username: "Bob", ID: 12743}
	expectedResult := ProjectInfo{Name: "Project1", Owner: GithubUser{Username: "Bob", ID: 12743}}
	var result ProjectInfo
	if err := json.Unmarshal([]byte(jsonBody), &result); err != nil {
		t.Error(err)
		t.Error("Unmarshalling did not succeed")
		log.Println(string(jsonBody))
		log.Println(err)
	}
	if result != expectedResult {
		t.Errorf("Result:\t%s\tnot matching expected:\t%s\n", result, expectedResult)
	}
}

func Test_JSONMarshallResponse(t *testing.T) {
	const expectedResult = `{"project":"kafka","owner":"apache","committer":"enothereska","commits":19,"language":["Java","Scala","Python","Shell","Batchfile"]}`
	// exp := strings.TrimSpace(expectedResult)
	responseBody := ResponseData{
		Project:   "kafka",
		Owner:     "apache",
		Committer: "enothereska",
		Commits:   19,
		Languages: []string{"Java", "Scala", "Python", "Shell", "Batchfile"}}
	result, err := json.Marshal(responseBody)
	if err != nil {
		t.Error(err)
	}
	if string(result) != string(expectedResult) {
		t.Errorf("Result:\t%s\tnot matching expected:\t%s\n", result, expectedResult)
	}
}

func Test_getRepoPathSuccess(t *testing.T) {
	testPath := "/projectinfo/v1/github.com/Stektpotet/Amazeking"
	expectedRepoPath := "repos/Stektpotet/Amazeking"
	resultRepoPath := getRepoPath(testPath)
	if resultRepoPath != expectedRepoPath {
		t.Errorf("%s not matching expected %s", resultRepoPath, expectedRepoPath)
	}

}

func Test_getRepoPathFail(t *testing.T) {
	testPath := "/Stektpotet/Amazeking"
	resultRepoPath := getRepoPath(testPath)
	if resultRepoPath != DEFAULT_REPO_PATH {
		t.Errorf("%s not matching defaulting path %s", resultRepoPath, DEFAULT_REPO_PATH)
	}

}
