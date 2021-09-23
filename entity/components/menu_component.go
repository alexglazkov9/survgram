package components

import (
	"github.com/alexglazkov9/survgram/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MenuComponent struct {
	BaseComponent `bson:"-" json:"-"`

	Menus Stack
}

type Menu struct {
	Msg         tgbotapi.Chattable
	MenuOptions map[string]func(*entity.Entity) interface{}
}

type Stack []Menu

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(kb Menu) {
	*s = append(*s, kb) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (Menu, bool) {
	if s.IsEmpty() {
		return Menu{}, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

// Return top element of stack. Return false if stack is empty.
func (s *Stack) Top() (Menu, bool) {
	if s.IsEmpty() {
		return Menu{}, false
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
