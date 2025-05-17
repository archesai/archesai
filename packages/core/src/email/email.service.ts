import type { SendMailOptions } from 'nodemailer'

import nodemailer from 'nodemailer'

import type { ConfigService } from '#config/config.service'

/**
 * Service to send emails
 */
export class EmailService {
  private readonly configService: ConfigService
  private readonly nodemailerTransport: nodemailer.Transporter

  constructor(configService: ConfigService) {
    this.configService = configService
    this.nodemailerTransport = nodemailer.createTransport({
      auth: {
        pass: this.configService.get('email.password'),
        user: this.configService.get('email.user')
      },
      service: this.configService.get('email.service')
    })
  }

  public async sendMail(options: SendMailOptions): Promise<void> {
    await this.nodemailerTransport.sendMail(options)
  }
}
