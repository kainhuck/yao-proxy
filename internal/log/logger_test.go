package log

import "testing"

func TestLogger(t *testing.T) {
	l := NewLogger(false)
	l.Info(1, 2, 3)
	l.Debug(1, 2, 3)
	l.Warn(1, 2, 3)
	l.Error(1, 2, 3)
	l.Debugf("name: %s", "xiaor")
}
