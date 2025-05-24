// import { getPlans } from '@archesai/client'
// import { Badge } from '@archesai/ui/components/shadcn/badge'
// import { Button } from '@archesai/ui/components/shadcn/button'
// import { Check } from '@archesai/ui/components/custom/icons'

// import {
//   Card,
//   CardContent,
//   CardDescription,
//   CardFooter,
//   CardHeader,
//   CardTitle
// } from '@archesai/ui/components/shadcn/card'

// const pricingList = [
//   {
//     benefitList: [
//       '1 Team member',
//       '2 GB Storage',
//       'Upto 4 pages',
//       'Community support',
//       'lorem ipsum dolor'
//     ]
//   },
//   {
//     benefitList: [
//       '4 Team member',
//       '4 GB Storage',
//       'Upto 6 pages',
//       'Priority support',
//       'lorem ipsum dolor'
//     ]
//   },
//   {
//     benefitList: [
//       '10 Team member',
//       '8 GB Storage',
//       'Upto 10 pages',
//       'Priority support',
//       'lorem ipsum dolor'
//     ]
//   }
// ]

// export const Pricing = async () => {
//   const response = await getPlans()
//   return (
//     <section
//       className='container py-24 sm:py-32'
//       id='pricing'
//     >
//       <h2 className='text-center text-3xl font-bold md:text-4xl'>
//         Get
//         <span className='from-primary/60 to-primary bg-gradient-to-b bg-clip-text text-transparent'>
//           {' '}
//           Unlimited{' '}
//         </span>
//         Access
//       </h2>
//       <h3 className='text-muted-foreground pb-8 pt-4 text-center text-xl'>
//         Lorem ipsum dolor sit amet consectetur adipisicing elit. Alias
//         reiciendis.
//       </h3>
//       <div className='grid gap-8 md:grid-cols-2 lg:grid-cols-3'>
//         {response.data.data.map((plan, i) => (
//           <Card
//             className={
//               plan.attributes.metadata.key === 'STANDARD'
//                 ? 'shadow-black/10 drop-shadow-xl dark:shadow-white/10'
//                 : ''
//             }
//             key={plan.attributes.name}
//           >
//             <CardHeader>
//               <CardTitle className='item-center flex justify-between'>
//                 {plan.attributes.name}
//                 {plan.attributes.metadata.key === 'STANDARD' ? (
//                   <Badge
//                     className='text-primary text-sm'
//                     variant='secondary'
//                   >
//                     Most popular
//                   </Badge>
//                 ) : null}
//               </CardTitle>
//               <div>
//                 <span className='text-3xl font-bold'>
//                   ${plan.attributes.unitAmount! / 100}
//                 </span>
//                 <span className='text-muted-foreground'> /month</span>
//               </div>

//               <CardDescription>{plan.attributes.description}</CardDescription>
//             </CardHeader>

//             <CardContent>
//               <Button className='w-full'>Choose plan</Button>
//             </CardContent>

//             <hr className='m-auto mb-4 w-4/5' />

//             <CardFooter className='flex'>
//               <div className='space-y-4'>
//                 {pricingList[i]?.benefitList.map((benefit: string) => (
//                   <span
//                     className='flex'
//                     key={benefit}
//                   >
//                     <Check className='text-green-500' />{' '}
//                     <h3 className='ml-2'>{benefit}</h3>
//                   </span>
//                 ))}
//               </div>
//             </CardFooter>
//           </Card>
//         ))}
//       </div>
//     </section>
//   )
// }
