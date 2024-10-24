import { Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";
import * as math from "mathjs";
import fetch from "node-fetch";
import { AudioContext } from "web-audio-api";

import { retry } from "../common/retry";

@Injectable()
export class KeyframesService {
  private readonly logger = new Logger(KeyframesService.name);

  async getKeyframes(
    url: string,
    framerate: number,
    fn: string,
    isTranslation: boolean
  ) {
    this.installPolyfill();
    const context = new AudioContext();

    const arrayBuffer = await retry(
      this.logger,
      () => fetch(url).then((response) => response.arrayBuffer()),
      5
    );

    const audioBuffer = await context.decodeAudioData(arrayBuffer);
    // Average between channels. Take abs so we don't have phase issues (and we eventually want absolute value anyway, for volume).
    function addAbsArrayElements(a, b) {
      return a.map((e, i) => Math.abs(e) + Math.abs(b[i]));
    }
    const channels = [];
    for (let i = 0; i < audioBuffer.numberOfChannels; i++) {
      channels.push(audioBuffer.getChannelData(i));
    }
    const rawData = channels
      .reduce(addAbsArrayElements)
      .map((x) => x / audioBuffer.numberOfChannels);
    // const rawData = audioBuffer.getChannelData(0); // We only need to work with one channel of data
    const samples = audioBuffer.duration * framerate; //rawData.length; // Number of samples we want to have in our final data set
    const blockSize = Math.floor(rawData.length / samples); // Number of samples in each subdivision
    let filteredData = [];
    for (let i = 0; i < samples; i++) {
      const chunk = rawData.slice(i * blockSize, (i + 1) * blockSize - 1);
      const sum = chunk.reduce((a, b) => a + b, 0);
      filteredData.push(sum / chunk.length);
    }
    const max = Math.max(...filteredData); // Normalise - maybe not ideal.
    // const Parser = require('expr-eval').Parser;
    // const parser = new Parser();
    // let expr = parser.parse(fn.value);
    filteredData = filteredData
      .map((x) => x / max)
      .map((x, ind) =>
        math.evaluate(
          fn.replace("x", x.toString()).replace("y", ind.toString())
        )
      );
    const string = this.getString(filteredData, isTranslation);

    return string;
  }

  getString(arr, isTranslation: boolean) {
    let string = "";
    for (const ind of Object.keys(arr)) {
      let sample = parseFloat(arr[ind]);
      if (sample > 1.01 && isTranslation) {
        sample = sample + 6;
      }
      string = string.concat(`${ind}: (${sample.toFixed(2)})`);
      if (parseInt(ind) < arr.length - 1) {
        string = string.concat(", ");
      }
    }
    return string;
  }

  installPolyfill() {
    function decodeAudioData_polyfill(
      audioData,
      successCallback,
      errorCallback
    ) {
      if (arguments.length > 1) {
        // Callback
        this.decodeAudioData(audioData, successCallback, errorCallback);
      } else {
        // Promise
        return new Promise((success, reject) =>
          this.decodeAudioData_original(audioData, success, reject)
        );
      }
    }
    if (!AudioContext.prototype.decodeAudioData_original) {
      AudioContext.prototype.decodeAudioData_original =
        AudioContext.prototype.decodeAudioData;
      AudioContext.prototype.decodeAudioData = decodeAudioData_polyfill;
    }
  }
}
