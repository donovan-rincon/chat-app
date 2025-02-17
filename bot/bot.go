package bot

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

type StockMessage struct {
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func GetRabbitMQURL() string {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}
	return rabbitMQURL
}

func ProcessStockRequest(StockCommand string) {
	stockCommandParams := strings.Split(StockCommand, "=")
	if len(stockCommandParams) != 2 {
		sendInvalidCommandResponse(StockCommand)
		return
	}
	stockCode := stockCommandParams[1]
	response, err := http.Get(fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", stockCode))
	if err != nil {
		log.Printf("Failed to call stock API: %v", err)
		sendErrorResponse(stockCode)
		return
	}
	body, _ := io.ReadAll(response.Body)
	response.Body.Close()
	reader := csv.NewReader(strings.NewReader(string(body)))
	records, _ := reader.ReadAll()

	if len(records) == 0 || len(records[0]) == 0 {
		sendInvalidCommandResponse(stockCode)
		return
	}

	symbol := records[1][0] // share symbol from response
	price := records[1][6]  // share price from response (considering close is the share price)
	msg := fmt.Sprintf("Stock %s quote is $%s per share", symbol, price)
	stockMessage := StockMessage{
		Username:  "stock_bot",
		Message:   msg,
		Timestamp: "now",
	}
	jsonMessage, err := json.Marshal(stockMessage)
	if err != nil {
		log.Fatalf("Failed to marshal stock message: %v", err)
		sendErrorResponse(err.Error())
	}
	sendStockMessage(jsonMessage)
}

func sendStockMessage(stockMessage []byte) {
	conn, err := amqp.Dial(GetRabbitMQURL())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"chatroom_messages",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        stockMessage,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
	log.Printf("Sent message: %s", stockMessage)
}

func sendErrorResponse(stockCode string) {
	sendStockMessage([]byte(fmt.Sprintf("Error: Unable to fetch stock for command '%s'. Please try again later.", stockCode)))
}

func sendInvalidCommandResponse(stockCode string) {
	sendStockMessage([]byte(fmt.Sprintf("Error: Invalid stock command '%s'. Please check the available commands and try again.", stockCode)))
}
