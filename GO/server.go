package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	LISTENER = "localhost:8081"
	BUFFER   = 10000
)

// Variables globales de uso general.
var (
	workingDirectory, _ = os.Getwd()
	projectDir          = filepath.Dir(workingDirectory)
	songsPath           = filepath.Join(projectDir, "mp3_files")
	metaPLS             = Playlists{}
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitCh     chan struct{}
	msgCh      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
		msgCh:      make(chan Message, BUFFER),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {

		}
	}(ln)
	s.ln = ln

	go s.acceptLoop()

	<-s.quitCh
	close(s.msgCh)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		fmt.Println("new connection to the server:", conn.RemoteAddr())

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	buf := make([]byte, BUFFER)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		// Obtener el mensaje enviado por el cliente.
		message := string(buf[:n])
		parts := strings.Split(message, "::")
		code := parts[0] //(PLAY_S,ADD_S,DEL_S,ADD_PL,DEL_PL)
		playlist := parts[1]
		name := parts[2]

		// Crear una instancia con la dirección remota del cliente.
		s.msgCh <- Message{
			from:    conn.RemoteAddr().String(),
			payload: []byte(message),
		}

		fmt.Println(parts)

		// Ejecutar acción dependiendo del contenido del mensaje.
		switch code {
		case "PLAY_S": // Envía los datos de la canción que se desea reproducir.
			sendMP3(conn, parts[2])
		case "ADD_S": // Añade una canción en la playlist deseada.

		default:
			_, err = conn.Write([]byte("Invalid Input"))
			if err != nil {
				return
			}
		}

		_, err = conn.Write([]byte("Message Received"))
		if err != nil {
			return
		}

	}
}

// sendMP3 Recibe la conexión del cliente y el nombre de la canción y envía los datos de la misma a través del servidor.
func sendMP3(conn net.Conn, name string) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	mp3File, err := os.Open(GetSongPath(name))
	if err != nil {
		fmt.Println("Error opening .mp3 file:", err)
		return
	}
	defer func(mp3File *os.File) {
		var err = mp3File.Close()
		if err != nil {
		}
	}(mp3File)

	_, err = io.Copy(conn, mp3File)
	if err != nil {
		fmt.Println("Error sending .mp3 file:", err)
		return
	}

	fmt.Printf(".mp3 file sent to %s\n", conn.RemoteAddr())
}

func main() {
	server := NewServer(LISTENER)
	if err := server.Start(); err != nil {
		log.Fatal("Error starting server:", err)
	}

	// The server will handle incoming connections and readLoop in the background.
	fmt.Printf("Server listening on %s\n", LISTENER)

	// Keep the main function running to allow the server to handle incoming connections.
	select {}
}
