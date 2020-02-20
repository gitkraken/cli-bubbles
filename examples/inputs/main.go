package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/tea"
	"github.com/charmbracelet/teaparty/input"
	te "github.com/muesli/termenv"
)

var (
	color               = te.ColorProfile().Color
	focusedPrompt       = te.String("> ").Foreground(color("205")).String()
	blurredPrompt       = "> "
	focusedSubmitButton = "[ " + te.String("Submit").Foreground(color("205")).String() + " ]"
	blurredSubmitButton = "[ " + te.String("Submit").Foreground(color("240")).String() + " ]"
)

func main() {
	if err := tea.NewProgram(
		initialize,
		update,
		view,
		subscriptions,
	).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

type Model struct {
	index         int
	nameInput     input.Model
	nickNameInput input.Model
	emailInput    input.Model
	submitButton  string
}

func initialize() (tea.Model, tea.Cmd) {
	name := input.DefaultModel()
	name.Placeholder = "Name"
	name.Focus()
	name.Prompt = focusedPrompt

	nickName := input.DefaultModel()
	nickName.Placeholder = "Nickname"
	nickName.Prompt = blurredPrompt

	email := input.DefaultModel()
	email.Placeholder = "Email"
	email.Prompt = blurredPrompt

	return Model{0, name, nickName, email, blurredSubmitButton}, nil
}

func update(msg tea.Msg, model tea.Model) (tea.Model, tea.Cmd) {
	m, ok := model.(Model)
	if !ok {
		panic("could not perform assertion on model")
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		// Cycle between inputs
		case "tab":
			fallthrough
		case "shift+tab":
			fallthrough
		case "enter":
			fallthrough
		case "up":
			fallthrough
		case "down":
			inputs := []input.Model{
				m.nameInput,
				m.nickNameInput,
				m.emailInput,
			}

			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.index == len(inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.index--
			} else {
				m.index++
			}

			if m.index > len(inputs) {
				m.index = 0
			} else if m.index < 0 {
				m.index = len(inputs)
			}

			for i := 0; i <= len(inputs)-1; i++ {
				if i == m.index {
					inputs[i].Focus()
					inputs[i].Prompt = focusedPrompt
					continue
				}
				inputs[i].Blur()
				inputs[i].Prompt = blurredPrompt
			}

			m.nameInput = inputs[0]
			m.nickNameInput = inputs[1]
			m.emailInput = inputs[2]

			if m.index == len(inputs) {
				m.submitButton = focusedSubmitButton
			} else {
				m.submitButton = blurredSubmitButton
			}

			return m, nil

		default:
			m.nameInput, _ = input.Update(msg, m.nameInput)
			m.nickNameInput, _ = input.Update(msg, m.nickNameInput)
			m.emailInput, _ = input.Update(msg, m.emailInput)
			return m, nil
		}

	default:
		m.nameInput, _ = input.Update(msg, m.nameInput)
		m.nickNameInput, _ = input.Update(msg, m.nickNameInput)
		m.emailInput, _ = input.Update(msg, m.emailInput)
		return m, nil
	}
}

func subscriptions(model tea.Model) tea.Subs {
	return tea.Subs{
		"blink": func(model tea.Model) tea.Msg {
			m, _ := model.(Model)
			return input.Blink(m.nameInput)
		},
	}
}

func view(model tea.Model) string {
	m, ok := model.(Model)
	if !ok {
		return "[error] could not perform assertion on model"
	}

	s := "\n"

	inputs := []string{
		input.View(m.nameInput),
		input.View(m.nickNameInput),
		input.View(m.emailInput),
	}

	for i := 0; i < len(inputs); i++ {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}

	s += "\n\n" + m.submitButton + "\n"

	return s
}
