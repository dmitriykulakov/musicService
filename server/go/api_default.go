package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"music_service/config"
	"net/http"
	"strings"
	"time"
)

const (
	DoesntExist         = "the song doesn't exist"
	SongAdded           = "the song is added"
	SongDeleted         = "the song is deleted"
	SongFilterEmpty     = "there are no songs with the filter"
	SongFilterOK        = "song filter is OK"
	InternalServerError = "internalServerError"
	SongChanged         = "the song is changed"
)

func Add(w http.ResponseWriter, r *http.Request) {
	const op = "add"

	log := SetupLogger()
	log = log.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := http.Post("http://"+config.NewRemoteServerConfig().Address+"/info", "application/json", r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Debug(err.Error())
		return
	}
	if resp.StatusCode == http.StatusOK {
		var song SongWithDetail
		json.NewDecoder(resp.Body).Decode(&song)
		defer resp.Body.Close()

		log = log.With(slog.String("song", fmt.Sprintln(song.Group, song.Song)))

		AddChan <- song
		if err := <-RespErrChan; err != nil {
			log.Info(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Debug(err.Error())
				return
			}
		} else {
			err := json.NewEncoder(w).Encode(SongAdded)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Debug(err.Error())
				return
			}
			log.Info(SongAdded)
		}
	}
	if resp.StatusCode == http.StatusBadRequest {
		w.WriteHeader(http.StatusBadRequest)
		log.Info(DoesntExist)
		err := json.NewEncoder(w).Encode(DoesntExist)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Debug(err.Error())
		}
	}
	if resp.StatusCode == http.StatusInternalServerError {
		log.Debug(InternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	const op = "delete"

	log := SetupLogger()
	log = log.With(slog.String("op", op))

	var song SongData
	json.NewDecoder(r.Body).Decode(&song)
	defer r.Body.Close()

	log = log.With(slog.String("song", fmt.Sprintln(song.Group, song.Song)))

	DeleteChan <- song
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	select {
	case err := <-RespErrChan:
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info(err.Error())
			err := json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Debug(err.Error())
			}
		} else {
			w.WriteHeader(http.StatusOK)
			log.Info(SongDeleted)
			err := json.NewEncoder(w).Encode(SongDeleted)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Debug(err.Error())
			}
		}
	case <-time.After(time.Second * 5):
		w.WriteHeader(http.StatusInternalServerError)
		log.Debug(InternalServerError)
	}
}

func GetSongText(w http.ResponseWriter, r *http.Request) {
	const op = "get_song_text"

	log := SetupLogger()
	log = log.With(slog.String("op", op))

	var song SongData
	json.NewDecoder(r.Body).Decode(&song)
	defer r.Body.Close()
	log = log.With(slog.String("song", fmt.Sprintln(song.Group, song.Song)))
	SongTextChan <- song
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	select {
	case resp := <-RespSongTextChan:
		if resp.Err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(resp.Err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Debug(err.Error())
			}
		} else {
			w.WriteHeader(http.StatusOK)
			var textSlice []SongTextPag
			for i, found := 1, true; found; i++ {
				var before string
				var after string
				before, after, found = strings.Cut(resp.Text, "\n\n")
				textSlice = append(textSlice, SongTextPag{i, before})
				resp.Text = after
			}
			err := json.NewEncoder(w).Encode(textSlice)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Debug(err.Error())
			}
		}
	case <-time.After(time.Second * 5):
		w.WriteHeader(http.StatusInternalServerError)
		log.Debug(InternalServerError)
	}

}

func GetSongs(w http.ResponseWriter, r *http.Request) {
	const op = "get_songs"

	log := SetupLogger()
	log = log.With(slog.String("op", op))

	var song SongWithDetail
	json.NewDecoder(r.Body).Decode(&song)
	defer r.Body.Close()

	log = log.With(slog.String("song", fmt.Sprintln(song.Group, song.Song)))

	GetSongsChan <- song
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	select {
	case resp := <-RespSongWithDetail:
		if len(resp) == 0 {
			log.Info(SongFilterEmpty)
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(SongFilterEmpty)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Debug(err.Error())
			}
		} else {
			log.Info(SongFilterOK)
			w.WriteHeader(http.StatusOK)
			var response []SongWithDetailPag
			for i, songWithDetail := range resp {
				response = append(response, SongWithDetailPag{i + 1, songWithDetail})
			}
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Debug(err.Error())
			}
		}
	case <-time.After(time.Second * 5):
		w.WriteHeader(http.StatusInternalServerError)
		log.Debug(InternalServerError)
	}
}

func Change(w http.ResponseWriter, r *http.Request) {
	const op = "change"

	log := SetupLogger()
	log = log.With(slog.String("op", op))

	var song SongWithDetail
	json.NewDecoder(r.Body).Decode(&song)
	defer r.Body.Close()

	log = log.With(slog.String("song", fmt.Sprintln(song.Group, song.Song)))

	ChangeChan <- song
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	select {
	case err := <-RespErrChan:
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info(err.Error())
			err := json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Debug(err.Error())
			}
		} else {
			w.WriteHeader(http.StatusOK)
			log.Info(SongChanged)
			err := json.NewEncoder(w).Encode(SongChanged)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Debug(err.Error())
			}
		}
	case <-time.After(time.Second * 5):
		w.WriteHeader(http.StatusInternalServerError)
		log.Debug(InternalServerError)
	}
}
