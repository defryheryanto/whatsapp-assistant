package assistant

import (
	"testing"

	"github.com/onsi/gomega"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func Test_getMessage(t *testing.T) {
	m := gomega.NewWithT(t)
	t.Run("should get message from conversation", func(t *testing.T) {
		message := getMessage(&events.Message{
			Message: &waE2E.Message{
				Conversation: proto.String("Hello World!"),
			},
		})
		m.Expect(message).To(gomega.Equal("Hello World!"))
	})
	t.Run("should get message from extended text message if conversation is empty", func(t *testing.T) {
		message := getMessage(&events.Message{
			Message: &waE2E.Message{
				ExtendedTextMessage: &waE2E.ExtendedTextMessage{
					Text: proto.String("Hello World!!!"),
				},
			},
		})
		m.Expect(message).To(gomega.Equal("Hello World!!!"))
	})
}

func Test_extractCommands(t *testing.T) {
	m := gomega.NewWithT(t)
	t.Run("should extract commands delimited by enter", func(t *testing.T) {
		commands := extractCommands("%assign\ntest @628123456789 @628213458697")
		m.Expect(commands).To(gomega.BeEquivalentTo([]string{"assign"}))
	})
	t.Run("should extract commands delimited by space", func(t *testing.T) {
		commands := extractCommands("%assign test @628123456789 @628213458697")
		m.Expect(commands).To(gomega.BeEquivalentTo([]string{"assign"}))
	})
}
