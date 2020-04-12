package check

import (
	"testing"
)

func TestResult(t *testing.T) {
	result := &Result{
		JobID:    "test",
		Previous: "aaa\nbbb",
		Current:  "aaa\nbbb\nccc",
	}
	if result.Diff().String() != `  aaa
- bbb
+ bbb
+ ccc
` {
		t.Fail()
	}
}
