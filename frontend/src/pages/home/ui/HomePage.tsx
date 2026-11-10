import { useState } from 'react'
import { ChatList } from './ChatList'
import { ConversationView } from './ConversationView'
import { NavRail } from './NavRail'
import { conversations as initialConversations, currentUser } from '@/shared/lib/messenger-data'
import { cn } from '@/shared/lib'

export function HomePage() {
    const [chats, setChats] = useState(initialConversations)
      const [activeId, setActiveId] = useState(initialConversations[0].id)
      const [showConversationOnMobile, setShowConversationOnMobile] = useState(false)
    
      const active = chats.find((chat) => chat.id === activeId) ?? chats[0]
    
      function selectChat(id: string) {
        setActiveId(id)
        setShowConversationOnMobile(true)
        setChats((prev) => prev.map((chat) => (chat.id === id ? { ...chat, unread: 0 } : chat)))
      }
    
      function sendMessage(text: string) {
        const time = new Date().toLocaleTimeString('ru-RU', {
          hour: '2-digit',
          minute: '2-digit',
        })
    
        setChats((prev) =>
          prev.map((chat) =>
            chat.id === active.id
              ? {
                  ...chat,
                  preview: `Вы: ${text}`,
                  previewOutgoing: true,
                  time,
                  messages: [
                    ...chat.messages,
                    {
                      id: `local-${Date.now()}`,
                      authorId: currentUser.name,
                      text,
                      time,
                      outgoing: true,
                      status: 'sent' as const,
                    },
                  ],
                }
              : chat,
          ),
        )
      }

    return (
        <>
            <main className="flex h-dvh overflow-hidden bg-background">
                <NavRail />

                <div
                    className={cn(
                        'h-full w-full shrink-0 md:block md:w-auto',
                        showConversationOnMobile ? 'hidden' : 'block',
                    )}
                >
                    <ChatList conversations={chats} activeId={active.id} onSelect={selectChat} />
                </div>

                <div
                    className={cn(
                        'h-full min-w-0 flex-1 md:block',
                        showConversationOnMobile ? 'block' : 'hidden',
                    )}
                >
                    <ConversationView
                        conversation={active}
                        onSend={sendMessage}
                        onBack={() => setShowConversationOnMobile(false)}
                    />
                </div>
            </main>
        </>
    )
}