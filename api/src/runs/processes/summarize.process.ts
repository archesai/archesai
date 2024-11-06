import { Logger } from "@nestjs/common";

import { retry } from "../../common/retry";
import { ContentService } from "../../content/content.service";
import { ContentEntity } from "../../content/entities/content.entity";
import { LLMService } from "../../llm/llm.service";

export const processSummarize = async (
  runId: string,
  runInputContents: ContentEntity[],
  logger: Logger,
  contentService: ContentService,
  llmService: LLMService
): Promise<ContentEntity[]> => {
  logger.log(`Summarizing content for run ${runId}`);
  const start = Date.now();
  const c = runInputContents
    .map((x) => x.text)
    .filter((x) => x)
    .join(" ");
  logger.log(`Got first tokens for content for run ${runId}`);
  const { summary } = await retry(
    logger,
    async () => await llmService.createSummary(c),
    3
  );
  logger.log(`Got summary for content for run ${runId}`);

  logger.log(
    "Summary saved. Completed in " + (Date.now() - start) / 1000 + "s"
  );

  const summaryContent = await contentService.create(
    runInputContents[0].orgname,
    {
      name: "Summary Tool - " + runInputContents.map((x) => x.name).join(", "),
      text: summary,
    }
  );

  return [new ContentEntity(summaryContent)];
};
