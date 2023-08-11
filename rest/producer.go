package rest

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func StartRESTProducingCountingInfo(networkHost string, networkPort int) {
	client := resty.New()
	rand.Seed(time.Now().UnixNano())

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
		url := "http://" + networkHost + ":" + strconv.Itoa(networkPort) + "/counting"

		// json post형식
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(collection_data).
			Post(url)

		if err != nil {
			log.Fatal("REST client 에러 발생:", err)
			return
		}

		if resp.StatusCode() != 200 {
			log.Println("서버에서 에러 발생. 응답 코드:", resp.Status())
		} else {
			log.Println("서버 응답:", resp.Status())
		}

		time.Sleep(1 * time.Second)
	}
}
