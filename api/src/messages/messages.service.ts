import { Inject, Injectable, Logger } from "@nestjs/common";
import { NotFoundException } from "@nestjs/common";
import { Chatbot, Content, Message, Thread } from "@prisma/client";
import GPT3Tokenizer from "gpt3-tokenizer";

import { ChatbotsService } from "../chatbots/chatbots.service";
import { retry } from "../common/retry";
import { SortByField, SortDirection } from "../common/search-query";
import { OpenAiCompletionsService } from "../completions/completions.openai.service";
import { ContentService } from "../content/content.service";
import { OpenAiEmbeddingsService } from "../embeddings/embeddings.openai.service";
import { MessageQueryDto } from "../messages/dto/message-query.dto";
import { MessageEntity } from "../messages/entities/message.entity";
import { OrganizationsService } from "../organizations/organizations.service";
import { ThreadsService } from "../threads/threads.service";
import {
  VECTOR_DB_SERVICE,
  VectorDBService,
} from "../vector-db/vector-db.service";
import { VectorRecordService } from "../vector-records/vector-record.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateMessageDto } from "./dto/create-message.dto";
import { MessageRepository } from "./message.repository";

@Injectable()
export class MessagesService {
  private readonly logger: Logger = new Logger("Messages Service");

  constructor(
    private messageRepository: MessageRepository,
    private threadsService: ThreadsService,
    private chatbotsService: ChatbotsService,
    private websocketsService: WebsocketsService,
    private organizationsService: OrganizationsService,
    @Inject(VECTOR_DB_SERVICE)
    private vectorDBService: VectorDBService,
    private openAiCompletionsService: OpenAiCompletionsService,
    private openAiEmbeddingsService: OpenAiEmbeddingsService,
    private contentService: ContentService,
    private vectorRecordService: VectorRecordService
  ) {}

  async answerQuestion(
    chatbot: Chatbot,
    thread: Thread,
    messages: Message[],
    orgname: string,
    content: Content[],
    createMessageDto: CreateMessageDto,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    signal: AbortSignal
  ) {
    this.logger.log("Answering question: " + createMessageDto.question);

    // Get memory
    const memory = [] as { answer: string; question: string }[];
    messages.forEach((i) => {
      memory.push({
        answer: i.answer,
        question: i.question,
      });
    });

    const emitAnswer = (answer: string) => {
      this.websocketsService.socket.to(orgname).emit("chat", {
        answer,
        chatbotId: chatbot.id,
        message: new MessageEntity({
          answer: answer,
          answerLength: createMessageDto.answerLength,
          citations: [],
          contextLength: createMessageDto.contextLength,
          createdAt: new Date(),
          credits: 0,
          id: "random",
          question: createMessageDto.question,
          similarityCutoff: createMessageDto.similarityCutoff,
          temperature: createMessageDto.temperature,
          threadId: thread.id,
          topK: createMessageDto.topK,
        }),
        orgname,
        threadId: thread.id,
      });
    };

    const { citations, context } = await this.createContext(
      orgname,
      content,
      createMessageDto
    );

    const completionMessages = [
      {
        content: chatbot.description,
        role: "assistant",
      },
    ];
    for (const m of memory) {
      completionMessages.push({
        content: m.question,
        role: "user",
      });
      completionMessages.push({
        content: m.answer,
        role: "assistant",
      });
    }
    if (context.length > 0) {
      completionMessages.push({
        content: context,
        role: "user",
      });
    }
    completionMessages.push({
      content: createMessageDto.question,
      role: "user",
    });

    const answer = await retry(
      this.logger,
      async () =>
        await this.openAiCompletionsService.createChatCompletion(
          {
            max_tokens:
              createMessageDto.contextLength + createMessageDto.answerLength,
            messages: completionMessages as any,
            temperature: createMessageDto.temperature,
          },
          (answer) => emitAnswer(answer)
        ),
      1
    );

    this.logger.log("Got final answer: " + answer);

    const tokenizer = new GPT3Tokenizer({ type: "gpt3" });
    const tokens =
      tokenizer.encode(context).text.length +
      tokenizer.encode(answer).text.length +
      tokenizer.encode(createMessageDto.question).text.length;

    this.logger.log("Used: " + tokens + " tokens in answering question");
    return {
      answer,
      citations,
      tokens: tokens,
    };
  }

