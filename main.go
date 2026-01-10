package main

import (
	"fmt"
	"os"

	"github.com/Br1an6/go-pedalboard/pkg/pedalboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	state          state
	inputDevices   []item
	outputDevices  []item
	effects        []item
	
	list           list.Model
	
	selectedInput  string
	selectedOutput string
	selectedEffect string
	
	stream         *pedalboard.AudioStream
	statusMsg      string
	err            error
}

type state int

const (
	selectInput state = iota
	selectOutput
	selectEffect
	playing
)

func initialModel() model {
	// Load devices
	inputs, err := pedalboard.GetInputDevices()
	if err != nil {
		fmt.Printf("Error getting inputs: %v\n", err)
		os.Exit(1)
	}
	outputs, err := pedalboard.GetOutputDevices()
	if err != nil {
		fmt.Printf("Error getting outputs: %v\n", err)
		os.Exit(1)
	}

	inputItems := make([]item, len(inputs))
	for i, v := range inputs {
		inputItems[i] = item{title: v, desc: "Input Device"}
	}
	// Add default option
	inputItems = append([]item{{title: "Default", desc: "System Default"}}, inputItems...)

	outputItems := make([]item, len(outputs))
	for i, v := range outputs {
		outputItems[i] = item{title: v, desc: "Output Device"}
	}
	outputItems = append([]item{{title: "Default", desc: "System Default"}}, outputItems...)

	effectItems := []item{
		{title: "Gain", desc: "Simple volume control"},
		{title: "Reverb", desc: "Room simulation"},
		{title: "Distortion", desc: "Waveshaping distortion"},
		{title: "Delay", desc: "Echo effect"},
		{title: "Chorus", desc: "Modulation effect"},
		{title: "Phaser", desc: "Phase shifting"},
	}

	l := list.New(convertItems(inputItems), list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Input Device"

	return model{
		state:         selectInput,
		inputDevices:  inputItems,
		outputDevices: outputItems,
	effects:       effectItems,
		list:          l,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			if m.stream != nil {
				m.stream.Stop()
				m.stream.Close()
			}
			return m, tea.Quit
		}

		if m.state == playing {
			if msg.String() == "q" || msg.String() == "esc" {
				if m.stream != nil {
					m.stream.Stop()
					m.stream.Close()
					m.stream = nil
				}
			return m, tea.Quit
			}
			return m, nil
		}
	
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	switch m.state {
	case selectInput:
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selectedInput = i.title
				if m.selectedInput == "Default" { m.selectedInput = "" }
				
				// Move to output
				m.state = selectOutput
				m.list.SetItems(convertItems(m.outputDevices))
				m.list.Title = "Select Output Device"
				m.list.ResetSelected()
			}
			return m, nil
		}
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case selectOutput:
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selectedOutput = i.title
				if m.selectedOutput == "Default" { m.selectedOutput = "" }

				// Move to effect
				m.state = selectEffect
				m.list.SetItems(convertItems(m.effects))
				m.list.Title = "Select Effect"
				m.list.ResetSelected()
			}
			return m, nil
		}
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case selectEffect:
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selectedEffect = i.title
				// Start Audio
				err := m.startAudio()
				if err != nil {
					m.err = err
					m.statusMsg = fmt.Sprintf("Error: %v", err)
				} else {
					m.statusMsg = fmt.Sprintf("Playing... Using %s -> %s with %s", 
						orDefault(m.selectedInput), orDefault(m.selectedOutput), m.selectedEffect)
				}
				m.state = playing
			}
			return m, nil
		}
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func convertItems(in []item) []list.Item {
	out := make([]list.Item, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

func orDefault(s string) string {
	if s == "" { return "Default" }
	return s
}

func (m *model) startAudio() error {
	proc, err := pedalboard.NewInternalProcessor(m.selectedEffect)
	if err != nil {
		return err
	}

	// Create stream with selected devices
	stream, err := pedalboard.NewAudioStreamWithDevices(proc, m.selectedInput, m.selectedOutput)
	if err != nil {
		return err
	}
	
	stream.Start()
	m.stream = stream
	return nil
}

func (m model) View() string {
	if m.state == playing {
		if m.err != nil {
			return docStyle.Render(fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err))
		}
		return docStyle.Render(fmt.Sprintf("%s\n\nPress q to stop and quit.", m.statusMsg))
	}
	return docStyle.Render(m.list.View())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
