import { Controller, Get } from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'
import { BillingService } from '@/src/billing/billing.service'
import { PlanEntity } from '@/src/billing/entities/plan.entity'

@ApiTags('Billing - Plans')
@Controller('billing/plans')
export class PlansController {
  constructor(private billingService: BillingService) {}

  /**
   * Get plans
   * @remarks This endpoint will return a list of available billing plans
   */
  @Get()
  async findAll(): Promise<PlanEntity[]> {
    const plans = await this.billingService.listPlans()
    return plans.map((plan) => new PlanEntity(plan))
  }
}
