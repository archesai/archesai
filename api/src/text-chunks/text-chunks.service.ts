import { Injectable, Logger } from "@nestjs/common";
import { TextChunk } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { TextChunkQueryDto } from "./dto/text-chunk-query.dto";
import { TextChunkRepository } from "./text-chunk.repository";

@Injectable()
export class TextChunksService
  implements BaseService<TextChunk, undefined, TextChunkQueryDto, undefined>
{
  private logger = new Logger(TextChunksService.name);
  constructor(private textChunkRepository: TextChunkRepository) {}

  async findAll(
    orgname: string,
    textChunkQueryDto: TextChunkQueryDto,
    contentId?: string
  ) {
    return this.textChunkRepository.findAll(
      orgname,
      textChunkQueryDto,
      contentId
    );
  }

  async findOne(id: string) {
    return this.textChunkRepository.findOne(id);
  }

  async query(
    orgname: string,
    embedding: number[],
    topK: number,
    contentIds?: string[]
  ) {
    return this.textChunkRepository.query(orgname, embedding, topK, contentIds);
  }

  async remove(orgname: string, id: string) {
    return this.textChunkRepository.remove(orgname, id);
  }

  async removeMany(orgname: string, ids: string[]) {
    return this.textChunkRepository.removeMany(orgname, ids);
  }

  async upsertTextChunks(
    orgname: string,
    contentId: string,
    records: {
      text: string;
    }[]
  ): Promise<void> {
    return this.textChunkRepository.upsertTextChunks(
      orgname,
      contentId,
      records
    );
  }

  async upsertVectors(
    orgname: string,
    contentId: string,
    records: {
      embedding: number[];
      textChunkId: string;
    }[]
  ): Promise<void> {
    return this.textChunkRepository.upsertVectors(orgname, contentId, records);
  }
}
