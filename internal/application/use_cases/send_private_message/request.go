package send_private_message

import "bytes"

type Request struct {
	RecipientID string
	MessageBody string
	ResponseToMessageID string
	Attachments []*bytes.Buffer
}
