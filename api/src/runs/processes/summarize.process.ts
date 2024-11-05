import { Logger } from "@nestjs/common";

import { retry } from "../../common/retry";
import { ContentService } from "../../content/content.service";
import { ContentEntity } from "../../content/entities/content.entity";
import { LLMService } from "../../llm/llm.service";
import { LoaderService } from "../../loader/loader.service";

export const processSummarize = async (
  content: ContentEntity,
  logger: Logger,
  loaderService: LoaderService,
  contentService: ContentService,
  llmService: LLMService
) => {
  const textContent = await contentService.findAll(content.id, {});
  const start = Date.now();
  const c = loaderService.getFirstTokens(
    textContent.results.map((x) => x.text),
    3000
  );
  const { summary } = await retry(
    logger,
    async () => await llmService.createSummary(c),
    3
  );

  logger.log(`Got summary for content for ${content.name}`);

  await contentService.updateRaw(content.orgname, content.id, {
    description: summary,
  });

  logger.log(
    "Summary saved. Completed in " + (Date.now() - start) / 1000 + "s"
  );
};
