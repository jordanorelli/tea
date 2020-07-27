package tea

import (
	"testing"
)

type nameless struct{}

func (nameless) Run(t *testing.T) {}

func TestNames(t *testing.T) {
	testTests := []struct {
		testToTest   Test
		expectedName string
	}{
		{
			testToTest:   nameless{},
			expectedName: "nameless",
		},
		{
			testToTest:   &nameless{},
			expectedName: "nameless",
		},
	}

	for _, doubleTest := range testTests {
		name := parseName(doubleTest.testToTest)
		if name != doubleTest.expectedName {
			t.Errorf("saw name %q expecting %q", name, doubleTest.expectedName)
		} else {
			t.Logf("test %v has expected name %q", doubleTest.testToTest, name)
		}
	}
}
