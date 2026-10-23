package files

import (
	"bytes"
	"io"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestStreamIO(t *testing.T) {
	for _, size := range []uint8{0, 1, 2, 7, 8} {
		f := New()
		var output bytes.Buffer
		f.OverrideStdout(&output)
		const value = uint64(0x0123456789abcdef)
		if n := f.WriteFile(uint64(Stdout), size, value); n != uint64(size) {
			t.Fatalf("size %d: wrote %d bytes", size, n)
		}
		want := []byte{1, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}[8-int(size):]
		if !bytes.Equal(output.Bytes(), want) {
			t.Fatalf("size %d: bytes = %x, want %x", size, output.Bytes(), want)
		}
		f.OverrideStdin(bytes.NewReader(output.Bytes()))
		var expected uint64
		for _, b := range want {
			expected = expected<<8 | uint64(b)
		}
		if got := f.ReadFile(uint64(Stdin), size); got != expected {
			t.Fatalf("size %d: read = %x, want %x", size, got, expected)
		}
	}
}

func expectPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	fn()
}

func TestFailedOpenPreservesDescriptors(t *testing.T) {
	f := New()
	path := filepath.Join(t.TempDir(), "existing.bin")
	fd := f.OpenFile(path, FileWriter)
	f.CloseFile(fd)
	expectPanic(t, func() { f.OpenFile(path, 255) })
	expectPanic(t, func() { f.OpenFile(path+"-missing", FileReader) })
	rd := f.OpenFile(path, FileReader)
	defer f.CloseFile(rd)
	if rd != fd {
		t.Fatal("failed open consumed a reusable descriptor")
	}
	other := f.OpenFile(path, FileReader)
	defer f.CloseFile(other)
	if other == rd || f.table[rd] == f.table[other] {
		t.Fatal("opening another file replaced an existing descriptor")
	}
}

type shortWriter struct{}

func (shortWriter) Write(buf []byte) (int, error) { return len(buf) - 1, nil }

func TestShortWrite(t *testing.T) {
	f := New()
	f.OverrideStdout(shortWriter{})
	defer func() {
		if got := recover(); got != io.ErrShortWrite {
			t.Errorf("short write panic = %v, want %v", got, io.ErrShortWrite)
		}
	}()
	f.WriteFile(uint64(Stdout), 2, 0x4142)
}

func TestInvalidSizesAndOffsets(t *testing.T) {
	f := New()
	for _, size := range []uint8{9, 255} {
		expectPanic(t, func() { f.ReadFile(uint64(Stdin), size) })
		expectPanic(t, func() { f.ReadFileAt(uint64(Stdin), 0, size) })
		expectPanic(t, func() { f.WriteFile(uint64(Stdout), size, 0) })
		expectPanic(t, func() { f.WriteFileAt(uint64(Stdout), 0, size, 0) })
	}
	expectPanic(t, func() { f.ReadFileAt(uint64(Stdin), math.MaxUint64, 1) })
	expectPanic(t, func() { f.WriteFileAt(uint64(Stdout), math.MaxUint64, 1, 0) })
}

type hostStream struct {
	bytes.Buffer
	closed bool
}

func (s *hostStream) Close() error {
	s.closed = true
	return nil
}

func TestStandardStreamsRemainHostOwned(t *testing.T) {
	f := New()
	stream := &hostStream{}
	f.OverrideStdin(stream)
	f.OverrideStdout(stream)
	f.OverrideStderr(stream)
	for _, fd := range []uint8{Stdin, Stdout, Stderr} {
		if f.CloseFile(uint64(fd)) != 0 || stream.closed {
			t.Fatal("VM closed a borrowed host stream")
		}
	}
	fd := f.OpenFile(filepath.Join(t.TempDir(), "owned.bin"), FileWriter)
	defer f.CloseFile(fd)
	if fd < 3 {
		t.Fatal("standard stream descriptor was recycled")
	}
	f.WriteFile(uint64(Stdout), 1, 65)
	if stream.String() != "A" {
		t.Fatal("opening an owned file replaced the host stream")
	}
}

func TestPartialReadAndEOF(t *testing.T) {
	f := New()
	f.OverrideStdin(bytes.NewReader([]byte{0x41, 0x42}))
	if got := f.ReadFile(uint64(Stdin), 8); got != 0x4142 {
		t.Fatalf("partial read = %x", got)
	}
	if got := f.ReadFile(uint64(Stdin), 8); got != 0 {
		t.Fatalf("EOF read = %x", got)
	}
}

func TestFileOffsetsAndHandleReuse(t *testing.T) {
	f := New()
	path := filepath.Join(t.TempDir(), "io.bin")
	fd := f.OpenFile(path, FileWriter)
	if n := f.WriteFileAt(fd, 2, 2, 0x4142); n != 2 {
		t.Fatalf("WriteAt = %d", n)
	}
	if f.CloseFile(fd) != 1 || f.CloseFile(fd) != 0 {
		t.Fatal("close status is incorrect")
	}
	rd := f.OpenFile(path, FileReader)
	if rd != fd {
		t.Fatal("closed descriptor was not reused")
	}
	defer f.CloseFile(rd)
	if got := f.ReadFileAt(rd, 2, 8); got != 0x4142 {
		t.Fatalf("partial ReadAt = %x", got)
	}
	if got := f.ReadFileAt(rd, 4, 1); got != 0 {
		t.Fatalf("ReadAt EOF = %x", got)
	}
	data, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(data, []byte{0, 0, 0x41, 0x42}) {
		t.Fatalf("file = %x, error = %v", data, err)
	}
}
