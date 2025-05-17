export interface CrudTestCases<CreateDto, UpdateDto> {
  create?: {
    accessToken: string
    createDto: CreateDto
    expectedStatus: number
    name: string
    orgname: string
  }[]
  delete?: {
    accessToken: string
    createDto: CreateDto
    expectedStatus: number
    name: string
    orgname: string
  }[]
  findMany?: {
    accessToken: string
    createDtos: CreateDto[]
    expectedStatus: number
    name: string
    orgname: string
    query: Record<string, unknown>
  }[]
  findOne?: {
    accessToken: string
    createDto: CreateDto
    expectedStatus: number
    name: string
    orgname: string
    query: Record<string, unknown>
  }[]
  update?: {
    accessToken: string
    createDto: CreateDto
    expectedStatus: number
    name: string
    orgname: string
    updateDto: UpdateDto
  }[]
}

export interface HttpOperations<Entity, CreateDto, UpdateDto> {
  create?: (
    accessToken: string,
    orgname: string,
    createDto: CreateDto
  ) => Promise<{ body: Entity; status: number }>
  delete?: (
    accessToken: string,
    orgname: string,
    id: string
  ) => Promise<{ body: { message: string }; status: number }>
  findMany?: (
    accessToken: string,
    orgname: string,
    query: Record<string, unknown>
  ) => Promise<{ body: object; status: number }>
  findOne?: (
    accessToken: string,
    orgname: string,
    id: string
  ) => Promise<{ body: Entity; status: number }>
  update?: (
    accessToken: string,
    orgname: string,
    id: string,
    updateDto: UpdateDto
  ) => Promise<{ body: Entity; status: number }>
}
