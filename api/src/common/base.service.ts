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
    RawUpdateInput
  >,
  PrismaModel,
  Include = any,
  RawUpdateInput = any,
> {
  constructor(protected readonly repository: Repo) {}

  async create(
    orgname: string,
    createDto: CreateDto,
    additionalData?: object
  ): Promise<Entity> {
    const entity = await this.repository.create(
      orgname,
      createDto,
      additionalData
    );

    this.emitMutationEvent(orgname);
    return this.toEntity(entity);
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
    await this.repository.remove(orgname, id);
    this.emitMutationEvent(orgname);
  }

  async update(
    orgname: string,
    id: string,
    updateDto: UpdateDto
  ): Promise<Entity> {
    const updated = await this.repository.update(orgname, id, updateDto);
    this.emitMutationEvent(orgname);
    return this.toEntity(updated);
  }

  /**
   * Convert the raw repository result to an Entity.
   * Override this method in the concrete service if necessary.
   */
  protected abstract emitMutationEvent(orgname: string): void;

  /**
   * Convert the raw repository result to an Entity.
   * Override this method in the concrete service if necessary.
   */
  protected abstract toEntity(model: PrismaModel): Entity;
}
