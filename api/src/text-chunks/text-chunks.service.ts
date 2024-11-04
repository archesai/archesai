import { Injectable, Logger } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { PaginatedDto } from "../common/paginated.dto";
import { TextChunkQueryDto } from "./dto/text-chunk-query.dto";
import { TextChunkEntity } from "./entities/text-chunk.entity";
import { TextChunkRepository } from "./text-chunk.repository";

@Injectable()
export class TextChunksService
  implements
    BaseService<TextChunkEntity, undefined, TextChunkQueryDto, undefined>
{
  private logger = new Logger(TextChunksService.name);
  constructor(private textChunkRepository: TextChunkRepository) {}

  async findAll(
    orgname: string,
    textChunkQueryDto: TextChunkQueryDto,
    contentId?: string
  ) {
    const { count, results } = await this.textChunkRepository.findAll(
      orgname,
      textChunkQueryDto,
      contentId
    );
    const textChunkEntities = results.map(
      (textChunk) => new TextChunkEntity(textChunk)
    );
    return new PaginatedDto<TextChunkEntity>({
      metadata: {
        limit: textChunkQueryDto.limit,
        offset: textChunkQueryDto.offset,
        totalResults: count,
      },
      results: textChunkEntities,
    });
  }

  async findOne(id: string) {
    return new TextChunkEntity(await this.textChunkRepository.findOne(id));
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
