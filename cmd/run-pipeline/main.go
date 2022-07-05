package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runs/pkg/convert"
	"runs/pkg/rest"
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
	printid := flag.String("printid", "false", "set to true if you wish to only output the pipeline run ID")
	parameters := flag.String("parameters", "{}", "add template parameters in object notation (json)")
	flag.Parse()
	// convert the printid *string to a boolean type
	o := convert.StringPointerToBool(printid)
	// initialize the URIParameters struct fo constructing the URL
	p := rest.URIParamaters{
		Organization: organization,
		Project:      project,
		PipelineID:   convert.StringToInt32(pipelineid),
	}
	// construct the URL for the RUNS request
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/pipelines/%d/runs?api-version=6.0-preview.1",
		p.Organization, p.Project, p.PipelineID)

	// building the requestbody and add the refname
	requestBody := rest.RequestBody{}
	requestBody.Resources.Repositories.Self.RefName = "refs/heads/main"

	// unmarshall the json parameter input into the requestbody struct type
	//example-input := "{\"param1\": \"myvalue-1\", \"param2\": \"myvalue-x\"}"
	var params interface{}
	json.Unmarshal([]byte(*parameters), &params)
	requestBody.TemplateParameters = params
	// marshall the requestbody struct type into a json object
	requestBodyJSON, err := json.Marshal(&requestBody)
	if err != nil {
		log.Fatal(err)
	}

	//add the personal access token to the header for basic authentication
	auth := rest.BasicAuth{Username: username, Password: password}
	// make authenticated post request to the RUNS API
	body, _ := rest.Call(url, "POST", requestBodyJSON, auth)

	// if printid == truw only print the runid to stdout
	// else print the entire json response object to stdout
	if o == true {
		response := rest.Response{}
		json.Unmarshal(body, &response)
		fmt.Printf("%d", response.ID)
	} else {
		fmt.Printf("%s", body)
	}

}
