import { PaginatedDto } from '@/src/common/dto/paginated.dto'
import { SearchQueryDto } from '@/src/common/dto/search-query.dto'

/**
 * Interface for your CRUD operations.
 * Each function can receive an optional 'input' object
 * if you need to pass different payloads.
 */
export interface CrudOperations<Entity, CreateDto, UpdateDto> {
  create?: (
    accessToken: string,
    orgname: string,
    createDto: CreateDto
  ) => Promise<{ status: number; body: Entity }>
  findOne?: (
    accessToken: string,
    orgname: string,
    id: string
  ) => Promise<{ status: number; body: Entity }>
  update?: (
    accessToken: string,
    orgname: string,
    id: string,
    updateDto: UpdateDto
  ) => Promise<{ status: number; body: Entity }>
  delete?: (
    accessToken: string,
    orgname: string,
    id: string
  ) => Promise<{ status: number; body: { message: string } }>
  findAll?: (
    accessToken: string,
    orgname: string,
    searchQueryDto: SearchQueryDto
  ) => Promise<{ status: number; body: PaginatedDto<Entity> }>
}

/**
 * Scenarios for each operation.
 * For example, 'create' might have multiple variations of data to test.
 */
export interface CrudTestCases<CreateDto, UpdateDto> {
  create?: Array<{
    name: string
    accessToken: string
    orgname: string
    createDto: CreateDto
    expectedStatus: number
  }>
  findOne?: Array<{
    name: string
    accessToken: string
    orgname: string
    createDto: CreateDto
    searchQueryDto: SearchQueryDto
    expectedStatus: number
  }>
  update?: Array<{
    name: string
    accessToken: string
    orgname: string
    createDto: CreateDto
    updateDto: UpdateDto
    expectedStatus: number
  }>
  findAll?: Array<{
    name: string
    accessToken: string
    orgname: string
    createDtos: CreateDto[]
    searchQueryDto: SearchQueryDto
    expectedStatus: number
  }>
  delete?: Array<{
    name: string
    accessToken: string
    orgname: string
    createDto: CreateDto
    expectedStatus: number
  }>
}

/**
 * Main test suite runner. For each CRUD operation (create, read, etc.),
 * we look at the array of scenarios and run 'it.each' on them.
 */
export function runCrudTestSuite<
  Entity extends { id: string },
  CreateDto,
  UpdateDto
