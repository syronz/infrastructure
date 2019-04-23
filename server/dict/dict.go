// Dict package implement dictionery plus translate messages to the proper language.
// For now three languages have been supported. If message not exist in words.toml
// parameter will being returned without any changes but undefined message appended to 
// the dict/not-translated.
package dict

import (
	"os"
	"bufio"
	"runtime"
	"path"
	"log"

	"github.com/syronz/infrastructure/server/app"
	"github.com/syronz/infrastructure/server/utils/debug"
	"github.com/BurntSushi/toml"
	"github.com/syronz/ozzo-routing"
)

type (
	// Params is used to replace placeholders inside messages. for example "User %s successfully
	// added to the database." %s will be replaced with param
	Params map[string]interface{}

	// Languages which have been supported
	dicTemplate struct {
		En	string `toml:"en"`
		Ku	string `toml:"ku"`
		Ar	string `toml:"ar"`
	}
)

var Words map[string]dicTemplate

// This part will be usefull inside _test files
func init() {
	LoadWords()
}

// Load messages from words.toml file and add them to the defined map.
// In case of words.toml not exist or can't be opened the program will halted
func LoadWords() {
	_, file, _, _ := runtime.Caller(0)
	dir, _ := path.Split(file)
	file = dir + "words.toml"
	Words = map[string]dicTemplate{}
	_, err := toml.DecodeFile(file, &Words)
	if err != nil {
		debug.Log(err)
		log.Fatal(err)
	}
}

// Translate message to the default language.
//
// TODO: after user completed this part must translate based on user's language.
func T(c *routing.Context,str string) string{
	defaultLanguage := app.Config.DefaultLanguage
	userLang := app.GetRequestScope(c).UserLanguage()

	_, ok := Words[str]
	if !ok {
		addToNotTranslated(str)
	} else {
		switch userLang {
		case "en":
			return Words[str].En
		case "ku":
			return Words[str].Ku
		case "ar":
			return Words[str].Ar

			default :
			switch defaultLanguage {
			case "en":
				return Words[str].En
			case "ku":
				return Words[str].Ku
			case "ar":
				return Words[str].Ar
			default:
				return Words[str].En
			}
		}
	}

	return str
}

// If message not defined inside in the words.toml the message will added to the 
// not-translated.txt for extending words.toml afterward
func addToNotTranslated(str string){
	_, file, _, _ := runtime.Caller(0)
	dir, _ := path.Split(file)
	dir = path.Clean(dir + "/..")

	_, isExist := os.Stat(dir + "/log/not-translated.txt")
	var err error
	var f *os.File

	if os.IsNotExist(isExist){
		f, err = os.OpenFile(dir + "/log/not-translated.txt", os.O_CREATE, 0644)
	} else {
		f, err = os.OpenFile(dir + "/log/not-translated.txt", os.O_RDWR, 0644)
	}


	if err != nil {
		debug.Log(dir,err)
		log.Fatal(err)
	}
	defer f.Close()

	var isInserted bool

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == str {
			isInserted = true
			break
		}
	}

	// is file inserted before to the file didn't add it
	if !isInserted {
		if _, err = f.WriteString(str + "\n"); err != nil {
			debug.Log(err)
		}
	}

}
