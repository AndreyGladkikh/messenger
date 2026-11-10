'use client'

import { Archive, MessagesSquare, Settings, Users } from 'lucide-react'
import { useState } from 'react'
import { Avatar } from './Avatar'
import { currentUser } from '@/shared/lib/messenger-data'
import { cn } from '@/shared/lib'

const items = [
  { id: 'chats', label: 'Чаты', icon: MessagesSquare },
  { id: 'contacts', label: 'Контакты', icon: Users },
  { id: 'archive', label: 'Архив', icon: Archive },
  { id: 'settings', label: 'Настройки', icon: Settings },
]

export function NavRail() {
  const [active, setActive] = useState('chats')

  return (
    <nav
      aria-label="Основная навигация"
      className="hidden w-16 shrink-0 flex-col items-center gap-1 border-r border-border bg-sidebar py-4 md:flex"
    >
      <span className="mb-4 flex size-9 items-center justify-center rounded-xl bg-primary font-mono text-sm font-semibold text-primary-foreground">
        m
      </span>

      {items.map((item) => {
        const isActive = active === item.id
        return (
          <button
            key={item.id}
            type="button"
            onClick={() => setActive(item.id)}
            aria-current={isActive ? 'page' : undefined}
            className={cn(
              'group relative flex size-10 items-center justify-center rounded-xl transition-colors',
              isActive
                ? 'bg-accent text-accent-foreground'
                : 'text-muted-foreground hover:bg-secondary hover:text-foreground',
            )}
          >
            <item.icon className="size-5" strokeWidth={1.75} aria-hidden="true" />
            <span className="sr-only">{item.label}</span>
            {isActive && (
              <span className="absolute -left-2.5 h-5 w-1 rounded-r-full bg-primary" aria-hidden="true" />
            )}
          </button>
        )
      })}

      <div className="mt-auto">
        <Avatar initials={currentUser.initials} size="sm" online />
      </div>
    </nav>
  )
}
