// Package cpu implements the MOS Technology 6502 microprocessor, as used in
// the NES.

package cpu

import (
	"log"
	"time"

	"gone/mask"
	"gone/mem"
)

// https://www.nesdev.org/wiki/CPU#Frequencies
// https://www.nesdev.org/wiki/Cycle_reference_chart#Clock_rates

// CPU_TICK  = time.Nanosecond * 5590
// CPU_TICK2 = time.Duration(time.Second).Nanoseconds() / 1789773

var (
	tick = 10e9 / 1789773 // cannot be inlined into time.Duration, even with cast
	Tick = time.Nanosecond * time.Duration(tick)
)

// An AddressingMode tells the Cpu where to access (look for) a given byte of
// memory. There are 13 possible modes.
//
// Most Instructions can index the full 64 kB range of memory, that is, 256
// pages of 256 bytes. The exception is ZeroPage, which is confined to the
// first page of 256 bytes.
//
// Note: OLC uses dedicated Cpu methods (specifically, func () byte) for each
// of the 13 modes, but I think we might get away with an enum + switch block.
type AddressingMode int

// https://problemkaputt.de/everynes.htm#cpumemoryaddressing
// https://www.middle-engine.com/blog/posts/2020/06/23/programming-the-nes-the-6502-in-detail#addressing-modes
// https://www.nesdev.org/wiki/CPU_addressing_modes
// https://www.youtube.com/watch?v=TGcjn8zMhfM

const (
	// 0 increments

	Implied     AddressingMode = iota // does not increment ProgramCounter
	Accumulator                       // use Cpu.Accumulator

	// 1 increment, 1 (or 3) read

	Immediate // use the ProgramCounter itself
	ZeroPage  // 0x0000-0x00ff
	ZeroPageX
	ZeroPageY // LDX, STX
	IndirectX // rarely used

	IndirectY // 3 reads, may involve page crossing
	Relative  // 3 reads

	// 2 increments, 2 reads

	Absolute
	AbsoluteX // may involve page crossing
	AbsoluteY // may involve page crossing

	// 2 increments, 4 reads

	Indirect // JMP
)

// func checkByteAddr(b uint16) {
// 	if b < 0 || b > 0xffff {
// 		panic(1024)
// 	}
// }

// representing flags as an enum takes 8 bytes (when it should really just be
// one byte). it also necessitates additional getter/setter, when impls
// (instructions) could just check/modify the struct field(s) directly
//
// type Flag byte
//
// const (
// 	Carry Flag = iota
// 	Zero
// 	Interrupt
// 	Decimal
// 	B
// 	Unused
// 	Overflow
// 	Negative
// )

// The Cpu has no memory of its own (aside from a number of small registers
// which amount to about 7 bytes). Instead, the Cpu interfaces with a Bus that
// provides memory.
type Cpu struct {
	Bus *mem.Bus

	// Flags Flags

	// https://problemkaputt.de/everynes.htm#cpuregistersandflags
	// https://www.nesdev.org/wiki/CPU_ALL#CPU_2
	// https://www.nesdev.org/wiki/Status_flags#Flags

	// Flags are 8 bits that make up the status register (aka P register).
	// B and Decimal are unused.
	//
	// 7654 3210
	// NV1B DIZC
	Flags struct {
		Negative  bool // bit 7; only if signed ints are used
		Overflow  bool // bit 6; only if signed ints are used
		Unused    bool // bit 5
		Interrupt bool // bit 2; if true, disables interrupts
		Zero      bool // bit 1
		Carry     bool // bit 0
		B         bool // bit 4; unused
		Decimal   bool // bit 3; inherited from 6502, but unused by NES
		// note: if numeric indexing is required, switch to `Flags byte`
	}

	// Status         byte // equivalent to Flags, compacted in a single byte

	Accumulator byte // The Accumulator represents a byte value for immediate use, similar to a local variable
	X           byte
	Y           byte

	// Stack instructions (PHA, PLA, PHP, PLP, JSR, RTS, BRK, RTI) always
	// access the 01 page (0x0100-0x01ff). The Cpu can store a low byte in
	// this register.
	Stack byte

	// The ProgramCounter is a 2-byte (word) memory address that increments
	// (almost) continuously. The byte located at this address should
	// provide the CPU with an Opcode that specifies the next instruction
	// to execute.
	ProgramCounter uint16
	// https://en.wikipedia.org/wiki/Program_counter
	// https://www.youtube.com/watch?v=Z5JC9Ve1sfI
	// when is PC ever decremented/reset?

	M           byte // after AddressingMode
	AbsAddress  uint16
	RelAddress  int8 // relative to current PC, used exclusively in brancing instructions
	PageCrossed bool
	Cycles      byte // decrements to 0, at which point a new instruction is executed
	// Opcode     Opcode // current opcode (not really necessary? maybe for interrupt purposes)
}

