package models

const (
	MessageTypeKefu    = "kefu"
	MessageTypeVisitor = "visitor"

	MessageStatusRead   = "read"
	MessageStatusUnread = "unread"
)

func InitialMessageStatus(mesType string) string {
	if mesType == MessageTypeVisitor {
		return MessageStatusUnread
	}
	return MessageStatusRead
}
