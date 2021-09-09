package components

type PlayerComponent struct {
	BaseComponent `bson:"-" json:"-"`

	TelegramID      int
	ChatID          int64
	CurrentLocation int
}

func (pc PlayerComponent) Update(dt float64) {

}

func (pc PlayerComponent) Save(dt float64) {

}
