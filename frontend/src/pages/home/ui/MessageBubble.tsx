import { Check, CheckCheck, FileText, Play } from 'lucide-react'
import { Avatar } from './Avatar'
import type { Message } from '@/shared/lib/messenger-data'
import { cn } from '@/shared/lib'

function initialsFrom(name: string) {
  return name
    .split(' ')
    .slice(0, 2)
    .map((part) => part[0])
    .join('')
    .toUpperCase()
}

type MessageBubbleProps = {
  message: Message
  showAuthor: boolean
  isGroup: boolean
}

export function MessageBubble({ message, showAuthor, isGroup }: MessageBubbleProps) {
  const { outgoing, attachment } = message

  return (
    <div className={cn('flex items-end gap-2', outgoing ? 'justify-end' : 'justify-start')}>
      {!outgoing && isGroup && (
        <span className={cn('mb-1', showAuthor ? 'opacity-100' : 'invisible')}>
          <Avatar initials={initialsFrom(message.authorId)} size="sm" />
        </span>
      )}

      <div
        className={cn(
          'flex max-w-[min(34rem,78%)] flex-col gap-1.5 px-3.5 py-2.5 text-sm leading-relaxed',
          outgoing
            ? 'rounded-2xl rounded-br-md bg-primary text-primary-foreground'
            : 'rounded-2xl rounded-bl-md bg-card text-card-foreground ring-1 ring-border',
        )}
      >
        {!outgoing && isGroup && showAuthor && (
          <span className="text-xs font-medium text-primary">{message.authorId}</span>
        )}

        {message.text && <p className="text-pretty">{message.text}</p>}

        {attachment?.kind === 'file' && (
          <div
            className={cn(
              'flex items-center gap-2.5 rounded-xl px-2.5 py-2',
              outgoing ? 'bg-primary-foreground/12' : 'bg-secondary',
            )}
          >
            <FileText className="size-4 shrink-0" strokeWidth={1.75} aria-hidden="true" />
            <span className="flex min-w-0 flex-col">
              <span className="truncate text-xs font-medium">{attachment.name}</span>
              <span
                className={cn(
                  'font-mono text-[11px]',
                  outgoing ? 'text-primary-foreground/70' : 'text-muted-foreground',
                )}
              >
                {attachment.meta}
              </span>
            </span>
          </div>
        )}

        {attachment?.kind === 'voice' && (
          <div className="flex items-center gap-2.5">
            <span
              className={cn(
                'flex size-8 shrink-0 items-center justify-center rounded-full',
                outgoing ? 'bg-primary-foreground/15' : 'bg-secondary',
              )}
            >
              <Play className="size-3.5" strokeWidth={2} aria-hidden="true" />
            </span>
            <span className="flex items-end gap-0.5" aria-hidden="true">
              {[8, 14, 20, 11, 24, 16, 9, 18, 22, 12, 7, 15, 19, 10].map((height, index) => (
                <span
                  key={index}
                  style={{ height }}
                  className={cn(
                    'w-0.5 rounded-full',
                    outgoing ? 'bg-primary-foreground/50' : 'bg-muted-foreground/45',
                  )}
                />
              ))}
            </span>
            <span
              className={cn(
                'font-mono text-[11px]',
                outgoing ? 'text-primary-foreground/70' : 'text-muted-foreground',
              )}
            >
              {attachment.meta}
            </span>
          </div>
        )}

        <span
          className={cn(
            'flex items-center gap-1 self-end font-mono text-[10px]',
            outgoing ? 'text-primary-foreground/70' : 'text-muted-foreground',
          )}
        >
          {message.time}
          {outgoing &&
            (message.status === 'read' || message.status === 'delivered' ? (
              <CheckCheck
                className={cn('size-3.5', message.status === 'read' && 'text-primary-foreground')}
                strokeWidth={2}
                aria-hidden="true"
              />
            ) : (
              <Check className="size-3.5" strokeWidth={2} aria-hidden="true" />
            ))}
          {outgoing && (
            <span className="sr-only">
              {message.status === 'read'
                ? 'Прочитано'
                : message.status === 'delivered'
                  ? 'Доставлено'
                  : 'Отправлено'}
            </span>
          )}
        </span>
      </div>
    </div>
  )
}
