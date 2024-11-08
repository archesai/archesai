import { BaseRepository } from "./base.repository";
import { PaginatedDto } from "./dto/paginated.dto";
import { SearchQueryDto } from "./dto/search-query.dto";

export abstract class BaseService<
  Entity,
  CreateDto,
  UpdateDto,
  Repo extends BaseRepository<
    PrismaModel,
    CreateDto,
    UpdateDto,
    Include,
    Select,
    RawUpdateInput
  >,
  PrismaModel,
  Include = any,
  Select = any,
  RawUpdateInput = any,
> {
  constructor(protected readonly repository: Repo) {}

  async create(
    orgname: string,
    createDto: CreateDto,
    additionalData?: object
  ): Promise<Entity> {
    return this.toEntity(
      await this.repository.create(orgname, createDto, additionalData)
    );
  }

  async findAll(
    orgname: string,
    queryDto: SearchQueryDto
  ): Promise<PaginatedDto<Entity>> {
    const { count, results } = await this.repository.findAll(orgname, queryDto);
    const entities = results.map((result) => this.toEntity(result));
    return new PaginatedDto<Entity>({
      metadata: {
        limit: queryDto.limit,
        offset: queryDto.offset,
        totalResults: count,
      },
      results: entities,
    });
  }

  async findOne(orgname: string, id: string): Promise<Entity> {
    const result = await this.repository.findOne(orgname, id);
    return this.toEntity(result);
  }

  async remove(orgname: string, id: string): Promise<void> {
    return this.repository.remove(orgname, id);
  }

  async update(
    orgname: string,
    id: string,
    updateDto: UpdateDto
  ): Promise<Entity> {
    const updated = await this.repository.update(orgname, id, updateDto);
    return this.toEntity(updated);
  }

  async updateRaw(
    orgname: string,
    id: string,
    data: RawUpdateInput,
    options?: { include?: Include; select?: Select }
  ): Promise<Entity> {
    return this.toEntity(
      await this.repository.updateRaw(orgname, id, data, options)
    );
  }

  /**
   * Convert the raw repository result to an Entity.
   * Override this method in the concrete service if necessary.
   */
  protected abstract toEntity(model: PrismaModel): Entity;
}
