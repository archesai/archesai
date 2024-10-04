import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { JwtService } from "@nestjs/jwt";
import { ApiToken } from "@prisma/client";
import { v4 } from "uuid";

import { BaseService } from "../common/base.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { ApiTokenRepository } from "./api-token.repository";
import { ApiTokenQueryDto } from "./dto/api-token-query.dto";
import { CreateApiTokenDto } from "./dto/create-api-token.dto";
import { UpdateApiTokenDto } from "./dto/update-api-token.dto";

@Injectable()
export class ApiTokensService
  implements
    BaseService<
      ApiToken,
      CreateApiTokenDto,
      ApiTokenQueryDto,
      UpdateApiTokenDto
    >
{
  constructor(
    private apiTokenRepository: ApiTokenRepository,
    private configService: ConfigService,
    private jwtService: JwtService,
    private websocketsService: WebsocketsService
  ) {}

  async create(orgname: string, createTokenDto: CreateApiTokenDto) {
    const id = v4();
    const token = this.jwtService.sign(
      {
        domains: createTokenDto.domains,
        id,
        orgname,
        role: createTokenDto.role,
      },
      {
        expiresIn: `${this.configService.get(
          "JWT_API_TOKEN_EXPIRATION_TIME"
        )}s`,
        secret: this.configService.get("JWT_API_TOKEN_SECRET"),
      }
    );
    const snippet = "*********" + token.slice(-5);
    const apiToken = await this.apiTokenRepository.create(
      orgname,
      createTokenDto,
      snippet
    );
    this.websocketsService.socket.to(orgname).emit("update");

    apiToken.key = token;

    return apiToken;
  }

  async findAll(orgname: string, apiTokenQueryDto: ApiTokenQueryDto) {
    return this.apiTokenRepository.findAll(orgname, apiTokenQueryDto);
  }

  async findOne(orgname: string, id: string) {
    return this.apiTokenRepository.findOne(orgname, id);
  }

  async remove(orgname: string, id: string) {
    await this.apiTokenRepository.remove(orgname, id);
    this.websocketsService.socket.to(orgname).emit("update");
  }

  async update(
    orgname: string,
    id: string,
    updateApiTokenDto: UpdateApiTokenDto
  ) {
    const apiToken = await this.apiTokenRepository.update(
      orgname,
      id,
      updateApiTokenDto
    );
    this.websocketsService.socket.to(orgname).emit("update");
    return apiToken;
  }
}
