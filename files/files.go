package files

import (
	"io"
	"os"
)

// Files is a file table.
type Files struct {
	table map[uint64]interface{}
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
		table: map[uint64]interface{}{
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

// OpenFileReader opens a file.
func (f *Files) OpenFile(path string, mod uint8) uint64 {
	index := uint64(len(f.table))
	if len(f.tombs) > 0 {
		index = f.tombs[len(f.tombs)-1]
		f.tombs = f.tombs[:len(f.tombs)-1]
	}
	var file *os.File
	var err error
	switch mod {
	case FileReader:
		file, err = os.Open(path)
	case FileWriter:
		file, err = os.Create(path)
	case FileAppender:
		file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	}
	if err != nil {
		panic(err)
	}
	f.table[index] = file
	return index
}

// CloseFile closes a file.
func (f *Files) CloseFile(index uint64) uint64 {
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

// ReadFile reads a file.
func (f *Files) ReadFile(index uint64, size uint8) uint64 {
	size = size & 0b00000111
	file, ok := f.table[index]
	if !ok {
		return 0
	}
	reader, ok := file.(io.Reader)
	if !ok {
		panic("FD is not Reader")
	}
	buf := make([]byte, size)
	_, err := reader.Read(buf)
	if err != nil {
		panic(err)
	}
	result := uint64(0)
	for i := uint8(0); i < size; i++ {
		result = result<<8 | uint64(buf[i])
	}
	return result
}

// ReadFileAt reads a file at a specified offset.
func (f *Files) ReadFileAt(index uint64, offset uint64, size uint8) uint64 {
	size = size & 0b00000111
	file, ok := f.table[index]
	if !ok {
		return 0
	}
	reader, ok := file.(io.ReaderAt)
	if !ok {
		panic("FD is not ReaderAt")
	}
	buf := make([]byte, size)
	_, err := reader.ReadAt(buf, int64(offset))
	if err != nil {
		panic(err)
	}
	result := uint64(0)
	for i := uint8(0); i < size; i++ {
		result = result<<8 | uint64(buf[i])
	}
	return result
}

// WriteFile writes a file.
func (f *Files) WriteFile(index uint64, size uint8, value uint64) uint64 {
	size = size & 0b00000111
	file, ok := f.table[index]
	if !ok {
		return 0
	}
	writer, ok := file.(io.Writer)
	if !ok {
		panic("FD is not Writer")
	}
	buf := make([]byte, size)
	for i := size - 1; i >= 0; i++ {
		buf[i] = byte(value & 0xFF)
		value = value >> 8
	}
	_, err := writer.Write(buf)
	if err != nil {
		panic(err)
	}
	return 1
}

// WriteFileAt writes a file at a specified offset.
func (f *Files) WriteFileAt(index uint64, offset uint64, size uint8, value uint64) uint64 {
	size = size & 0b00000111
	file, ok := f.table[index]
	if !ok {
		return 0
	}
	writer, ok := file.(io.WriterAt)
	if !ok {
		panic("FD is not WriterAt")
	}
	buf := make([]byte, size)
	for i := size - 1; i >= 0; i++ {
		buf[i] = byte(value & 0xFF)
		value = value >> 8
	}
	_, err := writer.WriteAt(buf, int64(offset))
	if err != nil {
		panic(err)
	}
	return 1
}
