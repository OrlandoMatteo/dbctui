package can

import "fmt"

type Message struct {
	CanId           uint64    `json:"can_id"`
	Pgn             uint64    `json:"pgn"`
	Source          uint64    `json:"source"`
	Name            string    `json:"name"`
	Priority        uint64    `json:"priority"`
	Label           string    `json:"label"`
	IsExtendedFrame bool      `json:"isExtendedFrame"`
	Dlc             uint64    `json:"dlc"`
	Comment         string    `json:"comment"`
	Signals         []*Signal `json:"signals"`
	LineInDbc       int       `json:"lineInDbc"`
}

func PrintMessage(canMessage *Message) {
	fmt.Printf("CanId %d, Name %s\n", canMessage.CanId, canMessage.Name)
	for _, signal := range canMessage.Signals {
		PrintSignal(signal, true)
	}
}

func PrintMessages(messages []*Message) {
	for _, message := range messages {
		PrintMessage(message)
	}
}
