package tests

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func CheckErr(t *testing.T, err error) {
	if err != nil {
		t.Log(err, caller())
		t.FailNow()
	}
}

func Expect(t *testing.T, condition bool) {
	if !condition {
		t.Log("Expectation failed", caller())
		t.FailNow()
	}
}

func ExpectEquals(t *testing.T, left, right interface{}) {
	if !reflect.DeepEqual(left, right) {
		t.Logf("Not Equals:\n%#v\n%#v\n@%s", left, right, caller())
		t.FailNow()
	}
}

func ExpectNotEquals(t *testing.T, left, right interface{}) {
	if reflect.DeepEqual(left, right) {
		t.Logf("Equals (should not be):\n%#v\n%#v\n@%s", left, right, caller())
		t.FailNow()
	}
}

func ExpectNotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Log("Unexpected NIL value", caller())
		t.FailNow()
	}
}

func caller() string {
	var b bytes.Buffer
	for i := 2; i < 5; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok &&
			!strings.HasSuffix(file, "/src/testing/testing.go") &&
			!strings.Contains(file, "/src/runtime/") {
			b.WriteString(fmt.Sprint("\n", file, ":", line))
		}
	}
	return b.String()
}
