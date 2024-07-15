package podsmanagement

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var podList = &PodList{
	pods: []string{
		"vwap",
		"twap",
	},
}

func TestCreatePod(t *testing.T) {
	tests := []struct {
		testName   string
		testInput  string
		testOutput error
	}{
		{
			testName:   "without errors",
			testInput:  "hft",
			testOutput: nil,
		},
		{
			testName:   "with errors",
			testInput:  "vwap",
			testOutput: errors.New("pod is already created"),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := podList.CreatePod(test.testInput)
			require.Equal(t, test.testOutput, err)
		})
	}
}

func TestDeletePod(t *testing.T) {
	tests := []struct {
		testName   string
		testInput  string
		testOutput error
	}{
		{
			testName:   "without errors",
			testInput:  "vwap",
			testOutput: nil,
		},
		{
			testName:   "with errors",
			testInput:  "test",
			testOutput: errors.New("pod does not exist"),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := podList.DeletePod(test.testInput)
			require.Equal(t, test.testOutput, err)
		})
	}
}
