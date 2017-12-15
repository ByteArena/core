package utils

import (
	"fmt"
	"strings"
	"time"
)

type Stopwatch struct {
	name          string
	running       map[string]time.Time
	completedKeys []string
	completed     []int64
}

func MakeStopwatch(name string) Stopwatch {
	return Stopwatch{
		name:          name,
		running:       make(map[string]time.Time),
		completedKeys: make([]string, 0),
		completed:     make([]int64, 0),
	}
}

func (w *Stopwatch) Start(key string) *Stopwatch {
	w.running[key] = time.Now()
	return w
}

func (w *Stopwatch) Stop(key string) int64 {
	w.completed = append(w.completed, time.Now().UnixNano()-w.running[key].UnixNano())
	w.completedKeys = append(w.completedKeys, key)
	delete(w.running, key)

	return w.completed[len(w.completed)-1]
}

func (w *Stopwatch) String() string {

	res := make([]string, len(w.running))
	res = append(res, "- "+w.name+" -----------------")

	for i, watch := range w.completed {
		res = append(res, fmt.Sprintf("* %s took %f ms", w.completedKeys[i], float64(watch)/1000000))
	}

	return strings.Join(res, "\n") + "\n"
}
