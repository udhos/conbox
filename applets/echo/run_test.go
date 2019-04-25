package echo

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/udhos/conbox/common"
)

type testEcho struct {
	input  string
	status int
	output string
}

var testTable = []testEcho{
	{"", 0, "\n"},
	{"   ", 0, "\n"},
	{"a b", 0, "a b\n"},
	{"   a    b   ", 0, "a b\n"},
	{"-n", 0, ""},
	{"-n   ", 0, ""},
	{"-n a b", 0, "a b"},
	{"-n   a    b   ", 0, "a b"},
}

func TestEcho1(t *testing.T) {
	for _, data := range testTable {
		args := common.Tokenize(data.input)
		echo1(t, args, data.status, data.output)
	}
}

func echo1(t *testing.T, args []string, wantStatus int, wantOutput string) {

	tmp, errTmp := ioutil.TempFile("", "echo1-*.txt")
	if errTmp != nil {
		t.Errorf("tmp file: %v", errTmp)
		return
	}

	stdout := os.Stdout // save

	os.Stdout = tmp

	status := Run(nil, args)

	os.Stdout = stdout // restore

	if status != wantStatus {
		t.Errorf("unexpected exit status: want=%d got=%d", wantStatus, status)
	}

	// flush
	if errClose := tmp.Close(); errClose != nil {
		t.Errorf("tmp close: %v", errClose)
	}

	out, errRead := ioutil.ReadFile(tmp.Name())
	if errRead != nil {
		t.Errorf("tmp read: %v", errRead)
	}

	if errRemove := os.Remove(tmp.Name()); errRemove != nil {
		t.Errorf("tmp remove: %v", errRemove)
	}

	str := string(out)
	if str != wantOutput {
		t.Errorf("unexpected output: want=[%s] got=[%s]", wantOutput, str)
	}
}
