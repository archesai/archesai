import { HttpService } from "@nestjs/axios";
import { Inject, Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { InternalServerErrorException } from "@nestjs/common";
import axios from "axios";
import { AxiosError } from "axios";
import ffmpeg from "fluent-ffmpeg";
import * as fs from "fs";
import * as os from "os";
import * as ospath from "path";
import { catchError, firstValueFrom } from "rxjs";
import { v4 } from "uuid";

import { retry } from "../common/retry";
import { STORAGE_SERVICE, StorageService } from "../storage/storage.service";
import { KeyframesService } from "./keyframes.service";

@Injectable()
export class AudioService {
  private readonly logger: Logger = new Logger("Audio Service");

  getKeyframes = async (
    audioUrl: string,
    framerate: number,
    fn: string,
    isTranslation: boolean
  ) => {
    return this.keyframesService.getKeyframes(
      audioUrl,
      framerate,
      fn,
      isTranslation
    );
  };

  splitAudio = async (audioUrl: string) => {
    this.logger.log("Hitting moises' API...");
    const { data: moisesResponse } = await retry(
      this.logger,
      () =>
        firstValueFrom(
          this.httpService
            .post(
              "https://developer-api.moises.ai/api/job",
              {
                name: v4(),
                params: {
                  inputUrl: audioUrl,
                },
                workflow: "arches-workflow",
              },
              {
                headers: {
                  Authorization: "5fa360fa-9974-47fc-bcb8-39142bf4dcea",
                  "Content-Type": "application/json",
                },
              }
            )
            .pipe(
              catchError((err: AxiosError) => {
                this.logger.error("Could not hit moises endpoint", err.message);
                throw new InternalServerErrorException(err.message);
              })
            )
        ),
      5
    );
    this.logger.log("Moises response: " + JSON.stringify(moisesResponse));
    const moisesJobId = moisesResponse.id;
    while (true) {
      this.logger.log("Checking moises job status...");
      await new Promise((resolve) => setTimeout(resolve, 5000));
      const { data: moisesCheckJobResponse } = await retry(
        this.logger,
        () =>
          firstValueFrom(
            this.httpService
              .get("https://developer-api.moises.ai/api/job/" + moisesJobId, {
                headers: {
                  Authorization: "5fa360fa-9974-47fc-bcb8-39142bf4dcea",
                },
              })
              .pipe(
                catchError((err: AxiosError) => {
                  this.logger.error(
                    "Could not hit moises endpoint",
                    err.message
                  );
                  throw new InternalServerErrorException(err.message);
                })
              )
          ),
        5
      );
      this.logger.log("Got status from moises" + moisesCheckJobResponse.status);
      if (moisesCheckJobResponse.status == "SUCCEEDED") {
        const bassSrc = moisesCheckJobResponse.result.Bass;
        const drumsSrc = moisesCheckJobResponse.result.Drums;
        this.logger.log("Bass src: " + bassSrc);
        this.logger.log("Drums src: " + drumsSrc);

        return { bassSrc, drumsSrc };
      } else if (moisesCheckJobResponse.status === "FAILED") {
        throw new Error("Moises' job failed");
      }
    }
  };

  trimAudio = async (
    orgname: string,
    audioUrl: string,
    startTime: number,
    duration: number
  ): Promise<string> => {
    const inputTmpPath = ospath.join(os.tmpdir(), "original.mp3");
    const outputTmpPath = ospath.join(os.tmpdir(), "trimmed.mp3");
    const response = await axios.get(audioUrl, {
      responseType: "arraybuffer",
    });
    fs.writeFileSync(inputTmpPath, response.data);

    // Create a temporary file for the output
    return new Promise<string>((resolve, reject) => {
      // check if startTime + duration is less than the length of the audio
      ffmpeg.ffprobe(inputTmpPath, (err, data) => {
        if (err) {
          reject(err);
        }
        if (data.format.duration < startTime + duration) {
          fs.unlinkSync(inputTmpPath);
          resolve(audioUrl);
        }
      });

      ffmpeg(inputTmpPath)
        .setStartTime(startTime)
        .setDuration(duration)
        .output(outputTmpPath)
        .on("end", async () => {
          const url = await this.storageService.upload(
            orgname,
            "audio/" + new Date().valueOf().toString() + ".mp3",
            {
              buffer: fs.readFileSync(outputTmpPath),
              encoding: "7bit",
              fieldname: "audio",
              mimetype: "audio/mp3", // You may want to detect this automatically, e.g. with the `file-type` package
              originalname: ospath.basename(outputTmpPath),
              size: fs.statSync(outputTmpPath).size,
            } as Express.Multer.File
          );

          fs.unlinkSync(outputTmpPath);
          fs.unlinkSync(inputTmpPath);
          resolve(url);
        })
        .on("error", reject)
        .run();
    });
  };

  constructor(
    @Inject(STORAGE_SERVICE) private storageService: StorageService,
    private httpService: HttpService,
    private keyframesService: KeyframesService
  ) {}
}
