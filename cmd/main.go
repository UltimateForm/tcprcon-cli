package cmd

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/UltimateForm/tcprcon-cli/internal/ansi"
	"github.com/UltimateForm/tcprcon/pkg/logger"
	"github.com/UltimateForm/tcprcon/pkg/common_rcon"
	"github.com/UltimateForm/tcprcon/pkg/packet"
	"github.com/UltimateForm/tcprcon/pkg/rcon"
	"golang.org/x/term"
)

var addressParam string
var portParam uint
var passwordParam string
var logLevelParam uint
var inputCmdParam string

func init() {
	flag.StringVar(&addressParam, "address", "localhost", "RCON address, excluding port")
	flag.UintVar(&portParam, "port", 7778, "RCON port")
	flag.StringVar(&passwordParam, "pw", "", "RCON password, if not provided will attempt to load from env variables, if unavailable will prompt")
	flag.UintVar(&logLevelParam, "log", logger.LevelWarning, "sets log level (syslog serverity tiers) for execution")
	flag.StringVar(&inputCmdParam, "cmd", "", "command to execute, if provided will not enter into interactive mode")
}

func determinePassword() (string, error) {
	if len(passwordParam) > 0 {
		logger.Debug.Println("using password from parameter")
		return passwordParam, nil
	}
	envPassword := os.Getenv("rcon_password")
	var password string
	if len(envPassword) > 0 {
		logger.Debug.Println("using password from os env")
		r := ""
		for r == "" {
			fmt.Print("RCON password found in environment variables, use for authentication? (y/n) ")
			stdinread := bufio.NewReader(os.Stdin)
			stdinbytes, _isPrefix, err := stdinread.ReadLine()
			if err != nil {
				return "", err
			}
			if _isPrefix {
				logger.Err.Println("prefix not supported")
				continue
			}
			r = string(stdinbytes)
		}
		if strings.ToLower(r) == "y" {
			password = envPassword
		}
	}
	if len(password) == 0 {
		fmt.Print("RCON PASSWORD: ")
		stdinbytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}
		password = string(stdinbytes)
	}
	return password, nil
}

func execInputCmd(rcon *rcon.Client) error {
	logger.Debug.Println("executing input command: " + inputCmdParam)
	execPacket := packet.New(rcon.Id(), packet.SERVERDATA_EXECCOMMAND, []byte(inputCmdParam))
	fmt.Printf(
		"(%v): SND CMD %v\n",
		ansi.Format(strconv.Itoa(int(rcon.Id())), ansi.Green, ansi.Bold),
		ansi.Format(inputCmdParam, ansi.Blue),
	)
	rcon.Write(execPacket.Serialize())
	packetRes, err := packet.Read(rcon)
	if err != nil {
		return errors.Join(errors.New("error while reading from RCON client"), err)
	}
	fmt.Printf(
		"(%v): RCV PKT %v\n%v\n",
		ansi.Format(strconv.Itoa(int(rcon.Id())), ansi.Green, ansi.Bold),
		ansi.Format(strconv.Itoa(int(packetRes.Type)), ansi.Green, ansi.Bold),
		ansi.Format(strings.TrimRight(packetRes.BodyStr(), "\n\r"), ansi.Green),
	)
	return nil
}

func Execute() {
	flag.Parse()
	logLevel := uint8(logLevelParam)
	logger.Setup(logLevel)
	logger.Debug.Printf("parsed parameters: address=%v, port=%v, pw=%v, log=%v, cmd=%v\n", addressParam, portParam, passwordParam != "", logLevelParam, inputCmdParam)
	fullAddress := addressParam + ":" + strconv.Itoa(int(portParam))
	password, err := determinePassword()
	if err != nil {
		logger.Critical.Fatal(err)
	}
	logger.Debug.Printf("Dialing %v at port %v\n", addressParam, portParam)
	rconClient, err := rcon.New(fullAddress)
	if err != nil {
		logger.Critical.Fatal(err)
	}
	defer rconClient.Close()

	logger.Debug.Println("Building auth packet")
	auhSuccess, authErr := common_rcon.Authenticate(rconClient, password)
	if authErr != nil {
		logger.Critical.Println(errors.Join(errors.New("auth failure"), authErr))
		return
	}
	if !auhSuccess {
		logger.Critical.Println(errors.New("unknown auth error"))
		return
	}

	if inputCmdParam != "" {
		if err := execInputCmd(rconClient); err != nil {
			logger.Critical.Println(err)
		}
		return
	} else {
		// could just rely on early return but i feel anxious :D
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		runRconTerminal(rconClient, ctx, logLevel)
	}
}
