package lang

import (
	"io"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

type Parser struct {
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()
		op, err := parseCommand(commandLine)

		if err != nil {
			return res, err
		}

		res = append(res, op)
	}

	return res, nil
}

func parseCommand(commandLine string) (painter.Operation, error) {
	parsedCommand := strings.Fields(commandLine)
	commandName := parsedCommand[0]
	commandParams := parsedCommand[1:]

	switch commandName {
	case "white":
		return painter.Fill{Color: color.White}, nil
	case "green":
		return painter.Fill{Color: color.RGBA{G: 0xff, A: 0xff}}, nil
	case "update":
		return painter.UpdateOp, nil
	case "bgrect":
		params, err := parseParams(commandParams, 4)
		if err != nil {
			return nil, err
		}
		return painter.BgRect{
			X1: params[0],
			Y1: params[1],
			X2: params[2],
			Y2: params[3],
		}, nil
	case "figure":
		params, err := parseParams(commandParams, 2)
		if err != nil {
			return nil, err
		}
		return painter.Figure{X: params[0], Y: params[1]}, nil
	case "move":
		params, err := parseParams(commandParams, 2)
		if err != nil {
			return nil, err
		}
		return painter.Move{X: params[0], Y: params[1]}, nil
	case "reset":
		return painter.ResetOp, nil
	default:
		return nil, errors.New("unknown command")
	}
}

func parseParams(params []string, length int) ([]float32, error) {
	var res []float32

	if len(params) != length {
		return nil, errors.New("invalid params count")
	}

	for _, item := range params {
		floatNum, err := strconv.ParseFloat(item, 32)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		if floatNum < 0 || floatNum > 1 {
			return nil, errors.New("invalid coordinates")
		}

		res = append(res, float32(floatNum))
	}

	return res, nil
}
