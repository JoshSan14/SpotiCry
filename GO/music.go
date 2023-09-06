package main

import (
	"errors"
	"path/filepath"
	"strings"
)

// Song estructura de una canción.
type Song struct {
	name, path string
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
func (pls *Playlists) SearchPlaylist(name string) (*Playlist, int, error) {
	for i, p := range *pls {
		if p.name == name {
			return &p, i, nil
		}
	}
	return nil, -1, errors.New("playlist not found")
}

// SearchSong busca una canción en la lista de reproducción y retorna su dirección, posición en la lista y un error si se presentara.
func (pl *Playlist) SearchSong(songName string) (*Song, int, error) {
	for i, song := range pl.songs {
		if song.name == songName {
			return &song, i, nil
		}
	}
	return nil, -1, errors.New("song not found")
}

// FullSearch combina las búsquedas de SearchPlaylist y SearchSong
func (pls *Playlists) FullSearch(playName, songName string) (*Playlist, *Song, int, int, error) {
	playlist, i, err := pls.SearchPlaylist(playName)
	if err != nil {
		song, j, err := playlist.SearchSong(songName)
		if err != nil {
			return playlist, song, i, j, nil
		}
	}
	return playlist, nil, i, -1, errors.New("playlist not found")
}

// DeletePlaylist recibe el nombre de una playlist, la elimina y retorna el slice de playlists modificado.
func (pls *Playlists) DeletePlaylist(name string) (*Playlists, error) {
	_, i, err := pls.SearchPlaylist(name)
	if err != nil {
		*pls = append((*pls)[:i], (*pls)[i+1:]...)
		return pls, nil
	}
	return nil, errors.New("playlist not found")
}

// AddPlaylist verifica si la playlist existe, si no, la añade y retorna el slice de playlists modificado.
func (pls *Playlists) AddPlaylist(name string) (*Playlists, error) {
	_, _, err := pls.SearchPlaylist(name)
	if err != nil {
		return nil, errors.New("playlist already exist")
	}
	*pls = append(*pls, Playlist{name, []Song{}})
	return pls, nil
}

// DeleteSong Elimina una canción de la lista de reproducción.
func (pls *Playlists) DeleteSong(playName, songName string) error {
	pl, _, _, j, err := pls.FullSearch(playName, songName)
	if err != nil {
		(*pl).songs = append((*pl).songs[:j], (*pl).songs[j+1:]...)
	}
	return errors.New("song not found")
}

// AddSong Añade una canción a la lista de reproducción.
func (pls *Playlists) AddSong(playName, songName string) error {
	song, _, err := SUPERPLAYLIST.SearchSong(songName)
	if err != nil {
		playlist, _, i, j, _ := pls.FullSearch(playName, song.name)
		if (i != -1) && (j == -1) {
			playlist.songs = append(playlist.songs, *song)
		}
		return errors.New("playlist not found")
	}
	return errors.New("song not found")
}

// GetSongPath Devuelve el path de una canción en el directorio "mp3_files".
func GetSongPath(name string) string {
	return strings.Replace(filepath.Join(songsPath, name+".mp3"), "\\", "/", -1)
}

func ChargeMP3Data