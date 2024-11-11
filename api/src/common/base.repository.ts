import {
  FieldFieldQuery,
  Operator,
  SearchQueryDto,
} from "./dto/search-query.dto";

export abstract class BaseRepository<
  PrismaModel,
  CreateDto,
  UpdateDto,
  Include,
  Select,
  RawUpdateInput,
> {
  constructor(
    protected readonly delegate: PrismaDelegate<PrismaModel, Include, Select>,
    private readonly defaultInclude?: Include,
    private readonly defaultSelect?: Select
  ) {}

  async create(
    orgname: string,
    createDto: CreateDto,
    additionalData?: object,
    options?: { include?: Include; select?: Select }
  ): Promise<PrismaModel> {
    return this.delegate.create({
      data: {
        organization: { connect: { orgname } },
        ...createDto,
        ...additionalData,
      },
      include: options?.include ?? this.defaultInclude,
      select: options?.select ?? this.defaultSelect,
    });
  }

  async findAll(
    orgname: string,
    queryDto: SearchQueryDto,
    options?: { include?: Include; select?: Select }
  ): Promise<{ count: number; results: PrismaModel[] }> {
    const whereConditions: any = {
      createdAt: {
        gte: queryDto.startDate,
        lte: queryDto.endDate,
      },
      orgname,
    };

    if (queryDto.filters) {
      queryDto.filters.forEach((filter: FieldFieldQuery) => {
        if (filter.operator)
          if (
            // If this is a relation filter
            [Operator.EVERY, Operator.NONE, Operator.SOME].includes(
              filter.operator
            )
          ) {
            whereConditions[filter.field] = {
              [filter.operator]: {
                id: { equals: filter.value },
              },
            };
          } else {
            // This is a filter on a scalar field
            whereConditions[filter.field] = { [filter.operator]: filter.value };
          }
      });
    }

    const count = await this.delegate.count({ where: whereConditions });
    const results = await this.delegate.findMany({
      include: options?.include ?? this.defaultInclude,
      orderBy: {
        [queryDto.sortBy]: queryDto.sortDirection,
      },
      select: options?.select ?? this.defaultSelect,
      skip: queryDto.offset,
      take: queryDto.limit,
      where: whereConditions,
    });

    return { count, results };
  }

  async findOne(
    orgname: string,
    id: string,
    options?: { include?: Include; select?: Select }
  ): Promise<PrismaModel> {
    return this.delegate.findUniqueOrThrow({
      include: options?.include ?? this.defaultInclude,
      select: options?.select ?? this.defaultSelect,
      where: { id },
    });
  }

  async remove(
    orgname: string,
    id: string,
    options?: { include?: Include; select?: Select }
  ): Promise<void> {
    await this.delegate.delete({
      include: options?.include ?? this.defaultInclude,
      select: options?.select ?? this.defaultSelect,
      where: { id },
    });
  }

  async update(
    orgname: string,
    id: string,
    updateDto: UpdateDto,
    options?: { include?: Include; select?: Select }
  ): Promise<PrismaModel> {
    return this.delegate.update({
      data: updateDto,
      include: options?.include ?? this.defaultInclude,
      select: options?.select ?? this.defaultSelect,
      where: { id },
    });
  }

  async updateRaw(
    orgname: string,
    id: string,
    data: RawUpdateInput,
    options?: { include?: Include; select?: Select }
  ): Promise<PrismaModel> {
    return this.delegate.update({
      data,
      include: options?.include ?? this.defaultInclude,
      select: options?.select ?? this.defaultSelect,
      where: { id },
    });
  }
}

type PrismaDelegate<PrismaModel, Include, Select> = {
  count: (args: { where: any }) => Promise<number>;
  create: (args: {
    data: any;
    include?: Include;
    select?: Select;
  }) => Promise<any & PrismaModel>;
  delete: (args: {
    include?: Include;
    select?: Select;
    where: any;
  }) => Promise<any & PrismaModel>;
  findMany: (args: {
    include?: Include;
    orderBy: any;
    select?: Select;
    skip: number;
    take: number;
    where: any;
  }) => Promise<(any & PrismaModel)[]>;
  findUniqueOrThrow: (args: {
    include?: Include;
    select?: Select;
    where: any;
  }) => Promise<any & PrismaModel>;
  update: (args: {
    data: any;
    include?: Include;
    select?: Select;
    where: any;
  }) => Promise<any & PrismaModel>;
};
