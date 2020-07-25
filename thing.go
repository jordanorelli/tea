package thing

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidKey = errors.New("invalid key")
	ErrNotFound   = errors.New("not found")
)

func validateKey(key string) error {
	if strings.Count(key, " ") > 0 {
		return fmt.Errorf("%w: key cannot contain spaces", ErrInvalidKey)
	}
	return nil
}

type Thing struct {
	data map[string]string
}

func (t *Thing) Get(key string) string {
	return t.data[key]
}

func (t *Thing) Set(key, value string) error {
	if err := validateKey(key); err != nil {
		return fmt.Errorf("unable to set value for key %q: %w", key, err)
	}
	t.data = map[string]string{key: value}
	return nil
}

func (t *Thing) Has(key string) bool {
	_, ok := t.data[key]
	return ok
}
