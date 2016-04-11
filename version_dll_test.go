package windll

import (
	"path/filepath"
	"fmt"
	"os"
	"os/exec"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func runCommandsInSameDir(dir string,cmds ...*exec.Cmd) error {
	for _, c := range cmds {
        c.Dir = dir
		out, err := c.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Error in cmd (%s - %s). ERROR: %s. OUTPUT: %s", c.Path, c.Args, err.Error(), string(out))
		}
	}

	return nil
}

func TestExtractProductVersion(t *testing.T) {
	Convey("Testing extraction of the Product Version", t, func() {
		dummyDir := os.ExpandEnv("$GOPATH/src/github.com/golang-devops/windll/dummy_exe")
        exeFileName := "dummy.exe"
        dummyExeFullPath := filepath.Join(dummyDir, exeFileName)
        resourceSysoFullPath := filepath.Join(dummyDir, "resource.syso")

        err := runCommandsInSameDir(
            dummyDir,
            exec.Command("go", "generate"),
			exec.Command("go", "build", "-o", exeFileName),
		)
        defer os.Remove(dummyExeFullPath)
        defer os.Remove(resourceSysoFullPath)
		So(err, ShouldBeNil)

		productVersion, err := VersionDLL.ExtractProductVersion(dummyExeFullPath)
		So(err, ShouldBeNil)
		So(productVersion, ShouldEqual, "1.0.2.3")
	})
}
