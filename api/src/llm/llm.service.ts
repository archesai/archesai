import { Injectable } from '@nestjs/common'
import { Logger } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import OpenAI from 'openai'
import { ChatCompletionCreateParamsStreaming } from 'openai/resources'

import { CreateChatCompletionDto } from './dto/create-chat-completion.dto'

@Injectable()
export class LLMService {
  public openai: OpenAI

  private readonly logger: Logger = new Logger('LLMService')

  constructor(private configService: ConfigService) {
    this.openai = new OpenAI({
      apiKey: this.configService.get('LLM_TYPE') == 'openai' ? this.configService.get('OPEN_AI_KEY') : 'ollama',
      baseURL: this.configService.get('LLM_TYPE') == 'openai' ? undefined : this.configService.get('OLLAMA_ENDPOINT'),
      organization: 'org-uCtGHWe8lpVBqo5thoryOqcS'
    })
  }

  async createChatCompletion(createChatCompletionDto: CreateChatCompletionDto, emitAnswer: (answer: string) => void) {
    this.logger.log('Sending messages to OpenAI: ' + JSON.stringify(createChatCompletionDto, null, 2))

    let answer = ''
    const stream = await this.openai.chat.completions.create({
      ...(createChatCompletionDto as ChatCompletionCreateParamsStreaming),
      model: this.configService.get('LLM_TYPE') == 'openai' ? 'gpt-4o' : 'llama3.1',
      stream: true
    })

    for await (const part of stream) {
      const content = part.choices[0].delta.content
      if (content) {
        answer = answer.concat(content)
        emitAnswer(answer)
      }
    }

    this.logger.log('Received Answer: ' + answer)
    return answer
  }

  async createImageSummary(imageUrl: string) {
    const response = await this.openai.chat.completions.create({
      messages: [
        {
          content: [
            { text: 'What’s in this image?', type: 'text' },
            {
              image_url: {
                url: imageUrl
              },
              type: 'image_url'
            }
          ],
          role: 'user'
        }
      ],
      model: 'gpt-4o'
    })
    return response.choices[0]
  }

  async createSummary(text: string) {
    const { choices, usage } = await this.openai.completions.create({
      frequency_penalty: 0,
      max_tokens: 80,
      model: this.configService.get('LLM_TYPE') == 'openai' ? 'gpt-3.5-turbo-instruct' : 'llama3.1',
      presence_penalty: 0,
      prompt: `Write a very short one to two sentance summary describing what this document is based on a part of its content. It could be a book, a legal document, a textbook, a newspaper, a bank statement, or another document like this.\n\nContent:\n${text}\n\n---\n\nSummary:`,
      temperature: 0.3,
      top_p: 1
    })

    return {
      summary: (choices[0].text as string).trim(),
      tokens: usage.total_tokens as number
    }
  }
}
