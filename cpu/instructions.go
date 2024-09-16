package cpu

// all function signatures were automatically generated from
// https://www.nesdev.org/obelisk-6502-guide/reference.html

// 1-byte instr e.g. clc
// 2-byte instr e.g. 1-byte read: lda $41
// 3-byte instr e.g. 2-byte read: lda $0105
//
// what (func name), how many args, how long (cycles)
//
// read 1 byte (8 bits), split into 4+4 bits
// lower 4 bits = column
// upper 4 bits = row

// http://www.6502.org/tutorials/6502opcodes.html
// https://analog-hors.github.io/site/pones-p1/img/6502-opcode-table.png
// https://atariwiki.org/wiki/attach/OpCodes/OpCodes.jpg
// https://makingnesgames.com/Instruction_Set.html
// https://pbsandjay.github.io/
// https://problemkaputt.de/everynes.htm#cpuarithmeticlogicaloperations
// https://www.chibiakumas.com/book/CheatSheetCollection.pdf
// https://www.nesdev.org/obelisk-6502-guide/reference.html (best)

// how to read obelisk guide:
// A,Z,N = A&M
// [target],[flags...] = [op]

// helper funcs for flags

// Set Negative if M7 = 1
func (c *Cpu) setNegativeM7() { c.Flags.Negative = c.M&0x80 > 0 }

// Set Zero if Accumulator = 0
func (c *Cpu) setZero() { c.Flags.Zero = c.Accumulator == 0 }

// TODO: all instructions should PC++?

// ADC - Add with Carry
func (c *Cpu) ADC() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#ADC

	// c.Accumulator += c.M

	// in Go, overflow is possible at runtime, and results in wrapping
	// (e.g. 250+6=0), i.e. the sum will be less
	sum := c.Accumulator + c.M
	if sum < c.Accumulator {
		// C 	Carry Flag 	Set if overflow in bit 7
		// "Set if overflow in bit 7" -- bit 7 of what: A, M, or A+M? i
		// assume A+M for now
		c.Flags.Carry = true
	}

	c.Accumulator = sum
	if c.Flags.Carry {
		c.Accumulator += 1
	}

	// Z 	Zero Flag 	Set if A = 0
	c.setZero()
	// V 	Overflow Flag 	Set if sign bit is incorrect
	// c.Flags.Overflow = ???
	// N 	Negative Flag 	Set if bit 7 set
	c.Flags.Negative = c.Accumulator&0x80 > 0

	return 0
}

// AND - Logical AND
func (c *Cpu) AND() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#AND
	c.Accumulator &= c.M
	c.setZero()
	c.setNegativeM7()
	return 0
}

// ASL - Arithmetic Shift Left
func (c *Cpu) ASL() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#ASL
	c.Flags.Carry = c.M&0x80 > 0 // old bit 7
	c.M <<= 2
	c.setZero()
	c.setNegativeM7()
	return 0
}

// BCC - Branch if Carry Clear
func (c *Cpu) BCC() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#BCC
	if !c.Flags.Carry {
		c.ProgramCounter += uint16(c.RelAddress)
	}
	return 0
}

// BCS - Branch if Carry Set
func (c *Cpu) BCS() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#BCS
	if c.Flags.Carry {
		c.ProgramCounter += uint16(c.RelAddress)
	}
	return 0
}

// BEQ - Branch if Equal
func (c *Cpu) BEQ() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#BEQ
	if c.Flags.Zero {
		c.ProgramCounter += uint16(c.RelAddress)
	}
	return 0
}

// BIT - Bit Test
func (c *Cpu) BIT() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#BIT
	c.Flags.Zero = c.M&c.Accumulator > 0
	c.setNegativeM7()
	c.Flags.Overflow = c.M&0x40 > 0 // bit 6 set
	return 0
}

// BMI - Branch if Minus
func (c *Cpu) BMI() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#BMI
	if c.Flags.Negative {
		c.ProgramCounter += uint16(c.RelAddress)
	}
	return 0
}

