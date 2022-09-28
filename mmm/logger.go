package mmm

import "fmt"

type logger struct {
	maxDepth int
	isOn     bool
}

const depthMark = "--"

var globalLogger *logger

func getLogger() *logger {
	return globalLogger
}

func createLogger(maxDepth int) *logger {
	globalLogger = &logger{
		maxDepth: maxDepth,
		isOn:     true,
	}
	return globalLogger
}

func (l *logger) turnLoggerOff() {
	l.isOn = false
}

func (l *logger) debugAlways(depth int, m ...any) {
	if !l.isOn {
		return
	}
	depthMarker := ""

	for i := 0; i < l.maxDepth-depth; i++ {
		depthMarker += depthMark
	}

	fmt.Println(l.maxDepth-depth, depthMarker, m)
}

func (l *logger) debug(depth int, m ...any) {
	if l.maxDepth-depth > 6 {
		return
	}
	l.debugAlways(depth, m...)
}
