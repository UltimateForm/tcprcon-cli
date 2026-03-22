package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/UltimateForm/tcprcon-cli/internal/ansi"
	"github.com/UltimateForm/tcprcon-cli/internal/fullterm"
	"github.com/UltimateForm/tcprcon/pkg/logger"
	"github.com/UltimateForm/tcprcon/pkg/packet"
	"github.com/UltimateForm/tcprcon/pkg/rcon"
)

func stayAlive(ctx context.Context, client *rcon.Client, command string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			id := client.Id()
			logger.Debug.Printf("sending keepalive packet with id %v", id)
			pulsePacket := packet.New(client.Id(), packet.SERVERDATA_EXECCOMMAND, []byte(command))
			// note: potential race condition here, it is possible we send packet at the exact same time the user does, just fyi
			client.Write(pulsePacket.Serialize())
		case <-ctx.Done():
			return
		}
	}
}

func runRconTerminal(ctx context.Context, client *rcon.Client, logLevel uint8, profileName string, pulseCmd string, pulseInterval time.Duration) {
	signatureProfile := "rcon"
	if profileName != "" {
		signatureProfile = profileName
	}
	app := fullterm.CreateApp(fmt.Sprintf("%v@%v", signatureProfile, client.Address))
	// dont worry we are resetting the logger before returning
	logger.SetupCustomDestination(logLevel, app)

	appErrors := make(chan error, 1)
	connectionErrors := make(chan error, 1)

	appRun := func() {
		appErrors <- app.Run(ctx)
	}
	packetChannel := packet.CreateResponseChannel(client, ctx)
	packetReader := func() {
		for {
			select {
			case <-ctx.Done():
				return
			case streamedPacket := <-packetChannel:
				if streamedPacket.Error != nil {
					if errors.Is(streamedPacket.Error, os.ErrDeadlineExceeded) {
						logger.Debug.Println("read deadline reached; connection is idle or server is silent.")
						continue
					}
					if errors.Is(streamedPacket.Error, io.EOF) {
						connectionErrors <- io.EOF
						return
					}
					logger.Err.Println(errors.Join(errors.New("error while reading from RCON client"), streamedPacket.Error))
					continue
				}
				fmt.Fprintf(
					app,
					"(%v): RCV PKT %v\n%v\n",
					ansi.Format(strconv.Itoa(int(streamedPacket.Id)), ansi.Green, ansi.Bold),
					ansi.Format(strconv.Itoa(int(streamedPacket.Type)), ansi.Green, ansi.Bold),
					ansi.Format(strings.TrimRight(streamedPacket.BodyStr(), "\n\r"), ansi.Green),
				)
			}
		}
	}
	submissionChan := app.Submissions()
	submissionReader := func() {
		for {
			select {
			case <-ctx.Done():
				return
			case cmd := <-submissionChan:
				execPacket := packet.New(client.Id(), packet.SERVERDATA_EXECCOMMAND, []byte(cmd))
				fmt.Fprintf(
					app,
					"(%v): SND CMD %v\n",
					ansi.Format(strconv.Itoa(int(client.Id())), ansi.Green, ansi.Bold),
					ansi.Format(cmd, ansi.Blue),
				)
				client.Write(execPacket.Serialize())
			}
		}
	}

	go submissionReader()
	go packetReader()
	go appRun()
	if pulseCmd != "" {
		go stayAlive(ctx, client, pulseCmd, pulseInterval)
	}

	select {
	case <-ctx.Done():
		logger.Debug.Println("context done")
		break
	case err := <-connectionErrors:
		defer func() {
			logger.Critical.Println(errors.Join(errors.New("connection error"), err))
		}()
		break
	case err := <-appErrors:
		// lets do this because the app might be unrealiable at this point
		if err != nil {
			defer func() {
				logger.Critical.Println(errors.Join(errors.New("app error"), err))
			}()
		} else {
			defer func() {
				logger.Debug.Println("graceful app exit")
			}()
		}
		break
	}
	app.Close()
	logger.Setup(logLevel)
}
