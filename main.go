package main

import (
	"fmt"
	"laonproject0801/db"
	"laonproject0801/rest"
	"laonproject0801/tcp"
	"laonproject0801/udp"
	"laonproject0801/websocket"
	"log"

	"github.com/go-ini/ini"
)

func main() {
	// config.ini 파일로부터 설정 값을 불러옴
	configData, err := ini.Load("config.ini")
	if err != nil {

		log.Fatal("설정 값을 불러오는 중 에러 발생:", err)
	} else {
		log.Println("ini파일에서 설정 값 가져왔음!")
	}
	dbConfig := configData.Section("DB")
	networkConfig := configData.Section("Network")

	// db 키 읽어오기
	dbHost := dbConfig.Key("Host").String()
	dbUser := dbConfig.Key("User").String()
	dbPassword := dbConfig.Key("Password").String()
	dbName := dbConfig.Key("Database").String()
	dbPort := dbConfig.Key("Port").MustInt()

	// network 키 읽어오기
	method := networkConfig.Key("Method").MustInt()
	networkHost := networkConfig.Key("Host").String()
	networkPort := networkConfig.Key("Port").MustInt()

	// 데이터베이스 연결 초기화
	err = db.InitializeDB(dbHost, dbUser, dbPassword, dbName, dbPort)
	if err != nil {
		log.Fatal("데이터베이스 초기화 중 에러 발생:", err)
	}

	switch method {
	case 1:
		// TCP 시작 - 카운팅 정보를 생성하여 Consumer에게 보냄
		go tcp.StartTCPProducingCountingInfo(networkHost, networkPort)
		// TCP 시작 - 받은 카운팅 정보를 데이터베이스에 삽입
		go tcp.StartTCPReceivingCountingInfo(networkHost, networkPort)
	case 2:
		// UDP 시작 - 카운팅 정보를 생성하여 Consumer에게 보냄
		go udp.StartUDPProducingCountingInfo(networkHost, networkPort)
		// UDP 시작 - 받은 카운팅 정보를 데이터베이스에 삽입
		go udp.StartUDPReceivingCountingInfo(networkHost, networkPort)
	case 3:
		// websocket 시작 - 카운팅 정보를 생성하여 Consumer에게 보냄
		go websocket.StartProducingCountingInfo(networkHost, networkPort)
		// websocket 시작 - 받은 카운팅 정보를 데이터베이스에 삽입
		go websocket.StartReceivingCountingInfo(networkHost, networkPort)
	case 4:
		// REST 시작 - 카운팅 정보를 생성하여 Consumer에게 보냄
		go rest.StartRESTProducingCountingInfo(networkHost, networkPort)
		// REST 시작 - 받은 카운팅 정보를 데이터베이스에 삽입
		go rest.StartRESTReceivingCountingInfo(networkHost, networkPort)
	default:
		fmt.Println("Invalid method selected.")
		return
	}
	// 메인 스레드를 계속 실행 상태로 유지
	select {}

}
