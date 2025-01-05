import { Module } from '@nestjs/common'

import { EmailService } from './email.service'

@Module({
  exports: [EmailService],
  providers: [EmailService]
})
export class EmailModule {}
