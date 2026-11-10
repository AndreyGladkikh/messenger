import { isRouteErrorResponse, useRevalidator, useRouteError } from 'react-router'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'

export function RootErrorBoundary() {
  const error = useRouteError()
  const revalidator = useRevalidator()
  const message = isRouteErrorResponse(error)
    ? `${error.status} ${error.statusText || 'Request failed'}`
    : error instanceof Error
      ? error.message
      : 'An unexpected error occurred.'

  return (
    <main className="grid min-h-screen place-items-center p-6">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Sorry, somethig went wrong</CardTitle>
        </CardHeader>
        <CardContent>
          <div>
            {message}
          </div>
          <Button onClick={() => revalidator.revalidate()}>
            Retry
          </Button>
        </CardContent>
      </Card>
    </main>
  )
}
