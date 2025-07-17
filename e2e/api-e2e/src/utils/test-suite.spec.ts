import type { CrudTestCases, HttpOperations } from '#types/crud-test-cases'

export function runCrudTestSuite<
  Entity extends { id: string },
  CreateDto,
  UpdateDto
>(
  operations: HttpOperations<Entity, CreateDto, UpdateDto>,
  scenarios: CrudTestCases<CreateDto, UpdateDto>
) {
  // ---------------------------------------
  // CREATE
  // ---------------------------------------
  if (scenarios.create && scenarios.create.length > 0) {
    describe(`CREATE`, () => {
      it.each(scenarios.create ?? [])(
        'should create with $name',
        async (scenario) => {
          // Ensure the create operation is defined
          expect(operations.create).toBeDefined()
          if (!operations.create) {
            throw new Error('Create operation is not defined')
          }

          // Create the entity
          const created = await operations.create(
            scenario.accessToken,
            scenario.organizationId,
            scenario.createDto
          )

          // Log response for debugging if status is incorrect
          if (created.status !== scenario.expectedStatus) {
            console.error(
              `Unexpected status code for CREATE operation. Response:`,
              created.body
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
              scenario.organizationId,
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
      it.each(scenarios.findOne ?? [])(
        'should read with $name',
        async (scenario) => {
          // Ensure the create operation is defined
          expect(operations.create).toBeDefined()
          if (!operations.create) {
            throw new Error('Create operation is not defined')
          }

          // Create the entity
          const created = await operations.create(
            scenario.accessToken,
            scenario.organizationId,
            scenario.createDto
          )

          // Ensure the findOne operation is defined
          expect(operations.findOne).toBeDefined()
          if (!operations.findOne) {
            throw new Error('FindOne operation is not defined')
          }

          // Find the entity
          const found = await operations.findOne(
            scenario.accessToken,
            scenario.organizationId,
            created.body.id
          )

          // Log response for debugging if status is incorrect
          if (found.status !== scenario.expectedStatus) {
            console.error(
              `Unexpected status code for FIND ONE operation. Response:`,
              found.body
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
              scenario.organizationId,
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
      it.each(scenarios.update ?? [])(
        'should update with $name',
        async (scenario) => {
          // Ensure the create operation is defined
          expect(operations.create).toBeDefined()
          if (!operations.create) {
            throw new Error('Create operation is not defined')
          }

          // Create the entity
          const created = await operations.create(
            scenario.accessToken,
            scenario.organizationId,
            scenario.createDto
          )

          // Ensure the findOne operation is defined
          expect(operations.findOne).toBeDefined()
          if (!operations.findOne) {
            throw new Error('FindOne operation is not defined')
          }

          // Find the entity
          const found = await operations.findOne(
            scenario.accessToken,
            scenario.organizationId,
            created.body.id
          )

          // Ensure the update operation is defined
          expect(operations.update).toBeDefined()
          if (!operations.update) {
            throw new Error('Update operation is not defined')
          }

          // Update the entity
          const updated = await operations.update(
            scenario.accessToken,
            scenario.organizationId,
            found.body.id,
            scenario.updateDto
          )

          // Log response for debugging if status is incorrect
          if (updated.status !== scenario.expectedStatus) {
            console.error(
              `Unexpected status code for UPDATE operation. Response:`,
              updated.body
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
              scenario.organizationId,
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
  if (scenarios.findMany && scenarios.findMany.length > 0) {
    describe(`FIND ALL`, () => {
      it.each(scenarios.findMany ?? [])(
        'should list with $name',
        async (scenario) => {
          // Ensure the create operation is defined
          expect(operations.create).toBeDefined()
          if (!operations.create) {
            throw new Error('Create operation is not defined')
          }

          // Create the entities
          const createdIds: string[] = []
          for (const createDto of scenario.createDtos) {
            expect(operations.create).toBeDefined()
            const created = await operations.create(
              scenario.accessToken,
              scenario.organizationId,
              createDto
            )
            createdIds.push(created.body.id)
          }

          // Ensure the findMany operation is defined
          expect(operations.findMany).toBeDefined()
          if (!operations.findMany) {
            throw new Error('FindAll operation is not defined')
          }

          // Find the entities
          const allFound = await operations.findMany(
            scenario.accessToken,
            scenario.organizationId,
            scenario.query
          )

          // Log response for debugging if status is incorrect
          if (allFound.status !== scenario.expectedStatus) {
            console.error(
              `Unexpected status code for FIND ALL operation. Response:`,
              allFound.body
            )
          }

          // Validate the response
          expect(allFound).toBeDefined()
          expect(allFound.status).toBe(scenario.expectedStatus)
          expect(allFound).toSatisfyApiSpec()

          // Clean up by deleting the entities if necessary
          if (operations.delete) {
            for (const id of createdIds) {
              await operations.delete(
                scenario.accessToken,
                scenario.organizationId,
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
      it.each(scenarios.delete ?? [])(
        'should delete with $name',
        async (scenario) => {
          // Ensure the create operation is defined
          expect(operations.create).toBeDefined()
          if (!operations.create) {
            throw new Error('Create operation is not defined')
          }

          // Create the entity
          const created = await operations.create(
            scenario.accessToken,
            scenario.organizationId,
            scenario.createDto
          )

          // Ensure the delete operation is defined
          expect(operations.delete).toBeDefined()
          if (!operations.delete) {
            throw new Error('Delete operation is not defined')
          }

          // Delete the entity
          const deleted = await operations.delete(
            scenario.accessToken,
            scenario.organizationId,
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
