package compute

import "testing"

type expectHelper struct {
	test *testing.T
}

func expect(test *testing.T) expectHelper {
	return expectHelper{test}
}

func (expect expectHelper) notNil(description string, actual interface{}) {
	if actual == nil {
		expect.test.Fatalf("%s was nil.", description)
	}
}

func (expect expectHelper) equalsString(description string, expected string, actual string) {
	if actual != expected {
		expect.test.Fatalf("%s was '%s' (expected '%s').", description, actual, expected)
	}
}

func (expect expectHelper) equalsInt(description string, expected int, actual int) {
	if actual != expected {
		expect.test.Fatalf("%s was %d (expected %d).", description, actual, expected)
	}
}
