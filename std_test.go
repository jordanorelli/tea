// +build std

package thing

import (
	"testing"
)

func TestSet(t *testing.T) {
	tests := []struct {
		key   string
		value string
		bad   bool
	}{
		{"foo", "bar", false},
		{"foo", "a-value", false},
		{"foo", "a value", false},
		{"one-two-three", "whatever", false},
		{"one-two-three", "what ever", false},
		{"one two three", "whatever", true},
		{" one-two-three", "whatever", true},
	}

	for _, test := range tests {
		thing := new(Thing)
		err := thing.Set(test.key, test.value)
		if !test.bad && err != nil {
			t.Errorf("should be able to set %q=%q but saw error %v", test.key, test.value, err)
		}
		if test.bad && err == nil {
			t.Errorf("able to set bad values %q=%q", test.key, test.value)
		}
	}
}

func TestHas(t *testing.T) {
	var (
		key   = "foo"
		value = "bar"
		thing = &Thing{
			data: map[string]string{
				key: value,
			},
		}
	)

	if !thing.Has(key) {
		t.Errorf("missing expected key %q", key)
	}
}

func TestGet(t *testing.T) {
	var (
		key   = "foo"
		value = "bar"
		thing = &Thing{
			data: map[string]string{
				key: value,
			},
		}
	)

	if v := thing.Get(key); v != value {
		t.Errorf("read value %q, expected %q", v, value)
	}
}

func TestKeys(t *testing.T) {
	goodKeys := []string{"one", "two", "3", "one-two-three", "steve?"}
	badKeys := []string{"o n e", "one two three", " one", "one "}

	for _, key := range goodKeys {
		if err := validateKey(key); err != nil {
			t.Errorf("key %q should be valid but saw validation error %v", key, err)
		}
	}

	for _, key := range badKeys {
		if err := validateKey(key); err == nil {
			t.Errorf("key %q should be invalid but passed validation", key)
		}
	}
}
