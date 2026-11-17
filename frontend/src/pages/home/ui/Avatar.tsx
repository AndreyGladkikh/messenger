import { useState } from 'react'
import { cn } from '@/shared/lib'

type AvatarProps = {
  initials: string
  src?: string
  alt?: string
  online?: boolean
  className?: string
  size?: 'sm' | 'md' | 'lg'
}

const sizes = {
  sm: 'size-8 text-[11px]',
  md: 'size-10 text-xs',
  lg: 'size-11 text-sm',
}

export function Avatar({ initials, src, alt = '', online = false, className, size = 'md' }: AvatarProps) {
  const [failedSource, setFailedSource] = useState<string>()

  return (
    <span className={cn('relative inline-flex shrink-0', className)}>
      <span
        className={cn(
          'relative flex items-center justify-center overflow-hidden rounded-full bg-secondary font-medium tracking-wide text-secondary-foreground ring-1 ring-border ring-inset',
          sizes[size],
        )}
      >
        {initials}
        {src && src !== failedSource && (
          <img
            src={src}
            alt={alt}
            className="absolute inset-0 size-full object-cover"
            onError={() => setFailedSource(src)}
          />
        )}
      </span>
      {online && (
        <span className="absolute -right-0.5 -bottom-0.5 flex size-3 items-center justify-center rounded-full bg-background">
          <span className="size-2 rounded-full bg-primary" />
          <span className="sr-only">в сети</span>
        </span>
      )}
    </span>
  )
}
