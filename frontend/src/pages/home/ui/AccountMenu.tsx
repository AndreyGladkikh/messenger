import { LogOut } from 'lucide-react'
import { useNavigate } from 'react-router'
import { useSignOut } from '@/features/auth'
import { currentUser } from '@/shared/lib/messenger-data'
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
  const navigate = useNavigate()
  const signOut = useSignOut()

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
        <Avatar initials={currentUser.initials} size="sm" online />
      </DropdownMenuTrigger>
      <DropdownMenuContent side="top" align="start" sideOffset={12} className="w-56">
        <DropdownMenuGroup>
          <DropdownMenuLabel>
            <span className="block">{currentUser.name}</span>
            <span className="block">{currentUser.handle}</span>
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
