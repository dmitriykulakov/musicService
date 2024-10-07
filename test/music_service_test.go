package testPack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	srv "music_service/server/go"
	"net/http"
	"testing"
)

func TestAdd(t *testing.T) {
	testAdd := []struct {
		input      srv.SongData
		want       string
		statusCode int
	}{
		{srv.SongData{Group: "Muse", Song: "Supermassive Black Hole"}, "the song is added", 200},
		{srv.SongData{Group: "Mue", Song: "Supermassive Black Hole"}, "the song doesn't exist", 400},
		{srv.SongData{Group: "Muser", Song: "Supermassive Black Hole"}, "the song doesn't exist", 400},
		{srv.SongData{Group: "Muse", Song: "Supermassive Hole"}, "the song doesn't exist", 400},
		{srv.SongData{Group: "mUSe", Song: "won't stand down"}, "the song is added", 200},
		{srv.SongData{Group: "", Song: ""}, "the song doesn't exist", 400},
		{srv.SongData{Group: "sTONE sOUR", Song: "Through Glass"}, "the song is added", 200},
		{srv.SongData{Group: "sTONE sOUR", Song: "Through Glass"}, "the song Stone Sour - Through Glass is already exist", 400},
	}
	for i, r := range testAdd {
		songJson, err := json.Marshal(r.input)
		if err != nil {
			log.Fatal(err)
		}
		req := bytes.NewReader(songJson)

		resp, err := http.Post("http://localhost:3340/add", "Apllication/json", req)
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode != r.statusCode {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %d, result %d", i+1, r.input, r.statusCode, resp.StatusCode)
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %d, result %d: OK", i+1, r.input, r.statusCode, resp.StatusCode)
		}
		var response string
		json.NewDecoder(resp.Body).Decode(&response)
		if response != r.want {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %s, result %s", i+1, r.input, r.want, response)
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %s, result %s: OK", i+1, r.input, r.want, response)
		}
	}
}

func TestSongText(t *testing.T) {
	testSongText := []struct {
		input      srv.SongData
		want       string
		statusCode int
	}{
		{srv.SongData{Group: "Mue", Song: "Supermassive Black Hole"}, "the song Mue - Supermassive Black Hole is not exist", 400},
		{srv.SongData{Group: "Muser", Song: "Supermassive Black Hole"}, "the song Muser - Supermassive Black Hole is not exist", 400},
		{srv.SongData{Group: "Muse", Song: "Supermassive Hole"}, "the song Muse - Supermassive Hole is not exist", 400},
		{srv.SongData{Group: "", Song: ""}, "the song  -  is not exist", 400},
		{srv.SongData{Group: "sTONE sOUR", Song: " Through Glass"}, "the song sTONE sOUR -  Through Glass is not exist", 400},
	}
	for i, r := range testSongText {
		songJson, err := json.Marshal(r.input)
		if err != nil {
			log.Fatal(err)
		}
		req := bytes.NewReader(songJson)

		resp, err := http.Post("http://localhost:3340/song", "Apllication/json", req)
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode != r.statusCode {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %d, result %d", i+1, r.input, r.statusCode, resp.StatusCode)
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %d, result %d: OK", i+1, r.input, r.statusCode, resp.StatusCode)
		}
		var response string
		json.NewDecoder(resp.Body).Decode(&response)
		if response != r.want {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %s, result %s", i+1, r.input, r.want, response)
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %s, result %s: OK", i+1, r.input, r.want, response)
		}
	}

	testSongText2 := []struct {
		input      srv.SongData
		want       int
		statusCode int
	}{
		{srv.SongData{Group: "Muse", Song: "Supermassive Black Hole"}, 9, 200},
		{srv.SongData{Group: "mUSe", Song: "won't sTANd down"}, 13, 200},
		{srv.SongData{Group: "sTONE sOUR", Song: "Through Glass"}, 12, 200}}
	for i, r := range testSongText2 {
		songJson, err := json.Marshal(r.input)
		if err != nil {
			log.Fatal(err)
		}
		req := bytes.NewReader(songJson)

		resp, err := http.Post("http://localhost:3340/song", "Apllication/json", req)
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode != r.statusCode {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %d, result %d", i+1, r.input, r.statusCode, resp.StatusCode)
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %d, result %d: OK", i+1, r.input, r.statusCode, resp.StatusCode)
		}
		var response []srv.SongTextPag
		json.NewDecoder(resp.Body).Decode(&response)
		if len(response) != r.want {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %d, result %d", i+1, r.input, r.want, len(response))
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %d, result %d: OK", i+1, r.input, r.want, len(response))
		}
	}
}

