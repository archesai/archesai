import { ArgumentMetadata, BadRequestException } from '@nestjs/common'
import { CreateRunDto } from '../dto/create-run.dto'
import { CustomValidationPipe } from '@/src/common/pipes/custom-validation.pipe'

describe('CreateRunDto with CustomValidationPipe', () => {
  let validationPipe: CustomValidationPipe
  let metadata: ArgumentMetadata

  beforeEach(() => {
    validationPipe = new CustomValidationPipe()
    metadata = {
      type: 'body',
      metatype: CreateRunDto
    }
  })

  it('should make certain fields required', async () => {
    const dto = {}
    try {
      await validationPipe.transform(dto, metadata)
      fail('Expected a BadRequestException due to missing required fields')
    } catch (error: any) {
      expect(error).toBeInstanceOf(BadRequestException)
      expect(error.message).toContain(
        'runType must be one of the following values: TOOL_RUN, PIPELINE_RUN'
      )
    }
  })

  it('should allow valid values', async () => {
    const dto = {
      runType: 'TOOL_RUN',
      toolId: '1',
      contentIds: ['1'],
      text: 'This is the text to use as input for the run.'
    }
    const result = await validationPipe.transform(dto, metadata)

    expect(result).toEqual(dto)
  })

  it('should handle defaults', async () => {
    const dto = {
      runType: 'TOOL_RUN',
      toolId: '1',
      text: 'This is the text to use as input for the run.'
    }
    const instance = await validationPipe.transform(dto, metadata)
    expect(instance.runType).toBe(dto.runType)
    expect(instance.toolId).toBe(dto.toolId)
    expect(instance.pipelineId).toBeUndefined()
    expect(instance.text).toBe(dto.text)
    expect(instance.url).toBeUndefined()
  })
})
