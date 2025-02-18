package mysql

import (
	"database/sql"
	"encoding/json"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"hls-streamer/app/models"
	"log"
	"sync"
)

type Storage struct {
	db *sql.DB
}

var (
	mysqlConnectionString = flag.String("SQLConnPassword", "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4,utf8", "DB connection")
	storage               *Storage
	once                  sync.Once
)

func initMySQLConnection() {
	dbConn, err := sql.Open("mysql", *mysqlConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	dbConn.SetMaxIdleConns(0)

	storage = &Storage{
		db: dbConn,
	}
}

func GetConnection() *Storage {
	once.Do(func() {
		initMySQLConnection()
	})

	return storage
}

func (s *Storage) GetVideoContent() ([]models.Video, error) {
	query := `
		SELECT vf.id, vf.formats
		FROM video_formats vf
		JOIN files_j_video_formats fjvf ON vf.id = fjvf.video_format_id
		JOIN files f ON fjvf.file_id = f.id
		WHERE f.is_stream = 1
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Video

	for rows.Next() {
		var id int
		var formatsStr string
		if err := rows.Scan(&id, &formatsStr); err != nil {
			return nil, err
		}

		var formats []*models.VideoFormat
		if err := json.Unmarshal([]byte(formatsStr), &formats); err != nil {
			return nil, err
		}

		result = append(result, models.Video{
			ID:       id,
			Formats:  formats,
			Position: id,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
