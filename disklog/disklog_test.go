package disklog_test

import (
	"fmt"
	. "github.com/syncsynchalt/dime-a-tap/disklog"
	"github.com/syncsynchalt/dime-a-tap/test"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDumpPacketEmptyDir(t *testing.T) {
	err := DumpPacket("", "label", []byte("\x01\x02\x03"))
	test.Ok(t, err)
}

func TestDumpPacketNonExistDir(t *testing.T) {
	err := DumpPacket("/does/not/exist", "label", []byte("\x01\x02\x03"))
	test.Assert(t, err != nil, "error is not set")
	prefix := "unable to log packet: open /does/not/exist/label.20"
	test.Assert(t, strings.HasPrefix(err.Error(), prefix),
		"error [%s] does not have expected prefix [%s]", err, prefix)
}

func getFirstFile(t *testing.T, dir string) string {
	files, err := ioutil.ReadDir(dir)
	test.Ok(t, err)
	t.Log(files)
	return dir + "/" + files[0].Name()
}

func TestDumpPacket(t *testing.T) {
	mydir := fmt.Sprintf("/tmp/golang.test.%d", time.Now().UnixNano())
	os.Mkdir(mydir, 0755)
	defer os.RemoveAll(mydir)
	err := DumpPacket(mydir, "[::1]:80134", []byte("\x01\x02\x03"))

	file := getFirstFile(t, mydir)
	test.Assert(t, strings.HasPrefix(file, mydir + "/[::1]:80134.20"),
		"file %s doesn't start with expected prefix", file)

	data, err := ioutil.ReadFile(file)
	test.Ok(t, err)
	test.Equals(t, []byte("\x01\x02\x03"), data)
}
