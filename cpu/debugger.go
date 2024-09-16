package cpu

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

type model struct {
	cpu     *Cpu
	program []byte

	offset uint16 // only for drawing pageTable
	prevPC uint16
	error  error
}

const pages = 65536 / 16

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m model) Init() tea.Cmd {
	m.cpu.LoadProgram([]byte(m.program), m.offset)
	m.cpu.ProgramCounter = m.offset
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch s {
		case "q":
			return m, tea.Quit

		// case "j":
		// 	m.cpu.ProgramCounter++
		// 	return m, nil

		case " ", "j":
			// op := Opcodes[m.cpu.Bus.FakeRam[m.cpu.ProgramCounter]]
			// op.Instruction(m.cpu)
			m.prevPC = m.cpu.ProgramCounter
			err := m.cpu.tick()
			if err != nil {
				m.error = err
				return m, tea.Quit
			}

		}
	}
	return m, nil
}

// renderPage renders a single page as a line. The current PC is highlighted.
func (m model) renderPage(start uint16) string {
	if start%16 != 0 {
		panic("start must be a multiple of 16")
	}
	s := fmt.Sprintf("%04x | ", start)
	for i, b := range m.cpu.Bus.FakeRam[start : start+16] {
		if start+uint16(i) == m.cpu.ProgramCounter {
			s += fmt.Sprintf("[%02x] ", b)
		} else {
			s += fmt.Sprintf(" %02x  ", b)
		}
	}
	return s
}

func (m model) status() string {
	var flags string
	for _, flag := range []bool{
		m.cpu.Flags.Negative,
		m.cpu.Flags.Overflow,
		m.cpu.Flags.Unused,
		m.cpu.Flags.B,
		m.cpu.Flags.Decimal,
		m.cpu.Flags.DisableInterrupt,
		m.cpu.Flags.Zero,
		m.cpu.Flags.Carry,
	} {
		if flag {
			flags += "/ "
		} else {
			flags += "  "
		}
	}
	return fmt.Sprintf(`
PC: %x (%x)
 M: %x
 A: %x
 X: %x
 Y: %x
N V _ B D I Z C
`,
		m.cpu.ProgramCounter,
		m.prevPC,
		m.cpu.M,
		m.cpu.Accumulator,
		m.cpu.X,
		m.cpu.Y,
	) + flags
}

func (m model) pageTable() string {
	header := "page | "
	for b := range 16 {
		header += fmt.Sprintf("  %01x  ", b)
	}

	pages := []string{header}

	offsets := []int{
		0, 16, 32, 48, 64,
		int(m.offset),
		int(m.offset + 16*1),
		int(m.offset + 16*2),
		int(m.offset + 16*3),
		int(m.offset + 16*4),
	}
	for _, i := range offsets {
		pages = append(pages, m.renderPage(uint16(i)))
	}
	return strings.Join(pages, "\n")
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.pageTable(),
			m.status(),
		),
		"",
		// strconv.FormatInt(int64(m.cpu.ProgramCounter), 16),
		spew.Sdump(Opcodes[m.cpu.Bus.FakeRam[m.cpu.ProgramCounter]]),
	)
}

// Debug loads the program into memory at the given offset, then starts an
// interactive TUI.
func (c *Cpu) Debug(program []byte, offset uint16) {
	m, err := tea.NewProgram(model{
		cpu:     c,
		program: program,
		offset:  offset,
	}).Run()
	if err != nil {
		panic(err)
	}
	x := m.(model)
	if x.error != nil {
		fmt.Println("Error:", x.error)
	}
}
