import { Controller } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import { BaseController } from "../common/base.controller";
import { CreateLabelDto } from "./dto/create-label.dto";
import { UpdateLabelDto } from "./dto/update-label.dto";
import { LabelEntity } from "./entities/label.entity";
import { LabelsService } from "./labels.service";

@ApiBearerAuth()
@ApiTags("Labels")
@Controller("/organizations/:orgname/labels")
export class LabelsController extends BaseController<
  LabelEntity,
  CreateLabelDto,
  UpdateLabelDto,
  LabelsService
> {
  constructor(private readonly labelsService: LabelsService) {
    super(labelsService);
  }
}
