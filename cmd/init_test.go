package cmd

import (
	"os"
	"testing"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/util"
)

func TestGatherOne(t *testing.T) {
	Convey("收集用户终端输入", t, func() {
		Convey("带默认值，用户输入'go build'", func() {
			monkey.Patch(util.Scanln, func() (line string, err error) {
				return "go build", nil
			})
			defer monkey.Unpatch(util.Scanln)
			actual := gatherOne("tool", "go install")
			So(actual, ShouldEqual, "go build")
		})

		Convey("带默认值，用户输入空", func() {
			monkey.Patch(util.Scanln, func() (line string, err error) {
				return "", nil
			})
			defer monkey.Unpatch(util.Scanln)
			So(gatherOne("tool", "go install"), ShouldEqual, "go install")
		})

		Convey("不带默认值，用户输入'go build'", func() {
			monkey.Patch(util.Scanln, func() (line string, err error) {
				return "go build", nil
			})
			defer monkey.Unpatch(util.Scanln)
			So(gatherOne("tool", ""), ShouldEqual, "go build")
		})
	})
}

func TestGatherOneVar(t *testing.T) {
	Convey("收集用户输入的变量名及其值", t, func() {
		var (
			varName = "Commit"
			varVal  = "{{.GitCommit}}"
		)
		var i int
		monkey.Patch(util.Scanln, func() (line string, err error) {
			defer func() { i++ }()
			switch i {
			case 0:
				return varName, nil // variable's name
			case 1:
				return varVal, nil // variable's value
			default:
				panic("unreachable")
			}
		})
		defer monkey.Unpatch(util.Scanln)

		v := gatherOneVar()
		So(v, ShouldNotBeNil)
		So(v.Variable, ShouldEqual, varName)
		So(v.Value, ShouldEqual, varVal)
	})
}

func TestGather(t *testing.T) {
	Convey("收集用户的多次输入", t, func() {
		var (
			tool       = "go build"
			c1th       = "y"
			importpath = "github.com/voidint/gbb/build"
			varName    = "Commit"
			varVal     = "{{.GitCommit}}"
			c2th       = "n"
		)
		var i int
		monkey.Patch(util.Scanln, func() (line string, err error) {
			defer func() { i++ }()
			switch i {
			case 0:
				return tool, nil // tool
			case 1:
				return c1th, nil // continue? yes
			case 2:
				return importpath, nil // importpath
			case 3:
				return varName, nil // variable's name
			case 4:
				return varVal, nil // variable's value
			case 5:
				return c2th, nil // continue? no
			default:
				panic("unreachable")
			}
		})
		defer monkey.Unpatch(util.Scanln)

		c := gather()
		So(c, ShouldNotBeNil)
		So(c.Version, ShouldEqual, Version)
		So(c.Tool, ShouldEqual, tool)
		So(c.Importpath, ShouldEqual, importpath)
		So(len(c.Variables), ShouldEqual, 1)
		So(c.Variables[0].Variable, ShouldEqual, varName)
		So(c.Variables[0].Value, ShouldEqual, varVal)
	})
}

func TestGenConfigFile(t *testing.T) {
	Convey("在指定路径生成配置文件", t, func() {
		monkey.Patch(gather, func() (c *config.Config) {
			return &config.Config{
				Version:    Version,
				Tool:       "go build",
				Importpath: "github.com/voidint/gbb/build",
			}
		})

		monkey.Patch(util.Scanln, func() (line string, err error) {
			return "y", nil
		})
		defer monkey.Unpatch(util.Scanln)

		filename := "./gbb_test.json"
		So(genConfigFile(filename), ShouldBeNil)
		os.Remove(filename)
	})
}

func TestInitCmd(t *testing.T) {
	monkey.Patch(genConfigFile, func(_ string) error {
		return nil
	})

	defer monkey.Unpatch(genConfigFile)
	initCmd.Run(nil, nil)
}
