package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"music_service/config"
	sw "music_service/swaggerAPI/go"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Broadcast(ctx context.Context, wg *sync.WaitGroup) {
	db := connectToDB()
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case song := <-sw.ReqChan:
			sw.RespChan <- findSong(song, db)
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func connectToDB() *gorm.DB {
	const op = "connectToDB"

	log := sw.SetupLogger()
	log = log.With(slog.String("op", op))

	pg := config.NewRemotePgConfig()
	cfg := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", pg.Host, pg.Username, pg.Password, pg.Database, pg.Port)

	db, err := gorm.Open(postgres.Open(cfg), &gorm.Config{})
	for i := 0; i < 5 && err != nil; i++ {
		time.Sleep(time.Second * 5)
		db, err = gorm.Open(postgres.Open(cfg), &gorm.Config{})
		if err != nil {
			log.Info("ConnectToDB: error to connect, please wait %v", err)
		}
	}
	if err != nil {
		log.Debug("ConnectToDB: error to connect %v", err)
	}
	db.Table("songs").AutoMigrate(&sw.SongWithDetail{})
	log.Info("The PSQL is ready")
	return db
}

func findSong(song sw.SongData, db *gorm.DB) sw.RespSongs {
	var resp []sw.SongWithDetail
	db.Table("songs").Where("\"group\" ilike ?", song.Group).Where("song ilike ?", song.Song).Find(&resp)
	if len(resp) == 1 {
		return sw.RespSongs{SongWithDetail: resp[0], Err: nil}
	} else {
		return sw.RespSongs{SongWithDetail: sw.SongWithDetail{}, Err: errors.New("the song doesn't exist")}
	}
}
