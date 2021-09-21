package components

type PlayerLocationComponent struct {
	BaseComponent `bson:"-" json:"-"`

	CurrentLocation int
	Destination     *int
	TravelTime      float64
}
