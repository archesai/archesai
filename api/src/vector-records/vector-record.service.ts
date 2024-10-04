import { Injectable, Logger } from "@nestjs/common";
import { VectorRecord } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { VectorRecordQueryDto } from "./dto/vector-record-query.dto";
import { VectorRecordRepository } from "./vector-record.repository";

@Injectable()
export class VectorRecordService
  implements
    BaseService<VectorRecord, undefined, VectorRecordQueryDto, undefined>
{
  private logger = new Logger(VectorRecordService.name);
  constructor(private vectorRecordRepository: VectorRecordRepository) {}

  async findAll(
    orgname: string,
    vectorRecordQueryDto: VectorRecordQueryDto,
    contentId?: string
  ) {
    return this.vectorRecordRepository.findAll(
      orgname,
      vectorRecordQueryDto,
      contentId
    );
  }

  async findOne(id: string) {
    return this.vectorRecordRepository.findOne(id);
  }
}