// BNE - Branch if Not Equal
func (c *Cpu) BNE() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#BNE
	if !c.Flags.Zero {
		c.ProgramCounter += uint16(c.RelAddress)
	}
	return 0
}

// BPL - Branch if Positive
func (c *Cpu) BPL() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#BPL
	if !c.Flags.Negative {
		c.ProgramCounter += uint16(c.RelAddress)
	}
	return 0
}

// BRK - Force Interrupt
func (c *Cpu) BRK() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#BRK
	c.ProgramCounter++
	c.nmi()
	return 0
}

// BVC - Branch if Overflow Clear
func (c *Cpu) BVC() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#BVC
	if !c.Flags.Overflow {
		c.ProgramCounter += uint16(c.RelAddress)
	}
	return 0
}

// BVS - Branch if Overflow Set
func (c *Cpu) BVS() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#BVS
	if c.Flags.Overflow {
		c.ProgramCounter += uint16(c.RelAddress)
	}
	return 0
}

// CLC - Clear Carry Flag
func (c *Cpu) CLC() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#CLC
	c.Flags.Carry = false
	return 0
}

// CLD - Clear Decimal Mode
func (c *Cpu) CLD() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#CLD
	c.Flags.Decimal = false
	return 0
}

// CLI - Clear Interrupt Disable
func (c *Cpu) CLI() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#CLI
	c.Flags.Interrupt = false
	return 0
}

// CLV - Clear Overflow Flag
func (c *Cpu) CLV() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#CLV
	c.Flags.Overflow = false
	return 0
}

// CMP - Compare
func (c *Cpu) CMP() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#CMP
	c.Flags.Carry = c.Accumulator >= c.M
	c.Flags.Zero = c.Accumulator == c.M
	c.Flags.Negative = (c.Accumulator-c.M)&0x80 > 0
	return 0
}

// CPX - Compare X Register
func (c *Cpu) CPX() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#CPX
	c.Flags.Carry = c.X >= c.M
	c.Flags.Zero = c.X == c.M
	c.Flags.Negative = (c.X-c.M)&0x80 > 0
	return 0
}

// CPY - Compare Y Register
func (c *Cpu) CPY() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#CPY
	c.Flags.Carry = c.Y >= c.M
	c.Flags.Zero = c.Y == c.M
	c.Flags.Negative = (c.Y-c.M)&0x80 > 0
	return 0
}

// DEC - Decrement Memory
func (c *Cpu) DEC() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#DEC
	c.M--
	c.setNegativeM7()
	return 0
}

// DEX - Decrement X Register
func (c *Cpu) DEX() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#DEX
	c.X--
	c.Flags.Negative = c.X&0x80 > 0
	return 0
}

// DEY - Decrement Y Register
func (c *Cpu) DEY() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#DEY
	c.Y--
	c.Flags.Negative = c.Y&0x80 > 0
	return 0
}

// EOR - Exclusive OR
func (c *Cpu) EOR() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#EOR
	c.Accumulator ^= c.M
	c.setZero()
	c.Flags.Negative = c.Accumulator&0x80 > 0
	return 0
}

// INC - Increment Memory
func (c *Cpu) INC() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#INC
	c.M++
	c.setNegativeM7()
	return 0
}

// INX - Increment X Register
func (c *Cpu) INX() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#INX
	c.X++
	c.Flags.Negative = c.X&0x80 > 0
	return 0
}

// INY - Increment Y Register
func (c *Cpu) INY() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#INY
	c.Y++
	c.Flags.Negative = c.Y&0x80 > 0
	return 0
}

// JMP - Jump
func (c *Cpu) JMP() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#JMP
	c.ProgramCounter = uint16(c.M) // TODO: zero page? or wait for 2nd byte?
	return 0
}

// JSR - Jump to Subroutine
func (c *Cpu) JSR() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#JSR
	// TODO: haven't touched the stack yet
	return 0
}

