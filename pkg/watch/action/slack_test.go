package action

import (
	"os"
	"reflect"
	"testing"

	"github.com/uphy/watch-web/pkg/domain/value"

	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/domain"
	"gopkg.in/yaml.v2"
)

func TestSlackPayload(t *testing.T) {
	f, err := os.Open("testdata/slack.yml")
	if err != nil {
		t.Error("failed to open test data: ", err)
		return
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)

	var itemList1 value.ItemList
	if err := decoder.Decode(&itemList1); err != nil {
		t.Error("failed to decode itemList1: ", err)
		return
	}
	var itemList2 value.ItemList
	if err := decoder.Decode(&itemList2); err != nil {
		t.Error("failed to decode itemList2: ", err)
		return
	}

	logger := logrus.New()
	ctx := &domain.JobContext{Log: logger.WithFields(logrus.Fields{})}
	res := domain.NewResult(&domain.JobInfo{}, itemList1, itemList2)
	updates := res.Diff()
	for _, update := range updates {
		actual, err := slackPayload(ctx, res, update)
		if err != nil {
			t.Error("error occurred during run function: ", err)
			return
		}
		var payload interface{}
		if err := decoder.Decode(&payload); err != nil {
			t.Error("failed to decode payload: ", err)
			return
		}
		compareYAML(t, actual, payload)
	}
}

func compareYAML(t *testing.T, actual, expected interface{}) {
	actualBytes, _ := yaml.Marshal(actual)
	expectedBytes, _ := yaml.Marshal(expected)
	if reflect.DeepEqual(actualBytes, expectedBytes) {
		return
	}
	t.Errorf("---Expected---\n%s\n---Actual---\n%s", string(expectedBytes), string(actualBytes))
}
