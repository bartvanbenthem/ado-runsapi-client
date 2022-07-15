package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runs/pkg/rest"
	"strconv"
	"time"
)

var (
	// get environment variables for URI construction
	organization = os.Getenv("ORGANIZATION")
	project      = os.Getenv("PROJECT")
	pipelineid   = os.Getenv("PIPELINE_ID")
	// get environment variables for basic auth
	username = os.Getenv("USER")
	password = os.Getenv("PAT")
	// set global variables for storing flag input
	jstr  *string
	watch *string
)

func runPipeline(p rest.URIParamaters, jstr *string) ([]byte, error) {
	// construct the URL for the RUNS request
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/pipelines/%d/runs?api-version=6.0-preview.1",
		p.Organization, p.Project, p.PipelineID)

	// building the requestbody and add the refname
	requestBody := rest.RequestBody{}
	requestBody.Resources.Repositories.Self.RefName = "refs/heads/main"

	// unmarshall the json parameter input into the requestbody struct type
	//example-input := "{\"param1\": \"myvalue-1\", \"param2\": \"myvalue-x\"}"
	var params interface{}
	json.Unmarshal([]byte(*jstr), &params)
	requestBody.TemplateParameters = params
	// marshall the requestbody struct type into a json object
	requestBodyJSON, err := json.Marshal(&requestBody)
	if err != nil {
		log.Fatal(err)
	}

	//add the personal access token to the header for basic authentication
	auth := rest.BasicAuth{Username: username, Password: password}
	// make authenticated post request to the RUNS API
	body, err := rest.Call(url, "POST", requestBodyJSON, auth)
	if err != nil {
		return body, err
	}

	return body, nil
}

func watchPipeline(p rest.URIParamaters, runid int32) {
	// construct the URL for the RUNS request
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/pipelines/%d/runs/%d?api-version=6.0-preview.1",
		p.Organization, p.Project, p.PipelineID, runid)

	// initialize the response struct
	response := rest.Response{}
	//add the personal access token to the header for basic authentication
	auth := rest.BasicAuth{Username: username, Password: password}
	// make authenticated get request to the RUNS API
	body, _ := rest.Call(url, "GET", nil, auth)
	// unmarshall response body into the response struct
	json.Unmarshal(body, &response)
	// print the initial pipeline run status
	fmt.Printf("Pipeline \"%s\" \"%s\" run %d\"", response.Name, response.State, response.ID)
	// while the pipeline run is in progress make a new request every 5 seconds
	// and update the current state of the pipeline run
	for response.State == "inProgress" {
		body, _ = rest.Call(url, "GET", nil, auth)
		json.Unmarshal(body, &response)
		// progress bar
		fmt.Printf(".")
		time.Sleep(5 * time.Second)
	}
	// when the pipeline run is not in prorgress anymore
	// print the state and result to stdout
	fmt.Printf("\nPipeline %s with ID %d is %s and %s\n",
		response.Name, response.ID, response.State, response.Result)
}

func setParameters() {
	// azure devops flags
	token := flag.String("token", "", "set Azure DevOps personal access token")
	proj := flag.String("project", "", "set Azure DevOps project name")
	org := flag.String("organization", "", "set Azure DevOps organization")
	pipe := flag.String("pipelineid", "", "set Azure DevOps pipeline ID")
	// util flags
	j := flag.String("parameters", "{}", "add template parameters in serialized json text")
	w := flag.String("watch", "false", "set to true if you wish to track the run status")
	// parse flags
	flag.Parse()
	// set global variables
	jstr = j
	watch = w

	if len(*token) != 0 {
		password = *token
	}
	if len(password) == 0 {
		log.Fatal("ERROR: token parameter is required \n")
	}

	if len(*proj) != 0 {
		project = *proj
	}
	if len(project) == 0 {
		log.Fatal("ERROR: project parameter is required \n")
	}

	if len(*org) != 0 {
		organization = *org
	}
	if len(organization) == 0 {
		log.Fatal("ERROR: organization parameter is required \n")
	}

	if len(*pipe) != 0 {
		pipelineid = *pipe
	}
	if len(pipelineid) == 0 {
		log.Fatal("ERROR: pipelineid parameter is required \n")
	}

}

func StringToInt32(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}

	return int32(i)
}

func StringToBool(s string) bool {
	i, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatal(err)
	}

	return bool(i)
}

func main() {
	// set parameters and check env variables when flags are not specified
	setParameters()
	// initialize the URIParameters struct fo constructing the URL
	p := rest.URIParamaters{
		Organization: organization,
		Project:      project,
		PipelineID:   StringToInt32(pipelineid),
	}

	// run pipeline
	body, err := runPipeline(p, jstr)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	// if watch is set to true start the watch pipeline function
	// else print the response body
	// convert the watch parameter *string to a boolean type
	if StringToBool(*watch) == true {
		response := rest.Response{}
		json.Unmarshal(body, &response)
		watchPipeline(p, int32(response.ID))
	} else {
		fmt.Printf("%s", body)
	}

}