func TestFilter(t *testing.T) {
	testFilter := []struct {
		input      srv.SongWithDetail
		want       int
		statusCode int
	}{
		{srv.SongWithDetail{SongData: srv.SongData{Group: "Muse", Song: "Supermassive Black Hole"}, SongDetail: srv.SongDetail{}}, 1, 200},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "Mue", Song: "Supermassive Black Hole"}, SongDetail: srv.SongDetail{}}, 0, 400},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "Muse"}, SongDetail: srv.SongDetail{}}, 2, 200},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "US"}, SongDetail: srv.SongDetail{}}, 2, 200},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "us", Song: "Black"}, SongDetail: srv.SongDetail{ReleaseDate: "0", Text: "baby", Link: "you"}}, 1, 200},
		{srv.SongWithDetail{SongData: srv.SongData{}, SongDetail: srv.SongDetail{Link: "youtube"}}, 3, 200},
		{srv.SongWithDetail{SongData: srv.SongData{}, SongDetail: srv.SongDetail{ReleaseDate: "0"}}, 3, 200},
		{srv.SongWithDetail{SongData: srv.SongData{}, SongDetail: srv.SongDetail{ReleaseDate: "2025"}}, 0, 400},
		{srv.SongWithDetail{SongData: srv.SongData{}, SongDetail: srv.SongDetail{Text: "looking"}}, 1, 200},
	}
	for i, r := range testFilter {
		songJson, err := json.Marshal(r.input)
		if err != nil {
			log.Fatal(err)
		}
		req := bytes.NewReader(songJson)

		resp, err := http.Post("http://localhost:3340/data", "Apllication/json", req)
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode != r.statusCode {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %d, result %d", i+1, r.input, r.statusCode, resp.StatusCode)
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %d, result %d: OK", i+1, r.input, r.statusCode, resp.StatusCode)
		}
		var response []srv.SongWithDetail
		json.NewDecoder(resp.Body).Decode(&response)
		if len(response) != r.want {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %d, result %d", i+1, r.input, r.want, len(response))
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %d, result %d: OK", i+1, r.input, r.want, len(response))
		}
	}
}

func TestChange(t *testing.T) {
	testChange := []struct {
		input      srv.SongWithDetail
		want       string
		statusCode int
	}{
		{srv.SongWithDetail{SongData: srv.SongData{Group: "Muse", Song: "Supermassive Black Hole"}, SongDetail: srv.SongDetail{Text: "ok"}}, "the song is changed", 200},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "Muse", Song: "Supermassive Black Hole"}, SongDetail: srv.SongDetail{ReleaseDate: "today"}}, "the song is changed", 200},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "Muse", Song: "Supermassive Black Hole"}, SongDetail: srv.SongDetail{Link: "no link"}}, "the song is changed", 200},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "Mue", Song: "Supermassive Black Hole"}}, "the song Mue - Supermassive Black Hole is not exist", 400},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "Muser", Song: "Supermassive Black Hole"}}, "the song Muser - Supermassive Black Hole is not exist", 400},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "Muse", Song: "Supermassive Hole"}}, "the song Muse - Supermassive Hole is not exist", 400},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "mUSe", Song: "won't stand down"}, SongDetail: srv.SongDetail{ReleaseDate: "today", Text: "ok", Link: "don't know"}}, "the song is changed", 200},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "", Song: ""}}, "the song  -  is not exist", 400},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "sTONE sOUR", Song: "Through Glass"}}, "the song is changed", 200},
		{srv.SongWithDetail{SongData: srv.SongData{Group: "sTONE sOUR", Song: "Through Glass"}}, "the song is changed", 200},
	}
	for i, r := range testChange {
		songJson, err := json.Marshal(r.input)
		if err != nil {
			log.Fatal(err)
		}
		req := bytes.NewReader(songJson)

		resp, err := http.Post("http://localhost:3340/change", "Apllication/json", req)
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode != r.statusCode {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %d, result %d", i+1, r.input, r.statusCode, resp.StatusCode)
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %d, result %d: OK", i+1, r.input, r.statusCode, resp.StatusCode)
		}
		var response string
		json.NewDecoder(resp.Body).Decode(&response)
		if response != r.want {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %s, result %s", i+1, r.input, r.want, response)
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %s, result %s: OK", i+1, r.input, r.want, response)
		}
	}
}

func TestDelete(t *testing.T) {
	testDelete := []struct {
		input      srv.SongData
		want       string
		statusCode int
	}{
		{srv.SongData{Group: "Mue", Song: "Supermassive Black Hole"}, "the song Mue - Supermassive Black Hole is not exist", 400},
		{srv.SongData{Group: "Muser", Song: "Supermassive Black Hole"}, "the song Muser - Supermassive Black Hole is not exist", 400},
		{srv.SongData{Group: "Muse", Song: "Supermassive Hole"}, "the song Muse - Supermassive Hole is not exist", 400},
		{srv.SongData{Group: "Muse", Song: "Supermassive Black Hole"}, "the song is deleted", 200},
		{srv.SongData{Group: "mUSe", Song: "won't sTANd down"}, "the song is deleted", 200},
		{srv.SongData{Group: "", Song: ""}, "the song  -  is not exist", 400},
		{srv.SongData{Group: "sTONE sOUR", Song: " Through Glass"}, "the song sTONE sOUR -  Through Glass is not exist", 400},
		{srv.SongData{Group: "sTONE sOUR", Song: "Through Glass"}, "the song is deleted", 200},
	}
	for i, r := range testDelete {
		songJson, err := json.Marshal(r.input)
		if err != nil {
			log.Fatal(err)
		}
		req := bytes.NewReader(songJson)

		resp, err := http.Post("http://localhost:3340/delete", "Apllication/json", req)
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode != r.statusCode {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %d, result %d", i+1, r.input, r.statusCode, resp.StatusCode)
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %d, result %d: OK", i+1, r.input, r.statusCode, resp.StatusCode)
		}
		var response string
		json.NewDecoder(resp.Body).Decode(&response)
		if response != r.want {
			t.Errorf("Test %d: FAIL, the example \"%v\" is wrong, expected %s, result %s", i+1, r.input, r.want, response)
		} else {
			t.Logf("Test %d: OK, the example \"%v\", expected %s, result %s: OK", i+1, r.input, r.want, response)
		}
	}
}
