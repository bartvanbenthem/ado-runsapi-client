package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runs/pkg/convert"
	"runs/pkg/rest"
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
	jstr     *string
	watch    *string
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

func main() {
	// set parameters and / or env variables
	setParameters()

	// convert the printid *string to a boolean type
	//o := convert.StringPointerToBool(printid)
	w := convert.StringPointerToBool(watch)
	// initialize the URIParameters struct fo constructing the URL
	p := rest.URIParamaters{
		Organization: organization,
		Project:      project,
		PipelineID:   convert.StringToInt32(pipelineid),
	}

	// run pipeline
	body, err := runPipeline(p, jstr)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	// if printid == truw only print the runid to stdout
	// else print the entire json response object to stdout
	var runid int32

	if w == true {
		response := rest.Response{}
		json.Unmarshal(body, &response)
		runid = int32(response.ID)
		watchPipeline(p, runid)
	} else {
		fmt.Printf("%s", body)
	}

}