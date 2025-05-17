import type { Metadata } from 'next'

import UserForm from '#components/forms/user-form'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/profile/general')

export default function ProfilePage() {
  return <UserForm />
}
