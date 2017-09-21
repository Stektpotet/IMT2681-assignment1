package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

func Test_JSONUnmarshalProjectInfo(t *testing.T) {
	const jsonBody = `{
        "full_name": "Bob/Project1",
        "owner": {"login": "Bob"}
    }`
	expectedResult := ProjectInfo{Name: "Bob/Project1", Owner: GithubUser{Username: "Bob"}}
	var result ProjectInfo
	if err := json.Unmarshal([]byte(jsonBody), &result); err != nil {
		t.Error(err)
		t.Error("Unmarshalling did not succeed")
		log.Println(string(jsonBody))
		log.Println(err)
	}
	if result != expectedResult {
		t.Errorf("Result:\t%+v\tnot matching expected:\t%+v\n", result, expectedResult)
	}
}

func Test_JSONMarshallResponse(t *testing.T) {
	const expectedResult = `{"project":"kafka","owner":"apache","committer":"enothereska","commits":19,"language":["Java","Scala","Python","Shell","Batchfile"]}`
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

//==============================GET REPO PATH TESTS===========================

var GetRepoPathTest = func(in string, expectedOut string, expectError bool, t *testing.T) {
	resultRepoPath, err := getRepoPath(in)
	if resultRepoPath != expectedOut {
		t.Errorf("\"%s\" not matching expected \"%s\"", resultRepoPath, expectedOut)
	}
	if expectError {
		if err == nil {
			t.Errorf("expected error because %s is invalid as path", in)
		}
	} else {
		if err != nil {
			t.Errorf("unexpected error occured %s! %s is invalid as path", err.Error(), in)
		}
	}
}

func Test_GetRepoPathDefault(t *testing.T) {
	GetRepoPathTest(ServiceBasePath+DefaultRepoPath, DefaultRepoPath, false, t)
}

func Test_GetRepoPathEmpty(t *testing.T) {
	GetRepoPathTest("", DefaultRepoPath, true, t)
}

func Test_GetRepoPathNonValid(t *testing.T) {
	GetRepoPathTest("cov/fe/fe", DefaultRepoPath, true, t)
}

//==============================SERVICE HANDLER REQUEST TESTS =================

func Test_ServiceHandlerSEND(t *testing.T) {
	reqest, _ := http.NewRequest("SEND", ServiceBasePath+DefaultRepoPath, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serviceHandler)
	handler.ServeHTTP(rr, reqest)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func Test_ServiceHandlerGET(t *testing.T) {
	reqest, _ := http.NewRequest("GET", ServiceBasePath+DefaultRepoPath, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serviceHandler)
	handler.ServeHTTP(rr, reqest)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

//=================================FULL SERVICE TESTS==========================

var FullServiceTest = func(path string, devEnv bool, expectError bool, t *testing.T) {
	serviceResponse, err := runService(path, devEnv)
	if expectError {
		if err == nil {
			t.Errorf("Error did not occur when error was expected!\n response is instead %+v", serviceResponse)
		}
	} else {
		if err != nil {
			t.Errorf("unexpected error occured \"%s\" when running service.", err.Error())
		}
		if serviceResponse.Project == "" {
			t.Errorf("service returns project with no name!\n%+v", serviceResponse)
		}
		if serviceResponse.Owner == "" {
			t.Errorf("service returns project wit no owner!\n%+v", serviceResponse)
		}
	}
}

func Test_RunServiceRemote(t *testing.T) {
	FullServiceTest(ServiceBasePath+DefaultRepoPath, true, false, t)
}

func Test_RunServiceLocal(t *testing.T) {
	FullServiceTest(ServiceBasePath+DefaultRepoPath, true, false, t)
}
func Test_RunServiceLocalBadRequest(t *testing.T) {
	FullServiceTest("", true, true, t)
}
