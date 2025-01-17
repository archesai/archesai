import { ConfigService } from '@/src/config/config.service'
import { Injectable, OnModuleInit } from '@nestjs/common'
import { PrismaClient } from '@prisma/client'

@Injectable()
export class PrismaService extends PrismaClient implements OnModuleInit {
  constructor(private readonly configService: ConfigService) {
    super({
      datasourceUrl: configService.get('database.url')
    })
  }

  async onModuleInit() {
    await this.$connect()
  }
}
