package tcp

import (
	"encoding/json"
	"fmt"
	"io"
	"laonproject0801/db"
	"log"
	"net"
	"reflect"
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

func StartTCPReceivingCountingInfo(networkHost string, networkPort int) {
	conn, err := net.Listen("tcp", networkHost+":"+strconv.Itoa(networkPort))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	for {
		conn, err := conn.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		Handler(conn)
		conn.Close()
	}
}

func Handler(conn net.Conn) {
	fadata := make([]byte, 4096)
	for {
		n, err := conn.Read(fadata)
		if err != nil {
			if io.EOF == err {
				log.Println(err)
				return
			}
			log.Println(err)
			return
		}
		if 0 < n {
			data := fadata[:n]
			var countingInfo CountingInfo
			err := json.Unmarshal(data, &countingInfo)
			if err != nil {
				log.Println("JSON 해석 에러:", err)
				return
			}
			fmt.Println("!$@@!$", reflect.TypeOf(data))

			// 차종 변수
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

			// 차 방향 변수
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
			fmt.Println("머가 문제야 !!!", vehicleTypeStr, directionStr)

			createdAtStr := countingInfo.CreateAt.Format("2006-01-02 15:04:05")
			// 데이터베이스에 저장
			err = db.InsertCountingInfo(countingInfo.AccessSeq, vehicleTypeStr, countingInfo.Lane, directionStr, countingInfo.Speed, createdAtStr)
			if err != nil {
				log.Println("데이터베이스 저장 에러:", err)
			} else {
				log.Println("데이터베이스에 성공적으로 저장됨")
			}
		}
	}
}
