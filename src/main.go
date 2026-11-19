package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type profileStatusMsg = struct {
	active   string
	profiles []string
	err      error
}

type setResultMsg struct {
	err error
}

type model struct {
	active   string
	profiles []string
	cursor   int
	status   string
	err      error
}

func initialModel() model {
	return model{
		status: "Loading power profiles...",
	}
}

func fetchProfilesCmd() tea.Cmd {
	return func() tea.Msg {
		activeOut, err := exec.Command("powerprofilesctl", "get").Output()
		if err != nil {
			return profileStatusMsg{err: fmt.Errorf("get failed: %w", err)}
		}

		listOut, err := exec.Command("powerprofilesctl", "list").Output()
		if err != nil {
			return profileStatusMsg{err: fmt.Errorf("list failed: %w", err)}
		}

		var profiles []string
		for _, line := range strings.Split(string(listOut), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasSuffix(trimmed, ":") {
				name := strings.TrimSuffix(trimmed, ":")
				// powerprofilesctl list prefixes active profile with '*'
				name = strings.TrimSpace(strings.TrimPrefix(name, "*"))
				if name != "" {
					profiles = append(profiles, name)
				}
			}
		}

		return profileStatusMsg{
			active:   strings.TrimSpace(string(activeOut)),
			profiles: profiles,
		}
	}
}

func setProfileCmd(profile string) tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("powerprofilesctl", "set", profile).CombinedOutput()
		if err != nil {
			return setResultMsg{err: fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))}
		}
		return setResultMsg{err: nil}
	}
}

func (m model) Init() tea.Cmd {
	return fetchProfilesCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case profileStatusMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = "Error reading power profiles"
			return m, nil
		}
		m.err = nil
		m.active = msg.active
		m.profiles = msg.profiles
		m.status = "Ready"
		for i, p := range m.profiles {
			if p == m.active {
				m.cursor = i
				break
			}
		}
	case setResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = "Failed to switch profile"
			return m, nil
		}
		return m, fetchProfilesCmd()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		case "enter", "space":
			if len(m.profiles) > 0 {
				target := m.profiles[m.cursor]
				m.status = "Applying " + target + "..."
				return m, setProfileCmd(target)
			}
		case "r":
			m.status = "Refreshing..."
			return m, fetchProfilesCmd()
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	var b strings.Builder
	b.WriteString("Battertui Power Profile TUI\n\n")
	fmt.Fprintf(&b, "Status: %s\n\n", m.status)

	if m.err != nil {
		fmt.Fprintf(&b, "Error: %v\n\n", m.err)
	}

	for i, p := range m.profiles {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		activeMark := ""
		if p == m.active {
			activeMark = "[*]"
		} else {
			activeMark = "[ ]"
		}
		fmt.Fprintf(&b, "%s %s %s\n", cursor, activeMark, p)
	}

	b.WriteString("\n[up/down] move  [enter] apply  [r] refresh  [q] quit\n")

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running battertui: %v\n", err)
		os.Exit(1)
	}
}
