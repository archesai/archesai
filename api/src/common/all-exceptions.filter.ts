import {
  ArgumentsHost,
  Catch,
  HttpException,
  HttpStatus,
} from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { BaseExceptionFilter } from "@nestjs/core";

@Catch()
export class AllExceptionsFilter extends BaseExceptionFilter {
  private readonly logger: Logger = new Logger("All Exceptions Filter");

  catch(exception: any, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse();

    let statusCode =
      exception instanceof HttpException
        ? exception.getStatus()
        : HttpStatus.INTERNAL_SERVER_ERROR;

    let message =
      exception?.response?.message ||
      exception?.message ||
      "Internal server error";

    if (
      message.includes("Unique constraint failed") ||
      message.includes("already in use") ||
      exception?.name == "UniqueConstraintViolationException"
    ) {
      const rx = /Unique constraint failed on the fields: \(`(\w+)`\)/g;
      const matches = rx.exec(message);
      if (matches) {
        message = `Resource already exists with provided '${matches[1].toLocaleLowerCase()}'`;
      }
      statusCode = 409;
    }

    if (message.includes("not found") || exception?.name == "NotFoundError") {
      const rx = /No '(\w+)' record\(s\)/g;
      const matches = rx.exec(message);
      if (matches) {
        message =
          "No " +
          matches[1].toLocaleLowerCase() +
          " found with the provided identifer";
      }
      statusCode = 404;
    }

    // if (message.includes("Unknown argument")) {
    //   message = "Bad Request";
    //   statusCode = 400;
    // }

    const errorResponse: any = {
      message,
      statusCode,
    };
    if (host.getType() != "http") {
      const e = new HttpException(exception.name, statusCode);
      return e;
    }

    this.logger.error(errorResponse);

    response.status(statusCode).json(errorResponse);
  }
}
