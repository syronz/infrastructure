package app

import (
	"runtime"
	"path"
	"log"

	"github.com/BurntSushi/toml"
)

type aclTemplate struct {
	Read		[]string	`toml:"read"`
	Write		[]string	`toml:"write"`
}

var Acl map[string]aclTemplate

func LoadAcl() {
	_, file, _, _ := runtime.Caller(0)
	dir, _ := path.Split(file)
	dir = path.Clean(dir + "/..")
	file = dir + "/config/acl.toml"
	Acl = map[string]aclTemplate{}

	if _, err := toml.DecodeFile(file, &Acl); err != nil {
		log.Fatal(err)
	}
}

