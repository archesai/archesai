export interface CrudTestCases<CreateDto, UpdateDto> {
  create?: {
    accessToken: string
    createDto: CreateDto
    expectedStatus: number
    name: string
    organizationId: string
  }[]
  delete?: {
    accessToken: string
    createDto: CreateDto
    expectedStatus: number
    name: string
    organizationId: string
  }[]
  findMany?: {
    accessToken: string
    createDtos: CreateDto[]
    expectedStatus: number
    name: string
    organizationId: string
    query: Record<string, unknown>
  }[]
  findOne?: {
    accessToken: string
    createDto: CreateDto
    expectedStatus: number
    name: string
    organizationId: string
    query: Record<string, unknown>
  }[]
  update?: {
    accessToken: string
    createDto: CreateDto
    expectedStatus: number
    name: string
    organizationId: string
    updateDto: UpdateDto
  }[]
}

export interface HttpOperations<Entity, CreateDto, UpdateDto> {
  create?: (
    accessToken: string,
    organizationId: string,
    createDto: CreateDto
  ) => Promise<{ body: Entity; status: number }>
  delete?: (
    accessToken: string,
    organizationId: string,
    id: string
  ) => Promise<{ body: { message: string }; status: number }>
  findMany?: (
    accessToken: string,
    organizationId: string,
    query: Record<string, unknown>
  ) => Promise<{ body: object; status: number }>
  findOne?: (
    accessToken: string,
    organizationId: string,
    id: string
  ) => Promise<{ body: Entity; status: number }>
  update?: (
    accessToken: string,
    organizationId: string,
    id: string,
    updateDto: UpdateDto
  ) => Promise<{ body: Entity; status: number }>
}
