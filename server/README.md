# Infrastructure Restful


STEP 1
use similar command to rename the project to different project
Real Case: [in the server directory]
for d in $(grep -rnwl github.com/syronz/infrastructure); do sed -i 's/github.com\/syronz\/infrastructure/github.com\/syronz\/PROJECT_NAME/g' $d; done

STEP 2
dep ensure

STEP 3
create custome database

STEP 4
edit config/app.toml and adding custome database and connection credentials

STEP 5 
go run main.go -migrate -reset

SETP 6 (for watching edits)
realize start -run -no-config

STEP 7
use generated username and password in step 6 to login to the system

Note: 
inside helper directory you can find insomnia json for how to connect with api




