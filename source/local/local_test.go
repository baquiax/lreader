package local

import (
	"errors"
	"testing"
)

func TestNewErrors(t *testing.T) {
	testCases := map[string]struct {
		path          string
		expectedError error
	}{
		"Should fail When path is empty": {
			path:          "",
			expectedError: ErrEmptyPath,
		},
		"Should fail When file does not exist": {
			path:          "unknown",
			expectedError: errors.New("file unknown does not exist"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := New(test.path)

			if err == nil {
				t.Errorf("expected error, got nil")
			}

			if err.Error() != test.expectedError.Error() {
				t.Errorf("expected error `%v`, got `%v`", test.expectedError, err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	fd, err := New("local.go")
	if err != nil {
		t.Fatalf("expected nil error but got %v", err)
	}

	err = fd.Close()
	if err != nil {
		t.Fatalf("expected nil error but got %v", err)
	}
}
