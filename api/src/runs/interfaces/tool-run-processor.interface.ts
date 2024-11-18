import { Logger } from "@nestjs/common";

import { ContentService } from "../../content/content.service";
import { ContentEntity } from "../../content/entities/content.entity";

// A tool run processor should take in the runId, the input contents, a logger, and the content service
// It should return the output contents
export interface IToolRunProcess {
  (
    runId: string,
    runInputContents: ContentEntity[],
    logger: Logger,
    contentService: ContentService,
    ...args: any[]
  ): Promise<ContentEntity[]>;
}
