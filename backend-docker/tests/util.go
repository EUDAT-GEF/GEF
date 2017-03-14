package tests

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func checkMsg(t *testing.T, err error, msg string) {
	if err != nil {
		t.Error(msg, caller())
		t.Error(err)
		t.FailNow()
	}
}

func check(t *testing.T, err error) {
	checkMsg(t, err, "")
}

func expect(t *testing.T, condition bool, msg string) {
	if !condition {
		t.Error(msg, caller())
		t.FailNow()
	}
}

func caller() string {
	var b bytes.Buffer
	for i := 2; i < 5; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && !strings.HasSuffix(file, "src/testing/testing.go") {
			b.WriteString(fmt.Sprint("\n", file, ":", line))
		}
	}
	return b.String()
}
