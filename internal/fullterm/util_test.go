package fullterm

import "testing"

func TestConstructCmdLineBasic(t *testing.T) {
	currLine := []byte("hell")
	newByte := byte('o')
	newLine, isSubmission := constructCmdLine(newByte, currLine)
	if isSubmission {
		t.Error("did not expect submission")
	}
	if string(newLine) != "hello" {
		t.Errorf("expected 'hello' but got '%s'", newLine)
	}
}

func TestConstructCmdLineSubmission(t *testing.T) {
	currLine := []byte("hell")
	newByte := byte(13)
	newLine, isSubmission := constructCmdLine(newByte, currLine)
	if !isSubmission {
		t.Error("expected submission")
	}
	if string(newLine) != "hell" {
		t.Errorf("expected 'hell' but got '%s'", newLine)
	}
}

func TestConstructCmdLineBackspace(t *testing.T) {
	currLine := []byte("hell")
	newByte := byte(127)
	newLine, isSubmission := constructCmdLine(newByte, currLine)
	if isSubmission {
		t.Error("did not expect submission")
	}
	if string(newLine) != "hel" {
		t.Errorf("expected 'hel' but got '%s'", newLine)
	}
}

func TestConstructCmdLineBackspaceEmpty(t *testing.T) {
	currLine := []byte("")
	newByte := byte(127)
	newLine, isSubmission := constructCmdLine(newByte, currLine)
	if isSubmission {
		t.Error("did not expect submission")
	}
	if string(newLine) != "" {
		t.Errorf("expected empty string, got '%s'", newLine)
	}
}

func TestConstructCmdLineBackspaceByte8(t *testing.T) {
	currLine := []byte("hi")
	newByte := byte(8)
	newLine, isSubmission := constructCmdLine(newByte, currLine)
	if isSubmission {
		t.Error("did not expect submission")
	}
	if string(newLine) != "h" {
		t.Errorf("expected 'h', got '%s'", newLine)
	}
}

func TestConstructCmdLineSubmissionLF(t *testing.T) {
	currLine := []byte("test")
	newByte := byte(10)
	newLine, isSubmission := constructCmdLine(newByte, currLine)
	if !isSubmission {
		t.Error("expected submission on LF")
	}
	if string(newLine) != "test" {
		t.Errorf("expected 'test', got '%s'", newLine)
	}
}

func TestClamp(t *testing.T) {
	minBound := 4
	maxBound := 23
	n := 11
	r := clamp(minBound, n, maxBound)
	if r != n {
		t.Fatalf("expect r to equal %v but instead it is %v", n, r)
	}
}

func TestClampGreaterThanMax(t *testing.T) {
	minBound := 2
	maxBound := 5
	n := 12
	r := clamp(minBound, n, maxBound)
	if r != maxBound {
		t.Fatalf("expect r to equal %v but instead it is %v", maxBound, r)
	}
}

func TestClampLesserThanMin(t *testing.T) {
	minBound := 202
	maxBound := 582
	n := 105
	r := clamp(minBound, n, maxBound)
	if r != minBound {
		t.Fatalf("expect r to equal %v but instead it is %v", minBound, r)
	}
}
