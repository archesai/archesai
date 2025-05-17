import type { ChatCompletionCreateParamsStreaming } from 'openai/resources/chat/completions'

import OpenAI from 'openai'

import type { ConfigService, HealthCheck, HealthStatus } from '@archesai/core'

import { Logger } from '@archesai/core'

/**
 * Service for interacting with the OpenAI Language Model API.
 */
export class LlmService implements HealthCheck {
  private readonly configService: ConfigService
  private readonly health: HealthStatus
  private readonly logger = new Logger(LlmService.name)
  private readonly openai: OpenAI

  constructor(configService: ConfigService) {
    this.configService = configService
    this.health = {
      status: 'COMPLETED'
    }
    this.openai = new OpenAI({
      apiKey: this.configService.get('llm.token'),
      baseURL: this.configService.get('llm.endpoint'),
      organization: 'org-uCtGHWe8lpVBqo5thoryOqcS'
    })
  }

  public async createChatCompletion(
    createChatCompletionDto: ChatCompletionCreateParamsStreaming,
    emitAnswer: (answer: string) => void
  ) {
    this.logger.debug('Sending messages to OpenAI', { createChatCompletionDto })

    let answer = ''
    const stream = await this.openai.chat.completions.create({
      ...createChatCompletionDto,
      model:
        this.configService.get('llm.type') == 'openai' ? 'gpt-4o' : 'llama3.1',
      stream: true
    })

    for await (const part of stream) {
      const content = part.choices[0]?.delta.content
      if (content) {
        answer = answer.concat(content)
        emitAnswer(answer)
      }
    }

    this.logger.debug('Received Answer: ' + answer)
    return answer
  }

  /**
   * Generates embeddings for an array of input texts using the configured LLM model.
   * The method dynamically selects the model based on the `llm.type` configuration.
   * It logs the number of texts embedded, total tokens used, and the time taken for the operation.
   * @param texts - An array of strings to generate embeddings for.
   * @returns A promise that resolves to an array of objects, each containing:
   *   - `embedding`: The generated embedding vector for the corresponding input text.
   *   - `tokens`: The number of tokens used for the embedding, averaged per text.
   * @throws Will throw an error if the embedding generation fails.
   */
  public async createEmbeddings(texts: string[]) {
    const start = Date.now()
    const { data, usage } = await this.openai.embeddings.create({
      input: texts,
      model:
        this.configService.get('llm.type') == 'openai'
          ? 'text-embedding-ada-002'
          : 'mxbai-embed-large'
    })
    const response = data.map((d) => {
      return {
        embedding: d.embedding,
        tokens: Math.ceil(usage.total_tokens / texts.length)
      }
    })
    this.logger.debug(
      `embedded ${texts.length.toString()} texts with ${usage.total_tokens.toString()} in ${((Date.now() - start) / 1000).toString()}s`
    )
    return response
  }

  /**
   * Generates a summary of the content in the provided image URL using OpenAI's GPT-4 model.
   * @param url - The URL of the image to be analyzed.
   * @returns A promise that resolves to the first choice of the response from the OpenAI API.
   * @throws Will throw an error if the OpenAI API request fails.
   */
  public async createImageSummary(url: string) {
    const response = await this.openai.chat.completions.create({
      messages: [
        {
          content: [
            { text: 'What’s in this image?', type: 'text' },
            {
              image_url: {
                url: url
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

  /**
   * Generates a concise summary of the provided text using an AI language model.
   * @param text - The input text to summarize. This could be content from a book, legal document,
   * textbook, newspaper, bank statement, or other types of documents.
   * @returns An object containing:
   * - `summary`: A one to two sentence summary of the input text.
   * - `tokens`: The total number of tokens used during the operation.
   * @throws An error if the AI model does not return a valid summary.
   */
  public async createSummary(text: string) {
    const { choices, usage } = await this.openai.completions.create({
      frequency_penalty: 0,
      max_tokens: 80,
      model:
        this.configService.get('llm.type') == 'openai'
          ? 'gpt-3.5-turbo-instruct'
          : 'llama3.1',
      presence_penalty: 0,
      prompt: `Write a very short one to two sentance summary describing what this document is based on a part of its content. It could be a book, a legal document, a textbook, a newspaper, a bank statement, or another document like this.\n\nContent:\n${text}\n\n---\n\nSummary:`,
      temperature: 0.3,
      top_p: 1
    })

    if (!choices[0]) {
      throw new Error('summary returned empty')
    }

    return {
      summary: choices[0].text.trim(),
      tokens: usage?.total_tokens ?? 0
    }
  }

  public getHealth(): HealthStatus {
    return this.health
  }
}
