export const VECTOR_DB_SERVICE = "VECTOR_DB_SERVICE";

export interface VectorDBService {
  deleteMany(orgname: string, ids: string[]): Promise<void>;
  fetchAll(
    orgname: string,
    ids: string[]
  ): Promise<{
    vectors: {
      [vectorId: string]: any;
    };
  }>;
  getSimilarity(A: number[], B: number[]): number;
  query(
    orgname: string,
    embedding: number[],
    topK: number,
    content?: { contentId: string }[]
  ): Promise<
    {
      id: string;
      score: number;
    }[]
  >;
  upsert(
    orgname: string,
    contentId: string,
    embeddings: number[][],
    texts?: string[]
  ): Promise<void>;
}

export class BaseVectorDBService {
  getSimilarity(A: number[], B: number[]) {
    let dotproduct = 0;
    let mA = 0;
    let mB = 0;
    for (let i = 0; i < A.length; i++) {
      // here you missed the i++
      dotproduct += A[i] * B[i];
      mA += A[i] * A[i];
      mB += B[i] * B[i];
    }
    mA = Math.sqrt(mA);
    mB = Math.sqrt(mB);
    const similarity = dotproduct / (mA * mB); // here you needed extra brackets
    return similarity;
  }
}