// LDA - Load Accumulator
func (c *Cpu) LDA() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#LDA
	c.Accumulator = c.M
	c.setZero()
	c.Flags.Negative = c.Accumulator&0x80 > 0
	return 0
}

// LDX - Load X Register
func (c *Cpu) LDX() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#LDX
	c.X = c.M
	c.setZero()
	c.Flags.Negative = c.X&0x80 > 0
	return 0
}

// LDY - Load Y Register
func (c *Cpu) LDY() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#LDY
	c.Y = c.M
	c.setZero()
	c.Flags.Negative = c.Y&0x80 > 0
	return 0
}

// LSR - Logical Shift Right
func (c *Cpu) LSR() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#LSR
	c.Flags.Carry = c.M&0x01 > 0 // old bit 0
	c.M >>= 2
	c.Flags.Zero = c.M == 0
	c.setNegativeM7()
	return 0
}

// NOP - No Operation
func (c *Cpu) NOP() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#NOP
	return 0
}

// ORA - Logical Inclusive OR
func (c *Cpu) ORA() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#ORA
	c.Accumulator |= c.M
	c.setZero()
	c.Flags.Negative = c.Accumulator&0x80 > 0
	return 0
}

// PHA - Push Accumulator
func (c *Cpu) PHA() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#PHA
	stackAddr := 0x0100 | uint16(c.Stack) // TODO: seems often repeated
	c.Write(stackAddr, c.Accumulator)
	return 0
}

// PHP - Push Processor Status
func (c *Cpu) PHP() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#PHP
	stackAddr := 0x0100 | uint16(c.Stack)

	var F byte
	for i, f := range []bool{
		c.Flags.Carry,
		c.Flags.Zero,
		c.Flags.Interrupt,
		c.Flags.Decimal,
		c.Flags.B,
		c.Flags.Unused,
		c.Flags.Overflow,
		c.Flags.Negative,
	} {
		if f {
			F |= 1 << i
		}
	}

	c.Write(stackAddr, F)
	return 0
}

// PLA - Pull Accumulator
func (c *Cpu) PLA() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#PLA
	stackAddr := 0x0100 | uint16(c.Stack)
	c.Accumulator = c.Read(stackAddr)
	return 0
}

// PLP - Pull Processor Status
func (c *Cpu) PLP() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#PLP
	stackAddr := 0x0100 | uint16(c.Stack)
	newFlags := c.Read(stackAddr)

	// for i, f := range []bool{
	// 	c.Flags.Carry,
	// 	c.Flags.Zero,
	// 	c.Flags.Interrupt,
	// 	c.Flags.Decimal,
	// 	c.Flags.B,
	// 	c.Flags.Unused,
	// 	c.Flags.Overflow,
	// 	c.Flags.Negative,
	// } {
	// 	f = true
	// 	if newFlags&(1<<i) > 0 {
	// 		f = true
	// 	}
	// }

	c.Flags.Carry = newFlags&(1<<0) > 0
	c.Flags.Zero = newFlags&(1<<1) > 0
	c.Flags.Interrupt = newFlags&(1<<2) > 0
	c.Flags.Decimal = newFlags&(1<<3) > 0
	c.Flags.B = newFlags&(1<<4) > 0
	c.Flags.Unused = newFlags&(1<<5) > 0
	c.Flags.Overflow = newFlags&(1<<6) > 0
	c.Flags.Negative = newFlags&(1<<7) > 0

	return 0
}

// ROL - Rotate Left
func (c *Cpu) ROL() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#ROL
	// similar to ASL
	c.Flags.Carry = c.M&0x80 > 0 // old bit 7
	c.M <<= 2

	if c.Flags.Carry {
		c.M |= 0x01
	}

	c.setZero()
	c.setNegativeM7()
	return 0
}

