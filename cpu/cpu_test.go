package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gone/mem"
)

func TestLoadProgram(t *testing.T) {
	// unhelpfully, this test program is nowhere to be found on OLC's repo
	program := "A2 0A 8E 00 00 A2 03 8E 01 00 AC 00 00 A9 00 18 6D 01 00 88 D0 FA 8D 02 00 EA EA EA" // 28 bytes
	// 162 10 142 ...

	C := Cpu{Bus: &mem.Bus{}}
	C.LoadProgram([]byte(program), 0x8000)
	assert.Equal(t, C.Bus.FakeRam[0x8000], uint8(0xa2))
	assert.Equal(t, C.Bus.FakeRam[0x8001], uint8(0x0a))
	assert.Equal(t, C.Bus.FakeRam[0x8002], uint8(0x8e))
	assert.Equal(t, C.Bus.FakeRam[0x801b], uint8(0xea))
	assert.Equal(t, C.Bus.FakeRam[0x801c], uint8(0))

	assert.Equal(t, Opcodes[C.Bus.FakeRam[0x8000]].Name, "LDX")
	assert.Equal(t, Opcodes[C.Bus.FakeRam[0x8001]].Name, "ASL")
	assert.Equal(t, Opcodes[C.Bus.FakeRam[0x8002]].Name, "STX")
	assert.Equal(t, Opcodes[C.Bus.FakeRam[0x801b]].Name, "NOP")
	assert.Equal(t, Opcodes[C.Bus.FakeRam[0x801c]].Name, "BRK")
}

func TestThirty(t *testing.T) {
	// this program is supposed to multiply 10 (0xa) by 3. the end state
	// should be:
	//
	// A=1e (30), X=3, Y=0
	// page 0: [0a 03 1e] (10 3 30)
	//
	// once this is done, 3 noops are called, then a BRK, which triggers an
	// NMI, effectively (writing a bunch of stuff to the stack and) jumping
	// to 0x0.
	//
	// at that point, the cpu decodes 1e and executes ASL on 0 in an
	// infinite loop.
	program := "A2 0A 8E 00 00 A2 03 8E 01 00 AC 00 00 A9 00 18 6D 01 00 88 D0 FA 8D 02 00 EA EA EA" // 28 bytes

	C := Cpu{Bus: &mem.Bus{}}
	// C.Debug([]byte(program), 0x8000)

	offset := uint16(0x8000)
	C.LoadProgram([]byte(program), offset)
	C.Bus.FakeRam[0xfffc] = 0x00 // reset
	C.Bus.FakeRam[0xfffd] = 0x80 // ?
	C.ProgramCounter = offset

	assert.Equal(t, Opcodes[C.Bus.FakeRam[C.ProgramCounter]].Name, "LDX")

	for _, cpuState := range []struct {
		M        uint8
		A        uint8
		X        uint8
		Y        uint8
		InstName string
	}{
		{M: 0xa, A: 0, X: 0xa, Y: 0, InstName: "STX"},
		{M: 0xa, A: 0, X: 0xa, Y: 0, InstName: "LDX"},
		{M: 3, A: 0, X: 3, Y: 0, InstName: "STX"},
		{M: 3, A: 0, X: 3, Y: 0, InstName: "LDY"},
		{M: 0xa, A: 0, X: 3, Y: 0xa, InstName: "LDA"},
		{M: 0, A: 0, X: 3, Y: 0xa, InstName: "CLC"},

		{M: 0, A: 0, X: 3, Y: 0xa, InstName: "ADC"},
		{M: 3, A: 3, X: 3, Y: 0xa, InstName: "DEY"},
		{M: 3, A: 3, X: 3, Y: 9, InstName: "BNE"},

		{M: 0x6d, A: 3, X: 3, Y: 9, InstName: "ADC"}, // note: we jumped back
		{M: 0x03, A: 6, X: 3, Y: 9, InstName: "DEY"},
		{M: 0x03, A: 6, X: 3, Y: 8, InstName: "BNE"},

		// {{{
		{M: 0x6d, A: 6, X: 3, Y: 8, InstName: "ADC"},
		{M: 0x03, A: 9, X: 3, Y: 8, InstName: "DEY"},
		{M: 0x03, A: 9, X: 3, Y: 7, InstName: "BNE"},

		{M: 0x6d, A: 9, X: 3, Y: 7, InstName: "ADC"},
		{M: 0x03, A: 12, X: 3, Y: 7, InstName: "DEY"},
		{M: 0x03, A: 12, X: 3, Y: 6, InstName: "BNE"},

		{M: 0x6d, A: 12, X: 3, Y: 6, InstName: "ADC"},
		{M: 0x03, A: 15, X: 3, Y: 6, InstName: "DEY"},
		{M: 0x03, A: 15, X: 3, Y: 5, InstName: "BNE"},

		{M: 0x6d, A: 15, X: 3, Y: 5, InstName: "ADC"},
		{M: 0x03, A: 18, X: 3, Y: 5, InstName: "DEY"},
		{M: 0x03, A: 18, X: 3, Y: 4, InstName: "BNE"},

		{M: 0x6d, A: 18, X: 3, Y: 4, InstName: "ADC"},
		{M: 0x03, A: 21, X: 3, Y: 4, InstName: "DEY"},
		{M: 0x03, A: 21, X: 3, Y: 3, InstName: "BNE"},

		{M: 0x6d, A: 21, X: 3, Y: 3, InstName: "ADC"},
		{M: 0x03, A: 24, X: 3, Y: 3, InstName: "DEY"},
		{M: 0x03, A: 24, X: 3, Y: 2, InstName: "BNE"},

		{M: 0x6d, A: 24, X: 3, Y: 2, InstName: "ADC"},
		{M: 0x03, A: 27, X: 3, Y: 2, InstName: "DEY"},
		{M: 0x03, A: 27, X: 3, Y: 1, InstName: "BNE"},

		{M: 0x6d, A: 27, X: 3, Y: 1, InstName: "ADC"},
		{M: 0x03, A: 30, X: 3, Y: 1, InstName: "DEY"},
		{M: 0x03, A: 30, X: 3, Y: 0, InstName: "BNE"},
		// }}}

		{M: 0x6d, A: 30, X: 3, Y: 0, InstName: "STA"},
		{M: 0x1e, A: 30, X: 3, Y: 0, InstName: "NOP"},
		{M: 0x1e, A: 30, X: 3, Y: 0, InstName: "NOP"},
		{M: 0x1e, A: 30, X: 3, Y: 0, InstName: "NOP"},
		{M: 0x1e, A: 30, X: 3, Y: 0, InstName: "BRK"},

		// UB from here on
		{M: 0x1e, A: 30, X: 3, Y: 0, InstName: "ASL"},
		{M: 0x78, A: 30, X: 3, Y: 0, InstName: ""},
	} {
		_ = C.tick()
		currInst := Opcodes[C.Bus.FakeRam[C.ProgramCounter]].Name
		assert.Equal(t, C.M, cpuState.M, "incorrect M at %s", currInst)
		assert.Equal(t, C.Accumulator, cpuState.A, "incorrect A at %s", currInst)
		assert.Equal(t, C.X, cpuState.X, "incorrect X at %s", currInst)
		assert.Equal(t, C.Y, cpuState.Y, "incorrect Y at %s", currInst)
		assert.Equal(t, currInst, cpuState.InstName)
	}

	assert.Equal(t, C.Bus.FakeRam[0], uint8(10))
	assert.Equal(t, C.Bus.FakeRam[1], uint8(3))
	assert.Equal(t, C.Bus.FakeRam[2], uint8(30))
}
