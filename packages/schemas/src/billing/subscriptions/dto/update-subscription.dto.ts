import type { Static } from '@sinclair/typebox'

import { CreateSubscriptionDtoSchema } from '#billing/subscriptions/dto/create-subscription.dto'

export const UpdateSubscriptionDtoSchema = CreateSubscriptionDtoSchema

export type UpdateSubscriptionDto = Static<typeof UpdateSubscriptionDtoSchema>
