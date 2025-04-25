package lang

import (
	"errors"
	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/stretchr/testify/assert"
	"image/color"
	"strings"
	"testing"
)

func TestParser_ValidCommands(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected painter.Operation
	}{
		{
			name:     "white command",
			input:    "white",
			expected: painter.Fill{Color: color.White},
		},
		{
			name:     "green command",
			input:    "green",
			expected: painter.Fill{Color: color.RGBA{G: 0xff, A: 0xff}},
		},
		{
			name:     "update command",
			input:    "update",
			expected: painter.UpdateOp,
		},
		{
			name:     "bgrect command",
			input:    "bgrect 0.1 0.2 0.9 0.8",
			expected: painter.BgRect{X1: 0.1, Y1: 0.2, X2: 0.9, Y2: 0.8},
		},
		{
			name:     "figure command",
			input:    "figure 0.5 0.5",
			expected: painter.Figure{X: 0.5, Y: 0.5},
		},
		{
			name:     "move command",
			input:    "move 0.3 0.7",
			expected: painter.Move{X: 0.3, Y: 0.7},
		},
		{
			name:     "reset command",
			input:    "reset",
			expected: painter.ResetOp,
		},
	}

	parser := Parser{}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parser.Parse(strings.NewReader(tc.input))
			assert.NoError(t, err)
			assert.Equal(t, 1, len(res))
			assert.EqualValues(t, tc.expected, res[0])
		})
	}
}

func TestParser_InvalidCommands(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "unknown command",
			input:       "invalid",
			expectedErr: errors.New("unknown command"),
		},
		{
			name:        "bgrect with wrong params count",
			input:       "bgrect 0.1 0.2",
			expectedErr: errors.New("invalid params count"),
		},
		{
			name:        "figure with invalid coordinates",
			input:       "figure 1.5 -0.5",
			expectedErr: errors.New("invalid coordinates"),
		},
		{
			name:        "invalid number format",
			input:       "move abc 0.5",
			expectedErr: errors.New("invalid params"),
		},
	}

	parser := Parser{}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.Parse(strings.NewReader(tc.input))
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr.Error())
		})
	}
}

func TestParser_MultiCommand(t *testing.T) {
	input := strings.NewReader("white\nfigure 0.5 0.5\nupdate")
	parser := Parser{}
	res, err := parser.Parse(input)

	assert.NoError(t, err)
	assert.Len(t, res, 3)
	assert.IsType(t, painter.Fill{}, res[0])
	assert.IsType(t, painter.Figure{}, res[1])
	assert.Equal(t, painter.UpdateOp, res[2])
}

func TestParser_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected painter.Operation
	}{
		{
			name:     "min coordinates",
			input:    "bgrect 0 0 0 0",
			expected: painter.BgRect{X1: 0, Y1: 0, X2: 0, Y2: 0},
		},
		{
			name:     "max coordinates",
			input:    "move 1 1",
			expected: painter.Move{X: 1, Y: 1},
		},
	}

	parser := Parser{}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parser.Parse(strings.NewReader(tc.input))
			assert.NoError(t, err)
			assert.EqualValues(t, tc.expected, res[0])
		})
	}
}
