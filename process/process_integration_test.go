package process_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/snowmerak/smips/files"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/process"
	"github.com/snowmerak/smips/register"
)

func instruction(op opcode.Op, size uint8, a, b, d uint16) opcode.Opcode {
	return opcode.NewRegularOpCode(size, a, b, d, uint8(op))
}

func set(value uint64, d uint16) opcode.Opcode {
	return opcode.NewWideOpCode(value, d, uint8(opcode.OpSet))
}

func ExampleProcess() {
	child := process.New(instruction(opcode.OpAdd, 0, register.ParamsBase, 1, register.ReturnsBase))
	p := process.New(
		instruction(opcode.OpAdd, 0, register.ParamsBase, 0, register.ExtraBase),
		instruction(opcode.OpLib, 0, 0, 0, 0),
		instruction(opcode.OpAdd, 0, register.ExtraBase, 0, register.ReturnsBase),
	).AddArgs(41).AddLib(child)
	p.Execute()
	fmt.Println(p.GetReturn(0))
	// Output: 42
}

func TestFileSystemCallsAcrossFrames(t *testing.T) {
	path := filepath.Join(t.TempDir(), "한글-파일.bin")
	const hi = register.ExtraBase - 1
	child := process.New(
		set(0x41424344, register.TemporaryBase),
		instruction(opcode.WriteFile, 4, register.ParamsBase, 0, register.TemporaryBase),
		instruction(opcode.OpAdd, 0, hi, 0, register.ReturnsBase),
	)
	p := process.New(
		set(uint64(files.FileWriter), register.TemporaryBase),
		instruction(opcode.OpenFile, 0, register.ParamsBase, register.TemporaryBase, register.TemporaryBase+1),
		instruction(opcode.OpAdd, 0, register.TemporaryBase+1, 0, register.ExtraBase),
		instruction(opcode.OpLib, 0, 0, 0, 0),
		instruction(opcode.OpAdd, 0, register.ExtraBase, 0, register.ReturnsBase),
		instruction(opcode.CloseFile, 0, register.TemporaryBase+1, 0, register.ReturnsBase+1),
		set(uint64(files.FileReader), register.TemporaryBase),
		instruction(opcode.OpenFile, 0, register.ParamsBase, register.TemporaryBase, register.TemporaryBase+1),
		instruction(opcode.ReadFileAt, 8, register.TemporaryBase+1, 0, register.ReturnsBase+2),
		instruction(opcode.CloseFile, 0, register.TemporaryBase+1, 0, register.ReturnsBase+3),
	).AddStringArgs(process.String([]rune(path))).AddLib(child)
	p.Execute()
	if p.GetReturn(0) != 4 || p.GetReturn(1) != 1 || p.GetReturn(2) != 0x41424344 || p.GetReturn(3) != 1 {
		t.Fatalf("write/close/read/close results = %d/%d/%x/%d", p.GetReturn(0), p.GetReturn(1), p.GetReturn(2), p.GetReturn(3))
	}
	data, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(data, []byte("ABCD")) {
		t.Fatalf("host file = %q, error = %v", data, err)
	}
}
