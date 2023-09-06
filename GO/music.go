package main

import (
	"errors"
	"fmt"
	"github.com/bogem/id3v2"

	"github.com/hajimehoshi/go-mp3"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Song estructura de una canción.
type Song struct {
	path, title string
	year        int
	length      float32
}

// Playlist Estructura de una lista de reproducción (slice de canciones).
type Playlist struct {
	name  string
	songs []Song
}

var SUPERPLAYLIST = Playlist{name: "SUPERPLAYLIST", songs: []Song{}}

// Playlists Slice de Playlist
type Playlists []Playlist

// SearchPlaylist recibe el nombre de una playlist y retorna su dirección, posición en la lista y un error si se presentara.
func (pls *Playlists) SearchPlaylist(playName string) (*Playlist, int, error) {
	for i, p := range *pls {
		if p.name == playName {
			return &p, i, nil
		}
	}
	return nil, -1, errors.New("playlist not found")
}

// SearchSong busca una canción en la lista de reproducción y retorna su dirección, posición en la lista y un error si se presentara.
func (pl *Playlist) SearchSong(songTitle string) (*Song, int, error) {
	for i, song := range pl.songs {
		if song.title == songTitle {
			return &song, i, nil
		}
	}
	return nil, -1, errors.New("song not found")
}

// FullSearch combina las búsquedas de SearchPlaylist y SearchSong
func (pls *Playlists) FullSearch(playName, songTitle string) (*Playlist, *Song, int, int, error) {
	playlist, i, err := pls.SearchPlaylist(playName)
	if err != nil {
		song, j, err := playlist.SearchSong(songTitle)
		if err != nil {
			return playlist, song, i, j, nil
		}
	}
	return playlist, nil, i, -1, errors.New("playlist not found")
}

// DeletePlaylist recibe el nombre de una playlist, la elimina y retorna el slice de playlists modificado.
func (pls *Playlists) DeletePlaylist(playName string) (*Playlists, error) {
	_, i, err := pls.SearchPlaylist(playName)
	if err != nil {
		*pls = append((*pls)[:i], (*pls)[i+1:]...)
		return pls, nil
	}
	return nil, errors.New("playlist not found")
}

// AddPlaylist verifica si la playlist existe, si no, la añade y retorna el slice de playlists modificado.
func (pls *Playlists) AddPlaylist(playName string) (*Playlists, error) {
	_, _, err := pls.SearchPlaylist(playName)
	if err != nil {
		return nil, errors.New("playlist already exist")
	}
	*pls = append(*pls, Playlist{playName, []Song{}})
	return pls, nil
}

// DeleteSong Elimina una canción de la lista de reproducción.
func (pls *Playlists) DeleteSong(playName, songTitle string) error {
	pl, _, _, j, err := pls.FullSearch(playName, songTitle)
	if err != nil {
		(*pl).songs = append((*pl).songs[:j], (*pl).songs[j+1:]...)
	}
	return errors.New("song not found")
}

// AddSong Añade una canción a la lista de reproducción.
func (pls *Playlists) AddSong(playName, songTitle string) error {
	song, _, err := SUPERPLAYLIST.SearchSong(songTitle)
	if err != nil {
		playlist, _, i, j, _ := pls.FullSearch(playName, song.title)
		if (i != -1) && (j == -1) {
			playlist.songs = append(playlist.songs, *song)
		}
		return errors.New("playlist not found")
	}
	return errors.New("song not found")
}

// GetSongPath Devuelve el path de una canción en el directorio "mp3_files".
func GetSongPath(songTitle string) string {
	return strings.Replace(filepath.Join(songsPath, songTitle+".mp3"), "\\", "/", -1)
}

func processMP3File(path string, info os.FileInfo) (*Song, error) {
	var song Song

	// Abrir el archivo MP3 para lectura.
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	tag, err := id3v2.ParseReader(file, id3v2.Options{Parse: true})
	if err != nil {
		return nil, err
	}

	// Crear un decodificador de MP3
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return nil, err
	}

	// Llenar los datos de la canción con metadatos.
	song.path = path
	song.title = tag.Title()
	song.year, err = strconv.Atoi(tag.Year())
	// Calcular la duración de un archivo MP3 en segundos
	song.length = float32(decoder.Length()) / float32(decoder.SampleRate())

	return &song, nil
}

func getMP3Data(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".mp3") {
			metadata, err := processMP3File(path, info)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", path, err)
			} else {

			}
		}
	})
}
