package services

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Parser() {
	file, err := os.Open("weather.dat")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		agent := fiber.Post("http://localhost:3000/api/v1/weather")
		agent.Body([]byte(line)) 

		statusCode, body, errs := agent.Bytes()

		if len(errs) > 0 {
			log.Printf("Error sending POST data #%d: %v\n", lineNum, errs)
			}

		if(statusCode != 200) {
			log.Printf("[%d] ❌ Error sending POST request: Body: %s", statusCode, string(body))
		} else {
			log.Printf("[%d] ✅ Successfully sent Request #%d - Body: %s", statusCode, lineNum,  string(body))
		}
		

		time.Sleep(10 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		log.Print("Error reading file:", err)
	}

}