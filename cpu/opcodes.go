package cpu

// An Opcode is associated with a unique byte Value (0x00-0xff). There are 256
// possible opcodes (16x16), but only 56 correspond to a valid Cpu instruction.
//
// Importantly, the Opcode carries with it information on the AddressingMode
// and number of Cycles that should elapse before the corresponding Instruction
// completes.
//
// Multiple Opcodes may execute the same Instruction, differing only in how the
// data is to be retrieved; this is handled by the Cpu, not the Instruction
// itself.
type Opcode struct {
	AddressingMode AddressingMode

	// Clock cycles required; typically 2 to 7 (hence a byte). Longer
	// instructions require more cycles to fetch and decode memory.
	//
	// Presumably, modern hardware (with a modern language) will be fast
	// enough to execute instructions at the speed of an NES (or faster).
	// Thus, clock cycles are more important in guaranteeing
	// synchronisation between different components.
	//
	// https://www.nesdev.org/wiki/Cycle_counting#Instruction_timings
	Cycles byte

	// these fields are less than ideal:
	// variadic does not actually let you define impls as func()
	// Instruction         func(...byte) byte
	// this is encapsulated, but unclear/misleading
	// Instruction         func(Instruction) byte

	// Args are passed to the func via the M field of c, not func args.
	// Generally, at most one arg is required.
	//
	// The byte returned by the Instruction call is not memory data. It
	// only indicates how many extra Cycles to wait. This number varies
	// greatly, depending on the AddressingMode, and the Instruction
	// itself.
	//
	// Theoretically, the returned byte may not be needed, since it is
	// already stored in Opcode.Cycles.
	Instruction func(c *Cpu) byte

	// Value               byte // The Value received by the CPU
	// NumBytes            int  // Always 1 to 3 (needed?)
	// CrossesPageBoundary bool // if true, increases Cycles (i.e. wait more ticks)
}

