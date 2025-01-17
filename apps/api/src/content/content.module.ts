import { Module } from '@nestjs/common'

import { PrismaModule } from '../prisma/prisma.module'
import { StorageModule } from '../storage/storage.module'
import { ContentController } from './content.controller'
import { ContentRepository } from './content.repository'
import { ContentService } from './content.service'
import { ScraperModule } from '../scraper/scraper.module'
import { AuthModule } from '@/src/auth/auth.module'

@Module({
  controllers: [ContentController],
  exports: [ContentService],
  imports: [PrismaModule, StorageModule.forRoot(), ScraperModule, AuthModule],
  providers: [ContentService, ContentRepository]
})
export class ContentModule {}
