# go-whosonfirst-remap

## Important

This code is for internal use only and will be deleted soon.

## Misc tools

### Create all the repos

```
#!/bin/sh

# https://developer.github.com/v3/repos/#create

CURL=`which curl`
GITHUB_TOKEN=''

for REPO in $@
do

    REPO_NAME=`basename ${REPO}`
    COUNTRY=`echo ${REPO_NAME} | awk -F '-' '{ print $3 }' | tr 'a-z' 'A-Z'`

    echo '{ "name": "'${REPO_NAME}'", "description": "Who's On First admin data for ${COUNTRY}", "homepage": "https://whosonfirst.org'", "private": false, "has_issues": true, "has_projects": true, "has_wiki": true }' > repo.json

    ${CURL} -H "Authorization: token ${GITHUB_TOKEN}" -X POST --data @repo.json https://api.github.com/orgs/whosonfirst-data/repos    

done
```