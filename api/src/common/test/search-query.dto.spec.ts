import { FieldAggregate, transformValues } from '../dto/search-query.dto'

describe('SearchQueryDto', () => {
  describe('transformValues()', () => {
    it('should transform a valid JSON string array', () => {
      const input = JSON.stringify([
        { field: 'createdAt', type: 'count', granularity: 'day' }
      ])
      const result = transformValues(input, FieldAggregate)
      expect(result).toEqual([
        { field: 'createdAt', type: 'count', granularity: 'day' }
      ])
    })

    it('should wrap a valid JSON object in an array if not an array', () => {
      const input = JSON.stringify({
        field: 'createdAt',
        type: 'count',
        granularity: 'day'
      })
      const result = transformValues(input, FieldAggregate)
      expect(result).toEqual([
        { field: 'createdAt', type: 'count', granularity: 'day' }
      ])
    })
  })
})
