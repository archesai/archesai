import {
  FieldFilter,
  OperatorEnum,
  SearchQueryDto
} from './dto/search-query.dto'

type PrismaDelegate<PrismaModel, Include> = {
  count: (args: { where: any }) => Promise<number>
  create: (args: { data: any; include?: Include }) => Promise<any & PrismaModel>
  delete: (args: {
    include?: Include
    where: any
  }) => Promise<any & PrismaModel>
  findMany: (args: {
    include?: Include
    orderBy: any
    skip: number
    take: number
    where: any
  }) => Promise<(any & PrismaModel)[]>
  findUniqueOrThrow: (args: {
    include?: Include
    where: any
  }) => Promise<any & PrismaModel>
  update: (args: {
    data: any
    include?: Include
    where: any
  }) => Promise<any & PrismaModel>
}

export abstract class BaseRepository<
  PrismaModel,
  CreateDto,
  UpdateDto,
  Include,
  RawUpdateInput
> {
  constructor(
    protected readonly delegate: PrismaDelegate<PrismaModel, Include>,
    private readonly defaultInclude?: Include
  ) {}

  async create(
    orgname: string,
    createDto: CreateDto,
    additionalData?: object
  ): Promise<PrismaModel> {
    return this.delegate.create({
      data: {
        organization: { connect: { orgname } },
        ...createDto,
        ...additionalData
      },
      include: this.defaultInclude
    })
  }

  async findAll(
    orgname: string,
    queryDto: SearchQueryDto
  ): Promise<{ count: number; results: PrismaModel[] }> {
    console.log(queryDto)
    const whereConditions: any = {
      createdAt: {
        gte: queryDto.startDate,
        lte: queryDto.endDate
      },
      orgname
    }

    if (queryDto.filters) {
      queryDto.filters.forEach((filter: FieldFilter) => {
        if (filter.operator)
          if (
            // If this is a relation filter
            [OperatorEnum.EVERY, OperatorEnum.NONE, OperatorEnum.SOME].includes(
              filter.operator
            )
          ) {
            whereConditions[filter.field] = {
              [filter.operator]: {
                id: { equals: filter.value }
              }
            }
          } else {
            // This is a filter on a scalar field
            whereConditions[filter.field] = { [filter.operator]: filter.value }
          }
      })
    }

    const count = await this.delegate.count({ where: whereConditions })
    const results = await this.delegate.findMany({
      include: this.defaultInclude,
      orderBy: {
        [queryDto.sortBy]: queryDto.sortDirection
      },
      skip: queryDto.offset,
      take: queryDto.limit,
      where: whereConditions
    })

    return { count, results }
  }

  async findOne(orgname: string, id: string): Promise<PrismaModel> {
    return this.delegate.findUniqueOrThrow({
      include: this.defaultInclude,
      where: { id }
    })
  }

  async remove(orgname: string, id: string): Promise<void> {
    await this.delegate.delete({
      include: this.defaultInclude,
      where: { id }
    })
  }

  async update(
    orgname: string,
    id: string,
    updateDto: UpdateDto
  ): Promise<PrismaModel> {
    return this.delegate.update({
      data: updateDto,
      include: this.defaultInclude,
      where: { id }
    })
  }

  async updateRaw(
    orgname: string,
    id: string,
    data: RawUpdateInput
  ): Promise<PrismaModel> {
    return this.delegate.update({
      data,
      include: this.defaultInclude,
      where: { id }
    })
  }
}
