package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/UltimateForm/tcprcon-cli/internal/ansi"
	"github.com/UltimateForm/tcprcon/pkg/logger"
	"github.com/UltimateForm/tcprcon/pkg/packet"
	"github.com/UltimateForm/tcprcon/pkg/rcon"
)

func execInputCmd(ctx context.Context, rcon *rcon.Client) error {
	logger.Debug.Println("executing input command: " + inputCmdParam)
	id := rcon.Id()
	execPacket := packet.New(id, packet.SERVERDATA_EXECCOMMAND, []byte(inputCmdParam))
	fmt.Printf(
		"(%v): SND CMD %v\n",
		ansi.Format(strconv.Itoa(int(id)), ansi.Green, ansi.Bold),
		ansi.Format(inputCmdParam, ansi.Blue),
	)
	rcon.Write(execPacket.Serialize())
	// not totally effective since response reader has its own internally defined timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	responseChan := packet.CreateResponseChannel(rcon, timeoutCtx)
	for {
		select {
		case <-timeoutCtx.Done():
			return errors.New("timeout while waiting for response")
		case pkt := <-responseChan:
			if pkt.Id != id {
				logger.Warn.Printf("ignoring a packet with mismatched ID: %v vs %v; possibly overly chatty server\n", id, pkt.Id)
				continue
			}
			fmt.Printf(
				"(%v): RCV PKT %v\n%v\n",
				ansi.Format(strconv.Itoa(int(pkt.Id)), ansi.Green, ansi.Bold),
				ansi.Format(strconv.Itoa(int(pkt.Type)), ansi.Green, ansi.Bold),
				ansi.Format(strings.TrimRight(pkt.BodyStr(), "\n\r"), ansi.Green),
			)
			return nil
		}
	}
}
