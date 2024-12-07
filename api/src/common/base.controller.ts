// base.controller.mixin.ts
import { Type } from '@nestjs/common'
import { Body, Delete, Get, Param, Patch, Post, Query } from '@nestjs/common'
import { ApiBody, ApiOperation, ApiResponse } from '@nestjs/swagger'

import { BaseService } from './base.service'
import { ApiPaginatedResponse } from './decorators/paginated.decorator'
import { PaginatedDto } from './dto/paginated.dto'
import { SearchQueryDto } from './dto/search-query.dto'

export function BaseController<
  Entity,
  CreateDto,
  UpdateDto,
  Service extends BaseService<Entity, CreateDto, UpdateDto, any, any>
>(
  EntityClass: Type<Entity>,
  CreateDtoClass: Type<CreateDto>,
  UpdateDtoClass: Type<UpdateDto>,
  itemType: string = EntityClass.name.replace('Entity', '').toLowerCase()
) {
  class BaseController {
    constructor(public readonly service: Service) {}

    @ApiBody({ type: CreateDtoClass })
    @ApiOperation({ summary: `Create a new ${itemType}` })
    @ApiResponse({
      description: `${itemType} created successfully.`,
      status: 201,
      type: EntityClass
    })
    @Post()
    async create(
      @Param('orgname') orgname: string,
      @Body() createDto: CreateDto,
      ...additionalData: any[]
    ): Promise<Entity> {
      return this.service.create(orgname, createDto, additionalData)
    }

    @ApiOperation({ summary: `Get all ${itemType}s` })
    @ApiPaginatedResponse(EntityClass)
    @Get()
    async findAll(
      @Param('orgname') orgname: string,
      @Query() searchQueryDto: SearchQueryDto
    ): Promise<PaginatedDto<Entity>> {
      return this.service.findAll(orgname, searchQueryDto)
    }

    @ApiOperation({ summary: `Get a single ${itemType}` })
    @ApiResponse({
      description: `${itemType} found.`,
      status: 200,
      type: EntityClass
    })
    @Get(':id')
    async findOne(
      @Param('orgname') orgname: string,
      @Param('id') id: string
    ): Promise<Entity> {
      return this.service.findOne(orgname, id)
    }

    @ApiOperation({ summary: `Delete a ${itemType}` })
    @ApiResponse({
      description: `${itemType} deleted successfully.`,
      status: 200
    })
    @Delete(':id')
    async remove(
      @Param('orgname') orgname: string,
      @Param('id') id: string
    ): Promise<void> {
      return this.service.remove(orgname, id)
    }

    @ApiBody({ type: UpdateDtoClass })
    @ApiOperation({ summary: `Update a ${itemType}` })
    @ApiResponse({
      description: `${itemType} updated successfully.`,
      status: 200,
      type: EntityClass
    })
    @Patch(':id')
    async update(
      @Param('orgname') orgname: string,
      @Param('id') id: string,
      @Body() updateDto: UpdateDto
    ): Promise<Entity> {
      return this.service.update(orgname, id, updateDto)
    }
  }

  return BaseController
}
