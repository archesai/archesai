import { Test, TestingModule } from "@nestjs/testing";

import { OrganizationsController } from "../organizations.controller";
import { OrganizationsService } from "../organizations.service";

describe("OrganizationsController", () => {
  let controller: OrganizationsController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [OrganizationsController],
      providers: [
        {
          provide: OrganizationsService,
          useFactory: () => ({
            findAll: jest.fn(() => Promise.resolve([])),
          }),
        },
      ],
    }).compile();

    controller = module.get<OrganizationsController>(OrganizationsController);
  });

  it("should be defined", () => {
    expect(controller).toBeDefined();
  });
});
