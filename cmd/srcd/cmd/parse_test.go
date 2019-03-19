// +build integration

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	cmdtest "github.com/src-d/engine/cmd/test-utils"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ParseTestSuite struct {
	cmdtest.IntegrationSuite
	testDir string
}

func TestParseTestSuite(t *testing.T) {
	s := ParseTestSuite{}
	suite.Run(t, &s)
}

type testCase struct {
	path     string
	filename string
	lang     string
}

// Uses files from https://github.com/leachim6/hello-world
var testCases = []testCase{
	{
		path:     filepath.FromSlash("testdata/hello.py"),
		filename: "hello.py",
		lang:     "python",
	},
	{
		path:     filepath.FromSlash("testdata/hello-py3.py"),
		filename: "hello-py3.py",
		lang:     "python",
	},
	{
		path:     filepath.FromSlash("testdata/hello.cpp"),
		filename: "hello.cpp",
		lang:     "c++",
	},
	{
		path:     filepath.FromSlash("testdata/hello.java"),
		filename: "hello.java",
		lang:     "java",
	},
	{
		path:     filepath.FromSlash("testdata/hello.js"),
		filename: "hello.js",
		lang:     "javascript",
	},
	{
		path:     filepath.FromSlash("testdata/hello.bash"),
		filename: "hello.bash",
		lang:     "shell",
	},
	{
		path:     filepath.FromSlash("testdata/hello.rb"),
		filename: "hello.rb",
		lang:     "ruby",
	},
	{
		path:     filepath.FromSlash("testdata/hello.go"),
		filename: "hello.go",
		lang:     "go",
	},
	{
		path:     filepath.FromSlash("testdata/hello.cs"),
		filename: "hello.cs",
		lang:     "c#",
	},
	{
		path:     filepath.FromSlash("testdata/hello.php"),
		filename: "hello.php",
		lang:     "php",
	},
}

func (s *ParseTestSuite) SetupTest() {
}

func (s *ParseTestSuite) TearDownTest() {
	s.RunStop(context.Background())
}

func (s *ParseTestSuite) TestDriversList() {
	require := s.Require()

	out, err := s.RunCommand(context.TODO(), "parse", "drivers", "list")
	outStr := out.String()

	require.NoError(err, outStr)

	/* Example output:

	LANGUAGE	VERSION
	----------	----------
	python		v2.8.0
	cpp		v1.1.0
	java		v2.5.0
	javascript	v2.6.0
	bash		v2.4.0
	ruby		v2.9.0
	go		v2.5.0
	csharp		v1.4.0
	php		v2.7.0
	*/

	// Simple checks to see if it's the table, and contains a known driver
	expected := regexp.MustCompile(`LANGUAGE\s+VERSION`)
	require.Regexp(expected, outStr)
	expected = regexp.MustCompile(`javascript\s+v\S+`)
	require.Regexp(expected, outStr)
}

func (s *ParseTestSuite) TestLang() {
	for _, tc := range testCases {
		s.T().Run(tc.filename, func(t *testing.T) {
			require := require.New(t)

			// Check the language is detected
			out, err := s.RunCommand(context.TODO(), "parse", "lang", tc.path)
			require.NoError(err, out.String())
			require.Equal(tc.lang+"\n", out.String())
		})
	}
}

// same as RunCommand, but captures only stdout instead of stdout + stderr
func (s *ParseTestSuite) runCommandStdout(ctx context.Context, cmd string, args ...string) (*bytes.Buffer, error) {
	var out bytes.Buffer

	command := s.CommandContext(ctx, cmd, args...)
	command.Stdout = &out

	return &out, command.Run()
}

type arg []string
type args []arg

func getArgCombinations(tc testCase) args {
	uast := args{arg{"uast", tc.path}}
	modes := args{
		arg{},
		arg{"--mode", "semantic"},
		arg{"--mode", "annotated"},
		arg{"--mode", "native"},
	}
	langs := args{
		arg{},
		arg{"--lang", tc.lang},
	}
	queries := args{
		arg{},
		arg{"--query", "/"}, // Xpath query to get the root node
	}

	return combine(combine(combine(uast, modes), langs), queries)
}

func combine(a args, b args) args {
	var out args
	for _, v := range a {
		for _, w := range b {
			out = append(out, append(v, w...))
		}
	}

	return out
}

func (s *ParseTestSuite) TestUast() {
	for _, tc := range testCases {
		argsCombinations := getArgCombinations(tc)

		for _, args := range argsCombinations {
			testName := fmt.Sprintf("%s %s", tc.filename, strings.Join(args[2:], " "))
			s.T().Run(testName, func(t *testing.T) {
				require := require.New(t)

				// Intentionally left to help creating new test cases
				// t.Log("about to run parse " + strings.Join(args, " "))

				// RunCommand mixes stdout and stderr. To parse the UAST output properly
				// we need to read stdout only
				uastOut, err := s.runCommandStdout(context.TODO(), "parse", args...)

				// ----------------
				// TODO Temporary test skip, it fails for cpp, bash, and csharp.
				// See https://github.com/src-d/engine/issues/297
				if tc.lang == "c++" || tc.lang == "shell" || tc.lang == "c#" {
					// This Error assertion will fail when #297 is fixed, to remind us to remove this skip
					require.Error(err)
					t.Skip("TEST FAILURE IS A KNOWN ISSUE (#297): " + uastOut.String())
				}
				// ----------------

				extraInfo := fmt.Sprintf("srcd parse %s\n%s", strings.Join(args, " "), uastOut.String())
				require.NoError(err, extraInfo)

				// Check the UAST output is valid json

				var js interface{}
				err = json.Unmarshal([]byte(uastOut.String()), &js)
				require.NoError(err, extraInfo)
			})
		}
	}
}