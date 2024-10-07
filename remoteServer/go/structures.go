package swagger

type SongData struct {
	Group string
	Song  string
}

type RespSongs struct {
	SongWithDetail
	Err error
}

type SongWithDetail struct {
	SongData
	SongDetail
}

var ReqChan = make(chan SongData)
var RespChan = make(chan RespSongs)
