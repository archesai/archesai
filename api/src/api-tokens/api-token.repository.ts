import { Injectable } from "@nestjs/common";
import { ApiToken, Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateApiTokenDto } from "./dto/create-api-token.dto";
import { UpdateApiTokenDto } from "./dto/update-api-token.dto";

@Injectable()
export class ApiTokenRepository extends BaseRepository<
  ApiToken,
  CreateApiTokenDto,
  UpdateApiTokenDto,
  Prisma.ApiTokenInclude,
  Prisma.ApiTokenSelect
> {
  constructor(private prisma: PrismaService) {
    super(prisma.apiToken);
  }

  async create(
    orgname: string,
    createApiTokenDto: CreateApiTokenDto,
    additionalData: {
      id: string;
      snippet: string;
      uid: string;
    }
  ) {
    return this.prisma.apiToken.create({
      data: {
        ...createApiTokenDto,
        id: additionalData.id,
        key: additionalData.snippet,
        organization: {
          connect: {
            orgname,
          },
        },
        user: {
          connect: {
            id: additionalData.uid,
          },
        },
      },
    });
  }
}
