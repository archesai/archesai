import { forwardRef, Module } from '@nestjs/common'

import { OrganizationsModule } from '@/src/organizations/organizations.module'
import { BillingController } from '@/src/billing/controllers/billing.controller'
import { BillingService } from '@/src/billing/billing.service'
import { PaymentMethodsController } from '@/src/billing/controllers/payment-methods.controller'
import { PlansController } from '@/src/billing/controllers/plans.controller'
import { SubscriptionsController } from '@/src/billing/controllers/subscriptions.controller'
import { StripeController } from '@/src/billing/controllers/stripe.controller'

@Module({
  controllers: [
    BillingController,
    PaymentMethodsController,
    PlansController,
    SubscriptionsController,
    StripeController
  ],
  exports: [BillingService],
  imports: [forwardRef(() => OrganizationsModule)],
  providers: [BillingService]
})
export class BillingModule {}
