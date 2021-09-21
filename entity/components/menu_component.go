package components

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type MenuComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Menus Stack
}

type Menu struct {
	Msg         tgbotapi.Chattable
	MenuOptions map[string]interface{}
}

type Stack []tgbotapi.Chattable

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(kb tgbotapi.Chattable) {
	*s = append(*s, kb) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (tgbotapi.Chattable, bool) {
	if s.IsEmpty() {
		return tgbotapi.MessageConfig{}, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

// Return top element of stack. Return false if stack is empty.
func (s *Stack) Top() (tgbotapi.Chattable, bool) {
	if s.IsEmpty() {
		return tgbotapi.MessageConfig{}, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		return element, true
	}
}

// Return top element of stack. Return false if stack is empty.
func (s *Stack) Clear() {
	for !s.IsEmpty() {
		s.Pop()
	}
}
