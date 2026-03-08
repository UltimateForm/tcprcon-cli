package fullterm

type ansiState int

const (
	ansiStateIdle ansiState = iota
	ansiStateEscape
	ansiStateCSI
	ansiStateCSITerm
)

type stdinAnsi struct {
	buf   []byte
	state ansiState
}

func newStdinAnsi() stdinAnsi {
	return stdinAnsi{
		buf:   make([]byte, 0),
		state: ansiStateIdle,
	}
}

func (src *stdinAnsi) handle(b byte) ansiState {
	switch {
	case b == 27:
		src.buf = []byte{b}
		src.state = ansiStateEscape
	case b == 91 && src.state == ansiStateEscape:
		src.buf = append(src.buf, b)
		src.state = ansiStateCSI
	case src.state == ansiStateCSI:
		src.buf = append(src.buf, b)
		if b > 63 && b < 127 {
			defer func() {
				src.state = ansiStateIdle
			}()
			src.state = ansiStateCSITerm
		}
	}
	return src.state
}
