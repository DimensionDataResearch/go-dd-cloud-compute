package compute

import (
	"reflect"
	"testing"
)

type expectHelper struct {
	test *testing.T
}

func expect(test *testing.T) expectHelper {
	return expectHelper{test}
}

func (expect expectHelper) isTrue(description string, condition bool) {
	if !condition {
		expect.test.Errorf("Expression was false: %s", description)
	}
}

func (expect expectHelper) isFalse(description string, condition bool) {
	if condition {
		expect.test.Errorf("Expression was true: %s", description)
	}
}

func (expect expectHelper) notNil(description string, actual interface{}) {
	if reflect.ValueOf(actual).IsNil() {
		expect.test.Fatalf("%s was nil.", description)
	}
}

func (expect expectHelper) equalsString(description string, expected string, actual string) {
	if actual != expected {
		expect.test.Errorf("%s was '%s' (expected '%s').", description, actual, expected)
	}
}

func (expect expectHelper) equalsInt(description string, expected int, actual int) {
	if actual != expected {
		expect.test.Errorf("%s was %d (expected %d).", description, actual, expected)
	}
}
