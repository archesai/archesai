import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { EmailService } from "./email.service";

@Module({
  exports: [EmailService],
  imports: [ConfigModule],
  providers: [EmailService],
})
export class EmailModule {}
