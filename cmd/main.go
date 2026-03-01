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
	"github.com/UltimateForm/tcprcon-cli/internal/config"
	"github.com/UltimateForm/tcprcon/pkg/common_rcon"
	"github.com/UltimateForm/tcprcon/pkg/logger"
	"github.com/UltimateForm/tcprcon/pkg/packet"
	"github.com/UltimateForm/tcprcon/pkg/rcon"
	"golang.org/x/term"
)

var addressParam string
var portParam uint
var passwordParam string
var logLevelParam uint
var inputCmdParam string
var saveParam string
var profileParam string

func init() {
	flag.StringVar(&addressParam, "address", "localhost", "RCON address, excluding port")
	flag.UintVar(&portParam, "port", 7778, "RCON port")
	flag.StringVar(&passwordParam, "pw", "", "RCON password, if not provided will attempt to load from env variables, if unavailable will prompt")
	flag.UintVar(&logLevelParam, "log", logger.LevelWarning, "sets log level (syslog severity tiers) for execution")
	flag.StringVar(&inputCmdParam, "cmd", "", "command to execute, if provided will not enter into interactive mode")
	flag.StringVar(&saveParam, "save", "", "saves current connection parameters as a profile, value is the profile name")
	flag.StringVar(&profileParam, "profile", "", "loads a saved profile by name, overriding default flags but overridden by explicit flags")
}

func determinePassword(currentPw string) (string, error) {
	if len(currentPw) > 0 {
		logger.Debug.Println("using password from parameter or profile")
		return currentPw, nil
	}

	envPassword := os.Getenv("rcon_password")
	var password string

	if len(envPassword) > 0 {
		logger.Debug.Println("using password from os env")
		r := ""
		reader := bufio.NewReader(os.Stdin)
		for r == "" {
			fmt.Print("RCON password found in environment variables, use for authentication? [y/N] ")
			stdinbytes, _isPrefix, err := reader.ReadLine()
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
	configBasePath, err := os.UserConfigDir()
	if err != nil {
		logger.Critical.Fatal(err)
	}
	resolvedAddress, resolvedPort, resolvedPassword, err := config.Resolve(
		configBasePath,
		profileParam,
		addressParam,
		portParam,
		passwordParam,
	)
	if err != nil {
		logger.Critical.Fatal(err)
	}

	logger.Debug.Printf("resolved parameters: address=%v, port=%v, pw=%v, log=%v, cmd=%v\n", resolvedAddress, resolvedPort, resolvedPassword != "", logLevelParam, inputCmdParam)

	fullAddress := resolvedAddress + ":" + strconv.Itoa(int(resolvedPort))
	password, err := determinePassword(resolvedPassword)
	if err != nil {
		logger.Critical.Fatal(err)
	}

	// TODO: consider moving to config lib
	if saveParam != "" {
		cfg, loadErr := config.Load(configBasePath)
		if loadErr != nil {
			logger.Critical.Fatal(errors.Join(errors.New("failed to load config for saving"), loadErr))
		}

		newProfile := config.Profile{
			Address: resolvedAddress,
			Port:    resolvedPort,
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Save password to config file? (stored in %v with restricted file permissions) [y/N]: ", ansi.Format("plaintext", ansi.Red, ansi.Bold))
		savePwChoice, prefix, _ := reader.ReadLine()
		if prefix {
			logger.Err.Println("prefix handling not implemented")
		}
		if strings.ToLower(strings.TrimSpace(string(savePwChoice))) == "y" {
			newProfile.Password = password
		}

		cfg.SetProfile(saveParam, newProfile)
		if saveErr := cfg.Save(configBasePath); saveErr != nil {
			logger.Critical.Fatal(errors.Join(errors.New("failed to save config file"), saveErr))
		}
		configFilePath, pathErr := config.BuildConfigPath(configBasePath)
		if pathErr != nil {
			logger.Critical.Fatal(errors.Join(errors.New("failed to get config file path for display"), pathErr))
		}
		logger.Info.Printf("Profile '%s' saved successfully to %s\n", saveParam, configFilePath)
		// TODO: consider exiting here if user calls cli with "save" param, maybe he just setting profile and dont want to run the full thing idk
	}

	logger.Debug.Printf("Dialing %v at port %v\n", resolvedAddress, resolvedPort)
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
