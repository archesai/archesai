import { forwardRef, Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { BillingModule } from "../billing/billing.module";
import { PipelinesModule } from "../pipelines/pipelines.module";
import { PrismaModule } from "../prisma/prisma.module";
import { ToolsModule } from "../tools/tools.module";
import { OrganizationRepository } from "./organization.repository";
import { OrganizationsController } from "./organizations.controller";
import { OrganizationsService } from "./organizations.service";

@Module({
  controllers: [OrganizationsController],
  exports: [OrganizationsService],
  imports: [
    forwardRef(() => BillingModule),
    PrismaModule,
    ConfigModule,
    ToolsModule,
    PipelinesModule,
  ],
  providers: [OrganizationsService, OrganizationRepository],
})
export class OrganizationsModule {}
