import { Logger, Type, UsePipes } from '@nestjs/common'
import { Body, Delete, Get, Param, Patch, Post, Query } from '@nestjs/common'
import { ApiBody, ApiOperation, ApiResponse } from '@nestjs/swagger'

import { BaseService } from './base.service'
import { ApiPaginatedResponse } from './decorators/paginated.decorator'
import { PaginatedDto } from './dto/paginated.dto'
import { OperatorEnum, SearchQueryDto } from './dto/search-query.dto'
import { AbstractValidationPipe } from './pipes/abstract-validation.pipe'
import { CurrentUser } from '../auth/decorators/current-user.decorator'
import { UserEntity } from '../users/entities/user.entity'

export function BaseController<
  Entity,
  CreateDto extends Partial<Entity>,
  UpdateDto extends Partial<Entity>,
  Service extends BaseService<Entity, any, any>
>(
  EntityClass: Type<Entity>,
  CreateDtoClass: Type<CreateDto>,
  UpdateDtoClass: Type<UpdateDto>,
  itemType: string = EntityClass.name.replace('Entity', '').toLowerCase()
) {
  class BaseController {
    readonly logger = new Logger(this.constructor.name)
    readonly itemType = itemType
    constructor(public readonly service: Service) {}

    @ApiBody({ type: CreateDtoClass })
    @ApiOperation({ summary: `Create a new ${itemType}` })
    @ApiResponse({
      description: `${itemType} created successfully.`,
      status: 201,
      type: EntityClass
    })
    @UsePipes(new AbstractValidationPipe({ body: CreateDtoClass }))
    @Post()
    async create(
      @Param('orgname') orgname: string,
      @Body() createDto: CreateDto,
      @CurrentUser() currentUserDto?: UserEntity
    ): Promise<Entity> {
      this.logger.debug(`creating ${itemType}`, {
        orgname,
        createDto,
        currentUserDto
      })
      return this.service.create({
        ...createDto,
        orgname,
        username: currentUserDto?.username
      })
    }

    @ApiOperation({ summary: `Get all ${itemType}s` })
    @ApiPaginatedResponse(EntityClass)
    @UsePipes(new AbstractValidationPipe({}))
    @Get()
    async findAll(
      @Param('orgname') orgname: string,
      @Query() searchQueryDto: SearchQueryDto
    ): Promise<PaginatedDto<Entity>> {
      this.logger.debug(`fetching all ${itemType}`, {
        orgname,
        searchQueryDto
      })
      return this.service.findAll({
        ...searchQueryDto,
        filters: [
          ...(searchQueryDto.filters || []),
          {
            field: 'orgname',
            operator: OperatorEnum.EQUALS,
            value: orgname
          }
        ]
      })
    }

    @ApiOperation({ summary: `Get a single ${itemType}` })
    @ApiResponse({
      description: `${itemType} found.`,
      status: 200,
      type: EntityClass
    })
    @UsePipes(new AbstractValidationPipe({}))
    @Get(':id')
    async findOne(
      @Param('orgname') orgname: string,
      @Param('id') id: string
    ): Promise<Entity> {
      this.logger.debug(`fetching single ${itemType}`, {
        orgname,
        id
      })
      return this.service.findOne(id)
    }

    @ApiOperation({ summary: `Delete a ${itemType}` })
    @ApiResponse({
      description: `${itemType} deleted successfully.`,
      status: 200
    })
    @UsePipes(new AbstractValidationPipe({}))
    @Delete(':id')
    async remove(
      @Param('orgname') orgname: string,
      @Param('id') id: string
    ): Promise<Entity> {
      this.logger.debug(`deleting single ${itemType}`, {
        orgname,
        id
      })
      return this.service.remove(id)
    }

    @ApiBody({ type: UpdateDtoClass })
    @ApiOperation({ summary: `Update a ${itemType}` })
    @ApiResponse({
      description: `${itemType} updated successfully.`,
      status: 200,
      type: EntityClass
    })
    @UsePipes(
      new AbstractValidationPipe({
        body: UpdateDtoClass
      })
    )
    @Patch(':id')
    async update(
      @Param('orgname') orgname: string,
      @Param('id') id: string,
      @Body() updateDto: UpdateDto
    ): Promise<Entity> {
      this.logger.debug(`updating ${itemType}`, {
        orgname,
        id,
        updateDto
      })
      return this.service.update(id, updateDto)
    }
  }

  return BaseController
}
