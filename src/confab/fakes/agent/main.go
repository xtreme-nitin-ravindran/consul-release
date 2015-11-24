package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	// store information about this fake process into JSON
	signal.Ignore()

	// validate command line arguments
	// expect them to look like
	//   fake-thing agent -config-dir=/some/path/to/some/dir
	if len(os.Args[1:]) == 0 {
		log.Fatal("expecting command as first argment")
	}

	var configDir string
	flagSet := flag.NewFlagSet("", flag.ExitOnError)
	flagSet.StringVar(&configDir, "config-dir", "", "config directory")
	flagSet.Parse(os.Args[2:])
	if configDir == "" {
		log.Fatal("missing required config-dir flag")
	}

	ow := NewOutputWriter(filepath.Join(configDir, "fake-output.json"), os.Getpid(), os.Args[1:])

	// read input options provided to us by the test
	var inputOptions struct {
		Members       []string
		FailRPCServer bool
	}

	if optionsBytes, err := ioutil.ReadFile(filepath.Join(configDir, "options.json")); err == nil {
		json.Unmarshal(optionsBytes, &inputOptions)
	}

	tcpAddr := ""
	if !inputOptions.FailRPCServer {
		tcpAddr = "127.0.0.1:8400"
	}

	server := &Server{
		HTTPAddr:     "127.0.0.1:8500",
		TCPAddr:      tcpAddr,
		Members:      inputOptions.Members,
		OutputWriter: ow,
	}

	err := server.Serve()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err)
	}

	for {
		if server.DidLeave {
			err := server.Exit()
			if err != nil {
				log.Fatalf("Failed to close server: %s\n", err)
			}

			ow.Exit()
			os.Exit(0)
		}

		time.Sleep(100 * time.Millisecond)
	}
}