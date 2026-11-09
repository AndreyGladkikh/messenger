import { isRouteErrorResponse, useRevalidator, useRouteError } from 'react-router'
import { Button } from '@/shared/ui/pouf/Button'
import { ErrorNote } from '@/shared/ui/pouf/feedback'
import { Stack } from '@/shared/ui/pouf/layout'
import { Card } from '@/shared/ui/pouf/surface'
import { Heading, Text } from '@/shared/ui/pouf/text'

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
      <div className="w-full max-w-[520px]">
        <Card>
          <Stack gap={4}>
            <Heading level={2}>We couldn’t load this page</Heading>
            <Text size="sm" muted>Messenger remains locked until your account and profile can be verified.</Text>
            <ErrorNote>{message}</ErrorNote>
            <Button
              block
              loading={revalidator.state === 'loading'}
              onClick={() => revalidator.revalidate()}
            >
              Retry
            </Button>
          </Stack>
        </Card>
      </div>
    </main>
  )
}
