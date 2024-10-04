import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { ARTokensModule } from "../ar-tokens/ar-tokens.module";
import { AuthModule } from "../auth/auth.module";
import { EmailModule } from "../email/email.module";
import { UsersModule } from "../users/users.module";
import { EmailChangeController } from "./email-change.controller";
import { EmailChangeService } from "./email-change.service";

@Module({
  controllers: [EmailChangeController],
  imports: [EmailModule, UsersModule, ConfigModule, ARTokensModule, AuthModule],
  providers: [EmailChangeService],
})
export class EmailChangeModule {}
