import { HttpService } from "@nestjs/axios";
import { Injectable } from "@nestjs/common";
import { InternalServerErrorException, Logger } from "@nestjs/common";
import { AxiosError } from "axios";
import { catchError, firstValueFrom } from "rxjs";

import { retry } from "../common/retry";
import { ToolRunsService } from "../tool-runs/tool-runs.service";

@Injectable()
export class RunpodService {
  private readonly logger: Logger = new Logger("Runpod Service");

  constructor(
    private readonly httpService: HttpService,
    private readonly toolRunsService: ToolRunsService
  ) {}

  async runPod(jobId: string, podId: string, input: any) {
    // START RUNPOD JOB
    this.logger.log("Running runpod");
    const { data: runpodCreateResponse } = await retry(
      this.logger,
      () =>
        firstValueFrom(
          this.httpService
            .post(`https://api.runpod.ai/v2/${podId}/run`, input, {
              headers: {
                Authorization:
                  "Bearer E0L11W179IXEBVJ09878F0ICVMNZD6JFSTWE7MZP",
                "Content-Type": "application/json",
              },
            })
            .pipe(
              catchError((err: AxiosError) => {
                this.logger.error("Could not hit runpod endpoint", err.message);
                throw new InternalServerErrorException(err.message);
              })
            )
        ),
      5
    );

    // CHECK RUNPOD JOB STATUS
    const runpodJobId = runpodCreateResponse.id;
    let firstAttempt = true;
    while (true) {
      await new Promise((resolve) => setTimeout(resolve, 5000));
      const { data: rundpodCheckJobResponse } = await retry(
        this.logger,
        () =>
          firstValueFrom(
            this.httpService
              .get(`https://api.runpod.ai/v2/${podId}/status/` + runpodJobId, {
                headers: {
                  Authorization:
                    "Bearer E0L11W179IXEBVJ09878F0ICVMNZD6JFSTWE7MZP",
                  "Content-Type": "application/json",
                },
              })
              .pipe(
                catchError((err: AxiosError) => {
                  this.logger.error(
                    "Could not hit runpod endpoint",
                    err.message
                  );
                  throw new InternalServerErrorException(err.message);
                })
              )
          ),
        5
      );
      this.logger.log(
        "Got runpod response: " +
          JSON.stringify(rundpodCheckJobResponse, null, 2)
      );
      if (rundpodCheckJobResponse.status == "COMPLETED") {
        return rundpodCheckJobResponse.output;
      } else if (rundpodCheckJobResponse.status === "IN_PROGRESS") {
        if (firstAttempt) {
          // mark as processing
          await this.toolRunsService.setStatus(jobId, "PROCESSING");
          firstAttempt = false;
        }

        await this.toolRunsService.setProgress(
          jobId,
          Number(rundpodCheckJobResponse.output) || 0.5
        );
      } else if (rundpodCheckJobResponse.status === "FAILED") {
        throw new InternalServerErrorException("Runpod job failed");
      }
    }
  }
}
