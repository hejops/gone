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

	// An Instruction usually modifies or copies register(s). Args (usually
	// just a byte) are passed to the func implicitly via the M field of c,
	// -not- explicitly via func args.
	//
	// With the sole exception of BRK, Instructions -never- increment the
	// PC. Cycles are incurred -only- in branching Instructions, and only
	// if the condition succeeds.
	//
	// The byte returned by the Instruction call is not memory data. It
	// only indicates how many extra Cycles to wait. (This is from OLC, and
	// may be rethought) This number varies greatly, depending on the
	// AddressingMode, and the Instruction itself.
	//
	// Theoretically, the returned byte may not be needed, since it is
	// already stored in Opcode.Cycles.
	Instruction func(c *Cpu) byte

	Name string // for debugger

	// Value               byte // The Value received by the CPU
	// NumBytes            int  // Always 1 to 3 (needed?)
	// CrossesPageBoundary bool // if true, increases Cycles (i.e. wait more ticks)
}

// The Opcodes table lists all 151 (?) byte values recognised by the Cpu. These
// values are mapped to 56 unique instructions.
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

	0x69: {Instruction: (*Cpu).ADC, Name: "ADC", Cycles: 2, AddressingMode: Immediate},
	0x65: {Instruction: (*Cpu).ADC, Name: "ADC", Cycles: 3, AddressingMode: ZeroPage},
	0x75: {Instruction: (*Cpu).ADC, Name: "ADC", Cycles: 4, AddressingMode: ZeroPageX},
	0x6D: {Instruction: (*Cpu).ADC, Name: "ADC", Cycles: 4, AddressingMode: Absolute},
	0x7D: {Instruction: (*Cpu).ADC, Name: "ADC", Cycles: 4, AddressingMode: AbsoluteX},
	0x79: {Instruction: (*Cpu).ADC, Name: "ADC", Cycles: 4, AddressingMode: AbsoluteY},
	0x61: {Instruction: (*Cpu).ADC, Name: "ADC", Cycles: 6, AddressingMode: IndirectX},
	0x71: {Instruction: (*Cpu).ADC, Name: "ADC", Cycles: 5, AddressingMode: IndirectY},
	0x29: {Instruction: (*Cpu).AND, Name: "AND", Cycles: 2, AddressingMode: Immediate},
	0x25: {Instruction: (*Cpu).AND, Name: "AND", Cycles: 3, AddressingMode: ZeroPage},
	0x35: {Instruction: (*Cpu).AND, Name: "AND", Cycles: 4, AddressingMode: ZeroPageX},
	0x2D: {Instruction: (*Cpu).AND, Name: "AND", Cycles: 4, AddressingMode: Absolute},
	0x3D: {Instruction: (*Cpu).AND, Name: "AND", Cycles: 4, AddressingMode: AbsoluteX},
	0x39: {Instruction: (*Cpu).AND, Name: "AND", Cycles: 4, AddressingMode: AbsoluteY},
	0x21: {Instruction: (*Cpu).AND, Name: "AND", Cycles: 6, AddressingMode: IndirectX},
	0x31: {Instruction: (*Cpu).AND, Name: "AND", Cycles: 5, AddressingMode: IndirectY},
	0x0A: {Instruction: (*Cpu).ASL, Name: "ASL", Cycles: 2, AddressingMode: Accumulator},
	0x06: {Instruction: (*Cpu).ASL, Name: "ASL", Cycles: 5, AddressingMode: ZeroPage},
	0x16: {Instruction: (*Cpu).ASL, Name: "ASL", Cycles: 6, AddressingMode: ZeroPageX},
	0x0E: {Instruction: (*Cpu).ASL, Name: "ASL", Cycles: 6, AddressingMode: Absolute},
	0x1E: {Instruction: (*Cpu).ASL, Name: "ASL", Cycles: 7, AddressingMode: AbsoluteX},
	0x24: {Instruction: (*Cpu).BIT, Name: "BIT", Cycles: 3, AddressingMode: ZeroPage},
	0x2C: {Instruction: (*Cpu).BIT, Name: "BIT", Cycles: 4, AddressingMode: Absolute},
	0x00: {Instruction: (*Cpu).BRK, Name: "BRK", Cycles: 7, AddressingMode: Implied},
	0xC9: {Instruction: (*Cpu).CMP, Name: "CMP", Cycles: 2, AddressingMode: Immediate},
	0xC5: {Instruction: (*Cpu).CMP, Name: "CMP", Cycles: 3, AddressingMode: ZeroPage},
	0xD5: {Instruction: (*Cpu).CMP, Name: "CMP", Cycles: 4, AddressingMode: ZeroPageX},
	0xCD: {Instruction: (*Cpu).CMP, Name: "CMP", Cycles: 4, AddressingMode: Absolute},
	0xDD: {Instruction: (*Cpu).CMP, Name: "CMP", Cycles: 4, AddressingMode: AbsoluteX},
	0xD9: {Instruction: (*Cpu).CMP, Name: "CMP", Cycles: 4, AddressingMode: AbsoluteY},
	0xC1: {Instruction: (*Cpu).CMP, Name: "CMP", Cycles: 6, AddressingMode: IndirectX},
	0xD1: {Instruction: (*Cpu).CMP, Name: "CMP", Cycles: 5, AddressingMode: IndirectY},
	0xE0: {Instruction: (*Cpu).CPX, Name: "CPX", Cycles: 2, AddressingMode: Immediate},
	0xE4: {Instruction: (*Cpu).CPX, Name: "CPX", Cycles: 3, AddressingMode: ZeroPage},
	0xEC: {Instruction: (*Cpu).CPX, Name: "CPX", Cycles: 4, AddressingMode: Absolute},
	0xC0: {Instruction: (*Cpu).CPY, Name: "CPY", Cycles: 2, AddressingMode: Immediate},
	0xC4: {Instruction: (*Cpu).CPY, Name: "CPY", Cycles: 3, AddressingMode: ZeroPage},
	0xCC: {Instruction: (*Cpu).CPY, Name: "CPY", Cycles: 4, AddressingMode: Absolute},
	0xC6: {Instruction: (*Cpu).DEC, Name: "DEC", Cycles: 5, AddressingMode: ZeroPage},
	0xD6: {Instruction: (*Cpu).DEC, Name: "DEC", Cycles: 6, AddressingMode: ZeroPageX},
	0xCE: {Instruction: (*Cpu).DEC, Name: "DEC", Cycles: 6, AddressingMode: Absolute},
	0xDE: {Instruction: (*Cpu).DEC, Name: "DEC", Cycles: 7, AddressingMode: AbsoluteX},
	0x49: {Instruction: (*Cpu).EOR, Name: "EOR", Cycles: 2, AddressingMode: Immediate},
	0x45: {Instruction: (*Cpu).EOR, Name: "EOR", Cycles: 3, AddressingMode: ZeroPage},
	0x55: {Instruction: (*Cpu).EOR, Name: "EOR", Cycles: 4, AddressingMode: ZeroPageX},
	0x4D: {Instruction: (*Cpu).EOR, Name: "EOR", Cycles: 4, AddressingMode: Absolute},
	0x5D: {Instruction: (*Cpu).EOR, Name: "EOR", Cycles: 4, AddressingMode: AbsoluteX},
	0x59: {Instruction: (*Cpu).EOR, Name: "EOR", Cycles: 4, AddressingMode: AbsoluteY},
	0x41: {Instruction: (*Cpu).EOR, Name: "EOR", Cycles: 6, AddressingMode: IndirectX},
	0x51: {Instruction: (*Cpu).EOR, Name: "EOR", Cycles: 5, AddressingMode: IndirectY},
	0xE6: {Instruction: (*Cpu).INC, Name: "INC", Cycles: 5, AddressingMode: ZeroPage},
	0xF6: {Instruction: (*Cpu).INC, Name: "INC", Cycles: 6, AddressingMode: ZeroPageX},
	0xEE: {Instruction: (*Cpu).INC, Name: "INC", Cycles: 6, AddressingMode: Absolute},
	0xFE: {Instruction: (*Cpu).INC, Name: "INC", Cycles: 7, AddressingMode: AbsoluteX},
	0x4C: {Instruction: (*Cpu).JMP, Name: "JMP", Cycles: 3, AddressingMode: Absolute},
	0x6C: {Instruction: (*Cpu).JMP, Name: "JMP", Cycles: 5, AddressingMode: Indirect},
	0x20: {Instruction: (*Cpu).JSR, Name: "JSR", Cycles: 6, AddressingMode: Absolute},
	0xA9: {Instruction: (*Cpu).LDA, Name: "LDA", Cycles: 2, AddressingMode: Immediate},
	0xA5: {Instruction: (*Cpu).LDA, Name: "LDA", Cycles: 3, AddressingMode: ZeroPage},
	0xB5: {Instruction: (*Cpu).LDA, Name: "LDA", Cycles: 4, AddressingMode: ZeroPageX},
	0xAD: {Instruction: (*Cpu).LDA, Name: "LDA", Cycles: 4, AddressingMode: Absolute},
	0xBD: {Instruction: (*Cpu).LDA, Name: "LDA", Cycles: 4, AddressingMode: AbsoluteX},
	0xB9: {Instruction: (*Cpu).LDA, Name: "LDA", Cycles: 4, AddressingMode: AbsoluteY},
	0xA1: {Instruction: (*Cpu).LDA, Name: "LDA", Cycles: 6, AddressingMode: IndirectX},
	0xB1: {Instruction: (*Cpu).LDA, Name: "LDA", Cycles: 5, AddressingMode: IndirectY},
	0xA2: {Instruction: (*Cpu).LDX, Name: "LDX", Cycles: 2, AddressingMode: Immediate},
	0xA6: {Instruction: (*Cpu).LDX, Name: "LDX", Cycles: 3, AddressingMode: ZeroPage},
	0xB6: {Instruction: (*Cpu).LDX, Name: "LDX", Cycles: 4, AddressingMode: ZeroPageY},
	0xAE: {Instruction: (*Cpu).LDX, Name: "LDX", Cycles: 4, AddressingMode: Absolute},
	0xBE: {Instruction: (*Cpu).LDX, Name: "LDX", Cycles: 4, AddressingMode: AbsoluteY},
	0xA0: {Instruction: (*Cpu).LDY, Name: "LDY", Cycles: 2, AddressingMode: Immediate},
	0xA4: {Instruction: (*Cpu).LDY, Name: "LDY", Cycles: 3, AddressingMode: ZeroPage},
	0xB4: {Instruction: (*Cpu).LDY, Name: "LDY", Cycles: 4, AddressingMode: ZeroPageX},
	0xAC: {Instruction: (*Cpu).LDY, Name: "LDY", Cycles: 4, AddressingMode: Absolute},
	0xBC: {Instruction: (*Cpu).LDY, Name: "LDY", Cycles: 4, AddressingMode: AbsoluteX},
	0x4A: {Instruction: (*Cpu).LSR, Name: "LSR", Cycles: 2, AddressingMode: Accumulator},
	0x46: {Instruction: (*Cpu).LSR, Name: "LSR", Cycles: 5, AddressingMode: ZeroPage},
	0x56: {Instruction: (*Cpu).LSR, Name: "LSR", Cycles: 6, AddressingMode: ZeroPageX},
	0x4E: {Instruction: (*Cpu).LSR, Name: "LSR", Cycles: 6, AddressingMode: Absolute},
	0x5E: {Instruction: (*Cpu).LSR, Name: "LSR", Cycles: 7, AddressingMode: AbsoluteX},
	0xEA: {Instruction: (*Cpu).NOP, Name: "NOP", Cycles: 2, AddressingMode: Implied},
	0x09: {Instruction: (*Cpu).ORA, Name: "ORA", Cycles: 2, AddressingMode: Immediate},
	0x05: {Instruction: (*Cpu).ORA, Name: "ORA", Cycles: 3, AddressingMode: ZeroPage},
	0x15: {Instruction: (*Cpu).ORA, Name: "ORA", Cycles: 4, AddressingMode: ZeroPageX},
	0x0D: {Instruction: (*Cpu).ORA, Name: "ORA", Cycles: 4, AddressingMode: Absolute},
	0x1D: {Instruction: (*Cpu).ORA, Name: "ORA", Cycles: 4, AddressingMode: AbsoluteX},
	0x19: {Instruction: (*Cpu).ORA, Name: "ORA", Cycles: 4, AddressingMode: AbsoluteY},
	0x01: {Instruction: (*Cpu).ORA, Name: "ORA", Cycles: 6, AddressingMode: IndirectX},
	0x11: {Instruction: (*Cpu).ORA, Name: "ORA", Cycles: 5, AddressingMode: IndirectY},
	0x2A: {Instruction: (*Cpu).ROL, Name: "ROL", Cycles: 2, AddressingMode: Accumulator},
	0x26: {Instruction: (*Cpu).ROL, Name: "ROL", Cycles: 5, AddressingMode: ZeroPage},
	0x36: {Instruction: (*Cpu).ROL, Name: "ROL", Cycles: 6, AddressingMode: ZeroPageX},
	0x2E: {Instruction: (*Cpu).ROL, Name: "ROL", Cycles: 6, AddressingMode: Absolute},
	0x3E: {Instruction: (*Cpu).ROL, Name: "ROL", Cycles: 7, AddressingMode: AbsoluteX},
	0x6A: {Instruction: (*Cpu).ROR, Name: "ROR", Cycles: 2, AddressingMode: Accumulator},
	0x66: {Instruction: (*Cpu).ROR, Name: "ROR", Cycles: 5, AddressingMode: ZeroPage},
	0x76: {Instruction: (*Cpu).ROR, Name: "ROR", Cycles: 6, AddressingMode: ZeroPageX},
	0x6E: {Instruction: (*Cpu).ROR, Name: "ROR", Cycles: 6, AddressingMode: Absolute},
	0x7E: {Instruction: (*Cpu).ROR, Name: "ROR", Cycles: 7, AddressingMode: AbsoluteX},
	0x40: {Instruction: (*Cpu).RTI, Name: "RTI", Cycles: 6, AddressingMode: Implied},
	0x60: {Instruction: (*Cpu).RTS, Name: "RTS", Cycles: 6, AddressingMode: Implied},
	0xE9: {Instruction: (*Cpu).SBC, Name: "SBC", Cycles: 2, AddressingMode: Immediate},
	0xE5: {Instruction: (*Cpu).SBC, Name: "SBC", Cycles: 3, AddressingMode: ZeroPage},
	0xF5: {Instruction: (*Cpu).SBC, Name: "SBC", Cycles: 4, AddressingMode: ZeroPageX},
	0xED: {Instruction: (*Cpu).SBC, Name: "SBC", Cycles: 4, AddressingMode: Absolute},
	0xFD: {Instruction: (*Cpu).SBC, Name: "SBC", Cycles: 4, AddressingMode: AbsoluteX},
	0xF9: {Instruction: (*Cpu).SBC, Name: "SBC", Cycles: 4, AddressingMode: AbsoluteY},
	0xE1: {Instruction: (*Cpu).SBC, Name: "SBC", Cycles: 6, AddressingMode: IndirectX},
	0xF1: {Instruction: (*Cpu).SBC, Name: "SBC", Cycles: 5, AddressingMode: IndirectY},
	0x85: {Instruction: (*Cpu).STA, Name: "STA", Cycles: 3, AddressingMode: ZeroPage},
	0x95: {Instruction: (*Cpu).STA, Name: "STA", Cycles: 4, AddressingMode: ZeroPageX},
	0x8D: {Instruction: (*Cpu).STA, Name: "STA", Cycles: 4, AddressingMode: Absolute},
	0x9D: {Instruction: (*Cpu).STA, Name: "STA", Cycles: 5, AddressingMode: AbsoluteX},
	0x99: {Instruction: (*Cpu).STA, Name: "STA", Cycles: 5, AddressingMode: AbsoluteY},
	0x81: {Instruction: (*Cpu).STA, Name: "STA", Cycles: 6, AddressingMode: IndirectX},
	0x91: {Instruction: (*Cpu).STA, Name: "STA", Cycles: 6, AddressingMode: IndirectY},
	0x86: {Instruction: (*Cpu).STX, Name: "STX", Cycles: 3, AddressingMode: ZeroPage},
	0x96: {Instruction: (*Cpu).STX, Name: "STX", Cycles: 4, AddressingMode: ZeroPageY},
	0x8E: {Instruction: (*Cpu).STX, Name: "STX", Cycles: 4, AddressingMode: Absolute},
	0x84: {Instruction: (*Cpu).STY, Name: "STY", Cycles: 3, AddressingMode: ZeroPage},
	0x94: {Instruction: (*Cpu).STY, Name: "STY", Cycles: 4, AddressingMode: ZeroPageX},
	0x8C: {Instruction: (*Cpu).STY, Name: "STY", Cycles: 4, AddressingMode: Absolute},

	// clear, set
	0x18: {Instruction: (*Cpu).CLC, Name: "CLC", Cycles: 2, AddressingMode: Implied},
	0x38: {Instruction: (*Cpu).SEC, Name: "SEC", Cycles: 2, AddressingMode: Implied},
	0x58: {Instruction: (*Cpu).CLI, Name: "CLI", Cycles: 2, AddressingMode: Implied},
	0x78: {Instruction: (*Cpu).SEI, Name: "SEI", Cycles: 2, AddressingMode: Implied},
	0xB8: {Instruction: (*Cpu).CLV, Name: "CLV", Cycles: 2, AddressingMode: Implied},
	0xD8: {Instruction: (*Cpu).CLD, Name: "CLD", Cycles: 2, AddressingMode: Implied},
	0xF8: {Instruction: (*Cpu).SED, Name: "SED", Cycles: 2, AddressingMode: Implied},

	// increment, decrement, transfer
	0xAA: {Instruction: (*Cpu).TAX, Name: "TAX", Cycles: 2, AddressingMode: Implied},
	0x8A: {Instruction: (*Cpu).TXA, Name: "TXA", Cycles: 2, AddressingMode: Implied},
	0xCA: {Instruction: (*Cpu).DEX, Name: "DEX", Cycles: 2, AddressingMode: Implied},
	0xE8: {Instruction: (*Cpu).INX, Name: "INX", Cycles: 2, AddressingMode: Implied},
	0xA8: {Instruction: (*Cpu).TAY, Name: "TAY", Cycles: 2, AddressingMode: Implied},
	0x98: {Instruction: (*Cpu).TYA, Name: "TYA", Cycles: 2, AddressingMode: Implied},
	0x88: {Instruction: (*Cpu).DEY, Name: "DEY", Cycles: 2, AddressingMode: Implied},
	0xC8: {Instruction: (*Cpu).INY, Name: "INY", Cycles: 2, AddressingMode: Implied},

	// branch
	0x10: {Instruction: (*Cpu).BPL, Name: "BPL", Cycles: 2, AddressingMode: Relative},
	0x30: {Instruction: (*Cpu).BMI, Name: "BMI", Cycles: 2, AddressingMode: Relative},
	0x50: {Instruction: (*Cpu).BVC, Name: "BVC", Cycles: 2, AddressingMode: Relative},
	0x70: {Instruction: (*Cpu).BVS, Name: "BVS", Cycles: 2, AddressingMode: Relative},
	0x90: {Instruction: (*Cpu).BCC, Name: "BCC", Cycles: 2, AddressingMode: Relative},
	0xB0: {Instruction: (*Cpu).BCS, Name: "BCS", Cycles: 2, AddressingMode: Relative},
	0xD0: {Instruction: (*Cpu).BNE, Name: "BNE", Cycles: 2, AddressingMode: Relative},
	0xF0: {Instruction: (*Cpu).BEQ, Name: "BEQ", Cycles: 2, AddressingMode: Relative},

	// stack
	0x9A: {Instruction: (*Cpu).TXS, Name: "TXS", Cycles: 2, AddressingMode: Implied},
	0xBA: {Instruction: (*Cpu).TSX, Name: "TSX", Cycles: 2, AddressingMode: Implied},
	0x48: {Instruction: (*Cpu).PHA, Name: "PHA", Cycles: 3, AddressingMode: Implied},
	0x68: {Instruction: (*Cpu).PLA, Name: "PLA", Cycles: 4, AddressingMode: Implied},
	0x08: {Instruction: (*Cpu).PHP, Name: "PHP", Cycles: 3, AddressingMode: Implied},
	0x28: {Instruction: (*Cpu).PLP, Name: "PLP", Cycles: 4, AddressingMode: Implied},
}
