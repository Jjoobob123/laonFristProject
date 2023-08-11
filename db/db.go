package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitializeDB 함수는 데이터베이스 연결을 초기화
func InitializeDB(dbHost string, dbUser, dbPassword string, dbName string, dbPort int) error {
	fmt.Println(dbHost, dbName, dbPassword, dbPort, dbUser)
	connectionStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	var err error
	db, err = sql.Open("postgres", connectionStr)

	if err != nil {
		log.Fatal("데이터베이스 초기화 에러 :", err)
	}
	// defer db.Close()

	// 데이터베이스 연결 확인 방법
	err = db.Ping()
	if err != nil {
		log.Fatal("데이터베이스 연결 확인 에러:", err)
	}

	// 테이블 생성 코드 생성 ( 테이블이 없을 경우에만 생성 될수 있도록 )
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS collection_data (
			id SERIAL PRIMARY KEY,
			accessseq INT,
			vehicletypestr VARCHAR(255),
			lane INT,
			directionstr VARCHAR(255),
			speed INT,
			createAtStr VARCHAR(255)
		);
		`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("테이블 생성 오류:", err)
	}
	fmt.Println("데이터베이스에 연결되었습니다.")
	return nil
}

// InsertCountingInfo 함수는 카운팅 정보를 데이터베이스에 삽입
func InsertCountingInfo(accessSeq int, vehicleTypeStr string, lane int, directionStr string, speed int, createAtStr string) error {

	_, err := db.Exec("INSERT INTO collection_data (accessseq, vehicletypestr, lane, directionstr, speed, createAtStr) VALUES ($1, $2, $3, $4, $5, $6)",
		accessSeq, vehicleTypeStr, lane, directionStr, speed, createAtStr)
	if err != nil {
		return fmt.Errorf("데이터 삽입시 오류: %v", err)
	}
	return nil

}
