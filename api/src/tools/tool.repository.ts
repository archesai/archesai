import { Injectable } from "@nestjs/common";
import { Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateToolDto } from "./dto/create-tool.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";
import { ToolModel } from "./entities/tool.entity";

@Injectable()
export class ToolRepository extends BaseRepository<
  ToolModel,
  CreateToolDto,
  UpdateToolDto,
  Prisma.ToolInclude,
  Prisma.ToolUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.tool);
  }

  async createDefaultTools(orgname: string) {
    return this.prisma.tool.createMany({
      data: [
        {
          description:
            "Extract text from a file. This tool supports all file types.",
          inputType: "TEXT",
          name: "Extract Text",
          orgname,
          outputType: "TEXT",
          toolBase: "extract-text",
        },
        {
          description: "Create an image from text.",
          inputType: "TEXT",
          name: "Text to Image",
          orgname,
          outputType: "IMAGE",
          toolBase: "text-to-image",
        },
        {
          description: "Summarize text. This tool supports all languages.",
          inputType: "TEXT",
          name: "Summarize",
          orgname,
          outputType: "TEXT",
          toolBase: "summarize",
        },
        {
          description:
            "Create embeddings from text. This tool supports all languages.",
          inputType: "TEXT",
          name: "Create Embeddings",
          orgname,
          outputType: "TEXT", // FIXME make this none
          toolBase: "create-embeddings",
        },
        {
          description:
            "Convert text to speech. This tool supports all languages.",
          inputType: "TEXT",
          name: "Text to Speech",
          orgname,
          outputType: "AUDIO",
          toolBase: "text-to-speech",
        },
      ],
    });
  }
}
