package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BasicAuth struct {
	Username string
	Password string
}

type URIParamaters struct {
	Organization    string `json:"organization"`
	PipelineID      int32  `json:"pipelineId"`
	Project         string `json:"project"`
	APIVersion      string `json:"api-version"`
	PipelineVersion int32  `json:"pipelineVersion"`
}

type RequestBody struct {
	PreviewRun bool `json:"previewRun,omitempty"`
	Resources  struct {
		Builds       string `json:"builds,omitempty"`
		Containers   string `json:"containers,omitempty"`
		Packages     string `json:"packages,omitempty"`
		Pipelines    string `json:"pipelines,omitempty"`
		Repositories struct {
			Self struct {
				RefName string `json:"refName"`
			} `json:"self"`
		} `json:"repositories"`
	} `json:"resources"`
	StagesToSkip       []string    `json:"stagesToSkip,omitempty"`
	TemplateParameters interface{} `json:"templateParameters,omitempty"`
	Variables          interface{} `json:"variables"`
	YamlOverride       interface{} `json:"yamlOverride,omitempty"`
}

type Response struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Web struct {
			Href string `json:"href"`
		} `json:"web"`
		PipelineWeb struct {
			Href string `json:"href"`
		} `json:"pipeline.web"`
		Pipeline struct {
			Href string `json:"href"`
		} `json:"pipeline"`
	} `json:"_links"`
	Pipeline struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		Revision int    `json:"revision"`
		Name     string `json:"name"`
		Folder   string `json:"folder"`
	} `json:"pipeline"`
	State       string    `json:"state"`
	Result      string    `json:"result"`
	CreatedDate time.Time `json:"createdDate"`
	URL         string    `json:"url"`
	Resources   struct {
		Repositories struct {
			Self struct {
				Repository struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"repository"`
				RefName string `json:"refName"`
				Version string `json:"version"`
			} `json:"self"`
		} `json:"repositories"`
	} `json:"resources"`
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Call(url, method string, requestBody []byte, auth BasicAuth) ([]byte, error) {

	var body []byte

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return body, fmt.Errorf("Error: %s", err.Error())
	}

	req.SetBasicAuth(auth.Username, auth.Password)
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return body, fmt.Errorf("Error: %s", err.Error())
	}

	body, err = io.ReadAll(response.Body)

	defer response.Body.Close()

	return body, nil
}

func FmtJsonOutput(str string) (string, error) {
	var fmtJson bytes.Buffer
	if err := json.Indent(&fmtJson, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return fmtJson.String(), nil
}
