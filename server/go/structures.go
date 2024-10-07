package server

type SongData struct {
	Group string
	Song  string
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongWithDetail struct {
	SongData
	SongDetail
}

type SongWithDetailPag struct {
	Number int
	SongWithDetail
}

type SongText struct {
	Text string
	Err  error
}

type SongTextPag struct {
	Number int
	Text   string
}

var AddChan = make(chan SongWithDetail)
var DeleteChan = make(chan SongData)
var SongTextChan = make(chan SongData)
var GetSongsChan = make(chan SongWithDetail)
var ChangeChan = make(chan SongWithDetail)

var RespErrChan = make(chan error)
var RespSongTextChan = make(chan SongText)
var RespSongWithDetail = make(chan []SongWithDetail)
