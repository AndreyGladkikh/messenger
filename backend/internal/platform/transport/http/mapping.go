package http

import "messenger/messenger/internal/messaging/application/query/get_chat_list"

func mapChatList(chatListView get_chat_list.ChatListView) ChatListView {
	items := make([]ChatListViewItem, 0, len(chatListView.Chats))
	for _, item := range chatListView.Chats {
		items = append(items, ChatListViewItem{
			Id: item.ID.String(),
			LastMessage: struct {
				Body   string `json:"body"`
				SentAt string `json:"sentAt"`
			}{
				Body:   item.LastMessage.Body,
				SentAt: item.LastMessage.SentAt.String(),
			},
			Name:        item.Name,
			UnreadCount: item.UnreadCount,
		})
	}

	return ChatListView{
		Chats: items,
	}
}
