import { Outlet, useNavigation } from 'react-router'
import { PageLoading } from '@/pages/error'

export function ProtectedRouteBoundary() {
  const navigation = useNavigation()

  return navigation.state === 'idle' ? <Outlet /> : <PageLoading />
}
