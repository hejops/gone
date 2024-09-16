package mem

// A Bus is the central (global) object that connects multiple 'hardware'
// components together, enabling communication between them. Each Bus has an
// independent memory layout that begins at 0x0000.
//
// In the NES, there are 2 Buses. One has 64 kB, responsible for CPU, memory,
// audio and cartridge (0x0000-0xffff). The other has 8 (?) kB, responsible for
// graphics (0x2000-0x3fff?).
//
// One or more components (structs) can be connected to a Bus by means of a
// pointer; e.g. Cpu.Bus = &Bus{}.
type Bus struct {
	// no divisions/mirroring of memory yet; not meant to be used for now
	FakeRam [64 * 1024]byte // 64 kB (0xffff), zeroed on init
}

// CPU     MEM     APU     CART
//  |       |       |       |
//  |       |0000   |4000   |4020
//  |       |07ff   |4017   |ffff
//  |------------------------------------ BUS 1
//  |
// PPU     GFX     VRAM    PALETTE
//  |       |       |       |
//  |       |       |       |
//  |       |       |       |
//  |------------------------------------ BUS 2

func (b Bus) Write(
	addr uint16, // addresses are 2 bytes wide
	data byte,
) {
	b.FakeRam[addr] = data
}

func (b Bus) Read(addr uint16, readonly bool) byte { return b.FakeRam[addr] }

// func newBus() Bus {
// 	return Bus{}
// }