>(
  operations: CrudOperations<Entity, CreateDto, UpdateDto>,
  scenarios: CrudTestCases<CreateDto, UpdateDto>
) {
  // ---------------------------------------
  // CREATE
  // ---------------------------------------
  if (scenarios.create && scenarios.create.length > 0) {
    describe(`CREATE`, () => {
      it.each(scenarios.create!)(
        'should create with $name',
        async (scenario) => {
          // Ensure the create operation is defined
          expect(operations.create).toBeDefined()

          // Create the entity
          const created = await operations.create!(
            scenario.accessToken,
            scenario.orgname,
            scenario.createDto
          )

          // Log response for debugging if status is incorrect
          if (created.status !== scenario.expectedStatus) {
            console.error(
              `Unexpected status code for CREATE operation. Response:`,
              (created.body as any).message
            )
          }

          // Validate the response
          expect(created).toBeDefined()
          expect(created.status).toBe(scenario.expectedStatus)
          expect(created).toSatisfyApiSpec()

          // Clean up by deleting the entity if necessary
          if (created.body.id && operations.delete) {
            await operations.delete(
              scenario.accessToken,
              scenario.orgname,
              created.body.id
            )
          }
        }
      )
    })
  }

  // ---------------------------------------
  // FIND ONE
  // ---------------------------------------
  if (scenarios.findOne && scenarios.findOne.length > 0) {
    describe(`FIND ONE`, () => {
      it.each(scenarios.findOne!)(
        'should read with $name',
        async (scenario) => {
          // Ensure the create operation is defined
          expect(operations.create).toBeDefined()

          // Create the entity
          const created = await operations.create!(
            scenario.accessToken,
            scenario.orgname,
            scenario.createDto
          )

          // Ensure the findOne operation is defined
          expect(operations.findOne).toBeDefined()

          // Find the entity
          const found = await operations.findOne!(
            scenario.accessToken,
            scenario.orgname,
            created.body.id
          )

          // Log response for debugging if status is incorrect
          if (found.status !== scenario.expectedStatus) {
            console.error(
              `Unexpected status code for FIND ONE operation. Response:`,
              (found.body as any).message
            )
          }

          // Validate the response
          expect(found).toBeDefined()
          expect(found.status).toBe(scenario.expectedStatus)
          expect(found).toSatisfyApiSpec()

          // Clean up by deleting the entity if necessary
          if (created.body.id && operations.delete) {
            await operations.delete(
              scenario.accessToken,
              scenario.orgname,
              created.body.id
            )
          }
        }
      )
    })
  }

  // ---------------------------------------
  // UPDATE
  // ---------------------------------------
  if (scenarios.update && scenarios.update.length > 0) {
    describe(`UPDATE`, () => {
      it.each(scenarios.update!)(
        'should update with $name',
        async (scenario) => {
          // Ensure the create operation is defined
          expect(operations.create).toBeDefined()

          // Create the entity
          const created = await operations.create!(
            scenario.accessToken,
            scenario.orgname,
            scenario.createDto
          )

          // Ensure the findOne operation is defined
          expect(operations.findOne).toBeDefined()

          // Find the entity
          const found = await operations.findOne!(
            scenario.accessToken,
            scenario.orgname,
            created.body.id
          )

          // Ensure the update operation is defined
          expect(operations.update).toBeDefined()

          // Update the entity
          const updated = await operations.update!(
            scenario.accessToken,
            scenario.orgname,
            found.body.id,
            scenario.updateDto
          )

          // Log response for debugging if status is incorrect
          if (updated.status !== scenario.expectedStatus) {
            console.error(
              `Unexpected status code for UPDATE operation. Response:`,
              (updated.body as any).message
            )
          }

          // Validate the response
          expect(updated).toBeDefined()
          expect(updated.status).toBe(scenario.expectedStatus)
          expect(updated).toSatisfyApiSpec()

          // Clean up by deleting the entity if necessary
          if (updated.body.id && operations.delete) {
            await operations.delete(
              scenario.accessToken,
              scenario.orgname,
              updated.body.id
            )
          }
        }
      )
    })
  }

  // ---------------------------------------
  // FIND ALL
  // ---------------------------------------
  if (scenarios.findAll && scenarios.findAll.length > 0) {
    describe(`FIND ALL`, () => {
      it.each(scenarios.findAll!)(
        'should list with $name',
        async (scenario) => {
          // Ensure the create operation is defined
          expect(operations.create).toBeDefined()

          // Create the entities
          const createdIds: string[] = []
          for (const createDto of scenario.createDtos) {
            expect(operations.create).toBeDefined()
            const created = await operations.create!(
              scenario.accessToken,
              scenario.orgname,
              createDto
            )
            createdIds.push(created.body.id)
          }

          // Ensure the findAll operation is defined
          expect(operations.findAll).toBeDefined()

          // Find the entities
          const allFound = await operations.findAll!(
            scenario.accessToken,
            scenario.orgname,
            scenario.searchQueryDto
          )

          // Log response for debugging if status is incorrect
          if (allFound.status !== scenario.expectedStatus) {
            console.error(
              `Unexpected status code for FIND ALL operation. Response:`,
              (allFound.body as any).message
            )
          }

          // Validate the response
          expect(allFound).toBeDefined()
          expect(allFound.status).toBe(scenario.expectedStatus)
          expect(allFound).toSatisfyApiSpec()

          // Clean up by deleting the entities if necessary
          if (operations.delete) {
            for (const id of createdIds) {
              expect(operations.delete).toBeDefined()
              await operations.delete(
                scenario.accessToken,
                scenario.orgname,
                id
              )
            }
          }
        }
      )
    })
  }

  // ---------------------------------------
  // DELETE
  // ---------------------------------------
  if (scenarios.delete && scenarios.delete.length > 0) {
    describe(`DELETE`, () => {
      it.each(scenarios.delete!)(
        'should delete with $name',
        async (scenario) => {
          // Ensure the create operation is defined
          expect(operations.create).toBeDefined()

          // Create the entity
          const created = await operations.create!(
            scenario.accessToken,
            scenario.orgname,
            scenario.createDto
          )

          // Ensure the delete operation is defined
          expect(operations.delete).toBeDefined()

          // Delete the entity
          const deleted = await operations.delete!(
            scenario.accessToken,
            scenario.orgname,
            created.body.id
          )

          // Log response for debugging if status is incorrect
          if (deleted.status !== scenario.expectedStatus) {
            console.error(
              `Unexpected status code for DELETE operation. Response:`,
              deleted.body.message
            )
          }

          // Validate the response
          expect(deleted).toBeDefined()
          expect(deleted.status).toBe(scenario.expectedStatus)
          expect(deleted).toSatisfyApiSpec()
        }
      )
    })
  }
}
