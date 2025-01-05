import { atom } from 'jotai'

export type AuthStatus =
  | 'Authenticated'
  | 'Loading'
  | 'Refreshing'
  | 'Unauthenticated'

export const authStatusAtom = atom<AuthStatus>('Loading')

export const defaultOrgnameAtom = atom<string>('')
