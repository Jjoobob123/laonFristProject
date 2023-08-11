package udp

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

func StartUDPProducingCountingInfo(networkHost string, networkPort int) {
	udpAddr, err := net.ResolveUDPAddr("udp", networkHost+":"+strconv.Itoa(networkPort))
	if err != nil {
		log.Fatal("UDP 주소 해석 실패:", err)
		return
	}

	// 무한 루프 막기!
	maxMessages := 100
	for i := 0; i < maxMessages; i++ {
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
			continue // 에러가 발생하면 다음 루프로 넘어감
		}
		fmt.Println("TTTTTTTTTTTTTTTT", reflect.TypeOf(jsonData))

		// 데이터 전송
		err = sendUDPData(jsonData, udpAddr)
		if err != nil {
			log.Println("UDP 전송 실패:", err)
		}

		time.Sleep(1 * time.Second)
	}
}

func sendUDPData(data []byte, udpAddr *net.UDPAddr) error {
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}
