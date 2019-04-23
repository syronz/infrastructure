// Package debug used for dump variables. It depends on app.toml and check debug parameter.
// if debug was true the log and other functions do their jobs.
package debug

import (
	"github.com/syronz/infrastructure/server/app"
	"github.com/davecgh/go-spew/spew"
	"fmt"
	"os"
	"log"
	"runtime"
	"path"
	"time"
)

// Fucntion for dumpt any type of variable with address of file and line of log. It
// depends to the debug.
func Log(v ...interface{}){
	pc, file, line, ok := runtime.Caller(1)
	_, fileName := path.Split(file)

	saveToFile(file, line, v)

	if(app.Config.Debug) {
		if ok {
			fmt.Printf("───● %v:%v \nCalled from %s, line #%d\n───○ func: %v\n",
			fileName, line, file, line, runtime.FuncForPC(pc).Name())
		}
		spew.Dump(v)
	}


}

// log dumps inside file
func saveToFile(file string, line int, v ...interface{}){
	_, fileToWrite, _, _ := runtime.Caller(0)
	dir, _ := path.Split(fileToWrite)
	dir = path.Clean(dir + "/../..")

	_, isExist := os.Stat(dir + "/log/dump.log")
	var err error
	var f *os.File

	if os.IsNotExist(isExist){
		f, err = os.OpenFile(dir + "/log/dump.log", os.O_CREATE, 0644)
	} else {
		f, err = os.OpenFile(dir + "/log/dump.log", os.O_APPEND|os.O_WRONLY, 0644)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	currentTime := time.Now()
	f.WriteString("───○ " + currentTime.Format("2006-01-02 15:04:05") + "\n")
	f.WriteString("───● " + file + ": " + fmt.Sprint(line) + "\n")
	spew.Fdump(f, v)

}
