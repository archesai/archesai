import { Prisma } from '@prisma/client'
import { OperatorEnum, SearchQueryDto } from './dto/search-query.dto'

export abstract class BaseRepository<
  TDelegate extends IPrismaDelegate<
    TModel,
    TCreateInput,
    TWhereInput,
    TWhereUniqueInput,
    TUpdateInput,
    TInclude
  >,
  TInclude = Prisma.Args<TDelegate, 'findFirst'>['include'],
  TModel = Prisma.Result<
    TDelegate,
    {
      include: TInclude
    },
    'create'
  >,
  TCreateInput = Prisma.Args<TDelegate, 'create'>['data'],
  TWhereInput = Prisma.Args<TDelegate, 'findFirst'>['where'],
  TWhereUniqueInput = Prisma.Args<TDelegate, 'findUnique'>['where'],
  TUpdateInput = Prisma.Args<TDelegate, 'update'>['data']
> {
  constructor(
    protected readonly delegate: TDelegate,
    protected readonly include?: TInclude
  ) {}

  async create(data: TCreateInput): Promise<TModel> {
    return this.delegate.create({ data, include: this.include })
  }

  async findAll(
    queryDto: SearchQueryDto
  ): Promise<{ count: number; results: TModel[] }> {
    const whereConditions: TWhereInput = this.buildWhereConditions(queryDto)

    const count = await this.delegate.count({ where: whereConditions })
    const results = await this.delegate.findMany({
      where: whereConditions,
      include: this.include,
      skip: queryDto.offset,
      take: queryDto.limit,
      // Optional ordering
      orderBy:
        queryDto.sortBy && queryDto.sortDirection
          ? { [queryDto.sortBy]: queryDto.sortDirection }
          : undefined
    })

    return { count, results }
  }

  async findOne(id: string): Promise<TModel> {
    return this.delegate.findUniqueOrThrow({
      where: { id } as TWhereUniqueInput,
      include: this.include
    })
  }

  async remove(id: string): Promise<TModel> {
    return this.delegate.delete({
      where: { id } as TWhereUniqueInput,
      include: this.include
    })
  }

  async update(id: string, data: TUpdateInput): Promise<TModel> {
    return this.delegate.update({
      where: { id } as TWhereUniqueInput,
      data,
      include: this.include
    })
  }

  async deleteMany(queryDto: SearchQueryDto): Promise<{ count: number }> {
    const whereConditions: TWhereInput = this.buildWhereConditions(queryDto)
    return this.delegate.deleteMany({
      where: whereConditions
    })
  }
  /**
   * Example helper method to build a typed `where` object from the queryDto.
   */
  protected buildWhereConditions(queryDto: SearchQueryDto): TWhereInput {
    const where: any = {
      createdAt: {
        gte: queryDto.startDate,
        lte: queryDto.endDate
      }
    }

    if (queryDto.filters) {
      for (const filter of queryDto.filters) {
        if (filter.operator) {
          if (
            [OperatorEnum.EVERY, OperatorEnum.NONE, OperatorEnum.SOME].includes(
              filter.operator
            )
          ) {
            where[filter.field] = {
              [filter.operator]: {
                id: { equals: filter.value }
              }
            }
          } else {
            where[filter.field] = { [filter.operator]: filter.value }
          }
        }
      }
    }
    return where
  }
}

export interface IPrismaDelegate<
  TModel = any, // The shape of the returned model (e.g., `Content`)
  TCreateInput = any, // The shape for creating (e.g., `ContentCreateInput`)
  TWhereInput = any, // For filtering multiple records
  TWhereUniqueInput = any, // For identifying a single unique record
  TUpdateInput = any, // For updating
  TInclude = any // For including related models
> {
  count(args: { where?: TWhereInput }): Promise<number>

  create(args: {
    data: TCreateInput
    include?: TInclude
  }): Promise<TModel & any>

  delete(args: {
    where: TWhereUniqueInput
    include?: TInclude
  }): Promise<TModel & any>

  findMany(args: {
    where?: TWhereInput
    include?: TInclude
    skip?: number
    take?: number
    orderBy?: Record<string, 'asc' | 'desc'> // if you want orderBy
  }): Promise<(TModel & any)[]>

  findUniqueOrThrow(args: {
    where: TWhereUniqueInput
    include?: TInclude
  }): Promise<TModel & any>

  update(args: {
    data: TUpdateInput
    where: TWhereUniqueInput
    include?: TInclude
  }): Promise<TModel & any>

  deleteMany(args: { where: TWhereInput }): Promise<{
    count: number
  }>
}
