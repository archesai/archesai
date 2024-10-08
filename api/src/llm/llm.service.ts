import { Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import OpenAI from "openai";
import { zodResponseFormat } from "openai/helpers/zod";
import { ChatCompletionCreateParamsStreaming } from "openai/resources";
import { z } from "zod";

import { CreateChatCompletionDto } from "./dto/create-chat-completion.dto";

// Schema for Lesson Metadata
const LessonMetadataSchema = z.object({
  description: z.string(),
  lessonNumber: z.number().int().positive(),
  nativeLanguage: z.string(),
  targetLanguage: z.string(),
  title: z.string(),
});

// Schema for Dialogues
const DialogueSchema = z.object({
  audioResources: z.array(z.string()).optional(),
  dialogueId: z.string().uuid(),
  originalText: z.string(),
  translation: z.string(),
});

// Schema for Vocabulary
const VocabularyItemSchema = z.object({
  pronunciation: z.string().optional(),
  translation: z.string(),
  usageExamples: z.array(z.string()).optional(),
  word: z.string(),
});

// Schema for Grammar Explanations
const GrammarPointSchema = z.object({
  examples: z.array(z.string()),
  explanation: z.string(),
  grammarPoint: z.string(),
});

// Schema for Exercises
const ExerciseSchema = z.object({
  answers: z.array(z.string()).optional(),
  exerciseId: z.string().uuid(),
  instructions: z.string(),
  questions: z.array(z.string()),
  type: z.enum([
    "multiple-choice",
    "fill-in-the-blank",
    "matching",
    "short-answer",
  ]),
});

// Schema for Cultural Notes
const CulturalNoteSchema = z.object({
  content: z.string(),
  topic: z.string(),
});

// Schema for Summary
const SummarySchema = z.object({
  keyTakeaways: z.array(z.string()),
  reviewQuestions: z.array(z.string()).optional(),
});

// Schema for Additional Resources
const AdditionalResourceSchema = z.object({
  links: z.array(z.string().url()),
  references: z.array(z.string()),
});

// Complete Lesson Schema
const LessonSchema = z.object({
  additionalResources: AdditionalResourceSchema.optional(),
  culturalNotes: z.array(CulturalNoteSchema).optional(),
  dialogues: z.array(DialogueSchema),
  exercises: z.array(ExerciseSchema),
  grammar: z.array(GrammarPointSchema),
  metadata: LessonMetadataSchema,
  summary: SummarySchema,
  vocabulary: z.array(VocabularyItemSchema),
});

@Injectable()
export class LLMService {
  private readonly logger: Logger = new Logger("LLMService");

  public openai: OpenAI;

  constructor(private configService: ConfigService) {
    this.openai = new OpenAI({
      apiKey:
        this.configService.get("LLM_TYPE") == "openai"
          ? this.configService.get("OPEN_AI_KEY")
          : "ollama",
      baseURL:
        this.configService.get("LLM_TYPE") == "openai"
          ? undefined
          : this.configService.get("OLLAMA_ENDPOINT"),
      organization: "org-uCtGHWe8lpVBqo5thoryOqcS",
    });
  }

  async createAssimil() {
    this.logger.log("Creating Assimil Lesson");
    const completion = await this.openai.beta.chat.completions.parse({
      messages: [
        { content: "Extract the event information.", role: "system" },
        {
          content: "Create me an assimil lesson in Spanish",
          role: "user",
        },
      ],
      model: "gpt-4o",
      response_format: zodResponseFormat(LessonSchema, "lesson"),
    });

    const lesson = completion.choices[0].message;

    // If the model refuses to respond, you will get a refusal message
    if (lesson.refusal) {
      this.logger.error(lesson.refusal);
    } else {
      this.logger.log("Assimil Lesson: " + JSON.stringify(lesson, null, 2));
    }
  }

  async createChatCompletion(
    createChatCompletionDto: CreateChatCompletionDto,
    emitAnswer: (answer: string) => void
  ) {
    this.logger.log(
      "Sending messages to OpenAI: " +
        JSON.stringify(createChatCompletionDto, null, 2)
    );

    let answer = "";
    const stream = await this.openai.chat.completions.create({
      ...(createChatCompletionDto as ChatCompletionCreateParamsStreaming),
      model:
        this.configService.get("LLM_TYPE") == "openai" ? "gpt-4o" : "llama3.1",
      stream: true,
    });

    for await (const part of stream) {
      const content = part.choices[0].delta.content;
      if (content) {
        answer = answer.concat(content);
        emitAnswer(answer);
      }
    }

    this.logger.log("Received Answer: " + answer);
    return answer;
  }

  async createImageSummary(imageUrl: string) {
    const response = await this.openai.chat.completions.create({
      messages: [
        {
          content: [
            { text: "What’s in this image?", type: "text" },
            {
              image_url: {
                url: imageUrl,
              },
              type: "image_url",
            },
          ],
          role: "user",
        },
      ],
      model: "gpt-4o",
    });
    return response.choices[0];
  }

  async createSummary(text: string) {
    const { choices, usage } = await this.openai.completions.create({
      frequency_penalty: 0,
      max_tokens: 80,
      model:
        this.configService.get("LLM_TYPE") == "openai"
          ? "gpt-3.5-turbo-instruct"
          : "llama3.1",
      presence_penalty: 0,
      prompt: `Write a very short few word summary describing what this document is based on a part of its content. It could be a book, a legal document, a textbook, a newspaper, a bank statement, or another document like this.\n\nContent:\n${text}\n\n---\n\nSummary:`,
      temperature: 0.3,
      top_p: 1,
    });

    return {
      summary: (choices[0].text as string).trim(),
      tokens: usage.total_tokens as number,
    };
  }
}