// Read reads one byte from the given addr. The addr is typically supplied by
// the program.
func (c *Cpu) Read(addr uint16) byte {
	// note: we usually return byte, but Cpu typically has to cast
	// ('concats') bytes into uint16 to form mem addresses
	return c.Bus.Read(addr, true)
}

// Write passes data to the Bus, which actually performs the write.
func (c *Cpu) Write(
	addr uint16, // addresses are 2 bytes (16 bits) wide; see xxd
	data byte,
) {
	c.Bus.Write(addr, data)
}

func (c *Cpu) fetch(b byte) Opcode {
	// if illegal byte, do nothing?
	oc, legal := Opcodes[b]
	if !legal {
		log.Fatalln("Illegal byte supplied:", b)
	}
	return oc
}

// decode fetches a byte of data from memory, accounting for the addressing
// mode. c.ProgramCounter is incremented zero to three times.
//
// The retrieved byte is stored in c.M, so that it can be used by the following
// Instruction.
func (c *Cpu) decode(a AddressingMode) { // {{{

	// https://www.ascii-code.com/

	// TODO: is a pc increment equivalent to a clock tick?

	switch a {

	// using a byte in a register directly is always faster than a memory
	// read (c.read). similarly, reading from the zero page is faster than
	// reading from distant pages.

	// 0 reads

	case Implied:
		// no byte to fetch
		return // 0

	case Accumulator:
		// the byte -is- the Accumulator
		c.M = c.Accumulator
		return

	// note: fogleman does not increment PC inside this func (lookup), OLC
	// does.

	case Immediate:
		// "data is supplied as part of the instruction -- the next byte"
		c.ProgramCounter++
		c.AbsAddress = c.ProgramCounter

	// 1 read

	case ZeroPage:
		c.AbsAddress = uint16(c.Read(c.ProgramCounter))
		c.ProgramCounter++
		c.AbsAddress &= 0x00ff // clear high byte (go to page 0), keep low byte

	case ZeroPageX:
		// think struct ptr + offset. c.X is probably set by a prior
		// instruction
		c.AbsAddress = uint16(c.Read(c.ProgramCounter) + c.X)
		c.ProgramCounter++
		c.AbsAddress &= 0x00ff

	case ZeroPageY:
		c.AbsAddress = uint16(c.Read(c.ProgramCounter) + c.Y)
		c.ProgramCounter++
		c.AbsAddress &= 0x00ff

	case Relative:
		// fetch a byte somewhere up to half a byte away from current
		// absolute address (in either direction)

		// Since a single signed byte value is used to indicate the
		// adjustment, it can only be in the range -128 to +127
		// a signed byte has leftmost bit of 1
		// 0x01 0b0000_0001  1   1
		// 0x00 0b0000_0000  0   0
		// 0xff 0b1111_1111 -1 255
		// 0xfe 0b1111_1110 -2 254
		//
		// https://www.simonv.fr/TypesConvert/?integers

		rel := c.Read(c.ProgramCounter)
		c.ProgramCounter++
		// this comparison checks the leftmost bit of rel. in concrete
		// terms, &0x80 returns 128 for all rel>=128 (in which case
		// move back a page), 0 otherwise (in which case we use rel as
		// is and move forward)

		// if rel&0x80 == 0 {
		// 	c.AbsAddress += uint16(rel) - 0x0100
		// } else {
		// 	c.AbsAddress += uint16(rel)
		// }

		// https://github.com/fogleman/nes/blob/3880f3400500b1ff2e89af4e12e90be46c73ae07/nes/cpu.go#L469
		c.AbsAddress += uint16(rel)
		if rel&0x80 == 0 {
			c.AbsAddress -= 0x0100
		}
		c.RelAddress = int8(rel)

	// 2 reads

	case Absolute:
		// read pc twice to get a 2-byte addr (1st col, then page),
		// then go to (read data from) that new addr

		// The 6502 is little endian, so the number at the 1st address
		// read becomes the low byte (column).
		// https://stackoverflow.com/a/77683792

		col := c.Read(c.ProgramCounter) // 0xff
		c.ProgramCounter++
		page := c.Read(c.ProgramCounter) // 0xff00
		c.ProgramCounter++
		c.AbsAddress = mask.Word(page, col)

	case AbsoluteX:
		col := c.Read(c.ProgramCounter)
		c.ProgramCounter++
		page := c.Read(c.ProgramCounter)
		c.ProgramCounter++
		c.AbsAddress = mask.Word(page, col)

		c.AbsAddress += uint16(c.X)
		if c.AbsAddress&0xff00 != uint16(page)<<8 {
			c.PageCrossed = true
		}

	case AbsoluteY:
		col := c.Read(c.ProgramCounter)
		c.ProgramCounter++
		page := c.Read(c.ProgramCounter)
		c.ProgramCounter++
		c.AbsAddress = mask.Word(page, col)

		c.AbsAddress += uint16(c.Y)
		if c.AbsAddress&0xff00 != uint16(page)<<8 {
			c.PageCrossed = true
		}

	// 3 reads

	case IndirectX:

		// only 1 pc increment, but 3 reads
		ptr := c.Read(c.ProgramCounter)
		c.ProgramCounter++
		// we only jump once, to somewhere in page 0. once there, we
		// read 2 adjacent bytes, with a given X offset, and concat
		// those 2 bytes into a new word (addr), which is where we go
		// to.

		// note: we first cast into uint16 to avoid byte overflow, and
		// discard the high byte of the results
		page := c.Read(uint16(ptr+c.X) & 0x00ff)
		col := c.Read(uint16(ptr+1+c.X) & 0x00ff) // no 0xxxff bug, apparently
		c.AbsAddress = mask.Word(page, col)

	case IndirectY:

		// only 1 pc increment, but 3 reads
		ptr := c.Read(c.ProgramCounter)
		c.ProgramCounter++

		// unlike IndirectX, the Y increment is applied -after- the
		// indirection, not before. this means that a page cross is
		// possible, and must be checked
		page := c.Read(uint16(ptr) & 0x00ff)
		col := c.Read(uint16(ptr+1) & 0x00ff)
		c.AbsAddress = mask.Word(page, col)

		c.AbsAddress += uint16(c.Y)
		if c.AbsAddress&0xff00 != uint16(page)<<8 {
			c.PageCrossed = true
		}

	// 4 reads

	case Indirect:

		// first, we get a 2-byte addr (row,col), then go to that addr,
		// similar to Absolute mode. however, unlike Absolute mode, we
		// don't stop there, because the 2 bytes we read are not data,
		// but a pointer to an address, which we must jump to in order
		// to get the actual data.
		//
		// as a result, 4 reads are performed in total. however, the pc
		// is still only incremented twice.

		ptrCol := c.Read(c.ProgramCounter)
		c.ProgramCounter++
		ptrPage := c.Read(c.ProgramCounter)
		ptr := mask.Word(ptrPage, ptrCol)
		c.ProgramCounter++

		// important: no pc increment or extra cycles required

		// now that we have the pointer, get the contents of the addr,
		// and its neighbour
		realCol := c.Read(ptr)

		var realPage byte
		if ptrCol == 0xff {
			// bug: while reading the bytes for the ptr, a page
			// cross may have occurred. if so, read from 1st byte
			// of the same page (0xYY00)
			// http://www.6502.org/tutorials/6502opcodes.html#JMP
			// https://atariwiki.org/wiki/Wiki.jsp?page=6502%20bugs
			realPage = c.Read(ptr & 0xff00)
		} else {
			// note: +1 since we need to fetch a new addr, which is
			// always 2 bytes! otherwise, we would be reading the
			// same byte twice, and the result would always be
			// something like 0xabab (assuming the value remains
			// unchanged), which is silly
			realPage = c.Read(ptr + 1)
		}

		c.AbsAddress = mask.Word(realPage, realCol)

	}

	c.M = c.Read(c.AbsAddress)
} // }}}

