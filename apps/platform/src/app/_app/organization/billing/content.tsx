import { useState } from 'react'

import {
  useCancelSubscription,
  useCreateCheckoutSession,
  useGetOneOrganizationSuspense,
  useGetPlansSuspense,
  useUpdateSubscription
} from '@archesai/client'
import { Loader2 } from '@archesai/ui/components/custom/icons'
import { Button } from '@archesai/ui/components/shadcn/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@archesai/ui/components/shadcn/card'
import { toast } from '@archesai/ui/components/shadcn/sonner'
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow
} from '@archesai/ui/components/shadcn/table'

export default function BillingPageContent() {
  const defaultOrgname = 'Arches Platform'
  const [clickedButtonIndex, setClickedButtonIndex] = useState<null | number>(
    -1
  )

  const { data: plans } = useGetPlansSuspense()

  const { data: organizationResponse } =
    useGetOneOrganizationSuspense(defaultOrgname)

  const {
    isPending: createCheckoutSessionLoading,
    mutateAsync: createCheckoutSesseion
  } = useCreateCheckoutSession({
    mutation: {
      onError: (error) => {
        console.error('Error creating checkout session:', error)
        toast('Error', {
          description: 'Could not create checkout session'
        })
      },
      onSuccess: () => {
        toast('Checkout session created', {
          description: 'The checkout session has been successfully created.'
        })
      }
    }
  })
  const {
    isPending: switchSubscriptionLoading,
    mutateAsync: switchSubscriptionPlan
  } = useUpdateSubscription()
  const {
    isPending: cancelSubscriptionLoading,
    mutateAsync: cancelSubscription
  } = useCancelSubscription()

  const organization = organizationResponse.data

  return (
    <div className='flex flex-col gap-3'>
      {/* New Card for Available Plans */}
      <Card>
        <CardHeader>
          <CardTitle>Available Plans</CardTitle>
          <CardDescription>
            Subscribe to a plan to unlock additional features.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <>
            {plans.length > 0 ?
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableCell>Plan Name</TableCell>
                    <TableCell>Description</TableCell>
                    <TableCell>Price</TableCell>
                    <TableCell>Interval</TableCell>
                    <TableCell>Actions</TableCell>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {plans.toReversed().map((plan) => (
                    <TableRow key={plan.id}>
                      <TableCell>{plan.name}</TableCell>
                      <TableCell>{plan.description ?? '-'}</TableCell>
                      <TableCell>
                        {plan.unitAmount ?
                          `$${(plan.unitAmount / 100).toFixed(2)} ${plan.currency.toUpperCase()}`
                        : 'Free'}
                      </TableCell>
                      <TableCell>
                        {plan.recurring ?
                          `${plan.recurring.interval_count.toString()} ${plan.recurring.interval}(s)`
                        : 'One-time'}
                      </TableCell>
                      <TableCell>
                        {organization.attributes.plan === plan.metadata.key ?
                          <Button
                            className='flex gap-2'
                            disabled={
                              clickedButtonIndex === plans.indexOf(plan) &&
                              cancelSubscriptionLoading
                            }
                            onClick={async () => {
                              setClickedButtonIndex(plans.indexOf(plan))
                              await cancelSubscription({
                                id: defaultOrgname
                              })
                              toast('Success', {
                                description: 'Plan canceled successfully.'
                              })
                            }}
                            size='sm'
                            variant='destructive'
                          >
                            {clickedButtonIndex === plans.indexOf(plan) &&
                              cancelSubscriptionLoading && (
                                <Loader2 className='h-5 w-5 animate-spin' />
                              )}
                            <span>Cancel Plan</span>
                          </Button>
                        : organization.attributes.plan === 'FREE' ?
                          <Button
                            className='flex gap-2'
                            disabled={
                              clickedButtonIndex === plans.indexOf(plan) &&
                              createCheckoutSessionLoading
                            }
                            onClick={async () => {
                              const data = await createCheckoutSesseion({
                                data: {
                                  priceId: plan.id
                                }
                              })
                              window.location.href = data.url
                            }}
                            size='sm'
                          >
                            {clickedButtonIndex === plans.indexOf(plan) &&
                              createCheckoutSessionLoading && (
                                <Loader2 className='h-5 w-5 animate-spin' />
                              )}
                            <span>Subscribe</span>
                          </Button>
                        : <Button
                            className='flex gap-2'
                            disabled={
                              clickedButtonIndex === plans.indexOf(plan) &&
                              switchSubscriptionLoading
                            }
                            onClick={async () => {
                              setClickedButtonIndex(plans.indexOf(plan))
                              await switchSubscriptionPlan(
                                {
                                  data: {
                                    planId: plan.id
                                  },
                                  id: defaultOrgname
                                },
                                {
                                  onError: (error) => {
                                    console.error(
                                      'Error switching subscription plan:',
                                      error
                                    )
                                    toast('Could not switch plan')
                                  },
                                  onSuccess: () => {
                                    toast('Success', {
                                      description: 'Plan switched successfully.'
                                    })
                                  }
                                }
                              )
                            }}
                            size='sm'
                          >
                            {clickedButtonIndex === plans.indexOf(plan) &&
                              switchSubscriptionLoading && (
                                <Loader2 className='h-5 w-5 animate-spin' />
                              )}
                            <span>Subscribe</span>
                          </Button>
                        }
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            : <p>No plans available.</p>}
          </>
        </CardContent>
      </Card>
    </div>
  )
}
