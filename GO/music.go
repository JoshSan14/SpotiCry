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

// SearchSong busca una canción en la lista de reproducción y retorna su dirección, posición en la lista y un error si se presentara.
func (pls *Playlists) SearchSong(playName, songName string) (*Playlist, *Song, int, error) {
	pl, _, err := pls.SearchPlaylist(playName)
	if err != nil {
		for i, song := range pl.songs {
			if song.name == songName {
				return pl, &song, i, nil
			}
		}
		return pl, nil, -1, errors.New("song not found")
	}
	return nil, nil, -1, errors.New("playlist not found")
}

// DeleteSong Elimina una canción de la lista de reproducción.
func (pls *Playlists) DeleteSong(playName, songName string) (*Playlist, error) {
	pl, _, i, err := pls.SearchSong(playName, songName)
	if err != nil {
		pl.songs = append((pl.songs)[:i], (pl.songs)[i+1:]...)
		return pl, nil
	}
	return nil, errors.New("song not found")
}

// AddSong Añade una canción a la lista de reproducción.
func (pls *Playlists) AddSong(name, path string) {
	p.songs = append(p.songs, Song{name, path})
}

// GetSongPath Devuelve el path de una canción en el directorio "mp3_files".
func GetSongPath(name string) string {
	return strings.Replace(filepath.Join(songsPath, name+".mp3"), "\\", "/", -1)
}

// Buscar una playlist en la lista de playlists.
