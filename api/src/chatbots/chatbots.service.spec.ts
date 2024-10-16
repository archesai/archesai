import { TestBed } from "@automock/jest";
import { ForbiddenException, NotFoundException } from "@nestjs/common";
import { PlanType } from "@prisma/client";

import { ContentService } from "../content/content.service";
import { OrganizationEntity } from "../organizations/entities/organization.entity";
import { OrganizationsService } from "../organizations/organizations.service";
import { ChatbotsService } from "./chatbots.service";

describe("ChatbotsService unit spec", () => {
  const mockOrganization = {
    billingEmail: "billingEmail",
    createdAt: new Date(),
    credits: 100,
    id: "id",
    members: [],
    orgname: "orgname",
    plan: "UNLIMITED" as const,
    stripeCustomerId: "stripeCustomerId",
    updatedAt: new Date(),
  } as OrganizationEntity;

  let chatbotsService: ChatbotsService;
  let mockedContentService: jest.Mocked<ContentService>;
  let mockedOrganizationsService: jest.Mocked<OrganizationsService>;

  beforeAll(() => {
    const { unit, unitRef } = TestBed.create(ChatbotsService).compile();

    chatbotsService = unit;

    mockedContentService = unitRef.get(ContentService);
    mockedOrganizationsService = unitRef.get(OrganizationsService);
  });

  describe("create", () => {
    it("should throw a 404 when the organization does not exist", async () => {
      mockedOrganizationsService.findOneByName.mockImplementationOnce(
        async () => {
          throw new NotFoundException();
        }
      );

      await expect(
        chatbotsService.create("nonexistentOrgname", {
          description: "description",
          llmBase: "GPT_3_5_TURBO_16_K",
          name: "name",
        })
      ).rejects.toThrow(NotFoundException);
    });

    it("should throw a 404 error when provided a non-existent document ID", async () => {
      mockedOrganizationsService.findOneByName.mockResolvedValueOnce(
        mockOrganization
      );
      mockedContentService.findOne.mockImplementationOnce(async () => {
        throw new NotFoundException();
      });

      await expect(
        chatbotsService.create("orgname", {
          description: "description",

          llmBase: "GPT_3_5_TURBO_16_K",
          name: "name",
        })
      ).rejects.toThrow(NotFoundException);
    });

    it("should throw a 403 error when the user tries to use GPT-4 and is not on API plan", async () => {
      for (const plan of ["FREE", "PRO"]) {
        mockedOrganizationsService.findOneByName.mockResolvedValueOnce({
          ...mockOrganization,
          plan: plan as PlanType,
        });

        await expect(
          chatbotsService.create("orgname", {
            description: "description",

            llmBase: "GPT_4",
            name: "name",
          })
        ).rejects.toThrow(ForbiddenException);
      }
    });
  });
});
