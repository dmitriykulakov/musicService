package database

import (
	"context"
	"fmt"
	"log/slog"
	"music_service/config"
	sw "music_service/server/go"
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
		case song := <-sw.AddChan:
			err := addSong(&song, db)
			sw.RespErrChan <- err
		case song := <-sw.DeleteChan:
			err := deleteSong(&song, db)
			sw.RespErrChan <- err
		case song := <-sw.SongTextChan:
			resp := textSong(&song, db)
			sw.RespSongTextChan <- resp
		case filter := <-sw.GetSongsChan:
			resp := getSongs(&filter, db)
			sw.RespSongWithDetail <- resp
		case song := <-sw.ChangeChan:
			err := changeSongs(&song, db)
			sw.RespErrChan <- err
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func connectToDB() *gorm.DB {
	const op = "connectToDB"

	log := sw.SetupLogger()
	log = log.With(slog.String("op", op))
	pg := config.NewInternalPgConfig()
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

func addSong(song *sw.SongWithDetail, db *gorm.DB) error {
	var songTmp []sw.SongWithDetail
	db.Table("songs").Where("\"group\" ilike ?", song.Group).Where("song ilike ?", song.Song).Find(&songTmp)
	if len(songTmp) == 0 {
		result := db.Table("songs").Create(*song)
		return result.Error
	}
	err := fmt.Errorf("the song %s - %s is already exist", song.Group, song.Song)
	return err
}

func deleteSong(song *sw.SongData, db *gorm.DB) error {
	var songTmp []sw.SongWithDetail
	db.Table("songs").Where("\"group\" ilike ?", song.Group).Where("song ilike ?", song.Song).Find(&songTmp)
	if len(songTmp) == 1 {
		result := db.Table("songs").Where("\"group\" ilike ?", song.Group).Where("song ilike ?", song.Song).Delete(song)
		return result.Error
	}
	err := fmt.Errorf("the song %s - %s is not exist", song.Group, song.Song)
	return err
}

func textSong(song *sw.SongData, db *gorm.DB) sw.SongText {
	var text string
	db.Table("songs").Where("\"group\" ilike ?", song.Group).Where("song ilike ?", song.Song).Select("text").Find(&text)
	if len(text) > 0 {
		return sw.SongText{Text: text, Err: nil}
	}
	err := fmt.Errorf("the song %s - %s is not exist", song.Group, song.Song)
	return sw.SongText{Text: "", Err: err}
}

func getSongs(filter *sw.SongWithDetail, db *gorm.DB) []sw.SongWithDetail {
	var songTmp []sw.SongWithDetail
	db.Table("songs").
		Where("\"group\" ilike ?", "%"+filter.Group+"%").
		Where("song ilike ?", "%"+filter.Song+"%").
		Where("release_date ilike ?", "%"+filter.ReleaseDate+"%").
		Where("text ilike ?", "%"+filter.Text+"%").
		Where("link ilike ?", "%"+filter.Link+"%").
		Find(&songTmp)
	return songTmp
}

func changeSongs(song *sw.SongWithDetail, db *gorm.DB) error {
	var songTmp []sw.SongWithDetail
	db.Table("songs").Where("\"group\" ilike ?", song.Group).Where("song ilike ?", song.Song).Find(&songTmp)
	if len(songTmp) == 1 {
		if song.ReleaseDate != "" {
			db.Table("songs").Where("\"group\" ilike ?", song.Group).Where("song ilike ?", song.Song).Update("release_date", song.ReleaseDate)
		}
		if song.Text != "" {
			db.Table("songs").Where("\"group\" ilike ?", song.Group).Where("song ilike ?", song.Song).Update("text", song.Text)
		}
		if song.Link != "" {
			db.Table("songs").Where("\"group\" ilike ?", song.Group).Where("song ilike ?", song.Song).Update("link", song.Link)
		}
		return nil
	}
	err := fmt.Errorf("the song %s - %s is not exist", song.Group, song.Song)
	return err
}
