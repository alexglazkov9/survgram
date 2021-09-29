package components

type PlayerComponent struct {
	BaseComponent `bson:"-" json:"-"`

	TelegramID int
	ChatID     int64
}
