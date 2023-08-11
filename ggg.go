// // package communication

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	// "reflect"
// 	"strconv"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// func SendCountingInfoViaWebSocket(accessSeq, vehicleType, lane, direction, speed int, networkHost string, networkPort int) {
// 	// 웹소켓 서버 주소
// 	websocketURL := "ws://" + networkHost + ":" + strconv.Itoa(networkPort) + "/counting"

// 	// 웹소켓 연결 설정
// 	conn, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
// 	if err != nil {
// 		log.Fatal("웹소켓 연결 실패:", err)
// 	}
// 	defer conn.Close()

// 	// 카운팅 정보를 맵으로 생성
// 	collection_data := map[string]interface{}{
// 		"accessSeq":   accessSeq,
// 		"vehicleType": vehicleType,
// 		"lane":        lane,
// 		"direction":   direction,
// 		"speed":       speed,
// 	}
// 	fmt.Println(collection_data)
// 	// 맵을 JSON으로 변환하여 웹소켓으로 전송
// 	jsonData, err := json.Marshal(collection_data)
// 	if err != nil {
// 		log.Println("JSON 변환 에러:", err)
// 		return
// 	}
// 	fmt.Println("fufufufufuufufu", reflect.TypeOf(jsonData))
// 	// 웹소켓으로 데이터 전송
// 	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
// 		log.Println("웹소켓으로 데이터 전송 실패:", err)
// 		return
// 	}
// 	log.Println("카운팅 정보를 웹소켓으로 전송 완료.")

// 	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
// 		log.Println("write close:", err)
// 		return
// 	}
// }

// // StartProducingCountingInfo 함수는 카운팅 정보를 랜덤하게 생산하고 Consumer에게 보냅니다.
// func StartProducingCountingInfo(networkHost string, networkPort int) {

// 	rand.Seed(time.Now().UnixNano())

// 	for {
// 		accessSeq := rand.Intn(100) + 1
// 		vehicleType := rand.Intn(6) + 2
// 		lane := rand.Intn(5) + 1
// 		direction := rand.Intn(4) + 1
// 		speed := rand.Intn(41) + 20
// 		fmt.Println(networkHost, networkPort)
// 		if networkHost != "" {

// 			SendCountingInfoViaWebSocket(accessSeq, vehicleType, lane, direction, speed, networkHost, networkPort)
// 		} else {
// 			log.Fatal("networkHost가 비어있습니다.")
// 		}

// 		time.After(time.Duration(1) * time.Millisecond * 1000)
// 	}
// }

// // server.go 파일
// package communication

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"reflect"
// 	"strconv"

// 	"laonproject0801/db"

// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// // StartReceivingCountingInfo 함수는 선택된 통신 방식으로 받은 데이터의 유효성을 체크하고 DB에 삽입합니다.
// func StartReceivingCountingInfo(networkHost string, networkPort int) {

// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		upgrader = websocket.Upgrader{
// 			CheckOrigin: func(r *http.Request) bool {
// 				// Sec-WebSocket-Version 헤더를 확인하여 클라이언트가 요청한 프로토콜 버전 확인
// 				version := r.Header.Get("Sec-WebSocket-Version")
// 				supportedVersions := []string{"13"} // 서버가 지원하는 웹소켓 프로토콜 버전
// 				for _, v := range supportedVersions {
// 					if version == v {
// 						fmt.Println("성공")
// 						return true // 클라이언트가 지원하는 프로토콜 버전과 일치하는 경우
// 					}
// 				}
// 				fmt.Println("실패다 바보야!")
// 				return false // 클라이언트가 지원하지 않는 프로토콜 버전인 경우
// 			},
// 		}

// 		conn, err := upgrader.Upgrade(w, r, nil)
// 		if err != nil {
// 			log.Println("웹소켓 연결 실패:", err)
// 			return
// 		}
// 		// defer conn.Close()
// 		log.Println("웹소켓 연결 성공.")

// 		for {

// 			var data map[string]interface{}
// 			fmt.Println(data, "@@!@#!@#@!#!@#")
// 			if err := conn.ReadJSON(&data); err != nil {

