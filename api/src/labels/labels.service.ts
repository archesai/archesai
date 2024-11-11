import { Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";

import { BaseService } from "../common/base.service";
import { CreateLabelDto } from "./dto/create-label.dto";
import { UpdateLabelDto } from "./dto/update-label.dto";
import { LabelEntity, LabelModel } from "./entities/label.entity";
import { LabelRepository } from "./label.repository";

@Injectable()
export class LabelsService extends BaseService<
  LabelEntity,
  CreateLabelDto,
  UpdateLabelDto,
  LabelRepository,
  LabelModel
> {
  private readonly logger: Logger = new Logger("Labels Service");

  constructor(private labelRepository: LabelRepository) {
    super(labelRepository);
  }

  protected toEntity(model: LabelModel): LabelEntity {
    return new LabelEntity(model);
  }
}
