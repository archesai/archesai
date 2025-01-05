export class CreateChatCompletionDto {
  best_of?: number
  frequency_penalty?: number
  ignore_eos?: boolean
  logit_bias?: Record<string, number>
  max_tokens?: number
  messages: MessageDto[]
  model?: string
  n?: number = 1
  name?: string
  presence_penalty?: number
  skip_special_tokens?: boolean
  stop?: string[]
  stop_token_ids?: number[]
  stream?: boolean
  temperature?: number = 0.7
  top_k?: number = -1
  top_p?: number = 1
  use_beam_search?: boolean
  user?: string
}

export class MessageDto {
  content: string
  name?: string
  role: 'assistant' | 'function' | 'system' | 'user'
}
