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
    Select
  >,
  PrismaModel,
  Include,
  Select,
> {
  constructor(protected readonly repository: Repo) {}

  async create(
    orgname: string,
    createDto: CreateDto,
    additionalData?: object,
    options?: { include?: Include; select?: Select }
  ): Promise<Entity> {
    return this.toEntity(
      await this.repository.create(orgname, createDto, additionalData, options)
    );
  }

  async findAll(
    orgname: string,
    queryDto: SearchQueryDto,
    options?: { include?: Include; select?: Select }
  ): Promise<PaginatedDto<Entity>> {
    const { count, results } = await this.repository.findAll(
      orgname,
      queryDto,
      options
    );
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

  async findOne(
    orgname: string,
    id: string,
    options?: { include?: Include; select?: Select }
  ): Promise<Entity> {
    const result = await this.repository.findOne(orgname, id, options);
    return this.toEntity(result);
  }

  async remove(
    orgname: string,
    id: string,
    options?: { include?: Include; select?: Select }
  ): Promise<void> {
    return this.repository.remove(orgname, id, options);
  }

  async update(
    orgname: string,
    id: string,
    updateDto: UpdateDto,
    options?: { include?: Include; select?: Select }
  ): Promise<Entity> {
    const updated = await this.repository.update(
      orgname,
      id,
      updateDto,
      options
    );
    return this.toEntity(updated);
  }

  /**
   * Convert the raw repository result to an Entity.
   * Override this method in the concrete service if necessary.
   */
  protected abstract toEntity(model: PrismaModel): Entity;
}
