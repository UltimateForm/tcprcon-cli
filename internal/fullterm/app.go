package fullterm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/UltimateForm/tcprcon-cli/internal/ansi"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

type app struct {
	DisplayChannel   chan string
	submissionChan   chan string
	stdinChannel     chan byte
	fd               int
	prevState        *term.State
	content          []string
	commandSignature string
	once             sync.Once
	ansiMachine      stdinAnsi
	history          [][]byte
	historyCursor    int
}

func (src *app) Write(bytes []byte) (int, error) {
	src.DisplayChannel <- string(bytes)
	return len(bytes), nil
}

func (src *app) ListenStdin(context context.Context) {
	// we are only listening to the stdin bytes here, to see how we handle conversion to human readable characters go to util.go
	for {
		select {
		case <-context.Done():
			return
		default:
			b := make([]byte, 1)
			_, err := os.Stdin.Read(b)
			if err != nil {
				return
			}
			src.stdinChannel <- b[0]
		}
	}
}

func (src *app) setHistoryTail(b []byte) {
	src.history[len(src.history)-1] = b
}

func (src *app) historyTail() []byte {
	return src.history[len(src.history)-1]
}

func (src *app) currentCmd() []byte {
	return src.history[src.historyCursor]
}

func (src *app) Submissions() <-chan string {
	return src.submissionChan
}

func visibleContent(content []string, height int) []string {
	currentRows := len(content)
	// ngl i forgot why we adding plus 1.. oh well
	startRow := max(currentRows-(height+1), 0)
	return content[startRow:]
}

func formatCommandEcho(cmd string) string {
	return ansi.Format("> "+cmd+"\n", ansi.Blue)
}

func (src *app) DrawContent(finalDraw bool) error {
	_, height, err := term.GetSize(src.fd)
	if err != nil {
		return err
	}
	if !finalDraw {
		fmt.Print(ansi.ClearScreen + ansi.CursorHome)
	}
	drawableRows := visibleContent(src.content, height)
	for i := range drawableRows {
		fmt.Print(drawableRows[i])
	}

	if finalDraw {
		return nil
	}
	ansi.MoveCursorTo(height, 0)
	fmt.Printf(ansi.Format("%v> ", ansi.Blue), src.commandSignature)
	fmt.Print(string(src.currentCmd()))
	return nil
}

func (src *app) traverseHistory(delta int) {
	src.historyCursor = clamp(0, src.historyCursor+delta, len(src.history)-1)
}

func (src *app) Run(context context.Context) error {

	// this could be an argument but i aint feeling yet
	src.fd = int(os.Stdin.Fd())
	if !term.IsTerminal(src.fd) {
		return errors.New("expected to run in terminal")
	}

	prevState, err := term.MakeRaw(src.fd)
	fmt.Print(ansi.EnterAltScreen)
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT)

	if err != nil {
		return err
	}
	src.prevState = prevState
	defer src.Close()

	currFlags, err := unix.IoctlGetTermios(src.fd, unix.TCGETS)
	if err != nil {
		return err
	}
	currFlags.Lflag |= unix.ISIG
	currFlags.Oflag |= unix.ONLCR | unix.OPOST
	// fyi there's a TCSETS as well that applies the setting differently
	if err := unix.IoctlSetTermios(src.fd, unix.TCSETSF, currFlags); err != nil {
		return err
	}

	go src.ListenStdin(context)
	for {
		select {
		case <-sigch:
			return nil
		case <-context.Done():
			return nil
		case newStdinInput := <-src.stdinChannel:
			ansiSeq, ansiState := src.ansiMachine.handle(newStdinInput)
			// src.content = append(src.content, fmt.Sprintf("ansi machine handling %v, at state %v\n", ansiSeq, ansiState))
			switch ansiState {
			case ansiStateIdle:
				// no ansi sequence ongoing so its just presentation bytes
				newCmd, isSubmission := constructCmdLine(newStdinInput, src.currentCmd())
				if isSubmission {
					src.content = append(src.content, formatCommandEcho(string(newCmd)))
					if len(src.historyTail()) > 0 {
						src.history = append(src.history, []byte{})
					}
					src.submissionChan <- string(newCmd)
				} else {
					src.setHistoryTail(newCmd)
				}
				// feel a bit awkward doing this every input stroke, we can get back to it later
				src.historyCursor = len(src.history) - 1
			case ansiStateCSITerm:
				switch string(ansiSeq) {
				case ansi.ArrowKeyUp:
					src.traverseHistory(-1)
				case ansi.ArrowKeyDown:
					src.traverseHistory(1)
				case ansi.PageUpKey:
				//tbd
				case ansi.PageDownKey:
				//tbd
				default:
					// src.content = append(src.content, fmt.Sprintf("unhandled csi %v, %v\n", strconv.Itoa(int(ansiState)), ansiSeq))
				}
			default:
				// src.content = append(src.content, fmt.Sprintf("unhandled state %v, %v\n", strconv.Itoa(int(ansiState)), ansiSeq))
			}

			if err := src.DrawContent(false); err != nil {
				return err
			}
		case newDisplayInput := <-src.DisplayChannel:
			src.content = append(src.content, newDisplayInput)
			if err := src.DrawContent(false); err != nil {
				return err
			}
		}
	}
}

func (src *app) Close() {
	src.once.Do(func() {
		// note: consider closing channels
		fmt.Print(ansi.ExitAltScreen)
		src.DrawContent(true)
		term.Restore(src.fd, src.prevState)
		fmt.Println()
	})
}

func CreateApp(commandSignature string) *app {
	// buffered, so we don't block on input
	displayChannel := make(chan string, 10)
	displayChannel <- ansi.Format("##########\n", ansi.Yellow, ansi.Bold)
	stdinChannel := make(chan byte)
	submissionChan := make(chan string, 10)
	return &app{
		DisplayChannel:   displayChannel,
		stdinChannel:     stdinChannel,
		submissionChan:   submissionChan,
		content:          make([]string, 0),
		commandSignature: commandSignature,
		ansiMachine:      newStdinAnsi(),
		history:          [][]byte{[]byte{}},
		historyCursor:    0,
	}
}
