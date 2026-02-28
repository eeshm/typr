package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"terminal-wpm/internal/content"
	"terminal-wpm/internal/engine"
)

type Config struct {
	Mode      string
	TimeLimit time.Duration
}

func Run(cfg Config) error {
	text, err := content.RandomText(cfg.Mode)
	if err != nil {
		return err
	}
	m := newModel(cfg, text)
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
	target    string
	session   *engine.Session
	now       time.Time
	width     int
	height    int
	done      bool
	timedOut  bool
	cancelled bool
	final     engine.Metrics
	err       error
}

func newModel(cfg Config, target string) model {
	return model{
		cfg:     cfg,
		target:  target,
		session: engine.NewSession(target, cfg.TimeLimit),
		now:     time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
		return m, nil
	case tickMsg:
		m.now = time.Time(typed)
		if !m.done && m.session.IsTimedOut(m.now) {
			m.timedOut = true
			m.done = true
			m.final = m.session.Snapshot(m.now, true, false)
		}
		if m.done {
			return m, nil
		}
		return m, tickCmd()
	case tea.KeyMsg:
		if m.done {
			switch typed.String() {
			case "enter", "q", "esc", "ctrl+c":
				return m, tea.Quit
			default:
				return m, nil
			}
		}

		m.now = time.Now()
		switch typed.String() {
		case "ctrl+c":
			m.cancelled = true
			m.done = true
			m.final = m.session.Snapshot(m.now, false, true)
			return m, nil
		case "backspace", "ctrl+h":
			m.session.Backspace()
		default:
			runes := typed.Runes
			if len(runes) == 1 {
				r := runes[0]
				if r >= 32 && r <= 126 {
					m.session.ApplyRune(r, m.now)
				}
			}
		}

		if m.session.IsCompleted() {
			m.done = true
			m.final = m.session.Snapshot(m.now, false, false)
			return m, nil
		}
		if m.cfg.TimeLimit > 0 && m.session.IsTimedOut(m.now) {
			m.timedOut = true
			m.done = true
			m.final = m.session.Snapshot(m.now, true, false)
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.err != nil {
		return errorStyle.Render(m.err.Error())
	}
	if m.done {
		return m.viewSummary()
	}
	return m.viewLive()
}
