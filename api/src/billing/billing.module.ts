import { forwardRef, Module } from "@nestjs/common";

import { OrganizationsModule } from "../organizations/organizations.module";
import { BillingController } from "./billing.controller";
import { BillingService } from "./billing.service";

@Module({
  controllers: [BillingController],
  exports: [BillingService],
  imports: [forwardRef(() => OrganizationsModule)],
  providers: [BillingService],
})
export class BillingModule {}
