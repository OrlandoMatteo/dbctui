package can

import "fmt"

type Signal struct {
	Name           string    `json:"name"`
	Label          string    `json:"label"`
	StartBit       uint64    `json:"startBit"`
	BitLength      uint64    `json:"bitLength"`
	IsLittleEndian bool      `json:"isLittleEndian"`
	IsSigned       bool      `json:"isSigned"`
	Factor         float64   `json:"factor"`
	Offset         float64   `json:"offset"`
	Min            float64   `json:"min"`
	Max            float64   `json:"max"`
	SourceUnit     string    `json:"sourceUnit"`
	DataType       string    `json:"dataType"`
	Choking        bool      `json:"choking"`
	Visibility     bool      `json:"visibility"`
	Interval       uint64    `json:"interval"`
	Category       string    `json:"category"`
	LineInDbc      int       `json:"lineInDbc"`
	Problems       []Problem `json:"problems"`
	PostfixMetric  string    `json:"postfixMetric"`
	States         []State   `json:"states"`
	MsgID          uint64    `json:"msgId"`
	MsgName        string    `json:"msgName"`
	SigID          uint64    `json:"sig_id"`
}

func PrintSignal(i *Signal, tabbed bool) {
	if tabbed == false {
		fmt.Printf("Signal: %d\n", *i)
	} else {
		fmt.Printf("\tSignal: %s  with id %d\n", i.Name, i.SigID)
	}
}
