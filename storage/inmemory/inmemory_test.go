package inmemory

import (
	"testing"
)

func TestInMemory(t *testing.T) {
	storage := Storage{}
	expectValue := 1993

	err := storage.Write(1993)
	if err != nil {
		t.Fatalf("expected nil error but got %v", err)
	}

	value, err := storage.Read()
	if err != nil {
		t.Fatalf("expected nil error but got %v", err)
	}

	if value != int64(expectValue) {
		t.Fatalf("expected %d but got %d", expectValue, value)
	}

}
