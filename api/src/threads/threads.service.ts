import { Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";

import { WebsocketsService } from "../websockets/websockets.service";
import { CreateThreadDto } from "./dto/create-thread.dto";
import { ThreadQueryDto } from "./dto/thread-query.dto";
import { ThreadRepository } from "./thread.repository";

@Injectable()
export class ThreadsService {
  private readonly logger: Logger = new Logger("Threads Service");

  constructor(
    private threadRepository: ThreadRepository,
    private websocketsService: WebsocketsService
  ) {}
  async cleanupUnused() {
    return this.threadRepository.cleanupUnused();
  }

  async create(
    orgname: string,

    createThreadDto: CreateThreadDto
  ) {
    const thread = await this.threadRepository.create(orgname, createThreadDto);
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "threads"],
    });
    return thread;
  }

  async findAll(
    orgname: string,

    threadQueryDto: ThreadQueryDto
  ) {
    return this.threadRepository.findAll(orgname, threadQueryDto);
  }

  async findOne(orgname: string, threadId: string) {
    return this.threadRepository.findOne(threadId);
  }

  async incrementCredits(orgname: string, threadId: string, credits: number) {
    return this.threadRepository.incrementCredits(orgname, threadId, credits);
  }

  async remove(orgname: string, threadId: string) {
    await this.threadRepository.delete(threadId);
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "threads"],
    });
  }

  async updateThreadName(orgname: string, threadId: string, name: string) {
    return this.threadRepository.updateThreadName(orgname, threadId, name);
  }
}
