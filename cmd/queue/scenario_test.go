package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/ruslan-onishchenko/go-test-task/pkg/servelayer/queue/v1"
)

func prepareURI(t *testing.T, queueName string, timeout ...int) (put string, get string) {
	t.Helper()

	if len(queueName) == 0 {
		t.Fatal("queueName is empty")
	}

	var timeoutPrefix string
	if len(timeout) > 0 && timeout[0] > 0 {
		timeoutPrefix = fmt.Sprintf("?timeout=%v", timeout[0])
	}

	put = "/v1/queue/" + queueName
	get = put + timeoutPrefix

	return put, get
}

func reverseStringSlice(t *testing.T, s []string) []string {
	t.Helper()

	result := make([]string, len(s))
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = s[j], s[i]
	}

	if len(s)%2 == 1 {
		result[len(s)/2] = s[len(s)/2]
	}

	return result
}

func TestScenarioSimple(t *testing.T) {
	firstQueueSend := []string{"one", "two", "three"}
	firstQueueAnswer := reverseStringSlice(t, firstQueueSend)
	firstQueueEqual := []bool{true, true, true}

	secondQueueSend := []string{"one", "two", "three"}
	secondQueueAnswer := reverseStringSlice(t, firstQueueSend)
	secondQueueAnswer[0] = "aa"
	secondQueueAnswer[1] = "bb"
	secondQueueEqual := []bool{false, false, true}

	tests := []struct {
		name   string
		send   []string
		answer []string
		equal  []bool
	}{
		{
			name:   "first-queue",
			send:   firstQueueSend,
			answer: firstQueueAnswer,
			equal:  firstQueueEqual,
		},
		{
			name:   "second-queue",
			send:   secondQueueSend,
			answer: secondQueueAnswer,
			equal:  secondQueueEqual,
		},
	}

	host := fmt.Sprintf("http://localhost:%v", *port)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(
				t,
				len(test.send),
				len(test.answer),
				"incorrect test data length of send slice must be equal as answer slice",
			)
			require.Equal(
				t,
				len(test.send),
				len(test.equal),
				"incorrect test data length of send slice must be equal as equal slice",
			)

			put, get := prepareURI(t, test.name)
			for _, str := range test.send {
				push := pb.PushRequest{
					Message: str,
				}

				byteSlice, err := json.Marshal(&push)
				require.NoError(t, err)

				req, err := http.NewRequest(http.MethodPut, host+put, bytes.NewReader(byteSlice))
				require.NoError(t, err)

				req.Header.Set("Content-Type", "application/json")

				response, err := http.DefaultClient.Do(req)
				require.NoError(t, err)

				require.Equal(t, 200, response.StatusCode)
			}

			for index, str := range test.answer {
				req, err := http.NewRequest(http.MethodGet, host+get, nil)
				require.NoError(t, err)

				req.Header.Set("Content-Type", "application/json")

				response, err := http.DefaultClient.Do(req)
				require.NoError(t, err)

				require.Equal(t, response.StatusCode, 200)

				byteSlice, err := io.ReadAll(response.Body)
				require.NoError(t, err)

				var resp pb.GetResponse

				err = json.Unmarshal(byteSlice, &resp)
				require.NoError(t, err)

				if test.equal[index] {
					require.Equal(t, resp.GetMessage(), str)
				} else {
					require.NotEqual(t, resp.GetMessage(), str)
				}
			}
		})
	}
}
