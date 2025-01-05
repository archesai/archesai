import { Tool as _PrismaTool } from '@prisma/client'

import { BaseEntity } from '../../common/entities/base.entity'
import { IsEnum, IsString } from 'class-validator'
import { Expose } from 'class-transformer'

export type ToolModel = _PrismaTool

export enum ToolIOTypeEnum {
  TEXT = 'TEXT',
  IMAGE = 'IMAGE',
  AUDIO = 'AUDIO',
  VIDEO = 'VIDEO'
}

export class ToolEntity extends BaseEntity implements ToolModel {
  /**
   * The tool description
   * @example This tool converts a file to text, regardless of the file type.
   */
  @IsString()
  @Expose()
  description: string

  /**
   * The tools input type
   * @example TEXT
   */
  @IsEnum(ToolIOTypeEnum)
  @Expose()
  inputType: ToolIOTypeEnum

  /**
   * The tool name
   * @example extract-text
   */
  @IsString()
  @Expose()
  name: string

  /**
   * The organization name
   * @example my-organization
   */
  @IsString()
  @Expose()
  orgname: string

  /**
   * The tools output type
   * @example TEXT
   */
  @IsEnum(ToolIOTypeEnum)
  @Expose()
  outputType: ToolIOTypeEnum

  /**
   * The base of the tool
   * @example extract-text
   */
  @IsString()
  @Expose()
  toolBase: string

  constructor(tool: ToolModel) {
    super()
    Object.assign(this, tool)
  }
}