  async create(
    orgname: string,
    chatbotId: string,
    threadId: string,
    createMessageDto: CreateMessageDto,
    signal: AbortSignal
  ) {
    try {
      const chatbot = await this.chatbotsService.findOne(chatbotId);
      const thread = await this.threadsService.findOne(
        orgname,
        chatbotId,
        threadId
      );
      const organization =
        await this.organizationsService.findOneByName(orgname);

      // Ensure credits are available
      if (organization.credits <= 0) {
        const message = await this.messageRepository.create(
          threadId,
          createMessageDto,
          "You do not have enough credits to ask this question.",
          0,
          []
        );
        this.websocketsService.socket.to(orgname).emit("update");
        return message;
      }

      // Update Thread Name if still default
      if (thread.name == "New Thread") {
        await this.threadsService.updateThreadName(
          orgname,
          threadId,
          createMessageDto.question
        );
      }

      // Get messages
      const messages = await this.messageRepository.findAll(orgname, threadId, {
        limit: 5,
        sortBy: SortByField.CREATED,
        sortDirection: SortDirection.DESCENDING,
      });
      this.logger.log("Got messages");

      const documents = [];
      // Get answer
      const { answer, citations, tokens } = await this.answerQuestion(
        chatbot,
        thread,
        messages.results.reverse(),
        orgname,
        documents,
        createMessageDto,
        signal
      );
      this.logger.log("Completed question, saving message");

      // Create answer in db
      const message = await this.messageRepository.create(
        threadId,
        createMessageDto,
        answer,
        tokens,
        citations
      );

      let multiple = 1;
      if (chatbot.llmBase === "GPT_3_5_TURBO_16_K") {
        multiple =
          createMessageDto.contextLength + createMessageDto.answerLength < 3800
            ? 1
            : 2;
      } else {
        multiple =
          createMessageDto.contextLength + createMessageDto.answerLength < 7500
            ? 20
            : 45;
      }
      await this.organizationsService.removeCredits(orgname, multiple * tokens);

      // Add credits to thread total
      await this.threadsService.incrementCredits(
        orgname,
        threadId,
        multiple * tokens
      );
      this.websocketsService.socket.to(orgname).emit("update");

      return message;
    } catch (err) {
      if (err instanceof NotFoundException) {
        throw err;
      }
      this.logger.error(err);
      const message = await this.messageRepository.create(
        threadId,
        createMessageDto,
        "Sorry, but I could not process your request. Please contact support if this continues.",
        0,
        []
      );
      this.websocketsService.socket.to(orgname).emit("update");
      return message;
    }
  }

  async createContext(
    orgname: string,
    content: Content[],
    createMessageDto: CreateMessageDto
  ) {
    const [questionEmbedding] = await retry(
      this.logger,
      async () =>
        await this.openAiEmbeddingsService.createEmbeddings([
          createMessageDto.question,
        ]),
      3
    );
    const queryResult = await this.vectorDBService.query(
      orgname,
      questionEmbedding.embedding,
      createMessageDto.topK,
      content.map((content) => ({ contentId: content.id }))
    );
    this.logger.log("Got query result: " + JSON.stringify(queryResult));

    const discoveredContent = {} as { [key: string]: Content };
    const citations = [] as {
      contentId: string;
      similarity: number;
      text: string;
    }[];

    const tokenizer = new GPT3Tokenizer({ type: "gpt3" });
    let currentLen = 0;

    let highestSimilarity = 0;
    if (queryResult.length > 0) {
      highestSimilarity = queryResult[0].score;
    }

    for (const match of queryResult) {
      if (match.score < createMessageDto.similarityCutoff) {
        continue;
      }
      if (highestSimilarity - match.score > 0.025) {
        continue;
      }
      const contentId = match.id.split("__")[0];
      const vectorRecord = await this.vectorRecordService.findOne(match.id);
      const text = vectorRecord.text;
      if (!text) {
        continue;
      }
      currentLen += tokenizer.encode(text).text.length;
      // break if we are gonna go over the limit
      if (currentLen > 14000) {
        break;
      }
      this.logger.log("Adding text segment: " + text);
      const content = await this.contentService.findOne(contentId);
      discoveredContent[contentId] = content;
      citations.push({
        contentId,
        similarity: match.score,
        text: text,
      });
      if (currentLen > createMessageDto.contextLength) {
        break;
      }
    }

    let context = ""; // by default, we don't include any context
    if (citations.length > 0) {
      context =
        `**Context**: Below you'll find brief summaries and excerpts from a selection of content relevant to your inquiry.` +
        Object.values(discoveredContent).map((doc) => {
          const contentCitations = citations.filter(
            (c) => c.contentId === doc.id
          );
          return (
            `

- **Title**: ${doc.name}
- **Summary**: //FIXME{doc.summary}

**Excerpts from ${doc.name}**:
` +
            contentCitations.map((excerpts, i) => {
              return `

**Excerpt ${i + 1}**: ${excerpts.text.replace(/\n/g, " ")}
`;
            })
          );
        });
    }

    return {
      citations,
      context,
    };
  }

  findAll(
    orgname: string,
    chatbotId: string,
    threadId: string,
    messageQueryDto: MessageQueryDto
  ) {
    return this.messageRepository.findAll(orgname, threadId, messageQueryDto);
  }
}
