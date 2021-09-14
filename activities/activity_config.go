package activities

import (
	"log"
	"math/rand"
)

type ActivityType string

const (
	FISHING   ActivityType = "fishing"
	BATTLE    ActivityType = "battle"
	GATHERING ActivityType = "gathering"
)

type ActivityConfig struct {
	Type           ActivityType
	ActivityChance float64
	SpawnChances   []SpawnChance
}

type SpawnChance struct {
	Id     int
	Chance float64
}

func GetSpawneeId(spawn_chances []SpawnChance) int {
	cmltv := make([]float64, len(spawn_chances))
	for i, sc := range spawn_chances {
		if i == 0 {
			cmltv[i] = sc.Chance
			continue
		}
		cmltv[i] = cmltv[i-1] + sc.Chance
	}

	rnd_f := rand.Float64()
	log.Println(rnd_f)
	log.Println(cmltv)
	for i, v := range cmltv {
		log.Println(v)
		if rnd_f <= v {
			return spawn_chances[i].Id
		}
	}
	return -1
}
