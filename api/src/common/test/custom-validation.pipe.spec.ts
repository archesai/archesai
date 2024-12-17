import { ArgumentMetadata, BadRequestException } from '@nestjs/common'
import {
  FieldAggregate,
  FieldFilter,
  SearchQueryDto
} from '../dto/search-query.dto'
import { CustomValidationPipe } from '../pipes/custom-validation.pipe'
import { validate } from 'class-validator'

describe('SearchQueryDto with CustomValidationPipe', () => {
  let validationPipe: CustomValidationPipe
  let metadata: ArgumentMetadata

  beforeEach(() => {
    validationPipe = new CustomValidationPipe()
    metadata = {
      type: 'query',
      metatype: SearchQueryDto
    }
  })

  it('should apply defaults and transform input correctly', async () => {
    const dto = {}
    const instance = await validationPipe.transform(dto, metadata)
    const errors = await validate(instance)

    // Check that defaults were applied
    expect(errors).toHaveLength(0)
    expect(instance.limit).toBe(10)
    expect(instance.offset).toBe(0)
    expect(instance.sortBy).toBe('createdAt')
    expect(instance.sortDirection).toBe('desc')
    expect(instance.filters).toEqual([])
    expect(instance.aggregates).toEqual([])
    expect(instance.endDate).toBeUndefined()
    expect(instance.startDate).toBeUndefined()
  })

  it('should parse JSON fields from strings to objects', async () => {
    const input = {
      filters: JSON.stringify([
        { field: 'status', operator: 'equals', value: 'COMPLETE' }
      ]),
      aggregates: JSON.stringify([
        { field: 'createdAt', type: 'count', granularity: 'day' }
      ])
    }

    const result = await validationPipe.transform(input, metadata)
    const filters = result.filters as FieldFilter[]
    const aggregates = result.aggregates as FieldAggregate[]

    expect(filters).toEqual([
      { field: 'status', operator: 'equals', value: 'COMPLETE' }
    ])
    expect(aggregates).toEqual([
      { field: 'createdAt', type: 'count', granularity: 'day' }
    ])
  })

  it('should handle strings and arrays as values', async () => {
    const input = {
      filters: JSON.stringify([
        { field: 'status', operator: 'equals', value: 'COMPLETE' },
        { field: 'time', operator: 'in', value: ['COMPLETE'] }
      ])
    }
    const result = await validationPipe.transform(input, metadata)
    const filters = result.filters as FieldFilter[]

    expect(filters).toEqual([
      { field: 'status', operator: 'equals', value: 'COMPLETE' },
      { field: 'time', operator: 'in', value: ['COMPLETE'] }
    ])
  })

  it('should throw a BadRequestException for invalid aggregates or filters', async () => {
    for (const input of [
      { filters: 'not-valid-json', aggregates: 'not-valid-json' },
      {
        filters: JSON.stringify([
          { field: 'status', operator: 'asdf', value: ['active'] }
        ])
      },
      {
        aggregates: JSON.stringify([
          { field: 'createdAt', type: 'count', granularity: 'bad_granularity' }
        ])
      },
      {
        filters: JSON.stringify([
          { field: 'status', operator: 'equals', value: 'COMPLETE' },
          { field: 'time', operator: 'in' }
        ])
      }
    ]) {
      await expect(validationPipe.transform(input, metadata)).rejects.toThrow(
        BadRequestException
      )
    }
  })

  it('should throw a BadRequestException for invalid JSON strings', async () => {
    const input = {
      filters: 'not-valid-json'
    }

    await expect(validationPipe.transform(input, metadata)).rejects.toThrow(
      BadRequestException
    )
  })

  it('should fail validation for invalid enum values', async () => {
    const input = {
      sortDirection: 'INVALID'
    }

    try {
      await validationPipe.transform(input, metadata)
      fail('Expected a BadRequestException due to invalid enum')
    } catch (error) {
      expect(error).toBeInstanceOf(BadRequestException)

      const response = error.getResponse()
      // The response typically has a structure like:
      // {
      //   statusCode: 400,
      //   message: [ /* array of validation error messages */ ],
      //   error: "Bad Request"
      // }

      // Check the validation messages here
      expect(response).toHaveProperty('message')
      expect(Array.isArray(response['message'])).toBe(true)
      expect(response['message']).toContain(
        'sortDirection must be one of the following values: asc, desc'
      )
    }
  })
})
