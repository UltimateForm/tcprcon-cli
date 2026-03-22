package fullterm

func constructCmdLine(newByte byte, cmdLine []byte) ([]byte, bool) {
	isSubmission := false
	switch newByte {
	case 127, 8: // backspace, delete
		if len(cmdLine) > 0 {
			cmdLine = cmdLine[:len(cmdLine)-1]
		}
	case 13, 10: // enter
		isSubmission = true
	case 27: // escape
	default:
		cmdLine = append(cmdLine, newByte)
	}
	return cmdLine, isSubmission
}

func clamp(minBound int, n int, maxBound int) int {
	return max(minBound, min(maxBound, n))
}
