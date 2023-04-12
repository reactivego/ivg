package icons

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// overwriteTestdataFiles is temporarily set to true when adding new
// testdataTestCases.
const overwriteTestdataFiles = false

// TestOverwriteTestdataFilesIsFalse tests that any change to
// overwriteTestdataFiles is only temporary. Programmers are assumed to run "go
// test" before sending out for code review or committing code.
func TestOverwriteTestdataFilesIsFalse(t *testing.T) {
	if overwriteTestdataFiles {
		t.Errorf("overwriteTestdataFiles is true; do not commit code changes")
	}
}

func TestHashes(t *testing.T) {
	got := make(map[string]string)
	for i := range list {
		checksum := md5.Sum(list[i].data)
		got[list[i].name] = base64.RawURLEncoding.EncodeToString(checksum[:])
	}
	if overwriteTestdataFiles {
		out, err := os.Create("icons.json")
		if err != nil {
			t.Fatal(err)
		}
		enc := json.NewEncoder(out)
		enc.SetIndent("", " ")
		err = enc.Encode(got)
		out.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
	in, err := os.Open("icons.json")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	dec := json.NewDecoder(in)
	expect := make(map[string]string)
	err = dec.Decode(&expect)
	if err != nil {
		t.Fatal(err)
	}
	for ke, ve := range expect {
		if vg, ok := got[ke]; !ok {
			t.Errorf("icon:%q expected, but not present", ke)
		} else {
			if !strings.EqualFold(ve, vg) {
				t.Errorf("icon:%q, md5 expected:%q, got:%q", ke, ve, vg)
			}
		}
	}
	for kg := range got {
		if _, ok := expect[kg]; !ok {
			t.Errorf("icon:%q present, but not expected", kg)
		}
	}
}
