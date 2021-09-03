package components

type LogComponent struct {
	BaseComponent

	Log *[]string
}

func (ac *LogComponent) Update(dt float64) {
}

func (lc *LogComponent) AddRecord(record string) {
	*lc.Log = append(*lc.Log, record)
}