// func (c *Cpu) execute(i func() byte) byte {
// 	return i()
// }

// tick runs a single fetch/decode/execute cycle, setting c.Cycles to the
// appropriate number. The Cpu must 'wait' this number of cycles before the
// next tick call.
func (c *Cpu) tick() {
	// https://en.wikipedia.org/wiki/Instruction_cycle#Summary_of_stages

	// like OLC, this is not clock cycle accurate; we perform all the work
	// at once, and simply wait until the correct number of cycles has
	// elapsed. real hardware is slow and is always performing something
	// every cycle, thus requiring the full number of cycles for execution
	//
	// https://old.reddit.com/r/EmuDev/comments/pkgxws/what_cycles_really_are/hc3fqcf/

	b := c.Read(c.ProgramCounter)
	op := c.fetch(b)
	c.ProgramCounter++ // decoding the opcode always requires 1 cycle

	// x := c.ProgramCounter
	c.decode(op.AddressingMode)
	// elapsed := c.ProgramCounter - x // TODO: then what?
	// _ = elapsed

	// executing the opcode requires another ?-? cycles
	op.Instruction(c)
	// c.execute(op.Instruction)

	c.Cycles = op.Cycles
	if c.PageCrossed {
		c.Cycles++
		c.PageCrossed = false
	}
}

func (c *Cpu) loop() {
	for {
		if c.Cycles == 0 {
			c.tick()
		}
		time.Sleep(Tick)
		c.Cycles--

		// c.tick()
		// time.Sleep(Tick * time.Duration(c.Cycles))
	}
}

func (c *Cpu) irq()   {} // async interrupt (after curr instr; may be ignored)
func (c *Cpu) nmi()   {} // async interrupt (after curr instr; cannot be ignored)
func (c *Cpu) reset() {} // async interrupt
