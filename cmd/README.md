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
