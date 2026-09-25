package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const pageSize = 10

type pageLoaded struct {
	stations []Station
	offset   int
	playLast bool
}

type loadFailed struct{ err error }

type model struct {
	player   *Player
	stations []Station
	offset   int
	index    int
	loading  bool
	err      error
	width    int
	height   int
}

func newModel(player *Player) model {
	return model{player: player, loading: true}
}

func (m model) Init() tea.Cmd {
	return loadPage(0, false)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case pageLoaded:
		m.loading = false
		if len(msg.stations) == 0 {
			if len(m.stations) == 0 {
				m.err = fmt.Errorf("no j-rock stations")
			}
			return m, nil
		}
		m.stations = msg.stations
		m.offset = msg.offset
		m.err = nil
		i := 0
		if msg.playLast {
			i = len(m.stations) - 1
		}
		return m.play(i)
	case loadFailed:
		m.loading = false
		if len(m.stations) == 0 {
			m.err = msg.err
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.player.Stop()
			return m, tea.Quit
		}
		if m.loading || len(m.stations) == 0 {
			return m, nil
		}
		switch msg.String() {
		case "p":
			if err := m.player.Toggle(); err != nil {
				m.err = err
			}
			return m, nil
		case "up", "k":
			return m.prev()
		case "down", "j":
			return m.next()
		}
	}
	return m, nil
}

func (m model) View() string {
	status := m.status()
	body := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("jofii"),
		statusStyle.Render(status),
		hintStyle.Render("↑↓ station · p pause · q quit"),
	)
	box := panelStyle.Render(body)
	if m.width == 0 || m.height == 0 {
		return box
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func (m model) status() string {
	switch {
	case m.loading:
		return "loading..."
	case m.err != nil:
		return m.err.Error()
	case len(m.stations) == 0:
		return "no j-rock stations"
	}
	prefix := "playing"
	if m.player.Paused() {
		prefix = "paused"
	}
	return prefix + ": " + m.stations[m.index].Name
}

func (m model) play(i int) (tea.Model, tea.Cmd) {
	if i < 0 || i >= len(m.stations) {
		return m, nil
	}
	m.index = i
	if err := m.player.Play(m.stations[i].URL); err != nil {
		m.err = err
		return m, nil
	}
	m.err = nil
	return m, nil
}

func (m model) prev() (tea.Model, tea.Cmd) {
	if m.index > 0 {
		return m.play(m.index - 1)
	}
	if m.offset == 0 {
		return m, nil
	}
	m.loading = true
	return m, loadPage(m.offset-pageSize, true)
}

func (m model) next() (tea.Model, tea.Cmd) {
	if m.index < len(m.stations)-1 {
		return m.play(m.index + 1)
	}
	m.loading = true
	return m, loadPage(m.offset+pageSize, false)
}

func loadPage(offset int, playLast bool) tea.Cmd {
	return func() tea.Msg {
		stations, err := getStations(pageSize, offset)
		if err != nil {
			return loadFailed{err}
		}
		return pageLoaded{stations: stations, offset: offset, playLast: playLast}
	}
}

var (
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#e23d4a")).
			Padding(1, 2).
			Width(52)
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#e23d4a")).Bold(true)
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f4f0ea")).MarginTop(1)
	hintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#8a7f72")).MarginTop(1)
)

func main() {
	player := &Player{}
	defer player.Stop()

	p := tea.NewProgram(newModel(player), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "jofii: %v\n", err)
		os.Exit(1)
	}
}
