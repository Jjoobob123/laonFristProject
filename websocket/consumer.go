package websocket

import (
	"encoding/json"
	"fmt"
	"laonproject0801/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // 연결된 클라이언트 저장을 위한 맵
var broadcast = make(chan string)            // 클라이언트로부터 수신한 데이터를 모든 클라이언트에게 보내기 위한 채널
var upgrader = websocket.Upgrader{}          // 웹소켓 연결을 업그레이드하기 위한 업그레이더

// handleConnections은 클라이언트 연결을 처리하는 핸들러 함수입니다.
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 클라이언트를 웹소켓 연결로 업그레이드
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("웹소켓 연결 업그레이드 오류:", err)
		return
	}
	defer conn.Close()

	// 클라이언트를 맵에 추가
	clients[conn] = true

	// 클라이언트로부터 데이터를 읽어오는 루프
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("수신 오류:", err)
			delete(clients, conn) // 에러 발생 시 클라이언트 맵에서 삭제
			break
		}

		// JSON 데이터를 파싱하여 맵으로 변환
		var data map[string]interface{}
		err = json.Unmarshal(msg, &data)
		if err != nil {
			log.Println("JSON 파싱 오류:", err)
			continue
		}

		parseCountingInfo(data)

		// 수신한 데이터를 모든 클라이언트에게 보내기 위해 broadcast 채널로 전송
		broadcast <- string(msg)
		fmt.Println("boardcast:", broadcast)

	}
}

func parseCountingInfo(data map[string]interface{}) {

	// 필드 검증 및 형 변환
	accessSeqFloat, ok := data["accessSeq"].(float64)
	if !ok {

		return
	}
	accessSeq := int(accessSeqFloat)

	vehicleTypeFloat, ok := data["vehicleType"].(float64)
	if !ok {

		return
	}
	vehicleType := int(vehicleTypeFloat)

	laneFloat, ok := data["lane"].(float64)
	if !ok {

		return
	}
	lane := int(laneFloat)

	directionFloat, ok := data["direction"].(float64)
	if !ok {
		return
	}
	direction := int(directionFloat)

	speedFloat, ok := data["speed"].(float64)
	if !ok {

		return
	}
	speed := int(speedFloat)

	createString := data["createAt"].(string)
	if !ok {
		log.Println("타임스탬프 형식 오류")
		return
	}

	originalTimeLayout := "2006-01-02T15:04:05.999999999-07:00" // 원본 타임스탬프의 레이아웃
	originalTime, err := time.Parse(originalTimeLayout, createString)
	if err != nil {
		log.Println("타임스탬프 파싱 오류:", err)
		return
	}

	// 원하는 형식으로 타임스탬프 포맷팅
	desiredTimeLayout := "2006-01-02 15:04:05" // 원하는 레이아웃
	createAtStr := originalTime.Format(desiredTimeLayout)

	// 차종 변수
	var vehicleTypeStr string
	switch vehicleType {
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
	fmt.Println(vehicleTypeStr, "부르릉")
	// 차 방향 변수
	var directionStr string
	switch direction {
	case 1:
		directionStr = "Stright"
	case 2:
		directionStr = "LeftTurn"
	case 3:
		directionStr = "RightTurn"
	case 4:
		directionStr = "Uturn"
	}
	// 데이터베이스 삽입
	err = db.InsertCountingInfo(accessSeq, vehicleTypeStr, lane, directionStr, speed, createAtStr)
	if err != nil {
		log.Println("데이터베이스 삽입 오류:", err)
		return

	}

}

// handleMessages는 클라이언트로부터 수신한 데이터를 모든 클라이언트에게 보내는 핸들러 함수입니다.
func handleMessages() {
	for {
		// broadcast 채널로부터 수신한 데이터를 읽어오고 클라이언트에게 전송
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("전송 오류:", err)
				client.Close()
				delete(clients, client) // 에러 발생 시 클라이언트 맵에서 삭제
			}
		}
	}
}

func StartReceivingCountingInfo(networkHost string, networkPort int) {
	// 서버가 /counting 엔드포인트로 클라이언트의 웹소켓 연결 요청을 처리
	http.HandleFunc("/counting", handleConnections)

	// 클라이언트로부터 수신한 데이터를 모든 클라이언트에게 보내는 핸들러를 별도의 고루틴으로 실행
	go handleMessages()

	// 서버 실행
	// for {
	// 	err := http.ListenAndServe(":"+strconv.Itoa(networkPort), nil)
	// 	if err != nil {
	// 		if strings.Contains(err.Error(), "address already in use") {
	// 			fmt.Println("포트 번호", networkPort, "사용 중. 다른 포트 시도 중...")
	// 			networkPort++ // 다른 포트 번호 시도
	// 		} else {
	// 			log.Fatal("웹서버 실행 오류:", err)
	// 		}
	// 	} else {
	// 		break // 포트 충돌 없이 서버 실행되었을 때 루프 종료
	// 	}
	// }
	// for {
	// 	err := http.ListenAndServe(":"+strconv.Itoa(networkPort), nil)
	// 	if err != nil {
	// 		if strings.Contains(err.Error(), "address already in use") {
	// 			fmt.Println("포트 번호", networkPort, "사용 중. 기존 포트 재시도 중...")
	// 		} else {
	// 			log.Fatal("웹서버 실행 오류:", err)
	// 		}
	// 	} else {
	// 		break // 포트 충돌 없이 서버 실행되었을 때 루프 종료
	// 	}
	// }
	// log.Println("웹소켓 서버가  실행됩니다.")
	err := http.ListenAndServe(":"+strconv.Itoa(networkPort), nil)
	if err != nil {
		log.Fatal("웹서버 실행 오류:", err)
	}
}
