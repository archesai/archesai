import { ConfigService } from "@nestjs/config";
import { Test, TestingModule } from "@nestjs/testing";

import { CurrentUserDto } from "../../auth/decorators/current-user.decorator";
import { PrismaService } from "../../prisma/prisma.service";
import { StripeService } from "../../billing/billing.service";
import { CreateOrganizationDto } from "../dto/create-organization.dto";
import { OrganizationsService } from "../organizations.service";

describe("OrganizationsService", () => {
  let organizationsService: OrganizationsService;
  let prismaService: PrismaService;
  let stripeService: StripeService;
  let configService: ConfigService;

  const user: CurrentUserDto = {
    createdAt: new Date(),
    deactivated: false,
    defaultOrgname: "org1",
    email: "user1@test.com",
    emailVerified: true,
    firstName: "User",
    id: "userId",
    lastName: "One",
    memberships: [],
    photoUrl: "",
    uid: "",
    updatedAt: new Date(),
    username: "user1",
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        OrganizationsService,
        {
          provide: PrismaService,
          useValue: {
            organization: {
              create: jest.fn().mockResolvedValue({}),
            },
          },
        },
        {
          provide: StripeService,
          useValue: {
            createCustomer: jest.fn().mockResolvedValue({}),
          },
        },
        {
          provide: ConfigService,
          useValue: {
            get: jest.fn().mockReturnValue(false),
          },
        },
      ],
    }).compile();

    organizationsService =
      module.get<OrganizationsService>(OrganizationsService);
    prismaService = module.get<PrismaService>(PrismaService);
    stripeService = module.get<StripeService>(StripeService);
    configService = module.get<ConfigService>(ConfigService);
  });

  it("should be defined", () => {
    expect(organizationsService).toBeDefined();
  });

  describe("createAndInitialize", () => {
    it("should create an organization and initialize stripe if billing is enabled", async () => {
      jest.spyOn(configService, "get").mockReturnValue(true);
      jest
        .spyOn(stripeService, "createCustomer")
        .mockResolvedValue({ id: "stripeId" } as any);
      jest
        .spyOn(prismaService.organization, "create")
        .mockResolvedValue({ id: "orgId" } as any);

      const createOrgDto: CreateOrganizationDto = {
        billingEmail: "billing@test.com",
        orgname: "org1",
      };

      const result = await organizationsService.createAndInitialize(
        user,
        createOrgDto
      );

      expect(result).toEqual({ id: "orgId" });
      expect(stripeService.createCustomer).toHaveBeenCalledWith(
        "org1",
        "billing@test.com"
      );
      expect(prismaService.organization.create).toHaveBeenCalled();
    });

    it("should create an organization without stripe if billing is disabled", async () => {
      jest.spyOn(configService, "get").mockReturnValue(false);
      jest
        .spyOn(prismaService.organization, "create")
        .mockResolvedValue({ id: "orgId" } as any);

      const createOrgDto: CreateOrganizationDto = {
        billingEmail: "billing@test.com",
        orgname: "org1",
      };

      const result = await organizationsService.createAndInitialize(
        user,
        createOrgDto
      );

      expect(result).toEqual({ id: "orgId" });
      expect(stripeService.createCustomer).not.toHaveBeenCalled();
      expect(prismaService.organization.create).toHaveBeenCalled();
    });
  });

  // Similarly, other functions can be tested like above
});
