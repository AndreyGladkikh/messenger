import { Link, useNavigate } from 'react-router'
import { CredentialsForm } from '@/features/auth'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/shared/ui/card'

export function SignInPage() {
  const navigate = useNavigate()

  return (
    <main className="not-typeset flex min-h-svh flex-col items-center justify-center bg-muted/40 px-4 py-10 font-[family-name:var(--font-geist)]">
      <title>Sign in · Messenger</title>
      <div className="flex w-full max-w-sm flex-col gap-7">
        <div className="flex items-center justify-center gap-2.5 font-medium">
          <span aria-hidden="true" className="flex size-9 items-center justify-center rounded-xl bg-primary font-mono text-primary-foreground">m</span>
          Messenger
        </div>
        <Card className="gap-6 pt-6">
          <CardHeader className="gap-2 px-6">
            <CardTitle><h1 className="text-2xl font-semibold tracking-tight">Welcome back.</h1></CardTitle>
            <CardDescription>Sign in to pick up the conversation.</CardDescription>
          </CardHeader>
          <CardContent className="px-6">
            <CredentialsForm mode="sign-in" onSuccess={() => navigate('/', { replace: true })} />
          </CardContent>
          <CardFooter className="justify-center gap-1.5 py-5">
            <span className="text-muted-foreground">New to Messenger?</span>
            <Link className="rounded-sm font-medium underline-offset-4 hover:underline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-ring" to="/auth/register">Sign up</Link>
          </CardFooter>
        </Card>
      </div>
    </main>
  )
}
