import { forwardRef, Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { OrganizationsModule } from "../organizations/organizations.module";
import { BillingController } from "./billing.controller";
import { BillingService } from "./billing.service";

@Module({
  controllers: [BillingController],
  exports: [BillingService],
  imports: [ConfigModule, forwardRef(() => OrganizationsModule)],
  providers: [BillingService],
})
export class BillingModule {}
