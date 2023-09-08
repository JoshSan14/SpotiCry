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

const (
	SONGSPATH string = "/mp3_files"
)

// Song estructura de una canción.
type Song struct {
	path, title, artist, album, genre string
	year                              int
	length                            float32
}

// Playlist Estructura de una lista de reproducción (slice de canciones).
type Playlist struct {
	name  string
	songs []Song
}

// Playlists Slice de Playlist
type Playlists []*Playlist

var (
	SUPERPLAYLIST = Playlist{name: "SUPERPLAYLIST", songs: []Song{}}
	GENPLAYLISTS  = Playlists{&SUPERPLAYLIST}
)

// SearchPlaylist recibe el nombre de una playlist y retorna su dirección, posición en la lista y un error si se presentara.
func (pls *Playlists) SearchPlaylist(playName string) (*Playlist, int, error) {
	for i, p := range *pls {
		if p.name == playName {
			return p, i, nil
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
	return nil, -1, errors.New("song not found A")
}

// FullSearch combina las búsquedas de SearchPlaylist y SearchSong
func (pls *Playlists) FullSearch(playName, songTitle string) (*Playlist, *Song, int, int, error) {
	playlist, i, err := pls.SearchPlaylist(playName)
	if err != nil {
		return nil, nil, -1, -1, errors.New("playlist not found")
	}
	song, j, err := playlist.SearchSong(songTitle)
	if err != nil {
		return playlist, nil, i, -1, errors.New("song not found B")
	}
	return playlist, song, i, j, nil
}

// DeletePlaylist recibe el nombre de una playlist, la elimina y retorna el slice de playlists modificado.
func (pls *Playlists) DeletePlaylist(playName string) (*Playlists, error) {
	_, i, err := pls.SearchPlaylist(playName)
	if err != nil {
		return nil, errors.New("playlist not found")
	}
	*pls = append((*pls)[:i], (*pls)[i+1:]...)
	return pls, nil
}

// AddPlaylist verifica si la playlist existe, si no, la añade y retorna el slice de playlists modificado.
func (pls *Playlists) AddPlaylist(playName string) (*Playlists, error) {
	_, _, err := pls.SearchPlaylist(playName)
	if err != nil {
		*pls = append(*pls, &Playlist{playName, []Song{}})
		return pls, nil
	}
	return nil, errors.New("playlist already exist")
}

// DeleteSong Elimina una canción de la lista de reproducción.
func (pls *Playlists) DeleteSong(playName, songTitle string) (*Playlist, error) {
	pl, _, _, j, err := pls.FullSearch(playName, songTitle)
	if err != nil {
		return pl, errors.New("song not found B")
	}
	(*pl).songs = append((*pl).songs[:j], (*pl).songs[j+1:]...)
	return pl, nil
}

// AddSong Añade una canción a la lista de reproducción.
func (pls *Playlists) AddSong(playName, songTitle string) (*Playlist, error) {
	song, _, err := SUPERPLAYLIST.SearchSong(songTitle)
	if err != nil {
		return nil, errors.New("song not found")
	}
	playlist, _, i, j, _ := pls.FullSearch(playName, song.title)
	if (i != -1) && (j == -1) {
		fmt.Println("ACCOMPLIOSHJE")
		playlist.songs = append(playlist.songs, *song)
		return playlist, nil

	}
	return nil, errors.New("playlist not found")
}

// GetSongPath Devuelve el path de una canción en el directorio "mp3_files".
func GetSongPath(songTitle string) string {
	return strings.Replace(filepath.Join(SONGSPATH, songTitle+".mp3"), "\\", "/", -1)
}

// ProcessMP3File procesa un archivo MP3 dado y extrae sus metadatos, incluyendo la duración en segundos.
// Devuelve una estructura *Song que contiene los metadatos y un posible error si ocurre alguno.
func ProcessMP3File(path string) (*Song, error) {
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
	song.artist = tag.Artist()
	song.album = tag.Album()
	song.genre = tag.Genre()
	song.year, err = strconv.Atoi(tag.Year())

	// Calcular la duración de un archivo MP3 en segundos
	song.length = float32(decoder.Length()) / float32(decoder.SampleRate())

	return &song, nil
}

func GetMP3Data(dir string, playlist *Playlist) (*Playlist, error) {
	// Recorrer el directorio para encontrar todos los archivos MP3
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Verificar si el archivo es un MP3
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".mp3") {
			// Process the MP3 file and get its metadata
			song, err := ProcessMP3File(path)
			if err != nil {
				// Handle or log the error if processing the file fails.
				fmt.Printf("Error processing file %s: %v\n", path, err)
			} else {
				// Append the processed song to the list
				playlist.songs = append(playlist.songs, *song)
			}
		}
		return nil
	})

	if err != nil {
		return nil, errors.New("")
	}

	return playlist, nil

}

func main() {
	_, err := GetMP3Data("./mp3_files", &SUPERPLAYLIST)
	if err != nil {
		return
	}
	fmt.Println(SUPERPLAYLIST)

	// Add a playlist
	_, err = GENPLAYLISTS.AddPlaylist("MyPlaylist")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("TEST 1", GENPLAYLISTS[0].songs)

	// Add a song to MyPlaylist
	_, err = GENPLAYLISTS.AddSong("MyPlaylist", "In the End")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("TEST 2", GENPLAYLISTS[1].songs)

	// Search for a playlist
	playlist, _, err := GENPLAYLISTS.SearchPlaylist("MyPlaylist")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("TEST 3: Found playlist:", playlist.name)
	}

	// Search for a song in a playlist
	song, _, err := playlist.SearchSong("In the End")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Found song:", song.title)
	}

	// Delete a song from a playlist
	_, err = GENPLAYLISTS.DeleteSong("MyPlaylist", "In the End")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Deleted song from MyPlaylist")
	}
	fmt.Println(GENPLAYLISTS[1].songs)

	// Delete a playlist
	_, err = GENPLAYLISTS.DeletePlaylist("MyPlaylist")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Deleted MyPlaylist")
	}
	fmt.Println(GENPLAYLISTS)

}
