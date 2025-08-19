package dbc

import (
	"dbctui/can"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Parse(fileContent string) ([]*can.Message, []*can.Signal) {
	var messages []*can.Message
	var signals []*can.Signal
	// Split fileContent in lines
	lines := strings.Split(fileContent, "\n")
	currentMessage := &can.Message{}
	for i, line := range lines {
		// split line on whitespaces
		tokens := strings.Split(strings.TrimSpace(line), " ")
		if len(tokens) == 0 {
			continue
		}
		switch tokens[0] {
		case "BO_":
			message, err := parseMessage(tokens, i)
			if err != nil {
				fmt.Printf("Error parsing message at line %d : %s\n", i, err)
			} else {
				// a new message is being parsed
				if message.CanId != currentMessage.CanId {
					messages = append(messages, currentMessage)
				}
				currentMessage = &message
			}
		case "SG_":
			signal, err := parseSignal(tokens, i, currentMessage)
			if err != nil {
				fmt.Printf("Error parsing signal at line %d : %s\n", i, err)
			} else {
				signals = append(signals, &signal)
				currentMessage.Signals = append(currentMessage.Signals, &signal)
			}
		case "VAL_":
			pSignal, err := findSignal(tokens[2], signals)
			if err != nil {
				fmt.Printf("Cannot find signal %s to assign values: %s\n", tokens[2], err)
			} else {
				err = parseValData(tokens, pSignal)
				if err != nil {
					fmt.Printf("Error parsing signal values at line %d : %s\n", i, err)
				}
			}
		case "BA_":
			CI := strings.ReplaceAll(tokens[1], "\"", "")
			if CI == "CI_SigId" {
				pSignal, err := findSignal(tokens[4], signals)
				if err != nil {
					fmt.Printf("Cannot find signal %s to assign id: %s\n", tokens[2], err)
				} else {
					err = parseSignalId(tokens, pSignal)
					if err != nil {
						fmt.Printf("Error parsing signal id at line %d : %s\n", i, err)
					}
				}
			}
		default:
			continue
		}
	}
	return messages, signals
}

func parseMessage(tokens []string, lineInDbc int) (can.Message, error) {
	if len(tokens) != 5 {
		return can.Message{}, errors.New(fmt.Sprintf("Expected 5 tokens, found %d", len(tokens)))
	}
	message := can.Message{}
	err := error(nil)
	_canId, err := strconv.ParseUint(tokens[1], 10, 64)
	message.CanId = _canId & 0x1fffffff
	if message.CanId == 0 {
		return can.Message{}, errors.New(fmt.Sprintf("Can't parse message id %s", tokens[1]))
	}
	if err != nil {
		return can.Message{}, errors.New(fmt.Sprintf("Error parsing message id %s", tokens[1]))
	}
	message.Name = strings.Replace(tokens[2], ":", "", -1)
	message.Dlc, err = strconv.ParseUint(tokens[1], 10, 64)
	if err != nil {
		return can.Message{}, errors.New(fmt.Sprintf("Error parsing message dlc %s", tokens[2]))
	}
	splitCanId(&message)
	message.LineInDbc = lineInDbc

	return message, nil
}

func splitCanId(message *can.Message) {
	canId := message.CanId
	isExtendedFrame := canId > 0xffff

	if isExtendedFrame {
		message.Source = canId & 0xff
		message.Pgn = (canId >> 8) & 0xffff
		message.Priority = (canId >> 24) & 0xff
	} else {
		message.Pgn = canId
	}
}

func parseSignal(tokens []string, lineInDbc int, message *can.Message) (can.Signal, error) {
	signal := can.Signal{}

	signal.Name = tokens[1]
	startToken := 3
	if tokens[3] == ":" {
		startToken = 4
	}
	bitInfo := tokens[startToken]
	factorOffest := tokens[startToken+1]
	minMax := tokens[startToken+2]

	err := parseBitInfo(bitInfo, &signal)
	if err != nil {
		return signal, err
	}
	err = parseFactorOffset(factorOffest, &signal)
	if err != nil {
		return signal, err
	}
	err = parseMinMax(minMax, &signal)
	if err != nil {
		return signal, err
	}
	signal.MsgID = message.CanId
	signal.MsgName = message.Name
	signal.LineInDbc = lineInDbc
	signal.Label = message.Name

	return signal, nil
}

func parseBitInfo(token string, signal *can.Signal) error {
	re := regexp.MustCompile(`^(\d+)\|(\d+)@(\d+)`)
	matches := re.FindStringSubmatch(token)

	if matches == nil || len(matches) != 4 {
		return errors.New("Error parsing bit info for string " + token)
	}
	matches = matches[1:]
	startBit, err := strconv.ParseUint(matches[0], 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Can't parse signal startBit %s", token))
	}
	signal.StartBit = startBit

	bitLength, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Can't parse signal bitLength %s", token))
	}
	signal.BitLength = bitLength
	signal.IsLittleEndian, err = strconv.ParseBool(matches[2])
	if err != nil {
		return errors.New(fmt.Sprintf("Can't parse signal isLittleEndian %s", token))
	}
	return nil
}

func parseFactorOffset(token string, signal *can.Signal) error {
	re := regexp.MustCompile(`[+-]?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?`)
	matches := re.FindAllString(token, -1)
	if matches == nil || len(matches) != 2 {
		return errors.New("Error parsing factorOffset for string " + token)
	}
	factor, err := strconv.ParseFloat(matches[0], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Can't parse factor %s", token))
	}
	signal.Factor = factor

	offset, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		signal.Offset = offset
	}
	return nil
}

func parseMinMax(token string, signal *can.Signal) error {
	re := regexp.MustCompile(`-?(\d+)\|(\d+)`)
	matches := re.FindStringSubmatch(token)
	if matches == nil || len(matches) != 3 {
		return errors.New("Error parsing minMax for string " + token)
	}
	matches = matches[1:]
	minVal, err := strconv.ParseFloat(matches[0], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Can't parse min value %s", token))
	} else {
		signal.Min = minVal
	}
	maxVal, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Can't parse max value %s", token))
	}
	signal.Max = maxVal
	return nil
}

func findSignal(signalName string, signals []*can.Signal) (*can.Signal, error) {
	for _, signal := range signals {
		if signal.Name == signalName {
			return signal, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Can't find signal %s", signalName))
}

func parseValData(tokens []string, signal *can.Signal) error {
	line := strings.Join(tokens[3:], " ")
	re := regexp.MustCompile(`\d+|"([^"]*)"`)
	matches := re.FindAllString(line, -1)
	if matches == nil {
		return errors.New("Error parsing val data for string " + tokens[2])
	}
	for i := 0; i < len(matches); i = i + 2 {
		state := can.State{}
		value, err := strconv.ParseUint(matches[i], 10, 64)
		if err != nil {
			return errors.New(fmt.Sprintf("Can't parse value %s", matches[i]))
		}
		state.Value = value

		state.Name = strings.ReplaceAll(matches[i+1], "\"", "")
		signal.States = append(signal.States, state)
	}
	return nil
}

func parseSignalId(tokens []string, signal *can.Signal) error {

	sigIdStr := strings.ReplaceAll(tokens[5], ";", "")
	sigId, err := strconv.ParseUint(sigIdStr, 10, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Can't parse signal id %s", sigIdStr))
	}
	signal.SigID = sigId

	return nil
}
