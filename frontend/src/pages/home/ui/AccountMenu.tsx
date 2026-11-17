import { useQuery } from '@tanstack/react-query'
import { LogOut } from 'lucide-react'
import { useLoaderData, useNavigate } from 'react-router'
import { useSignOut } from '@/features/auth'
import { profileAvatarQuery, type Profile } from '@/entities/profile'
import { Alert, AlertDescription } from '@/shared/ui/alert'
import { Button } from '@/shared/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/ui/dropdown-menu'
import { Avatar } from './Avatar'

export function AccountMenu() {
  const profile = useLoaderData<Profile>()
  const navigate = useNavigate()
  const signOut = useSignOut()
  const displayName = profile.displayName
  const initials = displayName.trim().split(/\s+/).slice(0, 2)
    .map((word) => Array.from(word)[0]).join('').toUpperCase()
  const avatar = useQuery(profileAvatarQuery(profile.avatarMediaId))

  function handleSignOut() {
    if (signOut.isPending) return

    signOut.mutate(undefined, {
      onSuccess: () => {
        void navigate('/auth/login', { replace: true })
      },
    })
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        render={<Button variant="ghost" size="icon" className="size-10 rounded-full" />}
        aria-label="Account menu"
      >
        <Avatar initials={initials} src={avatar.data} alt={`${displayName} profile photo`} size="sm" online />
      </DropdownMenuTrigger>
      <DropdownMenuContent side="top" align="start" sideOffset={12} className="w-56">
        <DropdownMenuGroup>
          <DropdownMenuLabel>
            <span className="block break-words">{displayName}</span>
            <span className="block break-words">@{profile.handle}</span>
          </DropdownMenuLabel>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem
            variant="destructive"
            closeOnClick={false}
            disabled={signOut.isPending}
            onClick={handleSignOut}
          >
            <LogOut aria-hidden="true" />
            {signOut.isPending ? 'Signing out…' : 'Sign out'}
          </DropdownMenuItem>
        </DropdownMenuGroup>
        {signOut.isError && (
          <Alert variant="destructive" className="mt-1">
            <AlertDescription>Could not sign out. Please try again.</AlertDescription>
          </Alert>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
