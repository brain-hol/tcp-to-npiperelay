package main

import (
	"net"
	"os"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "12321"
	CONN_TYPE = "tcp"
)

func handleRequest(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Debug().Msg("Closing connection")
	}()

	log.Debug().Msg("Handling connection")
	cmd := exec.Command(".\\npiperelay.exe", "-ei", "-s", "//./pipe/openssh-ssh-agent")
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msg("Something went wrong")
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Debug().Msg("Starting server")

	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting TCP to listen")
		os.Exit(1)
	}
	defer listener.Close()

	log.Debug().Msgf("Listening on %s:%s", CONN_HOST, CONN_PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to accept a connection")
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}
