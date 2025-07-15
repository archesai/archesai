import type { ChatCompletionCreateParamsStreaming } from 'openai/resources/chat/completions'

import OpenAI from 'openai'

import type { ConfigService, Logger } from '@archesai/core'

export const createLlmService = (
  configService: ConfigService,
  logger: Logger
) => {
  const openai = new OpenAI({
    apiKey: configService.get('llm.token'),
    baseURL: configService.get('llm.endpoint'),
    organization: 'org-uCtGHWe8lpVBqo5thoryOqcS'
  })

  return {
    async createChatCompletion(
      createChatCompletionDto: ChatCompletionCreateParamsStreaming,
      emitAnswer: (answer: string) => void
    ): Promise<string> {
      logger.debug('Sending messages to OpenAI', { createChatCompletionDto })

      let answer = ''
      const stream = await openai.chat.completions.create({
        ...createChatCompletionDto,
        model:
          configService.get('llm.type') == 'openai' ? 'gpt-4o' : 'llama3.1',
        stream: true
      })

      for await (const part of stream) {
        const content = part.choices[0]?.delta.content
        if (content) {
          answer = answer.concat(content)
          emitAnswer(answer)
        }
      }

      logger.debug('Received Answer: ' + answer)
      return answer
    },

    async createEmbeddings(texts: string[]): Promise<
      {
        embedding: number[]
        tokens: number
      }[]
    > {
      const start = Date.now()
      const { data, usage } = await openai.embeddings.create({
        input: texts,
        model:
          configService.get('llm.type') == 'openai' ?
            'text-embedding-ada-002'
          : 'mxbai-embed-large'
      })
      const response = data.map((d) => {
        return {
          embedding: d.embedding,
          tokens: Math.ceil(usage.total_tokens / texts.length)
        }
      })
      logger.debug(
        `embedded ${texts.length.toString()} texts with ${usage.total_tokens.toString()} in ${((Date.now() - start) / 1000).toString()}s`
      )
      return response
    },

    async createImageSummary(
      url: string
    ): Promise<OpenAI.Chat.Completions.ChatCompletion.Choice | undefined> {
      const response = await openai.chat.completions.create({
        messages: [
          {
            content: [
              { text: "What's in this image?", type: 'text' },
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
    },

    async createSummary(text: string): Promise<{
      summary: string
      tokens: number
    }> {
      const { choices, usage } = await openai.completions.create({
        frequency_penalty: 0,
        max_tokens: 80,
        model:
          configService.get('llm.type') == 'openai' ?
            'gpt-3.5-turbo-instruct'
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
  }
}

export type LlmService = ReturnType<typeof createLlmService>
