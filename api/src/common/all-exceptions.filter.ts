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
  private readonly logger: Logger = new Logger("AllExceptionsFilter");

  catch(exception: any, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const request = ctx.getRequest();
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

    const prodErrorResponse: any = {
      message,
      statusCode,
    };
    if (host.getType() != "http") {
      const e = new HttpException(exception.name, statusCode);
      return e;
    }

    const devErrorResponse: any = {
      errorName: exception?.name,
      method: request.method,
      msg: message,
      path: request.url,
      statusCode,
      timestamp: new Date().toISOString(),
    };

    this.logger.error(devErrorResponse);

    response
      .status(statusCode)
      .json(
        process.env.NODE_ENV === "production"
          ? prodErrorResponse
          : devErrorResponse,
      );
  }
}
