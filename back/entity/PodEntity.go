package entity

import "time"

type PodEntity struct {
	PodName   string    `json:"podName"`
	NameSpace int       `json:"nameSpace"`
	Ready     int       `json:"ready"`
	Status    int       `json:"status"`
	Restart   int       `json:"restart"`
	Time      time.Time `json:"time"`
}
