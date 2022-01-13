package lreader_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"baquiax.dev/lreader"
	"baquiax.dev/lreader/storage/inmemory"
)

func TestNewReaderErrors(t *testing.T) {
	storage := inmemory.Storage{}
	storage.Write(20)

	brokenstorage := inmemory.Storage{
		ForcedReadError: errors.New("imposible to read/write offset"),
	}

	tests := map[string]struct {
		reader        io.Reader
		storage       lreader.OffsetReadWritter
		expectedError error
	}{
		"Should fail When reader arg is nil": {
			reader:        nil,
			expectedError: lreader.ErrNilSource,
		},
		"Should fail When storage arg is nil": {
			reader:        strings.NewReader(""),
			expectedError: lreader.ErrNilStorage,
		},
		"Should fail When the storage fails reading the latest offset": {
			reader:        strings.NewReader(""),
			storage:       &storage,
			expectedError: errors.New("error retoring last position on the file: EOF"),
		},
		"Should fail When the offset is out of range": {
			reader:        strings.NewReader(""),
			storage:       &brokenstorage,
			expectedError: errors.New("error reading the stored offset: imposible to read/write offset"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := lreader.New(test.reader, test.storage)

			if r != nil {
				t.Fatal("when error is present the reader should be nil")
			}

			if err == nil {
				t.Errorf("expected error, got nil")
			}

			if err.Error() != test.expectedError.Error() {
				t.Errorf("expected error `%v`, got `%v`", test.expectedError, err)
			}
		})
	}
}

func TestNewReader(t *testing.T) {
	reader := strings.NewReader("")
	storage := inmemory.Storage{}

	r, err := lreader.New(reader, &storage)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if r == nil {
		t.Error("expected a reader, got nil")
	}
}

func TestReadLineErrors(t *testing.T) {
	storage := inmemory.Storage{}

	tests := map[string]struct {
		reader        io.Reader
		storage       lreader.OffsetReadWritter
		expectedError error
	}{
		"Should fail When EOF is found": {
			reader:        strings.NewReader(""),
			storage:       &storage,
			expectedError: errors.New("error reading lines: EOF"),
		},
		"Should fail When offset is not possible to write": {
			reader: strings.NewReader("hi"),
			storage: &inmemory.Storage{
				ForcedWriteError: errors.New("imposible to read/write offset"),
			},
			expectedError: errors.New("error writing the offset: imposible to read/write offset"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := lreader.New(test.reader, test.storage)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			_, err = r.ReadLine()

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if err.Error() != test.expectedError.Error() {
				t.Fatalf("expected error `%v`, got `%v`", test.expectedError, err)
			}
		})
	}
}

func TestReadLine(t *testing.T) {
	// TODO: test lines > 4096 bytes

	r, err := lreader.New(strings.NewReader("hello\nworld\n"), &inmemory.Storage{})

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	line, err := r.ReadLine()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if string(line) != "hello" {
		t.Fatalf("expected `hello`, got `%v`", line)
	}
}