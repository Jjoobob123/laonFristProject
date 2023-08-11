package udp

import (
	"encoding/json"
	"laonproject0801/db"
	"log"
	"net"
	"strconv"
	"time"
)

type CountingInfo struct {
	AccessSeq   int       `json:"accessSeq"`
	VehicleType int       `json:"vehicleType"`
	Lane        int       `json:"lane"`
	Direction   int       `json:"direction"`
	Speed       int       `json:"speed"`
	CreateAt    time.Time `json:"createAt"`
}

func StartUDPReceivingCountingInfo(networkHost string, networkPort int) {
	// UDP 주소 해석
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(networkPort))
	if err != nil {
		log.Fatal("UDP 주소 해석 실패:", err)
		return
	}

	// UDP 수신 시작
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("UDP 수신 시작 실패:", err)
		return
	}
	defer conn.Close()

	for {
		data := make([]byte, 4096)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println("UDP 수신 실패:", err)
			continue
		}

		// 수신한 데이터 출력
		response := string(data[:n])
		log.Println("UDP 메시지 수신:", response)

		// 수신한 데이터 처리 예시 (이 부분을 변경하여 데이터베이스에 삽입하면 됩니다.)
		var countingInfo CountingInfo
		err = json.Unmarshal(data[:n], &countingInfo)
		if err != nil {
			log.Println("JSON 해석 에러:", err)
			continue
		}

		// 유효성 체크 및 차종, 차 방향 변수 설정
		var vehicleTypeStr string
		switch countingInfo.VehicleType {
		case 2:
			vehicleTypeStr = "Car"
		case 3:
			vehicleTypeStr = "SmallTruck"
		case 4:
			vehicleTypeStr = "LargeTruck"
		case 5:
			vehicleTypeStr = "SmallBus"
		case 6:
			vehicleTypeStr = "LargeBus"
		case 7:
			vehicleTypeStr = "Motorcycle"
		default:
			vehicleTypeStr = "Unknown"
		}

		var directionStr string
		switch countingInfo.Direction {
		case 1:
			directionStr = "Stright"
		case 2:
			directionStr = "LeftTurn"
		case 3:
			directionStr = "RightTurn"
		case 4:
			directionStr = "Uturn"
		default:
			directionStr = "Unknown"
		}

		createdAtStr := countingInfo.CreateAt.Format("2006-01-02 15:04:05")
		// 데이터베이스에 저장
		err = db.InsertCountingInfo(countingInfo.AccessSeq, vehicleTypeStr, countingInfo.Lane, directionStr, countingInfo.Speed, createdAtStr)
		if err != nil {
			log.Println("데이터베이스 저장 에러:", err)
		} else {
			log.Println("데이터베이스에 성공적으로 저장됨")
		}

		// 클라이언트에게 응답 전송 (예시로 "OK"를 전송)
		_, err = conn.WriteToUDP([]byte("OK"), addr)
		if err != nil {
			log.Println("UDP 응답 전송 실패:", err)
		}
	}
}
