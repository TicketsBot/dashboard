package chatreplica

import (
	"fmt"
	v2 "github.com/TicketsBot/logarchiver/pkg/model/v2"
)

func FromTranscript(transcript v2.Transcript, ticketId int) Payload {
	return Payload{
		Entities:    EntitiesFromTranscript(transcript.Entities),
		Messages:    MessagesFromTranscript(transcript.Messages),
		ChannelName: fmt.Sprintf("ticket-%d", ticketId),
	}
}
