package main

import (
	"strings"
	"time"
)

func GetMessageType(Message MessageStruct) string {
	if !Message.IsSystemMessage {
		return "N/A"
	} else if strings.Contains(Message.Body, "You have a new trade!") {
		if strings.Contains(Message.Subject, "countered your Trade") {
			return "TradeCountered"
		} else if strings.Contains(Message.Subject, "You have a Trade request") {
			return "TradeReceived"
		}
	} else if strings.Contains(Message.Body, "Trade declined.") {
		return "TradeDeclined"
	} else if strings.Contains(Message.Body, "Your Trade is complete.") {
		return "TradeAccepted"
	}

	return "N/A"
}

func HandleMessages(Messages []MessageStruct) []int {
	MessagesToDelete := make([]int, 0)

	for _, Message := range Messages {
		MessageType := GetMessageType(Message)

		if MessageType == "N/A" {
			continue
		}

		if (MessageType == "TradeReceived" && Settings.DeleteTradeReceived) || (MessageType == "TradeDeclined" && Settings.DeleteTradeDeclined) || (MessageType == "TradeAccepted" && Settings.DeleteTradeAccepted) || (MessageType == "TradeCountered" && Settings.DeleteTradeCountered) {
			MessagesToDelete = append(MessagesToDelete, Message.Id)
			continue
		}
	}

	return MessagesToDelete
}

func ChunkArray(Array []int, Size int) [][]int {
	var divided [][]int

	chunkSize := (len(Array) + Size - 1) / Size

	for i := 0; i < len(Array); i += chunkSize {
		end := i + chunkSize

		if end > len(Array) {
			end = len(Array)
		}

		divided = append(divided, Array[i:end])
	}

	return divided
}

func ScanMessages() {
	println("Running")

	MessagesToDelete := make([]int, 0)

	PreviousPageNumber := 0
	for {
		Completed := false

		for {
			Success, Response, IsEnd, PageNumber, Messages := FetchMessages(PreviousPageNumber)
			PreviousPageNumber = PageNumber + 1

			if !Success {
				StatusCode := Response.StatusCode

				if StatusCode == 429 {
					time.Sleep(time.Second * 10)
					break
				}

				println("Failed to fetch message!")
				println(StatusCode)
				break
			}

			MessagesToDelete = append(MessagesToDelete, HandleMessages(Messages)...)
			println(PageNumber)

			Completed = IsEnd

			if IsEnd {
				println("Went through all the messages!")
				break
			}
		}

		if Completed {
			break
		}
	}

	for _, MessageDeleteChunk := range ChunkArray(MessagesToDelete, 50) {
		DeleteMessages(MessageDeleteChunk)
	}

	println("Done!")
	time.Sleep(time.Second * 3)
}
