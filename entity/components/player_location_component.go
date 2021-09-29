package components

type PlayerLocationComponent struct {
	BaseComponent `bson:"-" json:"-"`

	CurrentLocation int
	Destination     *int    `bson:"-" json:"-"`
	TravelTime      float64 `bson:"-" json:"-"`
}
