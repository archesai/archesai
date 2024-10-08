import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { createTransport } from "nodemailer";
import * as Mail from "nodemailer/lib/mailer";

@Injectable()
export class EmailService {
  private nodemailerTransport: Mail;
  constructor(private readonly configService: ConfigService) {
    this.nodemailerTransport = createTransport({
      auth: {
        pass: this.configService.get("EMAIL_PASSWORD"),
        user: this.configService.get("EMAIL_USER"),
      },
      service: this.configService.get("EMAIL_SERVICE"),
    });
  }

  async sendMail(options: Mail.Options) {
    return this.nodemailerTransport.sendMail(options);
  }
}
