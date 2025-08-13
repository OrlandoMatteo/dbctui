package ui

import (
	"dbctui/can"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResultType int

const PageSize = 20

const (
	SIGNAL ResultType = iota
	MESSAGE
)

type Model struct {
	input          textinput.Model
	resultsSignals []*can.Signal
	resultMessages []*can.Message
	allMessages    []*can.Message
	allSignals     []*can.Signal
	selected       struct{}
	cursor         int
	pageIndex      int
	resultType     ResultType
	width, height  int
}

var (
	borderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	activeCursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)

	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true).
		PaddingBottom(1)

	emptyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)
	footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		PaddingTop(1)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		PaddingTop(0).
		PaddingBottom(0)
)

func InitialModel(msgs []*can.Message, sigs []*can.Signal) Model {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()
	return Model{
		input:          ti,
		resultsSignals: []*can.Signal{},
		resultMessages: []*can.Message{},
		allMessages:    msgs,
		allSignals:     sigs,
		cursor:         0,
		resultType:     SIGNAL,
		pageIndex:      0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("DBCtui")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab": // switch result type
			if m.resultType == SIGNAL {
				m.resultType = MESSAGE
			} else {
				m.resultType = SIGNAL
			}
			m.pageIndex = 0
			m.cursor = 0

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.resultType == SIGNAL && m.cursor < len(m.resultsSignals)-1 {
				m.cursor++
			}
			if m.resultType == MESSAGE && m.cursor < len(m.resultMessages)-1 {
				m.cursor++
			}
		case "left":
			newIndex := m.pageIndex - 1
			m.pageIndex = max(newIndex, 0)
		case "right":
			newIndex := m.pageIndex + 1
			if m.resultType == MESSAGE {
				pages := len(m.resultMessages) % 20
				m.pageIndex = min(newIndex, pages)
			}
			if m.resultType == SIGNAL {
				pages := len(m.resultsSignals) % 20
				m.pageIndex = min(newIndex, pages)
			}
		default:
			m.pageIndex = 0
		}

	}

	// Update search input
	m.input, cmd = m.input.Update(msg)

	// Filter results
	query := strings.ToLower(m.input.Value())
	m.resultsSignals = nil
	m.resultMessages = nil
	if query != "" {

		for _, sig := range m.allSignals {
			if strings.Contains(strings.ToLower(sig.Name), query) {
				m.resultsSignals = append(m.resultsSignals, sig)
			}
		}
		for _, msg := range m.allMessages {
			if strings.Contains(strings.ToLower(msg.Name), query) {
				m.resultMessages = append(m.resultMessages, msg)
			}
		}

		// Clamp cursor
		if m.resultType == SIGNAL && m.cursor >= len(m.resultsSignals) {
			m.cursor = max(0, len(m.resultsSignals)-1)
		}
		if m.resultType == MESSAGE && m.cursor >= len(m.resultMessages) {
			m.cursor = max(0, len(m.resultMessages)-1)
		}
	}

	return m, cmd

}
func (m Model) renderList() string {
	var out []string

	// Header
	mode := "Signals"
	if m.resultType == MESSAGE {
		mode = "Messages"
	}
	out = append(out, headerStyle.Render(mode))

	switch m.resultType {
	case SIGNAL:
		if len(m.resultsSignals) == 0 {
			out = append(out, emptyStyle.Render("No results"))
		}
		pageStart := min(PageSize*m.pageIndex, len(m.resultsSignals))
		pageStop := min(pageStart+PageSize, len(m.resultsSignals))

		for i, sig := range m.resultsSignals[pageStart:pageStop] {
			line := fmt.Sprintf("%s", sig.Name)
			if i == m.cursor {
				line = activeCursorStyle.Render("> " + sig.Name)
			} else {
				line = "  " + sig.Name
			}
			out = append(out, line)
		}
	case MESSAGE:

		if len(m.resultMessages) == 0 {
			out = append(out, emptyStyle.Render("No results"))
		}
		pageStart := min(PageSize*m.pageIndex, len(m.resultMessages))
		pageStop := min(pageStart+PageSize, len(m.resultMessages))
		for i, msg := range m.resultMessages[pageStart:pageStop] {
			line := fmt.Sprintf("%s", msg.Name)
			if i == m.cursor {
				line = activeCursorStyle.Render("> " + msg.Name)
			} else {
				line = "  " + msg.Name
			}
			out = append(out, line)
		}
	}
	total := 0
	if m.resultType == SIGNAL {
		total = len(m.resultsSignals)
	} else {
		total = len(m.resultMessages)
	}
	numPages := 0
	if m.resultType == SIGNAL {
		numPages = len(m.resultsSignals) % PageSize
	} else {
		numPages = len(m.resultMessages) % PageSize
	}
	itemIndex := m.cursor + 1 + m.pageIndex*PageSize
	footer := footerStyle.Render(fmt.Sprintf("Item %d of %d, Page %d of %d", itemIndex, total, m.pageIndex+1, numPages))
	out = append(out, "", footer)

	return borderStyle.Width(m.width/2 - 2).Render(strings.Join(out, "\n"))

}

func (m Model) renderDetail() string {
	var out string
	switch m.resultType {
	case SIGNAL:
		if len(m.resultsSignals) > 0 {
			sig := m.resultsSignals[m.cursor]
			out = fmt.Sprintf(
				"Signal\n\nName: %s\nStartBit: %d\nLength: %d\nMsg: %s\n",
				sig.Name, sig.StartBit, sig.BitLength, sig.Label,
			)
		} else {
			out = emptyStyle.Render("No signal selected")
		}
	case MESSAGE:
		if len(m.resultMessages) > 0 {
			msg := m.resultMessages[m.cursor]
			out = fmt.Sprintf(
				"Message\n\nName: %s\nID: %#X\nSource: %d",
				msg.Name, msg.CanId, msg.Source,
			)
		} else {
			out = emptyStyle.Render("No message selected")
		}
	}
	return borderStyle.Width(m.width/2 - 2).Render(out)
}

func (m Model) View() string {
	top := m.input.View()

	// Create the split view
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderList(),
		m.renderDetail(),
	)
	help := helpStyle.Render("↑/↓ Move ← → ChangePage  Tab Switch type  Q Quit  Enter Select")

	return top + "\n\n" + content + "\n\n" + help
}

func padRight(s string, w int) string {
	if len(s) >= w {
		return s[:w]
	}
	return s + strings.Repeat(" ", w-len(s))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
