# Introduction 
cli client to execute Azure pipelines over REST API with personal acces token authentication. when specified the pipeline run state is being tracked automatically and the result is being updated after every run. Integrates with Azure pipeline parameter specifications.

## Non required prereqs
for structured formatting the json output install the jq utillity 
```bash
# use the Snap package manager
sudo snap install jq
# or use the apt package manager
sudo apt install jq
```

## usage
```bash
Usage of ado-runsapi:
  -organization string
        set Azure DevOps organization
  -parameters string
        add template parameters in serialized json text (default "{}")
  -pipelineid string
        set Azure DevOps pipeline ID
  -project string
        set Azure DevOps project name
  -token string
        set Azure DevOps personal access token
  -watch string
        set to true if you wish to track the run status (default "false")
```

## Set environment variables (instead of flags)
```bash
# Linux Bash
export PAT='Azure DevOps personal access token'
export ORGANIZATION='ORG_NAME'
export PROJECT='PROJECT_NAME'
export PIPELINE_ID='999999'
```

```powershell
# PowerShell
$env:PAT = 'Azure DevOps personal access token'
$env:ORGANIZATION = 'ORG_NAME'
$env:ADOPROJECT = 'PROJECT_NAME'
$env:ADOPIPELINEID = '999999'
```

## Build and Run
```bash
# clone repository and change dir
git clone 'https://github.com/bartvanbenthem/ado-runsapi-client.git'
cd ado-runsapi-client

# build ado-runsapi binary
go build -o bin ./cmd/ado-runsapi

# set env variables (personal setup)
source ../00-ENV/env.sh

# execute pipeline with parameters and continue when the response has been received
./bin/ado-runsapi --parameters="{\"param1\": \"myvalue-1\", \"param2\": \"golang rules\"}" | jq .

# execute pipeline with parameters and track pipeline run state untill run is no longer in progress
./bin/ado-runsapi --watch='true' --parameters="{\"param1\": \"myvalue-1\", \"param2\": \"golang rules\"}"

# run ado-runsapi with flags instead of environment variables
./bin/ado-runsapi --token=$PAT --project=$PROJECT --organization=$ORGANIZATION --pipelineid=$PIPELINE_ID \
    --watch='true' --parameters="{\"param1\": \"myvalue-1\", \"param2\": \"golang rules\"}"

```