// 				log.Println("웹소켓으로부터 데이터 수신 실패:", err)
// 				break
// 			}
// 			// fmt.Println(data, "^^^^^^^^^^^^^^^^^^^6")
// 			// fmt.Println(data)
// 			// fmt.Println("accessSeq의 타입은 무엇일까요:", reflect.TypeOf(data["accessSeq"]))

// 			accessSeqFloat, ok := data["accessSeq"].(float64)
// 			fmt.Println(accessSeqFloat, "@@@@@@@@@@")
// 			if !ok {

// 				log.Println("유효하지 않은 접근로 시퀀스:", data["accessSeq"])
// 				continue
// 			}
// 			accessSeq := int(accessSeqFloat)
// 			fmt.Println("accessSeq의 타입은 무엇일까요:", reflect.TypeOf(accessSeq))
// 			fmt.Println(accessSeq)

// 			vehicleTypeFloat, ok := data["vehicleType"].(float64)
// 			if !ok {
// 				log.Println("유효하지 않은 차종:", data["vehicleType"])
// 				continue
// 			}
// 			vehicleType := int(vehicleTypeFloat)

// 			laneFloat, ok := data["lane"].(float64)
// 			if !ok {
// 				log.Println("유효하지 않은 검지된 차선:", data["lane"])
// 				continue
// 			}
// 			lane := int(laneFloat)

// 			directionFloat, ok := data["direction"].(float64)
// 			if !ok {
// 				log.Println("유효하지 않은 이동 방향:", data["direction"])
// 				continue
// 			}
// 			direction := int(directionFloat)

// 			speedFloat, ok := data["speed"].(float64)
// 			if !ok {
// 				log.Println("유효하지 않은 속도:", data["speed"])
// 				continue
// 			}
// 			speed := int(speedFloat)

// 			// 차종 변수
// 			var vehicleTypeStr string
// 			switch vehicleType {
// 			case 2:
// 				vehicleTypeStr = "Car"
// 			case 3:
// 				vehicleTypeStr = "SmallTruck"
// 			case 4:
// 				vehicleTypeStr = "LargeTruck"
// 			case 5:
// 				vehicleTypeStr = "SmallBus"
// 			case 6:
// 				vehicleTypeStr = "LargeBus"
// 			case 7:
// 				vehicleTypeStr = "Motorcycle"
// 			default:
// 				vehicleTypeStr = "Unknown"
// 			}
// 			fmt.Println(reflect.TypeOf(vehicleTypeStr), "안녕안녕안녕")
// 			// 차 방향 변수
// 			var directionStr string
// 			switch direction {
// 			case 1:
// 				directionStr = "Stright"
// 			case 2:
// 				directionStr = "LeftTurn"
// 			case 3:
// 				directionStr = "RightTurn"
// 			case 4:
// 				directionStr = "Uturn"
// 			}
// 			fmt.Println(reflect.TypeOf(directionStr), "안녕안녕안녕")
// 			fmt.Printf("카운팅 정보 수신: 접근로 시퀀스=%d, 차종=%d, 검지된 차선=%d, 이동 방향=%d, 속도=%d\n",
// 				accessSeq, vehicleType, lane, direction, speed)

// 			// 데이터 유효성 체크 및 DB에 삽입
// 			err := db.InsertCountingInfo(accessSeq, vehicleTypeStr, lane, directionStr, speed)
// 			if err != nil {
// 				log.Println("데이터 유효성 체크 에러 발생:", err)
// 			}
// 		}
// 	})

// 	http.ListenAndServe(":"+strconv.Itoa(networkPort), nil) // int를 문자열로 변환

// }

// 웹소켓 프로듀서 time.After
// select {
// 	case <-time.After(time.Duration(1) * time.Second): // 1초마다 실행되는 타이머 이벤트
// 		vehicleInfo := generateRandomVehicleCountingInfo()
// 		message := []byte(fmt.Sprintf("%d,%d,%d,%d,%d", vehicleInfo.AccessSeq, vehicleInfo.VehicleType, vehicleInfo.Lane, vehicleInfo.Direction, vehicleInfo.Speed))
// 		err := conn.WriteMessage(websocket.TextMessage, message) // 차량 카운팅 정보 전송
// 		if err != nil {
// 			log.Println("웹소켓 데이터 전송 오류:", err)
// 			return
// 		}
