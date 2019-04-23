package dict

import (
	"testing"
	"github.com/syronz/infrastructure/server/app"
	"fmt"
)

/* Words Sample
["City with this name entered before"]
en = "City with this name entered before"
ku = "Shar baw nawa peshter daxl krawa"
ar = "Al madina with this name already defined"

["Internal Error"]
en = "Internal Error"
ku = "Keshay nawxo"
ar = "Aldaxeli maslat"

[host3]
en = "Internal Error"
ku = "Keshay nawxo"
ar = "Aldaxeli maslat"
*/
func TestTranslateTableKu(t *testing.T) {
	if err := app.LoadConfig("../config"); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}
	app.Config.DefaultLanguage = "Ku"

	var tests = []struct {
		s, want string
	}{
		{"City with this name entered before", "Shar baw nawa peshter daxl krawa"},
		{"Internal Error", "Keshay nawxo"},
	}

	for _, c := range tests {
		//t.Log(c, T(c.s))
		t.Log("$$$$$$$$$$$$$$$$", app.Config.DefaultLanguage)
		if T(c.s) != c.want {
			t.Errorf("Supposed to return %q for %q but returned %q", c.want, c.s, T(c.s))
		}
	}
}
