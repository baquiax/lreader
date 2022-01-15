/*Package lreader is intended to be a "stateful" lines' reader. It
provides a way to read files line by line in a lazy way, storing the last
point read in an storage system provided by user.

The reason to make this a "stateful" component is to increase the "resilience"
of the application using it. Meaning that if for some reason the application
using it crashes, you can resume the work form the last point processed. This is
done by storing constantly the last read byte in a storage system provided as an
argument.
*/
package lreader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

const (
	// maxLineBytes is used to limit the amount of bytes
	// kept in memory when reading long lines
	//
	// Check: bufio.defaultBufSize and bufio.ReadLine (prefix)
	maxLineBytes = 4096 * 2
)

var (
	// ErrNilSource is returned when the first argument of the ReadLine
	// method is invalid
	ErrNilSource = errors.New("source is nil")

	// ErrNilStorage is returned when the second argument of the ReadLine
	// method is invalid
	ErrNilStorage = errors.New("storage is nil")

	// ErrLongLine is returned when one of the lines of the given source
	// is too long the be kept in memory. See maxLineBytes
	ErrLongLine = errors.New("line to log to be readed")

	// ErrInvalidReader is returned if the Reader instace does not contain
	// the needed dependencies. It could happen if it was not created using
	// the New method.
	//
	// This error aim to reduce the probability to have unexpected panics
	ErrInvalidReader = errors.New("the Reader instances was not created with the given constructor New")
)

// OffsetReadWritter defines the behavior of the storage system in charge to store
// and read the last read byte on a file.
type OffsetReadWritter interface {
	// Read returns the last read byte on a file.
	Read() (int64, error)

	// Write writes the last read byte on a file.
	Write(value int64) error
}

// Reader represents an instance of a line reader. Currently it
// does not export any useful field.
type Reader struct {
	storage       OffsetReadWritter
	reader        *bufio.Reader
	currentOffset int64
}

// New abstract the logic to build a valid Reader.
func New(source io.Reader, storage OffsetReadWritter) (*Reader, error) {
	if source == nil {
		return nil, ErrNilSource
	}

	if storage == nil {
		return nil, ErrNilStorage
	}

	offset, err := storage.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading the stored offset: %v", err)
	}

	reader := bufio.NewReader(source)

	_, err = reader.Discard(int(offset))
	if err != nil {
		return nil, fmt.Errorf("error retoring last position on the file: %v", err)
	}

	return &Reader{
		reader:  reader,
		storage: storage,
	}, nil
}

func (r *Reader) readLine() ([]byte, error) {
	var bytes []byte
	var error error

	for {
		partialBytes, isPrefix, err := r.reader.ReadLine()

		if len(bytes)+len(partialBytes) > maxLineBytes {
			return bytes, ErrLongLine
		}

		bytes = append(bytes, partialBytes...)

		if err != nil {
			error = err
			break
		}

		if !isPrefix {
			break
		}
	}

	return bytes, error
}

// ReadLine returns the bytes of the next line in the given sources
// that implements the io.Reader interface.
//
// This slice of bytes does not include the new line character (\n)
//
// An expeted error is got when the end of file is found (io.EOF). When
// this happens, the returned bytes MUST be readed
func (r *Reader) ReadLine() ([]byte, error) {
	if err := r.validateInstance(); err != nil {
		return []byte{}, err
	}

	bytes, err := r.readLine()

	if err != nil {
		if err == io.EOF {
			return bytes, io.EOF
		}

		return bytes, fmt.Errorf("error reading lines: %v", err)
	}

	// TODO: improve commit interval
	err = r.storage.Write(r.currentOffset + int64(len(bytes)+1))
	if err != nil {
		return bytes, fmt.Errorf("error writing the offset: %v", err)
	}

	r.currentOffset += int64(len(bytes) + 1)

	return bytes, err
}

func (r *Reader) validateInstance() error {
	switch {
	case r == nil:
		return ErrInvalidReader
	case r.reader == nil:
		return ErrNilSource
	case r.storage == nil:
		return ErrNilStorage
	default:
		return nil
	}
}
