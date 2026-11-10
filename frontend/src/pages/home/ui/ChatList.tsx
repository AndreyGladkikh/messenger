'use client'

import { BellOff, Pin, Plus, Search } from 'lucide-react'
import { useMemo, useState } from 'react'
import { Avatar } from './Avatar'
import type { Conversation } from '@/shared/lib/messenger-data'
import { cn } from '@/shared/lib'

type Filter = 'all' | 'unread' | 'groups'

const filters: { id: Filter; label: string }[] = [
  { id: 'all', label: 'Все' },
  { id: 'unread', label: 'Непрочитанные' },
  { id: 'groups', label: 'Группы' },
]

type ChatListProps = {
  conversations: Conversation[]
  activeId: string
  onSelect: (id: string) => void
}

export function ChatList({ conversations, activeId, onSelect }: ChatListProps) {
  const [query, setQuery] = useState('')
  const [filter, setFilter] = useState<Filter>('all')

  const visible = useMemo(() => {
    const normalized = query.trim().toLowerCase()
    return conversations.filter((chat) => {
      if (filter === 'unread' && chat.unread === 0) return false
      if (filter === 'groups' && chat.kind !== 'group') return false
      if (!normalized) return true
      return (
        chat.name.toLowerCase().includes(normalized) ||
        chat.preview.toLowerCase().includes(normalized)
      )
    })
  }, [conversations, filter, query])

  const totalUnread = conversations.reduce((sum, chat) => sum + chat.unread, 0)

  return (
    <div className="flex h-full w-full flex-col border-r border-border bg-sidebar md:w-80 lg:w-88">
      <header className="flex flex-col gap-3 px-4 pt-4 pb-3">
        <div className="flex items-center justify-between gap-2">
          <h1 className="text-lg font-semibold tracking-tight">Чаты</h1>
          <div className="flex items-center gap-2">
            <span className="rounded-full bg-accent px-2 py-0.5 font-mono text-[11px] text-accent-foreground">
              {totalUnread} новых
            </span>
            <button
              type="button"
              className="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground transition-opacity hover:opacity-90 focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-sidebar focus-visible:outline-none"
            >
              <Plus className="size-4" strokeWidth={2} aria-hidden="true" />
              <span className="sr-only">Новый чат</span>
            </button>
          </div>
        </div>

        <div className="relative">
          <Search
            className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground"
            strokeWidth={1.75}
            aria-hidden="true"
          />
          <input
            type="search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Поиск по чатам"
            aria-label="Поиск по чатам"
            className="h-10 w-full rounded-xl border border-border bg-background pr-3 pl-9 text-sm text-foreground placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/25 focus-visible:outline-none"
          />
        </div>

        <div className="flex items-center gap-1.5" role="tablist" aria-label="Фильтр чатов">
          {filters.map((item) => (
            <button
              key={item.id}
              type="button"
              role="tab"
              aria-selected={filter === item.id}
              onClick={() => setFilter(item.id)}
              className={cn(
                'rounded-lg px-2.5 py-1 text-xs font-medium transition-colors',
                filter === item.id
                  ? 'bg-foreground text-background'
                  : 'text-muted-foreground hover:bg-secondary hover:text-foreground',
              )}
            >
              {item.label}
            </button>
          ))}
        </div>
      </header>

      <div className="scrollbar-slim flex-1 overflow-y-auto px-2 pb-4">
        {visible.length === 0 ? (
          <p className="px-2 py-8 text-center text-sm text-muted-foreground">
            Ничего не найдено
          </p>
        ) : (
          <ul className="flex flex-col gap-0.5">
            {visible.map((chat) => {
              const isActive = chat.id === activeId
              return (
                <li key={chat.id}>
                  <button
                    type="button"
                    onClick={() => onSelect(chat.id)}
                    aria-current={isActive ? 'true' : undefined}
                    className={cn(
                      'flex w-full items-center gap-3 rounded-xl px-2 py-2.5 text-left transition-colors',
                      isActive ? 'bg-card shadow-xs ring-1 ring-border' : 'hover:bg-secondary',
                    )}
                  >
                    <Avatar initials={chat.initials} online={chat.online} />

                    <span className="flex min-w-0 flex-1 flex-col gap-0.5">
                      <span className="flex items-center gap-1.5">
                        <span className="truncate text-sm font-medium text-foreground">
                          {chat.name}
                        </span>
                        {chat.pinned && (
                          <Pin
                            className="size-3 shrink-0 text-muted-foreground"
                            strokeWidth={2}
                            aria-hidden="true"
                          />
                        )}
                        {chat.muted && (
                          <BellOff
                            className="size-3 shrink-0 text-muted-foreground"
                            strokeWidth={2}
                            aria-hidden="true"
                          />
                        )}
                        <span className="ml-auto shrink-0 font-mono text-[11px] text-muted-foreground">
                          {chat.time}
                        </span>
                      </span>

                      <span className="flex items-center gap-2">
                        <span
                          className={cn(
                            'truncate text-xs',
                            chat.typing ? 'text-primary' : 'text-muted-foreground',
                          )}
                        >
                          {chat.typing ? 'печатает…' : chat.preview}
                        </span>
                        {chat.unread > 0 && (
                          <span className="ml-auto flex h-5 min-w-5 shrink-0 items-center justify-center rounded-full bg-primary px-1.5 font-mono text-[11px] font-medium text-primary-foreground">
                            {chat.unread}
                          </span>
                        )}
                      </span>
                    </span>
                  </button>
                </li>
              )
            })}
          </ul>
        )}
      </div>
    </div>
  )
}
