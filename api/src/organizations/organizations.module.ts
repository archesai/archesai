import { forwardRef, Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { BillingModule } from "../billing/billing.module";
import { PrismaModule } from "../prisma/prisma.module";
import { OrganizationsController } from "./organizations.controller";
import { OrganizationsService } from "./organizations.service";

@Module({
  controllers: [OrganizationsController],
  exports: [OrganizationsService],
  imports: [forwardRef(() => BillingModule), PrismaModule, ConfigModule],
  providers: [OrganizationsService],
})
export class OrganizationsModule {}
