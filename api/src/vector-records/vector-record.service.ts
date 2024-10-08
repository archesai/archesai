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

  async query(
    orgname: string,
    embedding: number[],
    topK: number,
    contentIds?: string[]
  ) {
    return this.vectorRecordRepository.query(
      orgname,
      embedding,
      topK,
      contentIds
    );
  }

  async remove(orgname: string, id: string) {
    return this.vectorRecordRepository.remove(orgname, id);
  }

  async removeMany(orgname: string, ids: string[]) {
    return this.vectorRecordRepository.removeMany(orgname, ids);
  }

  async upsert(
    orgname: string,
    contentId: string,
    records: {
      embedding: number[];
      text: string;
    }[]
  ) {
    return this.vectorRecordRepository.upsert(orgname, contentId, records);
  }
}
