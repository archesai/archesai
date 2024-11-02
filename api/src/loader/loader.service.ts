import { HttpService } from "@nestjs/axios";
import { Injectable } from "@nestjs/common";
import { BadRequestException } from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { AxiosError } from "axios";
import * as fs from "fs";
import GPT3Tokenizer from "gpt3-tokenizer";
import { catchError, firstValueFrom } from "rxjs";
import { UnstructuredClient } from "unstructured-client";
import { Strategy } from "unstructured-client/sdk/models/shared";

const key = "XZL0F7FKDFYGGOQUpUVW4sSW9pRhTp";

@Injectable()
export class LoaderService {
  private readonly logger: Logger = new Logger("LoaderService");

  constructor(
    private readonly httpService: HttpService,
    private readonly configService: ConfigService
  ) {}

  async extractUrl(
    url: string,
    chunkSize: number,
    delimiter?: string
  ): Promise<{
    mimeType: string;
    preview: string;
    textContent: { page: number; text: string; tokens: number }[];
    title: string;
    totalTokens: number;
  }> {
    this.logger.log(
      `Extracting text from ${url.split("?")[0]} with ${chunkSize} chunk size and ${delimiter} delimiter`
    );

    try {
      const { data } = await firstValueFrom(
        this.httpService
          .post(this.configService.get("LOADER_ENDPOINT"), {
            chunkSize,
            delimiter,
            url,
          })
          .pipe(
            catchError((err: AxiosError) => {
              this.logger.error(
                "Error hitting loader endpoint: " + err.message
              );
              throw new BadRequestException();
            })
          )
      );

      const { contentType, preview, textContent, title } = data as {
        contentType: string;
        preview: string;
        textContent: { page: number; text: string; tokens: number }[];
        title: string;
      };

      const sanitizedTextContent = textContent.map((data) => ({
        ...data,
        text: data.text
          .replaceAll(/\0/g, "")
          .replaceAll(/[^ -~\u00A0-\uD7FF\uE000-\uFDCF\uFDF0-\uFFFD\n]/g, ""),
      }));

      return {
        mimeType: contentType,
        preview,
        textContent: sanitizedTextContent,
        title,
        totalTokens: textContent.reduce((acc, curr) => acc + curr.tokens, 0),
      };
    } catch (err) {
      this.logger.error("Error hitting loader endpoint: " + err.message);
      throw new BadRequestException();
    }
  }

  getFirstTokens(textContent: string[], maxLen: number) {
    const tokenizer = new GPT3Tokenizer({ type: "gpt3" });

    let tokens = "";
    let currentLen = 0;

    for (const text of textContent) {
      tokens += " " + text;
      currentLen += tokenizer.encode(text).text.length;
      if (currentLen > maxLen) {
        break;
      }
    }
    return tokens;
  }

  async ingestFile(filename: string) {
    const client = new UnstructuredClient({
      security: {
        apiKeyAuth: key,
      },
      serverURL: "https://api.unstructuredapp.io",
    });

    const data = fs.readFileSync(filename);

    try {
      const result = await client.general.partition({
        partitionParameters: {
          files: {
            content: data,
            fileName: filename,
          },
          strategy: Strategy.HiRes,
        },
      });

      if (result.statusCode == 200) {
        console.log(result.elements);
      }
    } catch (err) {
      console.log(err);
    }
  }
}
