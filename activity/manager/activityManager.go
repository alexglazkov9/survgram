package activitymanager

import (
	"log"

	"github.com/alexglazkov9/survgram/activity"
)

type ActivityManager struct {
	Activities []*activity.Activity
}

func (a *ActivityManager) Update(dt float64) {
	for i, activity := range a.Activities {
		if activity.IsComplete {
			a.Activities = append(a.Activities[:i], a.Activities[i+1:]...)
			log.Printf("activities len: %d", len(a.Activities))
			continue
		}
		activity.Update(dt)
	}
}

func (a *ActivityManager) Add(ac *activity.Activity) {
	a.Activities = append(a.Activities, ac)
	log.Printf("activities len: %d", len(a.Activities))
}
