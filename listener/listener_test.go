package listener

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/twitchscience/aws_utils/logger"
)

type mockSQSClient struct {
	sqsiface.SQSAPI
}

func (m *mockSQSClient) DeleteMessage(msg *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return nil, nil
}

type testSQSHandler struct {
	Called int
}

func (h *testSQSHandler) Handle(msg *sqs.Message) error {
	h.Called++
	return nil
}

func TestSQSListenerHandling(t *testing.T) {
	uniq := 1000
	handler := &testSQSHandler{Called: 0}
	listener := BuildSQSListener(handler, time.Second, &mockSQSClient{}, nil)
	for i := 0; i < uniq; i++ {
		s := fmt.Sprintf("%v", i)
		listener.handle(&sqs.Message{Body: &s}, nil)
	}
	for i := 0; i < uniq; i++ {
		s := fmt.Sprintf("%v", i)
		listener.handle(&sqs.Message{Body: &s}, nil)
	}
	if handler.Called != 2*uniq {
		t.Errorf("expected %d but got %d", 2*uniq, handler.Called)
	}
}

func TestSQSListenerHandlingDedup(t *testing.T) {
	logger.Init("warn")
	uniq := 1000
	handler := &testSQSHandler{Called: 0}
	filter := NewDedupSQSFilter(uniq, time.Hour)
	listener := BuildSQSListener(handler, time.Second, &mockSQSClient{}, filter)
	for i := 0; i < uniq; i++ {
		s := fmt.Sprintf("%v", i)
		listener.handle(&sqs.Message{Body: &s}, nil)
	}
	for i := 0; i < uniq; i++ {
		s := fmt.Sprintf("%v", i)
		listener.handle(&sqs.Message{Body: &s}, nil)
	}
	if handler.Called != uniq {
		t.Errorf("expected %d but got %d", uniq, handler.Called)
	}
}
