package files

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"
)

// Files is a file table.
type Files struct {
	table map[uint64]any
	tombs []uint64
}

const (
	Stdin uint8 = iota
	Stdout
	Stderr
	FileReader
	FileWriter
	FileAppender
)

func New() Files {
	return Files{
		table: map[uint64]any{
			uint64(Stdin):  os.Stdin,
			uint64(Stdout): os.Stdout,
			uint64(Stderr): os.Stderr,
		},
		tombs: []uint64{},
	}
}

// OverrideStdin overrides stdin.
func (f *Files) OverrideStdin(reader io.Reader) {
	f.table[uint64(Stdin)] = reader
}

// OverrideStdout overrides stdout.
func (f *Files) OverrideStdout(writer io.Writer) {
	f.table[uint64(Stdout)] = writer
}

// OverrideStderr overrides stderr.
func (f *Files) OverrideStderr(writer io.Writer) {
	f.table[uint64(Stderr)] = writer
}

// OpenFile opens a file and returns a reusable VM descriptor.
func (f *Files) OpenFile(path string, mod uint8) uint64 {
	var file *os.File
	var err error
	switch mod {
	case FileReader:
		file, err = os.Open(path)
	case FileWriter:
		file, err = os.Create(path)
	case FileAppender:
		file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	default:
		panic("files: invalid open mode")
	}
	if err != nil {
		panic(err)
	}
	// Do not consume a reusable descriptor until opening has succeeded.
	index := uint64(len(f.table))
	if len(f.tombs) > 0 {
		index = f.tombs[len(f.tombs)-1]
		f.tombs = f.tombs[:len(f.tombs)-1]
	}
	f.table[index] = file
	return index
}

// CloseFile closes an owned file. Standard streams are borrowed from the host
// and cannot be closed or recycled as ordinary descriptors.
func (f *Files) CloseFile(index uint64) uint64 {
	if index <= uint64(Stderr) {
		return 0
	}
	file, ok := f.table[index]
	if !ok {
		return 0
	}
	closer, ok := file.(io.Closer)
	if !ok {
		panic("FD is not Closer")
	}
	if err := closer.Close(); err != nil {
		panic(err)
	}
	f.tombs = append(f.tombs, index)
	delete(f.table, index)
	return 1
}

// ReadFile reads up to size bytes as a big-endian value. EOF returns zero.
func (f *Files) ReadFile(index uint64, size uint8) uint64 {
	checkSize(size)
	if size == 0 {
		return 0
	}
	file, ok := f.table[index]
	if !ok {
		return 0
	}
	reader, ok := file.(io.Reader)
	if !ok {
		panic("FD is not Reader")
	}
	buf := make([]byte, size)
	n, err := reader.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		panic(err)
	}
	return decode(buf[:n])
}

// ReadFileAt reads a file at a specified offset.
func (f *Files) ReadFileAt(index uint64, offset uint64, size uint8) uint64 {
	checkSize(size)
	checkOffset(offset)
	if size == 0 {
		return 0
	}
	file, ok := f.table[index]
	if !ok {
		return 0
	}
	reader, ok := file.(io.ReaderAt)
	if !ok {
		panic("FD is not ReaderAt")
	}
	buf := make([]byte, size)
	n, err := reader.ReadAt(buf, int64(offset))
	if err != nil && !errors.Is(err, io.EOF) {
		panic(err)
	}
	return decode(buf[:n])
}

// WriteFile writes size low bytes in big-endian order and returns the byte count.
func (f *Files) WriteFile(index uint64, size uint8, value uint64) uint64 {
	checkSize(size)
	if size == 0 {
		return 0
	}
	file, ok := f.table[index]
	if !ok {
		return 0
	}
	writer, ok := file.(io.Writer)
	if !ok {
		panic("FD is not Writer")
	}
	n, err := writer.Write(encode(value, size))
	if err != nil {
		panic(err)
	}
	if n != int(size) {
		panic(io.ErrShortWrite)
	}
	return uint64(n)
}

// WriteFileAt writes a file at a specified offset.
func (f *Files) WriteFileAt(index uint64, offset uint64, size uint8, value uint64) uint64 {
	checkSize(size)
	checkOffset(offset)
	if size == 0 {
		return 0
	}
	file, ok := f.table[index]
	if !ok {
		return 0
	}
	writer, ok := file.(io.WriterAt)
	if !ok {
		panic("FD is not WriterAt")
	}
	n, err := writer.WriteAt(encode(value, size), int64(offset))
	if err != nil {
		panic(err)
	}
	if n != int(size) {
		panic(io.ErrShortWrite)
	}
	return uint64(n)
}

func checkSize(size uint8) {
	if size > 8 {
		panic("files: byte count exceeds 8")
	}
}

func checkOffset(offset uint64) {
	if offset > math.MaxInt64 {
		panic("files: offset exceeds int64")
	}
}

func encode(value uint64, size uint8) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], value)
	return buf[8-int(size):]
}

func decode(buf []byte) uint64 {
	var value uint64
	for _, b := range buf {
		value = value<<8 | uint64(b)
	}
	return value
}
