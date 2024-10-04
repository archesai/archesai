import { Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { kmeans } from "ml-kmeans";

@Injectable()
export class KMeansService {
  private readonly logger: Logger = new Logger("KMeansService");

  getClusters(embeddings: number[][], clusters: number) {
    const res = kmeans(embeddings, clusters, {});
    return res;
  }
}
