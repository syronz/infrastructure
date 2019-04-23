# Infrastructure Restful

STEP 1
use similar command to rename the project to different project
Real Case: [in the server directory]
for d in $(grep -rnwl github.com/syronz/infrastructure); do sed -i 's/github.com\/syronz\/infrastructure/github.com\/syronz\/PROJECT_NAME/g' $d; done

STEP 2
dep ensure




