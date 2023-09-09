package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	LISTENER = "localhost:8081"
	BUFFER   = 10000
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr   string
	ln           net.Listener
	quitCh       chan struct{}
	msgCh        chan Message
	genPlaylists []*Playlists
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr:   listenAddr,
		quitCh:       make(chan struct{}),
		msgCh:        make(chan Message, BUFFER),
		genPlaylists: []*Playlists{},
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

		clientPlaylists := &Playlists{}
		s.genPlaylists = append(s.genPlaylists, clientPlaylists)

		go s.readLoop(conn, clientPlaylists)
	}
}

func (s *Server) readLoop(conn net.Conn, clientPlaylists *Playlists) {

	buf := make([]byte, BUFFER)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			return // Terminar la función si ocurre un error en la lectura
		}

		// Obtener el mensaje enviado por el cliente.
		message := string(buf[:n])
		parts := strings.Split(message, "::")
		code := parts[0] //(PLAY_S,ADD_S,DEL_S,ADD_P,DEL_P)
		playName := parts[1]
		songTitle := parts[2]
		extraCode := parts[3]

		// Crear una instancia con la dirección remota del cliente.
		s.msgCh <- Message{
			from:    conn.RemoteAddr().String(),
			payload: []byte(message),
		}

		// Ejecutar acción dependiendo del contenido del mensaje.
		switch code {

		// Ejecuta la búsqueda de una canción y envía sus datos.
		case "SRH_S":
			var songs = make([]*Song, 1)
			song, _, err := SUPERPLAYLIST.SearchSong(songTitle)
			if err != nil {
				fmt.Println(err)
			}
			songs = append(songs, song)
			SendSongData(conn, songs)

		// Ejecuta la búsqueda de una playlist y envía sus datos.
		case "SRH_P":
			playlist, _, err := clientPlaylists.SearchPlaylist(playName)
			if err != nil {
				fmt.Println(err)
			}
			SendSongData(conn, playlist.songs)

		// Envía los datos de todas las playlists del cliente.
		case "SHW_FP":
			clientPlaylists.SendFullPlaylistsData(conn)

		// Envía los bytes de un archivo MP3 por el servidor.
		case "PLAY_S":
			clientPlaylists.SendMP3(conn, playName, songTitle)
		// Añade una canción a una playlist.
		case "ADD_S":
			playlist, err := clientPlaylists.AddSong(playName, songTitle)
			if err != nil {
				fmt.Println(err)
			}
			SendSongData(conn, playlist.songs)
		// Elimina una canción de una playlist.
		case "DEL_S":
			playlist, err := clientPlaylists.DeleteSong(playName, songTitle)
			if err != nil {
				fmt.Println(err)
			}
			SendSongData(conn, playlist.songs)
		// Añade una nueva playlist al slice de playlists.
		case "ADD_P":
			_, err := clientPlaylists.AddPlaylist(playName)
			if err != nil {
				fmt.Println(err)
			}
			clientPlaylists.SendFullPlaylistsData(conn)
		// Elimina una playlist en el slice de playlists.
		case "DEL_P":
			_, err := clientPlaylists.DeletePlaylist(playName)
			if err != nil {
				fmt.Println(err)
			}
			clientPlaylists.SendFullPlaylistsData(conn)
		// Filtra la playlist por año.
		case "FLT_Y":
			year, _ := strconv.Atoi(extraCode)
			songs, err := clientPlaylists.FilterByYear(playName, year)
			if err != nil {
				fmt.Println(err)
			}
			SendSongData(conn, songs)
		// Filtra la playlist por duración.
		case "FLT_L":
			duration, _ := strconv.Atoi(extraCode)
			songs, err := clientPlaylists.FilterByDuration(playName, duration)
			if err != nil {
				fmt.Println(err)
			}
			SendSongData(conn, songs)
		// Filtra la playlist por album.
		case "FLT_A":
			songs, err := clientPlaylists.FilterByAlbum(playName, extraCode)
			if err != nil {
				fmt.Println(err)
			}
			SendSongData(conn, songs)
		// Caso base:
		default:
			_, err = conn.Write([]byte("Invalid Input"))
			if err != nil {
				fmt.Println(err)
			}
		}

	}
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
