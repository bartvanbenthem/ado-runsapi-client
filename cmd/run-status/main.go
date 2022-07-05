package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
)

func main() {
	// parse cli input parameters
	runid := flag.String("runid", "999", "set to true if you wish to only output the pipeline run ID")
	flag.Parse()
	// convert the printid *string to a Int32 type
	input := convert.StringPointerToInt32(runid)
	// initialize the URIParameters struct fo constructing the URL
	p := rest.URIParamaters{
		Organization: organization,
		Project:      project,
		PipelineID:   convert.StringToInt32(pipelineid),
	}
	// construct the URL for the RUNS request
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/pipelines/%d/runs/%d?api-version=6.0-preview.1",
		p.Organization, p.Project, p.PipelineID, input)

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
