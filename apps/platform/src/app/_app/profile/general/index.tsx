import { createFileRoute } from '@tanstack/react-router'

import UserForm from '#components/forms/user-form'
import { getRouteMeta } from '#lib/site-utils'

export const metadata = getRouteMeta('/profile/general')

export const Route = createFileRoute('/_app/profile/general/')({
  component: ProfilePage
})

export default function ProfilePage() {
  return <UserForm />
}
