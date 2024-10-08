import { Body, Controller, Param, Post } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import {
  ApiCrudOperation,
  Operation,
} from "../common/api-crud-operation.decorator";
import { BaseController } from "../common/base.controller";
import { ContentEntity } from "../content/entities/content.entity";
import { DocumentsService } from "./documents.service";
import { CreateDocumentDto } from "./dto/create-document.dto";

@ApiBearerAuth()
@ApiTags("Documents")
@Controller("organizations/:orgname/documents")
export class DocumentsController
  implements
    BaseController<ContentEntity, CreateDocumentDto, undefined, undefined>
{
  constructor(private readonly documentsService: DocumentsService) {}

  @ApiCrudOperation(Operation.CREATE, "document", ContentEntity, true)
  @Post("/")
  create(
    @Param("orgname") orgname: string,
    @Body() createDocumentDto: CreateDocumentDto
  ) {
    return this.documentsService.create(orgname, createDocumentDto);
  }
}
