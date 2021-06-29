// +build !unit

package templates_unit_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	renderedtDataPath = "local-render"
	goldenDataPath    = "golden"
	testDataPath      = "testdata"
)

var (
	// updateFixtures is set by the `-update` flag.
	updateFixtures bool

	// prettyDiff is set by the `-pretty-diff` flag.
	prettyDiff bool

	// write any rejected test data into this path
	rejectPath string
)

func TestMain(m *testing.M) {
	flag.BoolVar(&updateFixtures, "update", false, "update text fixtures in place")
	prettyDiff = os.Getenv("HELM_TEST_PRETTY_DIFF") != ""
	flag.BoolVar(&prettyDiff, "pretty-diff", prettyDiff, "display the full text when diffing")
	flag.StringVar(&rejectPath, "reject-path", "rejects", "write results for failed tests to this path (path is relative to the test location)")
	flag.Parse()
	os.Exit(m.Run())
}

func TestUnitTestTemplates(t *testing.T) {
	// for each folder in testdata/matchme
	matchDirs, err := ioutil.ReadDir(testDataPath)
	if err != nil {
		t.Errorf("reading directory %s failed: %v", testDataPath, err)
	}

	for _, d := range matchDirs {
		if d.IsDir() {
			err := checkValuesFile(d.Name(), t)
			if err != nil {
				t.Errorf("reading %s failed: %v", d.Name(), err)
			}
		}
	}
}

func checkValuesFile(dirName string, t *testing.T) error {
	valuesFileName := filepath.Join(testDataPath, dirName, "values.yaml")
	fmt.Printf("checking %s\n", valuesFileName)

	_, err := os.Stat(valuesFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} // simply ignore missing values.yaml
		return err
	}

	// for each file in folder "golden"
	baseDir := filepath.Join(testDataPath, dirName)
	goldenPath := filepath.Join(baseDir, goldenDataPath)
	actualPath := filepath.Join(baseDir, renderedtDataPath)

	err = filepath.Walk(goldenPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".yaml") {
				return nil
			}

			resourceName := strings.TrimPrefix(path, goldenPath)
			actualName := filepath.Join(actualPath, resourceName)
			actual := readTestdata(t, actualName)
			diffTestdata(t, path, baseDir, actual)

			return nil
		})
	if err != nil {
		t.Errorf("reading directory %s failed: %v", goldenPath, err)
	}
	return nil
}

func diffTestdata(t *testing.T, path, baseDir, actual string) {
	expected := readTestdata(t, path)
	if actual == expected {
		return
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expected, actual, true)
	diffs = dmp.DiffCleanupSemantic(diffs)
	var diff string
	if prettyDiff {
		diff = dmp.DiffPrettyText(diffs)
	} else {
		diff = dmp.PatchToText(dmp.PatchMake(diffs))
	}
	t.Errorf("mismatch: %s\n%s", path, diff)

	if updateFixtures {
		p := filepath.Join(baseDir, goldenDataPath, filepath.Base(path))
		writeTestdata(t, p, []byte(actual))
	}

	if rejectPath != "" {
		p := filepath.Join(baseDir, rejectPath, filepath.Base(path)+".rej")
		writeRejects(t, p, []byte(actual))
	}
}

func writeTestdata(t *testing.T, fileName string, data []byte) {
	if err := ioutil.WriteFile(fileName, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func writeRejects(t *testing.T, fileName string, data []byte) {
	if err := ioutil.WriteFile(fileName, data, 0644); err != nil {
		t.Fatal(err)
	}
}

// readTestdata reads a file and returns the contents of that file as a string.
func readTestdata(t *testing.T, fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("Failed to open expected input file: %v", err)
	}

	fixture, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}

	return string(fixture)
}
