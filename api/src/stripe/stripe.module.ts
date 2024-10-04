import { forwardRef, Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { OrganizationsModule } from "../organizations/organizations.module";
import { StripeController } from "./stripe.controller";
import { StripeService } from "./stripe.service";

@Module({
  controllers: [StripeController],
  exports: [StripeService],
  imports: [ConfigModule, forwardRef(() => OrganizationsModule), ConfigModule],
  providers: [StripeService],
})
export class StripeModule {}
