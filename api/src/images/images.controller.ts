import { Body, Controller, Param, Post } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { ContentEntity } from "../content/entities/content.entity";
import { CreateImageDto } from "./dto/create-image.dto";
import { ImagesService } from "./images.service";

@ApiBearerAuth()
@ApiTags("Images")
@Controller("organizations/:orgname/images")
export class ImagesController
  implements
    BaseController<ContentEntity, CreateImageDto, undefined, undefined>
{
  constructor(private readonly imagesService: ImagesService) {}

  @ApiCrudOperation(Operation.CREATE, "image", ContentEntity, true)
  @Post("/")
  create(
    @Param("orgname") orgname: string,
    @Body() createImageDto: CreateImageDto
  ) {
    return this.imagesService.create(orgname, createImageDto);
  }
}
