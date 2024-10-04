import { TestBed } from "@automock/jest";
import { NotFoundException } from "@nestjs/common";

import { ChatbotsService } from "../../chatbots/chatbots.service";
import { ContentService } from "../../content/content.service";
import { OrganizationsService } from "../../organizations/organizations.service";
import { ThreadsService } from "../../threads/threads.service";
import { WebsocketsService } from "../../websockets/websockets.service";
import { MessageRepository } from "../message.repository";
import { MessagesService } from "../messages.service";

describe("MessagesService unit spec", () => {
  let messagesService: MessagesService;
  let mockedChatbotsService: jest.Mocked<ChatbotsService>;
  let mockedOrganizationsService: jest.Mocked<OrganizationsService>;
  let mockedThreadsService: jest.Mocked<ThreadsService>;
  let mockedWebsocketsService: jest.Mocked<WebsocketsService>;
  let mockedMessageRepository: jest.Mocked<MessageRepository>;
  let mockedContentService: jest.Mocked<ContentService>;
  const mockEmit = jest.fn();

  const createMessageDto = {
    answerLength: 42,
    contextLength: 42,
    question: "question",
    similarityCutoff: 42,
    temperature: 42,
    topK: 42,
  };

  const mockOrganization = {
    billingEmail: "billingEmail",
    createdAt: new Date(),
    credits: 10000,
    id: "id",
    name: "orgname",
    orgname: "orgname",
    plan: "API" as const,
    stripeCustomerId: "stripeCustomerId",
    updatedAt: new Date(),
  };

  const mockChatbot = {
    accessScope: "ORGANIZATION" as const,
    createdAt: new Date(),
    description: "description",
    documents: [],
    id: "id",
    llmBase: "GPT_4" as const,
    name: "name",
    orgname: "orgname",
    updatedAt: new Date(),
  };

  const mockThread = {
    _count: {
      message: 0,
    },
    chatbotId: "id",
    createdAt: new Date(),
    credits: 0,
    documents: [],
    id: "id",
    name: "name",
    orgname: "orgname",
    updatedAt: new Date(),
  };

  const mockMessage = {
    answer: "answer",
    answerLength: 42,
    citations: [],
    contextLength: 42,
    createdAt: new Date(),
    credits: 42,
    id: "id",
    question: "question",
    similarityCutoff: 42,
    temperature: 42,
    threadId: "id",
    topK: 42,
    updatedAt: new Date(),
  };

  beforeEach(() => {
    const { unit, unitRef } = TestBed.create(MessagesService).compile();

    messagesService = unit;
    mockedChatbotsService = unitRef.get<ChatbotsService>(ChatbotsService);
    mockedOrganizationsService =
      unitRef.get<OrganizationsService>(OrganizationsService);
    mockedThreadsService = unitRef.get<ThreadsService>(ThreadsService);
    mockedWebsocketsService = unitRef.get<WebsocketsService>(WebsocketsService);
    mockedMessageRepository = unitRef.get<MessageRepository>(MessageRepository);
    mockedContentService = unitRef.get<ContentService>(ContentService);
    mockedOrganizationsService.findOneByName.mockImplementation(async () => {
      return mockOrganization;
    });
    mockedChatbotsService.findOne.mockImplementation(async () => {
      return mockChatbot;
    });
    mockedThreadsService.findOne.mockImplementation(async () => {
      return mockThread;
    });
    mockedMessageRepository.findAll.mockImplementation(async () => {
      return {
        count: 5,
        messages: [
          mockMessage,
          mockMessage,
          mockMessage,
          mockMessage,
          mockMessage,
        ],
      };
    });
    mockedContentService.findAll.mockImplementation(async (): Promise<any> => {
      return {
        aggregates: {},
        metadata: {
          limit: 100,
          offset: 0,
          total: 3,
        },
        results: [
          {
            id: "doc1",
          },
          {
            id: "doc2",
          },
          {
            id: "doc3",
          },
        ],
      };
    });
    mockedContentService.findOne.mockImplementation(
      async (id): Promise<any> => {
        return {
          id,
        };
      }
    );

    // Mock the socket methods
    mockedWebsocketsService.socket = {
      ...mockedWebsocketsService.socket,
      to: jest.fn().mockReturnValue({ emit: mockEmit }),
    } as any;
  });

  describe("create", () => {
    it("should throw a 404 when the organization does not exist", async () => {
      mockedOrganizationsService.findOneByName.mockImplementationOnce(
        async () => {
          throw new NotFoundException();
        }
      );

      await expect(
        messagesService.create(
          "nonexistentOrgname",
          "badChatbotId",
          "badThreadId",
          createMessageDto,
          null
        )
      ).rejects.toThrow(NotFoundException);
    });

    it("should return none when no credits", async () => {
      mockedOrganizationsService.findOneByName.mockImplementationOnce(
        async () => {
          return { ...mockOrganization, credits: 0 };
        }
      );
      await messagesService.create(
        "orgname",
        "id",
        "id",
        createMessageDto,
        null
      );

      expect(mockedWebsocketsService.socket.to).toHaveBeenCalledWith("orgname");
      expect(mockedMessageRepository.create).toHaveBeenCalledWith(
        "id", // this is threadId in your code
        createMessageDto,
        "You do not have enough credits to ask this question.",
        0,
        []
      );
      expect(mockEmit).toHaveBeenCalledWith("update");
    });

    it('should update the thread name if it is the default "New Thread"', async () => {
      mockedThreadsService.findOne.mockResolvedValueOnce({
        ...mockThread,
        name: "New Thread",
      });

      await messagesService.create(
        "orgname",
        "id",
        "id",
        createMessageDto,
        null
      );

      expect(mockedThreadsService.updateThreadName).toHaveBeenCalledWith(
        "orgname",
        "id",
        createMessageDto.question
      );
    });

    it("should generate a response for the provided question", async () => {
      await messagesService.create(
        "orgname",
        "id",
        "id",
        createMessageDto,
        null
      );

      expect(mockedMessageRepository.create).toHaveBeenCalled(); // Add more specific assertions based on the method's logic.
    });

    it("should handle errors and return a message indicating a processing issue", async () => {
      mockedChatbotsService.findOne.mockImplementationOnce(() => {
        throw new Error("An error occurred");
      });

      await messagesService.create(
        "orgname",
        "id",
        "id",
        createMessageDto,
        null
      );

      expect(mockedMessageRepository.create).toHaveBeenCalledWith(
        "id", // this is threadId in your code
        createMessageDto,
        "Sorry, but I could not process your request. Please contact support if this continues.",
        0,
        []
      );
    });

    it("should filter first thread specific documents", async () => {
      mockedThreadsService.findOne.mockImplementation(async () => {
        return {
          ...mockThread,
          documents: [
            {
              id: "doc1",
              name: "doc1",
            },
            {
              id: "doc2",
              name: "doc2",
            },
          ],
        };
      });

      const answerQuestionSpy = jest.spyOn(messagesService, "answerQuestion");
      await messagesService.create(
        "orgname",
        "id",
        "id",
        createMessageDto,
        null
      );

      expect(answerQuestionSpy).toHaveBeenCalledWith(
        mockChatbot,
        {
          ...mockThread,
          documents: [
            {
              id: "doc1",
              name: "doc1",
            },
            {
              id: "doc2",
              name: "doc2",
            },
          ],
        },
        [mockMessage, mockMessage, mockMessage, mockMessage, mockMessage],
        "orgname",
        [
          {
            id: "doc1",
          },
          {
            id: "doc2",
          },
        ],
        createMessageDto,
        null
      );
    });

    it("should filter first after on agent documents if document specific", async () => {
      mockedChatbotsService.findOne.mockImplementation(async () => {
        return {
          ...mockChatbot,
          accessScope: "DOCUMENT" as const,
          documents: [
            {
              id: "doc3",
              name: "doc3",
            },
            {
              id: "doc4",
              name: "doc4",
            },
          ],
        };
      });

      const answerQuestionSpy = jest.spyOn(messagesService, "answerQuestion");
      await messagesService.create(
        "orgname",
        "id",
        "id",
        createMessageDto,
        null
      );

      expect(answerQuestionSpy).toHaveBeenCalledWith(
        {
          ...mockChatbot,
          accessScope: "DOCUMENT" as const,
          documents: [
            {
              id: "doc3",
              name: "doc3",
            },
            {
              id: "doc4",
              name: "doc4",
            },
          ],
        },
        mockThread,
        [mockMessage, mockMessage, mockMessage, mockMessage, mockMessage],
        "orgname",
        [
          {
            id: "doc3",
          },
          {
            id: "doc4",
          },
        ],
        createMessageDto,
        null
      );
    });

    it("should filter first after on organization documents if organiztion access scope specific", async () => {
      mockedChatbotsService.findOne.mockImplementation(async () => {
        return {
          ...mockChatbot,
        };
      });

      const answerQuestionSpy = jest.spyOn(messagesService, "answerQuestion");
      await messagesService.create(
        "orgname",
        "id",
        "id",
        createMessageDto,
        null
      );

      expect(answerQuestionSpy).toHaveBeenCalledWith(
        {
          ...mockChatbot,
        },
        mockThread,
        [mockMessage, mockMessage, mockMessage, mockMessage, mockMessage],
        "orgname",
        [],
        createMessageDto,
        null
      );
    });

    it("should block all document matches if document level and no docs specified", async () => {
      mockedChatbotsService.findOne.mockImplementation(async () => {
        return {
          ...mockChatbot,
          accessScope: "DOCUMENT" as const,
        };
      });

      const answerQuestionSpy = jest.spyOn(messagesService, "answerQuestion");
      await messagesService.create(
        "orgname",
        "id",
        "id",
        createMessageDto,
        null
      );

      expect(answerQuestionSpy).toHaveBeenCalledWith(
        {
          ...mockChatbot,
          accessScope: "DOCUMENT" as const,
        },
        mockThread,
        [mockMessage, mockMessage, mockMessage, mockMessage, mockMessage],
        "orgname",
        [
          {
            id: "block-matches",
          },
        ],
        createMessageDto,
        null
      );
    });
  });
});
