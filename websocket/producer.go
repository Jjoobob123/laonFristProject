package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var done chan interface{}    // receiveHandler가 종료됨을 알리는 채널
var interrupt chan os.Signal // 인터럽트 시그널을 받기 위한 채널

// VehicleCountingInfo는 차량 카운팅 정보를 나타냅니다.
type collection_data struct {
	AccessSeq   int       // 접근로 시퀀스(1~100번 사이)
	VehicleType int       // 차종(2: 승용차, 3: 소형 트럭, 4: 대형 트럭, 5: 소형 버스, 6: 대형 버스, 7: 오토바이)
	Lane        int       // 검지된 차선(1~5번 사이)
	Direction   int       // 이동 방향(1: 직진, 2: 좌회전, 3: 우회전, 4: 유턴)
	Speed       int       // 속도(20~60)[정수로 할 것]
	CreateAt    time.Time // 생성 시간[DB에 INSERT할 때 시간으로 할 것]
}

// generateRandomVehicleCountingInfo는 랜덤한 차량 카운팅 정보를 생성합니다.
func generateRandomVehicleCountingInfo() collection_data {
	return collection_data{

		AccessSeq:   rand.Intn(100) + 1,
		VehicleType: rand.Intn(6) + 2,   // 2부터 7까지 (포함) 차종을 나타냄.
		Lane:        rand.Intn(5) + 1,   // 1부터 5까지 (포함) 차선 번호를 나타냄.
		Direction:   rand.Intn(4) + 1,   // 1부터 4까지 (포함) 이동 방향을 나타냄.
		Speed:       rand.Intn(41) + 20, // 20부터 60까지 (포함) 속도를 나타냄.
		CreateAt:    time.Now(),
	}
}

// receiveHandler는 웹소켓 연결에서 데이터를 수신하는 핸들러 함수입니다.
func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage() // 웹소켓으로부터 메시지를 읽어옴
		if err != nil {
			log.Println("수신 오류:", err)
			return
		}
		log.Printf("수신된 데이터: %s\n", msg) // 수신한 메시지를 로그로 출력
		fmt.Println("타입 무엇이드냐", reflect.TypeOf(msg))

	}
}

func StartProducingCountingInfo(networkHost string, networkPort int) {
	done = make(chan interface{})    // receiveHandler가 종료됨을 알리는 채널 생성
	interrupt = make(chan os.Signal) // 인터럽트 시그널을 받기 위한 채널 생성

	signal.Notify(interrupt, os.Interrupt) // SIGINT 시그널을 받을 때 interrupt 채널로 전달

	socketUrl := "ws://" + networkHost + ":" + strconv.Itoa(networkPort) + "/counting"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil) // 웹소켓 서버에 연결
	if err != nil {
		log.Fatal("웹소켓 서버 연결 오류:", err)
	}
	defer conn.Close()

	go receiveHandler(conn) // 데이터 수신을 위해 receiveHandler를 고루틴으로 실행

	// 메인 루프: 웹소켓 서버에 주기적으로 차량 카운팅 정보를 전송
	for {
		select {
		case <-time.After(time.Duration(1) * time.Second): // 1초마다 실행되는 타이머 이벤트
			vehicleInfo := generateRandomVehicleCountingInfo()
			// JSON 데이터 생성
			data := map[string]interface{}{
				"accessSeq":   vehicleInfo.AccessSeq,
				"vehicleType": vehicleInfo.VehicleType,
				"lane":        vehicleInfo.Lane,
				"direction":   vehicleInfo.Direction,
				"speed":       vehicleInfo.Speed,
				"createAt":    vehicleInfo.CreateAt,
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Println("JSON 마샬링 오류:", err)
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, jsonData) // 차량 카운팅 정보 전송
			if err != nil {
				log.Println("웹소켓 데이터 전송 오류:", err)
				return
			}

		case <-interrupt: // SIGINT (Ctrl + C) 시그널을 받으면 이 케이스가 실행됨
			log.Println("SIGINT 인터럽트 시그널 수신. 모든 연결 종료 중")

			// 웹소켓 연결 종료를 위해 CloseNormalClosure 메시지를 보냄
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("웹소켓 종료 오류:", err)
				return
			}

		}
	}
}
