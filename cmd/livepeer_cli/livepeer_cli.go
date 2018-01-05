package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "livepeer-cli"
	app.Usage = "interact with local Livepeer node"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "http",
			Usage: "local http port",
			Value: "8935",
		},
		cli.StringFlag{
			Name:  "rtmp",
			Usage: "local rtmp port",
			Value: "1935",
		},
		cli.StringFlag{
			Name:  "host",
			Usage: "host for the Livepeer node",
			Value: "localhost",
		},
		cli.IntFlag{
			Name:  "loglevel",
			Value: 4,
			Usage: "log level to emit to the screen",
		},
		cli.BoolFlag{
			Name:  "transcoder",
			Usage: "transcoder on off flag",
		},
	}
	app.Action = func(c *cli.Context) error {
		// Set up the logger to print everything and the random generator
		log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(c.Int("loglevel")), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))
		rand.Seed(time.Now().UnixNano())

		// Start the wizard and relinquish control
		w := &wizard{
			endpoint:   fmt.Sprintf("http://%v:%v/status", c.String("host"), c.String("http")),
			rtmpPort:   c.String("rtmp"),
			httpPort:   c.String("http"),
			host:       c.String("host"),
			transcoder: c.Bool("transcoder"),
			in:         bufio.NewReader(os.Stdin),
		}
		w.run()

		return nil
	}
	app.Run(os.Args)
}

type wizard struct {
	endpoint   string // Local livepeer node
	rtmpPort   string
	httpPort   string
	host       string
	transcoder bool
	in         *bufio.Reader // Wrapper around stdin to allow reading user input
}

func (w *wizard) run() {
	// Make sure there is a local node running
	_, err := http.Get(w.endpoint)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot find local node. Is your node running on http:%v and rtmp:%v?", w.httpPort, w.rtmpPort))
		return
	}

	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("| Welcome to livepeer-cli, your Livepeer command line tool  |")
	fmt.Println("|                                                           |")
	fmt.Println("| This tool lets you interact with a local Livepeer node    |")
	fmt.Println("| and participate in the Livepeer protocol without the	    |")
	fmt.Println("| hassle that it would normally entail.                     |")
	fmt.Println("|                                                           |")
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println()

	w.stats(w.transcoder)

	// Basics done, loop ad infinitum about what to do
	for {
		fmt.Println()
		fmt.Println("What would you like to do? (default = stats)")
		fmt.Println(" 1. Get node status")
		fmt.Println(" 2. Initialize round")
		fmt.Println(" 3. Bond")
		fmt.Println(" 4. Unbond")
		fmt.Println(" 5. Withdraw stake")
		fmt.Println(" 6. Withdraw fees")
		fmt.Println(" 7. Claim rewards and fees")
		fmt.Println(" 8. Get test LPT")
		fmt.Println(" 9. Get test ETH")
		fmt.Println(" 10. List registered transcoders")

		if w.transcoder {
			fmt.Println(" 11. Become a transcoder")
			fmt.Println(" 12. Set transcoder config")

			w.doCLIOpt(w.read(), true)
		} else {
			fmt.Println(" 11. Deposit ETH")
			fmt.Println(" 12. Withdraw ETH")
			fmt.Println(" 13. Broadcast video")
			fmt.Println(" 14. Stream video")
			fmt.Println(" 15. Set broadcast config")

			w.doCLIOpt(w.read(), false)
		}
	}
}

func (w *wizard) doCLIOpt(choice string, transcoder bool) {
	switch choice {
	case "1":
		w.stats(w.transcoder)
		return
	case "2":
		w.initializeRound()
		return
	case "3":
		w.bond()
		return
	case "4":
		w.unbond()
		return
	case "5":
		w.withdrawStake()
		return
	case "6":
		w.withdrawFees()
		return
	case "7":
		w.claimRewardsAndFees()
		return
	case "8":
		w.requestTokens()
		return
	case "9":
		fmt.Print("Go to eth-testnet.livepeer.org and use the faucet. (enter to continue)")
		w.read()
		return
	case "10":
		w.registeredTranscoderStats()
		return
	}

	if transcoder {
		switch choice {
		case "11":
			w.activateTranscoder()
			return
		case "12":
			w.setTranscoderConfig()
			return
		}
	} else {
		switch choice {
		case "11":
			w.deposit()
			return
		case "12":
			w.withdraw()
			return
		case "13":
			w.broadcast()
			return
		case "14":
			w.stream()
			return
		case "15":
			w.setBroadcastConfig()
			return
		}
	}

	log.Error("That's not something I can do")
}