// The Opcodes table lists all 128 (?) byte values recognised by the Cpu. 56
// unique instructions.
var Opcodes = map[byte]Opcode{
	// Generated from http://www.6502.org/tutorials/6502opcodes.html

	// OLC uses a 16x16 slice (Vec, in cpp) of opcodes, i.e. exhaustive.
	// Crucially, he chooses to 1) contain all opcodes within the (global)
	// Cpu struct, and 2) access them via slice index, also contained
	// within the Cpu struct. (AddressingMode and Instruction are also
	// contained in the Cpu)
	//
	// Not: in the event of illegal opcode (outside the 56), OLC invokes
	// either NOP or XXX; not sure of the significance of this.
	//
	// My approach is to 1) allow Opcodes to exist as entities of their
	// own, but allow them to reference an arbitrary (parent) Cpu via
	// pointer, and 2) access Opcode via a (global) map, and default to a
	// dummy opcode.

	// TODO: is map index or slice index faster? slice index requires slice
	// of len 256, map doesn't

	// 0x65: {Instruction: ADC}, // requires global (or at least packaged) func
	// 0x69: {Instruction: C.ADC}, // requires global C (var C Cpu)

	// 0x69: {Instruction: (*Cpu).ADC, Cycles: 2, AddressingMode: Immediate},
	// 0x65: {Instruction: (*Cpu).ADC, Cycles: 3, AddressingMode: ZeroPage},
	// 0x75: {Instruction: (*Cpu).ADC, Cycles: 4, AddressingMode: ZeroPageX},
	// 0x6d: {Instruction: (*Cpu).ADC, Cycles: 4, AddressingMode: Absolute},
	// 0x7d: {Instruction: (*Cpu).ADC, Cycles: 4, AddressingMode: AbsoluteX},
	// 0x79: {Instruction: (*Cpu).ADC, Cycles: 4, AddressingMode: AbsoluteY},
	// 0x61: {Instruction: (*Cpu).ADC, Cycles: 6, AddressingMode: IndirectX},
	// 0x71: {Instruction: (*Cpu).ADC, Cycles: 5, AddressingMode: IndirectY},

	0x69: {Instruction: (*Cpu).ADC, Cycles: 2, AddressingMode: Immediate},
	0x65: {Instruction: (*Cpu).ADC, Cycles: 3, AddressingMode: ZeroPage},
	0x75: {Instruction: (*Cpu).ADC, Cycles: 4, AddressingMode: ZeroPageX},
	0x6D: {Instruction: (*Cpu).ADC, Cycles: 4, AddressingMode: Absolute},
	0x7D: {Instruction: (*Cpu).ADC, Cycles: 4, AddressingMode: AbsoluteX},
	0x79: {Instruction: (*Cpu).ADC, Cycles: 4, AddressingMode: AbsoluteY},
	0x61: {Instruction: (*Cpu).ADC, Cycles: 6, AddressingMode: IndirectX},
	0x71: {Instruction: (*Cpu).ADC, Cycles: 5, AddressingMode: IndirectY},
	0x29: {Instruction: (*Cpu).AND, Cycles: 2, AddressingMode: Immediate},
	0x25: {Instruction: (*Cpu).AND, Cycles: 3, AddressingMode: ZeroPage},
	0x35: {Instruction: (*Cpu).AND, Cycles: 4, AddressingMode: ZeroPageX},
	0x2D: {Instruction: (*Cpu).AND, Cycles: 4, AddressingMode: Absolute},
	0x3D: {Instruction: (*Cpu).AND, Cycles: 4, AddressingMode: AbsoluteX},
	0x39: {Instruction: (*Cpu).AND, Cycles: 4, AddressingMode: AbsoluteY},
	0x21: {Instruction: (*Cpu).AND, Cycles: 6, AddressingMode: IndirectX},
	0x31: {Instruction: (*Cpu).AND, Cycles: 5, AddressingMode: IndirectY},
	0x0A: {Instruction: (*Cpu).ASL, Cycles: 2, AddressingMode: Accumulator},
	0x06: {Instruction: (*Cpu).ASL, Cycles: 5, AddressingMode: ZeroPage},
	0x16: {Instruction: (*Cpu).ASL, Cycles: 6, AddressingMode: ZeroPageX},
	0x0E: {Instruction: (*Cpu).ASL, Cycles: 6, AddressingMode: Absolute},
	0x1E: {Instruction: (*Cpu).ASL, Cycles: 7, AddressingMode: AbsoluteX},
	0x24: {Instruction: (*Cpu).BIT, Cycles: 3, AddressingMode: ZeroPage},
	0x2C: {Instruction: (*Cpu).BIT, Cycles: 4, AddressingMode: Absolute},
	0x00: {Instruction: (*Cpu).BRK, Cycles: 7, AddressingMode: Implied},
	0xC9: {Instruction: (*Cpu).CMP, Cycles: 2, AddressingMode: Immediate},
	0xC5: {Instruction: (*Cpu).CMP, Cycles: 3, AddressingMode: ZeroPage},
	0xD5: {Instruction: (*Cpu).CMP, Cycles: 4, AddressingMode: ZeroPageX},
	0xCD: {Instruction: (*Cpu).CMP, Cycles: 4, AddressingMode: Absolute},
	0xDD: {Instruction: (*Cpu).CMP, Cycles: 4, AddressingMode: AbsoluteX},
	0xD9: {Instruction: (*Cpu).CMP, Cycles: 4, AddressingMode: AbsoluteY},
	0xC1: {Instruction: (*Cpu).CMP, Cycles: 6, AddressingMode: IndirectX},
	0xD1: {Instruction: (*Cpu).CMP, Cycles: 5, AddressingMode: IndirectY},
	0xE0: {Instruction: (*Cpu).CPX, Cycles: 2, AddressingMode: Immediate},
	0xE4: {Instruction: (*Cpu).CPX, Cycles: 3, AddressingMode: ZeroPage},
	0xEC: {Instruction: (*Cpu).CPX, Cycles: 4, AddressingMode: Absolute},
	0xC0: {Instruction: (*Cpu).CPY, Cycles: 2, AddressingMode: Immediate},
	0xC4: {Instruction: (*Cpu).CPY, Cycles: 3, AddressingMode: ZeroPage},
	0xCC: {Instruction: (*Cpu).CPY, Cycles: 4, AddressingMode: Absolute},
	0xC6: {Instruction: (*Cpu).DEC, Cycles: 5, AddressingMode: ZeroPage},
	0xD6: {Instruction: (*Cpu).DEC, Cycles: 6, AddressingMode: ZeroPageX},
	0xCE: {Instruction: (*Cpu).DEC, Cycles: 6, AddressingMode: Absolute},
	0xDE: {Instruction: (*Cpu).DEC, Cycles: 7, AddressingMode: AbsoluteX},
	0x49: {Instruction: (*Cpu).EOR, Cycles: 2, AddressingMode: Immediate},
	0x45: {Instruction: (*Cpu).EOR, Cycles: 3, AddressingMode: ZeroPage},
	0x55: {Instruction: (*Cpu).EOR, Cycles: 4, AddressingMode: ZeroPageX},
	0x4D: {Instruction: (*Cpu).EOR, Cycles: 4, AddressingMode: Absolute},
	0x5D: {Instruction: (*Cpu).EOR, Cycles: 4, AddressingMode: AbsoluteX},
	0x59: {Instruction: (*Cpu).EOR, Cycles: 4, AddressingMode: AbsoluteY},
	0x41: {Instruction: (*Cpu).EOR, Cycles: 6, AddressingMode: IndirectX},
	0x51: {Instruction: (*Cpu).EOR, Cycles: 5, AddressingMode: IndirectY},
	0xE6: {Instruction: (*Cpu).INC, Cycles: 5, AddressingMode: ZeroPage},
	0xF6: {Instruction: (*Cpu).INC, Cycles: 6, AddressingMode: ZeroPageX},
	0xEE: {Instruction: (*Cpu).INC, Cycles: 6, AddressingMode: Absolute},
	0xFE: {Instruction: (*Cpu).INC, Cycles: 7, AddressingMode: AbsoluteX},
	0x4C: {Instruction: (*Cpu).JMP, Cycles: 3, AddressingMode: Absolute},
	0x6C: {Instruction: (*Cpu).JMP, Cycles: 5, AddressingMode: Indirect},
	0x20: {Instruction: (*Cpu).JSR, Cycles: 6, AddressingMode: Absolute},
	0xA9: {Instruction: (*Cpu).LDA, Cycles: 2, AddressingMode: Immediate},
	0xA5: {Instruction: (*Cpu).LDA, Cycles: 3, AddressingMode: ZeroPage},
	0xB5: {Instruction: (*Cpu).LDA, Cycles: 4, AddressingMode: ZeroPageX},
	0xAD: {Instruction: (*Cpu).LDA, Cycles: 4, AddressingMode: Absolute},
	0xBD: {Instruction: (*Cpu).LDA, Cycles: 4, AddressingMode: AbsoluteX},
	0xB9: {Instruction: (*Cpu).LDA, Cycles: 4, AddressingMode: AbsoluteY},
	0xA1: {Instruction: (*Cpu).LDA, Cycles: 6, AddressingMode: IndirectX},
	0xB1: {Instruction: (*Cpu).LDA, Cycles: 5, AddressingMode: IndirectY},
	0xA2: {Instruction: (*Cpu).LDX, Cycles: 2, AddressingMode: Immediate},
	0xA6: {Instruction: (*Cpu).LDX, Cycles: 3, AddressingMode: ZeroPage},
	0xB6: {Instruction: (*Cpu).LDX, Cycles: 4, AddressingMode: ZeroPageY},
	0xAE: {Instruction: (*Cpu).LDX, Cycles: 4, AddressingMode: Absolute},
	0xBE: {Instruction: (*Cpu).LDX, Cycles: 4, AddressingMode: AbsoluteY},
	0xA0: {Instruction: (*Cpu).LDY, Cycles: 2, AddressingMode: Immediate},
	0xA4: {Instruction: (*Cpu).LDY, Cycles: 3, AddressingMode: ZeroPage},
	0xB4: {Instruction: (*Cpu).LDY, Cycles: 4, AddressingMode: ZeroPageX},
	0xAC: {Instruction: (*Cpu).LDY, Cycles: 4, AddressingMode: Absolute},
	0xBC: {Instruction: (*Cpu).LDY, Cycles: 4, AddressingMode: AbsoluteX},
	0x4A: {Instruction: (*Cpu).LSR, Cycles: 2, AddressingMode: Accumulator},
	0x46: {Instruction: (*Cpu).LSR, Cycles: 5, AddressingMode: ZeroPage},
	0x56: {Instruction: (*Cpu).LSR, Cycles: 6, AddressingMode: ZeroPageX},
	0x4E: {Instruction: (*Cpu).LSR, Cycles: 6, AddressingMode: Absolute},
	0x5E: {Instruction: (*Cpu).LSR, Cycles: 7, AddressingMode: AbsoluteX},
	0xEA: {Instruction: (*Cpu).NOP, Cycles: 2, AddressingMode: Implied},
	0x09: {Instruction: (*Cpu).ORA, Cycles: 2, AddressingMode: Immediate},
	0x05: {Instruction: (*Cpu).ORA, Cycles: 3, AddressingMode: ZeroPage},
	0x15: {Instruction: (*Cpu).ORA, Cycles: 4, AddressingMode: ZeroPageX},
	0x0D: {Instruction: (*Cpu).ORA, Cycles: 4, AddressingMode: Absolute},
	0x1D: {Instruction: (*Cpu).ORA, Cycles: 4, AddressingMode: AbsoluteX},
	0x19: {Instruction: (*Cpu).ORA, Cycles: 4, AddressingMode: AbsoluteY},
	0x01: {Instruction: (*Cpu).ORA, Cycles: 6, AddressingMode: IndirectX},
	0x11: {Instruction: (*Cpu).ORA, Cycles: 5, AddressingMode: IndirectY},
	0x2A: {Instruction: (*Cpu).ROL, Cycles: 2, AddressingMode: Accumulator},
	0x26: {Instruction: (*Cpu).ROL, Cycles: 5, AddressingMode: ZeroPage},
	0x36: {Instruction: (*Cpu).ROL, Cycles: 6, AddressingMode: ZeroPageX},
	0x2E: {Instruction: (*Cpu).ROL, Cycles: 6, AddressingMode: Absolute},
	0x3E: {Instruction: (*Cpu).ROL, Cycles: 7, AddressingMode: AbsoluteX},
	0x6A: {Instruction: (*Cpu).ROR, Cycles: 2, AddressingMode: Accumulator},
	0x66: {Instruction: (*Cpu).ROR, Cycles: 5, AddressingMode: ZeroPage},
	0x76: {Instruction: (*Cpu).ROR, Cycles: 6, AddressingMode: ZeroPageX},
	0x6E: {Instruction: (*Cpu).ROR, Cycles: 6, AddressingMode: Absolute},
	0x7E: {Instruction: (*Cpu).ROR, Cycles: 7, AddressingMode: AbsoluteX},
	0x40: {Instruction: (*Cpu).RTI, Cycles: 6, AddressingMode: Implied},
	0x60: {Instruction: (*Cpu).RTS, Cycles: 6, AddressingMode: Implied},
	0xE9: {Instruction: (*Cpu).SBC, Cycles: 2, AddressingMode: Immediate},
	0xE5: {Instruction: (*Cpu).SBC, Cycles: 3, AddressingMode: ZeroPage},
	0xF5: {Instruction: (*Cpu).SBC, Cycles: 4, AddressingMode: ZeroPageX},
	0xED: {Instruction: (*Cpu).SBC, Cycles: 4, AddressingMode: Absolute},
	0xFD: {Instruction: (*Cpu).SBC, Cycles: 4, AddressingMode: AbsoluteX},
	0xF9: {Instruction: (*Cpu).SBC, Cycles: 4, AddressingMode: AbsoluteY},
	0xE1: {Instruction: (*Cpu).SBC, Cycles: 6, AddressingMode: IndirectX},
	0xF1: {Instruction: (*Cpu).SBC, Cycles: 5, AddressingMode: IndirectY},
	0x85: {Instruction: (*Cpu).STA, Cycles: 3, AddressingMode: ZeroPage},
	0x95: {Instruction: (*Cpu).STA, Cycles: 4, AddressingMode: ZeroPageX},
	0x8D: {Instruction: (*Cpu).STA, Cycles: 4, AddressingMode: Absolute},
	0x9D: {Instruction: (*Cpu).STA, Cycles: 5, AddressingMode: AbsoluteX},
	0x99: {Instruction: (*Cpu).STA, Cycles: 5, AddressingMode: AbsoluteY},
	0x81: {Instruction: (*Cpu).STA, Cycles: 6, AddressingMode: IndirectX},
	0x91: {Instruction: (*Cpu).STA, Cycles: 6, AddressingMode: IndirectY},
	// ptr): {Instruction: (*Cpu).X to, Cycles: 2, AddressingMode: TXS(Transf},
	// X): {Instruction: (*Cpu).Stac, Cycles: 2, AddressingMode: TSX(Transf},
	// Accumulator): {Instruction: (*Cpu).mula, Cycles: 3, AddressingMode: PHA(PusHA},
	// Accumulator): {Instruction: (*Cpu).mula, Cycles: 4, AddressingMode: PLA(PuLlA},
	// status): {Instruction: (*Cpu).esso, Cycles: 3, AddressingMode: PHP(PusHP},
	// status): {Instruction: (*Cpu).esso, Cycles: 4, AddressingMode: PLP(PuLlP},
	0x86: {Instruction: (*Cpu).STX, Cycles: 3, AddressingMode: ZeroPage},
	0x96: {Instruction: (*Cpu).STX, Cycles: 4, AddressingMode: ZeroPageY},
	0x8E: {Instruction: (*Cpu).STX, Cycles: 4, AddressingMode: Absolute},
	0x84: {Instruction: (*Cpu).STY, Cycles: 3, AddressingMode: ZeroPage},
	0x94: {Instruction: (*Cpu).STY, Cycles: 4, AddressingMode: ZeroPageX},
	0x8C: {Instruction: (*Cpu).STY, Cycles: 4, AddressingMode: Absolute},
}
