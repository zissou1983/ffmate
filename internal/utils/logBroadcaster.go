package utils

type LogBroadcaster struct {
	Callback func([]byte)
}

func (cw *LogBroadcaster) Write(p []byte) (n int, err error) {
	cw.Callback(p)
	return len(p), nil
}
