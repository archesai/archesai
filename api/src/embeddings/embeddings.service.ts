export interface EmbeddingsService {
  createEmbeddings(texts: string[]): Promise<
    {
      embedding: number[]
      tokens: number
    }[]
  >
}
