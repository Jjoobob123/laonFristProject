// package tcp

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net"
// 	"reflect"
// 	"strconv"
// 	"time"
// )

// func StartTCPProducingCountingInfo(networkHost string, networkPort int) {
// 	tcpURL := networkHost + ":" + strconv.Itoa(networkPort)
// 	// TCP 연결 부분
// 	conn, err := net.Dial("tcp", tcpURL)
// 	if err != nil {
// 		log.Println("TCP 연결 실패:", err)
// 		return
// 	}

// 	for {
// 		accessSeq := rand.Intn(100) + 1
// 		vehicleType := rand.Intn(6) + 2
// 		lane := rand.Intn(5) + 1
// 		direction := rand.Intn(4) + 1
// 		speed := rand.Intn(41) + 20

// 		collection_data := map[string]interface{}{
// 			"accessSeq":   accessSeq,
// 			"vehicleType": vehicleType,
// 			"lane":        lane,
// 			"direction":   direction,
// 			"speed":       speed,
// 		}
// 		fmt.Println(collection_data)

// 		jsonData, err := json.Marshal(collection_data)
// 		if err != nil {
// 			log.Println("JSON 변환 에러:", err)
// 			return // 에러가 발생하면 루프 탈출
// 		}
// 		fmt.Println("TTTTTTTTTTTTTTTT", reflect.TypeOf(jsonData))

// 		// jsondata를 TCP 연결을 통해 보냄
// 		_, err = conn.Write(jsonData)
// 		if err != nil {
// 			log.Println("TCP 전송 실패:", err)
// 			return // 에러가 발생하면 루프 탈출
// 		}
// 		fmt.Println("2222222222")

// 		// 서버로부터 응답 읽기
// 		data := make([]byte, 4096)
// 		n, err := conn.Read(data)
// 		if err != nil {
// 			log.Println(err)
// 			return // 에러가 발생하면 루프 탈출
// 		}

// 		response := string(data[:n])
// 		log.Println("Server send: " + response)

// 		// 응답 처리 예시
// 		if response == "OK" {
// 			log.Println("서버가 메시지를 성공적으로 처리했습니다.")
// 		} else {
// 			log.Println("서버가 메시지 처리에 실패했습니다.")
// 		}

// 		time.Sleep(1 * time.Second)
// 	}
// }

package tcp

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"time"
)

func sendData(conn net.Conn) {
	for {
		accessSeq := rand.Intn(100) + 1
		vehicleType := rand.Intn(6) + 2
		lane := rand.Intn(5) + 1
		direction := rand.Intn(4) + 1
		speed := rand.Intn(41) + 20
		createAt := time.Now()

		collection_data := map[string]interface{}{
			"accessSeq":   accessSeq,
			"vehicleType": vehicleType,
			"lane":        lane,
			"direction":   direction,
			"speed":       speed,
			"createAt":    createAt,
		}
		fmt.Println(collection_data)

		jsonData, err := json.Marshal(collection_data)
		if err != nil {
			log.Println("JSON 변환 에러:", err)
			return // 에러가 발생하면 함수 종료
		}
		fmt.Println("TTTTTTTTTTTTTTTT", reflect.TypeOf(jsonData))

		// jsondata를 TCP 연결을 통해 보냄
		_, err = conn.Write(jsonData)
		if err != nil {
			log.Println("TCP 전송 실패:", err)
			return // 에러가 발생하면 함수 종료
		}

		time.Sleep(1 * time.Second)
	}
}

func readResponse(conn net.Conn) {
	for {
		// 서버로부터 응답 읽기
		data := make([]byte, 4096)
		n, err := conn.Read(data)
		if err != nil {
			log.Println(err)
			return // 에러가 발생하면 함수 종료
		}

		response := string(data[:n])
		log.Println("Server send: " + response)

		// 응답 처리 예시
		if response == "OK" {
			log.Println("서버가 메시지를 성공적으로 처리했습니다.")
		} else {
			log.Println("서버가 메시지 처리에 실패했습니다.")
		}
	}
}

func StartTCPProducingCountingInfo(networkHost string, networkPort int) {
	tcpURL := networkHost + ":" + strconv.Itoa(networkPort)
	// TCP 연결 부분
	conn, err := net.Dial("tcp", tcpURL)
	if err != nil {
		log.Println("TCP 연결 실패:", err)
		return
	}
	defer conn.Close() // 작업이 끝나면 반드시 연결을 닫아줍니다.

	sendData(conn)     // 데이터 생성 및 전송 고루틴 시작
	readResponse(conn) // 응답 읽기 고루틴 시작

}
