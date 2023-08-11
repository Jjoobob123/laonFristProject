package rest

import (
	"fmt"
	"laonproject0801/db"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartRESTReceivingCountingInfo(networkHost string, networkPort int) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	port := strconv.Itoa(networkPort)

	r.POST("/counting", func(c *gin.Context) {
		var data map[string]interface{}

		// JSON 바인딩 검증
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
			return
		}

		// 필드 검증 및 형 변환
		accessSeqFloat, ok := data["accessSeq"].(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid accessSeq"})
			return
		}
		accessSeq := int(accessSeqFloat)

		vehicleTypeFloat, ok := data["vehicleType"].(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicleType"})
			return
		}
		vehicleType := int(vehicleTypeFloat)

		laneFloat, ok := data["lane"].(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lane"})
			return
		}
		lane := int(laneFloat)

		directionFloat, ok := data["direction"].(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid direction젠장"})
			return
		}
		direction := int(directionFloat)

		speedFloat, ok := data["speed"].(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid speed"})
			return
		}
		speed := int(speedFloat)

		createAtString, ok := data["createAt"].(string)
		if !ok {
			return
		}
		createAtStr := string(createAtString)
		fmt.Println(data["createAt"])
		fmt.Println(reflect.TypeOf(data["createAt"]))
		fmt.Println(reflect.TypeOf(createAtStr))

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
		fmt.Println(reflect.TypeOf(vehicleTypeStr), "안녕안녕안녕")
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
		err := db.InsertCountingInfo(accessSeq, vehicleTypeStr, lane, directionStr, speed, createAtStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "데이터베이스에 데이터 삽입 실패", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Fatal(r.Run(":" + port))
}
