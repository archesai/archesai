// src/token/token.service.ts

import { BadRequestException, Injectable } from "@nestjs/common";
import { ARTokenType } from "@prisma/client";
import * as bcrypt from "bcryptjs";
import * as crypto from "crypto";

import { PrismaService } from "../prisma/prisma.service";

@Injectable()
export class ARTokensService {
  constructor(private readonly prisma: PrismaService) {}

  /**
   * Generates and stores a token for a specific action type.
   * @param type - The type of token (e.g., PASSWORD_RESET, EMAIL_VERIFICATION, EMAIL_CHANGE).
   * @param userId - The ID of the user associated with the token.
   * @param expiresInHours - Token expiry time in hours.
   * @param additionalData - Optional additional data to store with the token (e.g., newEmail).
   * @returns The plain token string.
   */
  async createToken(
    type: ARTokenType,
    userId: string,
    expiresInHours: number,
    additionalData?: Record<string, any>
  ): Promise<string> {
    const token = crypto.randomBytes(32).toString("hex");
    const hashedToken = await bcrypt.hash(token, 10);
    const expiresAt = new Date();
    expiresAt.setHours(expiresAt.getHours() + expiresInHours);

    // Prepare data to store based on token type
    const tokenData: any = {
      expiresAt,
      token: hashedToken,
      type,
      user: { connect: { id: userId } },
    };

    if (type === ARTokenType.EMAIL_CHANGE && additionalData?.newEmail) {
      tokenData.newEmail = additionalData.newEmail;
    }

    // check if one token of this type already exists
    const existingToken = await this.prisma.aRToken.findFirst({
      where: { type },
    });
    if (existingToken) {
      await this.prisma.aRToken.delete({ where: { id: existingToken.id } });
    }

    await this.prisma.aRToken.create({
      data: tokenData,
    });

    return token;
  }

  /**
   * Verifies a token of a specific type.
   * @param type - The type of token to verify.
   * @param token - The plain token string to verify.
   * @returns An object containing the user ID and any additional data associated with the token.
   * @throws BadRequestException if token is invalid or expired.
   */
  async verifyToken(
    type: ARTokenType,
    token: string
  ): Promise<{ additionalData?: Record<string, any>; userId: string }> {
    const tokens = await this.prisma.aRToken.findMany({
      include: { user: true },
      where: { type },
    });

    // Iterate through tokens to find a match
    for (const tokenRecord of tokens) {
      const isMatch = await bcrypt.compare(token, tokenRecord.token);
      if (isMatch) {
        if (tokenRecord.expiresAt < new Date()) {
          throw new BadRequestException("Token has expired.");
        }

        const userId = tokenRecord.userId;
        const additionalData: Record<string, any> = {};

        // Extract additional data based on token type
        if (type === ARTokenType.EMAIL_CHANGE && tokenRecord.newEmail) {
          additionalData.newEmail = tokenRecord.newEmail;
        }

        // Delete the token after successful verification
        await this.prisma.aRToken.delete({ where: { id: tokenRecord.id } });

        return { additionalData, userId };
      }
    }

    throw new BadRequestException("Invalid or expired token.");
  }
}
