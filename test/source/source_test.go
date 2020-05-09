package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ghodss/yaml"

	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/watch/store"

	"github.com/uphy/watch-web/pkg/watch"

	_ "github.com/mattn/anko/packages"
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
		Previous interface{}       `json:"previous"`
		Expects  Expects           `json:"expects"`
	}
	Expects struct {
		Result  domain.ItemList `json:"result"`
		Changed *bool           `json:"changed"`
		Diff    domain.Updates  `json:"diff"`
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
		testDataPath := filepath.Join(dir, file.Name())
		loader := config.NewLoader(logger, testDataPath)
		testData := LoadTestData(testDataPath)
		for i, test := range testData.Tests {
			reporter.SetTestName(test.Name)
			ctx := loader.TemplateContext()
			for k, v := range test.Vars {
				ctx.Set(k, v)
			}
			source, err := loader.CreateSource(testData.Source)
			if err != nil {
				reporter.Error("failed to create source:", err)
				continue
			}
			jobID := fmt.Sprintf("%s-%d-%s", file.Name(), i, test.Name)
			store := store.NewMemoryStore()
			if s, ok := test.Previous.(string); ok {
				store.SetValue(jobID, s)
			} else {
				previousJSON, err := json.Marshal(test.Previous)
				if err != nil {
					reporter.Error("failed to marshal previous value")
					continue
				}
				store.SetValue(jobID, string(previousJSON))
			}
			exe := watch.NewExecutor(store, make([]domain.Action, 0), logger)
			job := watch.NewJob(&domain.JobInfo{
				ID: jobID,
			}, source)
			if err := exe.AddJob(job, nil); err != nil {
				reporter.Error("failed to add job:", err)
				continue
			}
			res, err := exe.Check(job)
			if err != nil {
				reporter.Error("failed to check update:", err)
				continue
			}
			if test.Expects.Result != nil {
				expected := test.Expects.Result
				compareResult(reporter, "Result", expected, res.Current)
			}
			if test.Expects.Diff != nil || test.Expects.Changed != nil {
				r := res.Diff()
				changed := r.Changes()
				if test.Expects.Changed != nil && changed != *test.Expects.Changed {
					reporter.Errorf("Diff changed property is wrong: expected=%v, actual=%v", *test.Expects.Changed, changed)
					continue
				}
				if test.Expects.Diff != nil {
					expected := test.Expects.Diff
					compareUpdates(reporter, "Diff", expected, r)
				}
			}
		}
	}
}

func compareResult(reporter *reporter, label string, expected, actual domain.ItemList) {
	if !reflect.DeepEqual(expected, actual) {
		reporter.Errorf(`%s wrong:
expected:
%s

actual:
%s
`, label, expected.YAML(), actual.YAML())
	}
}

func compareUpdates(reporter *reporter, label string, expected, actual domain.Updates) {
	y1 := expected.YAML()
	y2 := actual.YAML()
	if !reflect.DeepEqual(y1, y2) {
		reporter.Errorf(`%s wrong:
expected:
%s

actual:
%s
`, label, y1, y2)
	}
}

func LoadTestData(file string) *TestData {
	f, _ := os.Open(file)
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	var v TestData
	err := yaml.Unmarshal(b, &v)
	if err != nil {
		log.Fatal("failed to parse:", file, " ", err)
	}
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
