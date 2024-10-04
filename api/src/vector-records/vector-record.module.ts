import { Module } from "@nestjs/common";

import { PrismaModule } from "../prisma/prisma.module";
import { VectorRecordController } from "./vector-record.controller";
import { VectorRecordRepository } from "./vector-record.repository";
import { VectorRecordService } from "./vector-record.service";

@Module({
  controllers: [VectorRecordController],
  exports: [VectorRecordService],
  imports: [PrismaModule],
  providers: [VectorRecordService, VectorRecordRepository],
})
export class VectorRecordModule {}
