package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"terminal-wpm/internal/content"
	"terminal-wpm/internal/engine"
	"terminal-wpm/internal/history"
)

// phase tracks which screen the TUI is showing.
type phase int

const (
	phaseMenu   phase = iota // word-count selection
	phaseTyping              // active typing test
	phaseDone                // final results
)

// wordOption represents one selectable word-count choice.
type wordOption struct {
	label string
	count int
}

var wordOptions = []wordOption{
	{"30 words", 30},
	{"60 words", 60},
}

type Config struct {
	Mode      string
	TimeLimit time.Duration
	WordCount int
}

func Run(cfg Config) error {
	m := newModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	if resolved, ok := finalModel.(model); ok {
		return resolved.err
	}
	return nil
}

type tickMsg time.Time

type model struct {
	cfg       Config
	phase     phase
	menuIdx   int // currently highlighted menu option
	target    string
	session   *engine.Session
	now       time.Time
	width     int
	height    int
	timedOut  bool
	cancelled bool
	final     engine.Metrics
	history   []history.Record
	scrollY   int // scroll offset for results screen
	err       error
}

func newModel(cfg Config) model {
	return model{
		cfg:   cfg,
		phase: phaseMenu,
		now:   time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return nil // no tick needed during menu
}

func tickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// startTyping generates the text and transitions to the typing phase.
func (m *model) startTyping() tea.Cmd {
	chosen := wordOptions[m.menuIdx]
	m.cfg.WordCount = chosen.count

	text, err := content.RandomText(m.cfg.Mode, m.cfg.WordCount)
	if err != nil {
		m.err = err
		return nil
	}

	m.target = text
	m.session = engine.NewSession(text, m.cfg.TimeLimit)
	m.phase = phaseTyping
	m.now = time.Now()
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
		return m, nil

	case tickMsg:
		if m.phase != phaseTyping {
			return m, nil
		}
		m.now = time.Time(typed)
		if m.session.IsTimedOut(m.now) {
			m.timedOut = true
			m.phase = phaseDone
			m.final = m.session.Snapshot(m.now, true, false)
			m.saveHistory()
			return m, nil
		}
		return m, tickCmd()

	case tea.KeyMsg:
		switch m.phase {
		case phaseMenu:
			return m.updateMenu(typed)
		case phaseTyping:
			return m.updateTyping(typed)
		case phaseDone:
			return m.updateDone(typed)
		}
	}
	return m, nil
}

// --- menu phase input ---

func (m model) updateMenu(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.menuIdx > 0 {
			m.menuIdx--
		}
	case "down", "j":
		if m.menuIdx < len(wordOptions)-1 {
			m.menuIdx++
		}
	case "enter", " ":
		cmd := m.startTyping()
		return m, cmd
	}
	return m, nil
}

// --- typing phase input ---

func (m model) updateTyping(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.now = time.Now()
	switch key.String() {
	case "ctrl+c":
		m.cancelled = true
		m.phase = phaseDone
		m.final = m.session.Snapshot(m.now, false, true)
		m.saveHistory()
		return m, nil
	case "backspace", "ctrl+h":
		m.session.Backspace()
	default:
		runes := key.Runes
		if len(runes) == 1 {
			r := runes[0]
			if r >= 32 && r <= 126 {
				if !m.session.ApplyRune(r, m.now) {
					// Wrong key — emit terminal bell as error sound.
					fmt.Print("\a")
				}
			}
		}
	}

	if m.session.IsCompleted() {
		m.phase = phaseDone
		m.final = m.session.Snapshot(m.now, false, false)
		m.saveHistory()
		return m, nil
	}
	if m.cfg.TimeLimit > 0 && m.session.IsTimedOut(m.now) {
		m.timedOut = true
		m.phase = phaseDone
		m.final = m.session.Snapshot(m.now, true, false)
		m.saveHistory()
	}
	return m, nil
}

// --- done phase input ---

func (m model) updateDone(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "enter", "q", "esc", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.scrollY > 0 {
			m.scrollY--
		}
	case "down", "j":
		m.scrollY++
	}
	return m, nil
}

// saveHistory persists the current result and loads recent records for display.
func (m *model) saveHistory() {
	tier := performanceTier(m.final.WPM)
	rec := history.Record{
		Date:      time.Now(),
		Mode:      m.cfg.Mode,
		WordCount: m.cfg.WordCount,
		WPM:       m.final.WPM,
		RawWPM:    m.final.RawWPM,
		Accuracy:  m.final.Accuracy,
		Errors:    m.final.Errors,
		TimeTaken: m.final.TimeTaken.Seconds(),
		Completed: m.final.Completed,
		Tier:      tier,
	}
	_ = history.Save(rec) // best-effort; don't block on save errors
	m.history = history.Recent(5)
}

func (m model) View() string {
	if m.err != nil {
		errView := errorStyle.Render(m.err.Error()) + "\n" + hintStyle.Render("Press Ctrl+C to exit")
		if m.width > 0 && m.height > 0 {
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, errView)
		}
		return errView
	}
	switch m.phase {
	case phaseMenu:
		return m.viewMenu()
	case phaseDone:
		return m.viewSummary()
	default:
		return m.viewLive()
	}
}
