package source

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ghodss/yaml"

	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/watch/store"

	"github.com/uphy/watch-web/pkg/watch"

	"github.com/uphy/watch-web/pkg/domain"

	"github.com/uphy/watch-web/pkg/config"
)

type (
	TestData struct {
		Source *config.SourceConfig `json:"source"`
		Tests  []Test               `json:"tests"`
	}
	Test struct {
		Name     string            `json:"name"`
		Vars     map[string]string `json:"vars"`
		Previous string            `json:"previous"`
		Expects  Expects           `json:"expects"`
	}
	Expects struct {
		Result *string `json:"result"`
		Diff   *string `json:"diff"`
	}
	reporter struct {
		t        *testing.T
		fileName string
		testName string
	}
)

func TestAll(t *testing.T) {
	dir := "data"
	files, _ := ioutil.ReadDir("data")
	logger := logrus.New()
	reporter := &reporter{t: t}
	for _, file := range files {
		reporter.SetFileName(file.Name())
		testData := LoadTestData(filepath.Join(dir, file.Name()))
		for i, test := range testData.Tests {
			reporter.SetTestName(test.Name)
			ctx := domain.NewRootTemplateContext()
			for k, v := range test.Vars {
				ctx.Set(k, v)
			}
			source, err := testData.Source.Source(ctx)
			if err != nil {
				reporter.Error("failed to create source:", err)
				return
			}
			jobID := fmt.Sprintf("%s-%d", file.Name(), i)
			store := store.NewMemoryStore()
			store.SetValue(jobID, test.Previous)
			exe := watch.NewExecutor(store, make([]domain.Action, 0), logger)
			job := watch.NewJob(&domain.JobInfo{
				ID: jobID,
			}, source)
			if err := exe.AddJob(job, nil); err != nil {
				reporter.Error("failed to add job:", err)
				return
			}
			res, err := exe.Check(job)
			if err != nil {
				reporter.Error("failed to check update:", err)
				return
			}
			if test.Expects.Result != nil {
				expected := *test.Expects.Result
				compareString(reporter, "Result", expected, res.Current)
			}
			if test.Expects.Diff != nil {
				expected := *test.Expects.Diff
				r, err := res.Diff()
				if err != nil {
					reporter.Error("failed to diff:", err)
					return
				}
				compareString(reporter, "Diff", expected, r.String())
			}
		}
	}
}

func compareString(reporter *reporter, label, expected, actual string) {
	expected = strings.Trim(expected, " \t\n")
	actual = strings.Trim(actual, " \t\n")
	if actual != expected {
		reporter.Errorf(`%s wrong:
expected:
%s

actual:
%s
`, label, expected, actual)
	}
}

func LoadTestData(file string) *TestData {
	f, _ := os.Open(file)
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	var v TestData
	yaml.Unmarshal(b, &v)
	return &v
}

func (r *reporter) error(s string) {
	r.t.Error(fmt.Sprintf("[FAIL: file=%s, test=%s] %s", r.fileName, r.testName, s))
}

func (r *reporter) Error(v ...interface{}) {
	r.error(fmt.Sprint(v...))
}

func (r *reporter) Errorf(format string, v ...interface{}) {
	r.error(fmt.Sprintf(format, v...))
}

func (r *reporter) SetFileName(fileName string) {
	r.fileName = fileName
	r.testName = ""
}

func (r *reporter) SetTestName(testName string) {
	r.testName = testName
}
