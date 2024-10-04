import { CreateChatCompletionDto } from "./dto/create-chat-completion.dto";

export interface CompletionsService {
  createChatCompletion(
    createChatCompletionDto: CreateChatCompletionDto,
    // This is to send a notification as the answer is streamed
    emitAnswer: (answer: string) => void,
  ): Promise<string>;
}
