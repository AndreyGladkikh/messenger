import { redirect, type LoaderFunctionArgs, type MiddlewareFunction } from 'react-router'
import { currentUserQuery, type Account } from '@/features/auth'
import { profileQuery } from '@/entities/profile'
import { UnauthenticatedError } from '@/shared/api'
import { queryClientContext } from './context'

async function getAccount(context: LoaderFunctionArgs['context']): Promise<Account | null> {
  try {
    return await context.get(queryClientContext).fetchQuery(currentUserQuery())
  } catch (error) {
    if (error instanceof UnauthenticatedError) return null
    throw error
  }
}

async function getProfile(context: LoaderFunctionArgs['context'], account: Account) {
  return context.get(queryClientContext).fetchQuery(profileQuery(account.id))
}

export const protectedMiddleware: MiddlewareFunction = async ({ context }, next) => {
  console.log('protectedMiddleware');
  
  const account = await getAccount(context)
  console.log(account);
  
  if (!account) throw redirect('/auth/login')

  const profile = await getProfile(context, account)
  if (!profile) throw redirect('/onboarding/profile')

  return next()
}

export const publicAuthMiddleware: MiddlewareFunction = async ({ context }, next) => {
  const account = await getAccount(context)
  if (!account) return next()

  const profile = await getProfile(context, account)
  throw redirect(profile ? '/' : '/onboarding/profile')
}

export const onboardingMiddleware: MiddlewareFunction = async ({ context }, next) => {
  const account = await getAccount(context)
  if (!account) throw redirect('/auth/login')

  const profile = await getProfile(context, account)
  if (profile) throw redirect('/')

  return next()
}

export async function accountLoader({ context }: LoaderFunctionArgs): Promise<Account> {
  const account = context.get(queryClientContext).getQueryData(currentUserQuery().queryKey) as Account | undefined
  if (!account) throw redirect('/auth/login')
  return account
}
