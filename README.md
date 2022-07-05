# Introduction 
cli client to execute Azure pipelines over REST API. The pipeline state is being tracked automatically during every run and the result is being updated after every run. Integrates with Azure pipeline parameters.

## Non required prereqs
for structured formatting the json output install the jq utillity 
```bash
# use the Snap package manager
sudo snap install jq
# or use the apt package manager
sudo apt install jq
```

## Set environment variables (required)
```bash
# Linux Bash
export PAT='Azure DevOps personal access token'
export ORGANIZATION='ORG_NAME'
export PROJECT='PROJECT_NAME'
export PIPELINE_ID='999999'
```

```powershell
# Windows PowerShell
$env:PAT = 'Azure DevOps personal access token'
$env:ORGANIZATION = 'ORG_NAME'
$env:ADOPROJECT = 'PROJECT_NAME'
$env:ADOPIPELINEID = '999999'
```

## Build and Run
```bash
#clone repository and change dir
git clone 'https://github.com/bartvanbenthem/ado-runsapi-client.git'
cd ado-runsapi-client
# build run-pipeline
go build -o bin ./cmd/run-pipeline
# build get azure pipeline ID based on pipeline name
go build -o bin ./cmd/run-status


# run-pipeline and print full output
./bin/run-pipeline | jq .
# or run-pipeline and print only the run ID
./bin/run-pipeline --printid='true'

# check pipeline status for completion
./bin/run-status --runid='999999'

# Oneliner for Starting a new run of the pipeline and keep track of the current state
RUNID=$(./bin/run-pipeline \
            --printid='true' \
            --parameters="{\"param1\": \"myvalue-1\", \"param2\": \"golang rules\"}") \
&& ./bin/run-status --runid=$RUNID

```

### CURL examples
```bash
## Test if authentication working (List projects in your organization)
curl -s -u $USER:$PAT \
"https://dev.azure.com/$ORGANIZATION/_apis/projects?api-version=6.0" | jq .

# list all pipelines example
curl -s --request GET \
-u 'runs':$PAT \
--header "Content-Type: application/json" \
"https://dev.azure.com/$ORGANIZATION/$PROJECT/_apis/pipelines?api-version=7.1-preview.1" | jq .

# run pipeline example
curl -s --request POST \
-u 'runs':$PAT \
--header "Content-Type: application/json" \
--data '{
    "resources": {
        "repositories": {
            "self": {
                "refName": "refs/heads/main"
            }
        }
    },
   "templateParameters": {
        "param1": "myvalue1",
        "param2": "myvaluex"
    }
}' \
"https://dev.azure.com/$ORGANIZATION/$PROJECT/_apis/pipelines/$PIPELINE_ID/runs?api-version=6.0-preview.1" | jq .

# get pipeline run status
RUN_ID='999999'
curl -s --request GET \
-u 'runs':$PAT \
--header "Content-Type: application/json" \
"https://dev.azure.com/$ORGANIZATION/$PROJECT/_apis/pipelines/$PIPELINE_ID/runs/$RUN_ID?api-version=6.0-preview.1" | jq .

```