// ROR - Rotate Right
func (c *Cpu) ROR() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#ROR
	c.Flags.Carry = c.M&0x01 > 0 // old bit 0
	c.M >>= 2

	if c.Flags.Carry {
		c.M |= 0x80
	}

	c.setZero()
	c.setNegativeM7()
	return 0
}

// RTI - Return from Interrupt
func (c *Cpu) RTI() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#RTI
	// invoked at the end of an interrupt
	c.PLP()
	// TODO: something about the PC?
	return 0
}

// RTS - Return from Subroutine
func (c *Cpu) RTS() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#RTS
	stackAddr := 0x0100 | uint16(c.Stack)
	// The RTS instruction is used at the end of a subroutine to return to
	// the calling routine. It pulls the program counter (minus one) from
	// the stack. (so we correct it with +1?)
	c.ProgramCounter = uint16(c.Read(stackAddr)) + 1
	return 0
}

// SBC - Subtract with Carry
func (c *Cpu) SBC() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#SBC

	diff := c.Accumulator - c.M
	if diff > c.Accumulator {
		// C 	Carry Flag 	Set if overflow in bit 7
		// "Set if overflow in bit 7" -- bit 7 of what: A, M, or A+M? i
		// assume A+M for now
		c.Flags.Carry = true
	}

	c.Accumulator = diff
	if !c.Flags.Carry {
		c.Accumulator -= 1
	}

	c.setZero()
	// V 	Overflow Flag 	Set if sign bit is incorrect
	// c.Flags.Overflow = ???
	// N 	Negative Flag 	Set if bit 7 set
	c.Flags.Negative = c.Accumulator&0x80 > 0

	return 0
}

// SEC - Set Carry Flag
func (c *Cpu) SEC() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#SEC
	c.Flags.Carry = true
	return 0
}

// SED - Set Decimal Flag
func (c *Cpu) SED() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#SED
	c.Flags.Decimal = true
	return 0
}

// SEI - Set Interrupt Disable
func (c *Cpu) SEI() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#SEI
	c.Flags.Interrupt = true
	return 0
}

// STA - Store Accumulator
func (c *Cpu) STA() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#STA
	c.M = c.Accumulator
	return 0
}

// STX - Store X Register
func (c *Cpu) STX() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#STX
	c.M = c.X
	return 0
}

// STY - Store Y Register
func (c *Cpu) STY() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#STY
	c.M = c.Y
	return 0
}

// TAX - Transfer Accumulator to X
func (c *Cpu) TAX() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#TAX
	c.X = c.Accumulator
	c.Flags.Zero = c.X == 0
	c.Flags.Negative = c.X&0x80 > 0
	return 0
}

// TAY - Transfer Accumulator to Y
func (c *Cpu) TAY() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#TAY
	c.Y = c.Accumulator
	c.Flags.Zero = c.Y == 0
	c.Flags.Negative = c.Y&0x80 > 0
	return 0
}

// TSX - Transfer Stack Pointer to X
func (c *Cpu) TSX() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#TSX
	stackAddr := 0x0100 | uint16(c.Stack)
	c.X = c.Read(stackAddr)
	c.Flags.Zero = c.X == 0
	c.Flags.Negative = c.X&0x80 > 0
	return 0
}

// TXA - Transfer X to Accumulator
func (c *Cpu) TXA() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#TXA
	c.Accumulator = c.X
	c.Flags.Zero = c.X == 0
	c.Flags.Negative = c.Accumulator&0x80 > 0
	return 0
}

// TXS - Transfer X to Stack Pointer
func (c *Cpu) TXS() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#TXS
	stackAddr := 0x0100 | uint16(c.Stack)
	c.Write(stackAddr, c.X)
	return 0
}

// TYA - Transfer Y to Accumulator
func (c *Cpu) TYA() byte {
	// https://www.nesdev.org/obelisk-6502-guide/reference.html#TYA
	c.Accumulator = c.Y
	c.Flags.Zero = c.Y == 0
	c.Flags.Negative = c.Accumulator&0x80 > 0
	return 0
}
