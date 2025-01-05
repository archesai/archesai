import { BaseRepository, IPrismaDelegate } from './base.repository'
import { PaginatedDto } from './dto/paginated.dto'
import { SearchQueryDto } from './dto/search-query.dto'

export abstract class BaseService<
  TDto,
  TModel,
  TRepo extends BaseRepository<IPrismaDelegate<TModel>, any, TModel>
> {
  constructor(protected readonly repository: TRepo) {}

  async create(data: Partial<TDto>): Promise<TDto> {
    const model = await this.repository.create(data)
    const entity = this.toEntity(model)
    this.emitMutationEvent(entity)
    return entity
  }

  async findAll(queryDto: SearchQueryDto): Promise<PaginatedDto<TDto>> {
    const { count, results } = await this.repository.findAll(queryDto)
    const entities = results.map((result) => this.toEntity(result))
    const paginatedEntity = new PaginatedDto<TDto>({
      aggregates: [],
      metadata: {
        limit: queryDto.limit!,
        offset: queryDto.offset!,
        totalResults: count
      },
      results: entities
    })
    return paginatedEntity
  }

  async findOne(id: string): Promise<TDto> {
    const model = await this.repository.findOne(id)
    const entity = this.toEntity(model)
    return entity
  }

  async remove(id: string): Promise<TDto> {
    const deletedModel = await this.repository.remove(id)
    const deletedEntity = this.toEntity(deletedModel)
    this.emitMutationEvent(deletedEntity)
    return deletedEntity
  }

  async update(id: string, data: Partial<TDto>): Promise<TDto> {
    const updated = await this.repository.update(id, data)
    const entity = this.toEntity(updated)
    this.emitMutationEvent(entity)
    return entity
  }

  /**
   * Convert the raw repository result to an TDto.
   * Override this method in the concrete service if necessary.
   */
  protected abstract emitMutationEvent(entity: TDto): void

  /**
   * Convert the raw repository result to an TDto.
   * Override this method in the concrete service if necessary.
   */
  protected abstract toEntity(model: TModel): TDto
}
