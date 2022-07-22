package workers

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/apex/log"
)

// WorkQueue A buffered channel that we can send work requests on.
// The buffer size of the channel is completely arbitrary,
// but you want to set it high enough so that sending work requests over
// it does not fill up, and block the send operation: WorkQueue <- work.
var WorkQueue = make(chan *MessageRequest, 10000)

// MessageRequest message to be published on Google Pub/Sub
type MessageRequest struct {
	Ctx   context.Context
	Topic *pubsub.Topic
	Attr  map[string]string
}

// MessageRequests collection
type MessageRequests struct {
	MessageRequests []*MessageRequest
}

func publishMessageRequest(m *MessageRequest) error {
	result := m.Topic.Publish(m.Ctx, &pubsub.Message{
		Attributes: m.Attr,
	})
	// Block until the result is returned and a server-generated
	// Guid is returned for the published message.
	id, err := result.Get(m.Ctx)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"objectID": m.Attr["objectID"],
		}).Error("failed to publish message")
	} else {
		log.Infof("Published message Guid: %s", id)
	}

	return nil
}
