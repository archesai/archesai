import { BaseRepository } from "../base.repository";

class TestRepository extends BaseRepository<any, any, any, any, any> {
  constructor(delegate: any) {
    super(delegate);
  }
}

const mockDelegate = {
  count: jest.fn(),
  create: jest.fn(),
  delete: jest.fn(),
  findMany: jest.fn(),
  findUniqueOrThrow: jest.fn(),
  update: jest.fn(),
};

describe("BaseRepository", () => {
  let repository: TestRepository;

  beforeEach(() => {
    repository = new TestRepository(mockDelegate);
  });

  describe("create", () => {
    it("should create a new record", async () => {
      const orgname = "test-org";
      const createDto = { name: "Test Item" };
      const additionalData = { extraField: "extraValue" };
      const expectedResult = { id: "1", ...createDto, ...additionalData };

      mockDelegate.create.mockResolvedValue(expectedResult);

      const result = await repository.create(
        orgname,
        createDto,
        additionalData
      );

      expect(mockDelegate.create).toHaveBeenCalledWith({
        data: {
          organization: { connect: { orgname } },
          ...createDto,
          ...additionalData,
        },
        include: undefined, // or your defaultInclude if set
      });
      expect(result).toEqual(expectedResult);
    });
  });

  // Add similar tests for findAll, findOne, update, and remove methods
});
