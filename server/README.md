# Infrastructure Restful

use similar command to rename the project to different project
for d in $(grep -rnwl github.com/syronz/PRE_NAME); do sed -i 's/github.com\/syronz\/PRE_NAME/github.com\/syronz\/infrastructure\/server/g' $d; done

