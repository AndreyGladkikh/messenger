'use client'

import {
  ArrowLeft,
  EllipsisVertical,
  Mic,
  Paperclip,
  Phone,
  Search,
  SendHorizontal,
  Smile,
  Video,
} from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { Avatar } from './Avatar'
import { MessageBubble } from './MessageBubble'
import type { Conversation, Message } from '@/shared/lib/messenger-data'

type ConversationViewProps = {
  conversation: Conversation
  onSend: (text: string) => void
  onBack: () => void
}

export function ConversationView({ conversation, onSend, onBack }: ConversationViewProps) {
  const [draft, setDraft] = useState('')
  const scrollRef = useRef<HTMLDivElement>(null)
  const isGroup = conversation.kind === 'group'

  useEffect(() => {
    const node = scrollRef.current
    if (node) node.scrollTop = node.scrollHeight
  }, [conversation.id, conversation.messages.length])

  function submit() {
    const text = draft.trim()
    if (!text) return
    onSend(text)
    setDraft('')
  }

  return (
    <section className="flex h-full min-w-0 flex-1 flex-col bg-background">
      <header className="flex items-center gap-3 border-b border-border px-3 py-3 md:px-5">
        <button
          type="button"
          onClick={onBack}
          className="flex size-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground md:hidden"
        >
          <ArrowLeft className="size-5" strokeWidth={1.75} aria-hidden="true" />
          <span className="sr-only">Назад к списку чатов</span>
        </button>

        <Avatar initials={conversation.initials} online={conversation.online} size="lg" />

        <div className="flex min-w-0 flex-col">
          <h2 className="truncate text-sm font-semibold text-foreground">{conversation.name}</h2>
          <p className="truncate text-xs text-muted-foreground">
            {conversation.typing
              ? 'кто-то печатает…'
              : isGroup
                ? `${conversation.members} участников · ${conversation.lastSeen}`
                : conversation.lastSeen}
          </p>
        </div>

        <div className="ml-auto flex items-center gap-1">
          {[
            { icon: Search, label: 'Поиск в чате' },
            { icon: Phone, label: 'Аудиозвонок' },
            { icon: Video, label: 'Видеозвонок' },
            { icon: EllipsisVertical, label: 'Действия с чатом' },
          ].map((action) => (
            <button
              key={action.label}
              type="button"
              className="flex size-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground"
            >
              <action.icon className="size-[18px]" strokeWidth={1.75} aria-hidden="true" />
              <span className="sr-only">{action.label}</span>
            </button>
          ))}
        </div>
      </header>

      <div
        ref={scrollRef}
        className="scrollbar-slim flex-1 overflow-y-auto px-3 py-5 md:px-6"
        role="log"
        aria-label={`Переписка: ${conversation.name}`}
      >
        <div className="mx-auto flex max-w-3xl flex-col gap-2">
          <div className="mb-2 flex items-center gap-3">
            <span className="h-px flex-1 bg-border" />
            <span className="font-mono text-[11px] tracking-wide text-muted-foreground uppercase">
              сегодня
            </span>
            <span className="h-px flex-1 bg-border" />
          </div>

          {conversation.messages.map((message: Message, index) => {
            const previous = conversation.messages[index - 1]
            const showAuthor = !previous || previous.authorId !== message.authorId
            return (
              <MessageBubble
                key={message.id}
                message={message}
                showAuthor={showAuthor}
                isGroup={isGroup}
              />
            )
          })}

          {conversation.typing && (
            <div className="flex items-center gap-2">
              <span className="flex items-center gap-1 rounded-2xl rounded-bl-md bg-card px-3.5 py-3 ring-1 ring-border">
                {[0, 150, 300].map((delay) => (
                  <span
                    key={delay}
                    style={{ animationDelay: `${delay}ms` }}
                    className="size-1.5 animate-bounce rounded-full bg-muted-foreground"
                  />
                ))}
                <span className="sr-only">Собеседник печатает</span>
              </span>
            </div>
          )}
        </div>
      </div>

      <footer className="border-t border-border px-3 py-3 md:px-6 md:py-4">
        <div className="mx-auto flex max-w-3xl items-end gap-2">
          <button
            type="button"
            className="flex size-10 shrink-0 items-center justify-center rounded-xl text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground"
          >
            <Paperclip className="size-[18px]" strokeWidth={1.75} aria-hidden="true" />
            <span className="sr-only">Прикрепить файл</span>
          </button>

          <div className="flex min-w-0 flex-1 items-end rounded-2xl border border-border bg-card focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/25">
            <textarea
              value={draft}
              onChange={(event) => setDraft(event.target.value)}
              onKeyDown={(event) => {
                if (
                  event.key === 'Enter' &&
                  !event.shiftKey &&
                  !event.nativeEvent.isComposing &&
                  event.keyCode !== 229
                ) {
                  event.preventDefault()
                  submit()
                }
              }}
              rows={1}
              placeholder="Напишите сообщение…"
              aria-label="Текст сообщения"
              className="max-h-32 min-h-10 w-full resize-none bg-transparent px-3.5 py-2.5 text-sm leading-relaxed text-foreground placeholder:text-muted-foreground focus:outline-none"
            />
            <button
              type="button"
              className="mr-1 mb-1 flex size-8 shrink-0 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground"
            >
              <Smile className="size-[18px]" strokeWidth={1.75} aria-hidden="true" />
              <span className="sr-only">Эмодзи</span>
            </button>
          </div>

          {draft.trim() ? (
            <button
              type="button"
              onClick={submit}
              className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-primary text-primary-foreground transition-opacity hover:opacity-90"
            >
              <SendHorizontal className="size-[18px]" strokeWidth={2} aria-hidden="true" />
              <span className="sr-only">Отправить</span>
            </button>
          ) : (
            <button
              type="button"
              className="flex size-10 shrink-0 items-center justify-center rounded-xl text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground"
            >
              <Mic className="size-[18px]" strokeWidth={1.75} aria-hidden="true" />
              <span className="sr-only">Записать голосовое сообщение</span>
            </button>
          )}
        </div>
        <p className="mx-auto mt-2 hidden max-w-3xl font-mono text-[11px] text-muted-foreground md:block">
          Enter — отправить · Shift + Enter — новая строка
        </p>
      </footer>
    </section>
  )
}
