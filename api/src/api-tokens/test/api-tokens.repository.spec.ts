import { Test, TestingModule } from "@nestjs/testing";

import { PrismaService } from "../../prisma/prisma.service";
import { ApiTokenRepository } from "../api-token.repository";
import { CreateApiTokenDto } from "../dto/create-api-token.dto";
import { ApiTokenModel } from "../entities/api-token.entity";

describe("ApiTokenRepository", () => {
  let repository: ApiTokenRepository;
  let prismaService: PrismaService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ApiTokenRepository,
        {
          provide: PrismaService,
          useValue: {
            apiToken: {
              create: jest.fn(),
            },
          },
        },
      ],
    }).compile();

    repository = module.get<ApiTokenRepository>(ApiTokenRepository);
    prismaService = module.get<PrismaService>(PrismaService);
  });

  describe("create", () => {
    it("should create an API token", async () => {
      // Arrange
      const orgname = "test-org";
      const createApiTokenDto: CreateApiTokenDto = {
        domains: "*",
        name: "Test Token",
        role: "USER",
      };
      const additionalData = {
        id: "test-id",
        key: "test-key",
        username: "test-user",
      };

      const expectedResult: ApiTokenModel = {
        ...createApiTokenDto,
        createdAt: new Date(),
        id: additionalData.id,
        key: additionalData.key,
        orgname: "org-id",
        updatedAt: new Date(),
        username: "test-user",
      };

      prismaService.apiToken.create = jest
        .fn()
        .mockResolvedValue(expectedResult);

      // Act
      const result = await repository.create(
        orgname,
        createApiTokenDto,
        additionalData
      );

      // Assert
      expect(prismaService.apiToken.create).toHaveBeenCalledWith({
        data: {
          ...createApiTokenDto,
          id: additionalData.id,
          key: additionalData.key,
          organization: {
            connect: {
              orgname,
            },
          },
          user: {
            connect: {
              username: additionalData.username,
            },
          },
        },
      });
      expect(result).toEqual(expectedResult);
    });
  });
});
