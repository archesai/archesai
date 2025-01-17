import { Injectable } from '@nestjs/common'
import { createTransport } from 'nodemailer'
import Mail from 'nodemailer/lib/mailer'
import { ConfigService } from '../config/config.service'

@Injectable()
export class EmailService {
  private nodemailerTransport: Mail
  constructor(private readonly configService: ConfigService) {
    this.nodemailerTransport = createTransport({
      auth: {
        pass: this.configService.get('email.password'),
        user: this.configService.get('email.user')
      },
      service: this.configService.get('email.service')
    })
  }

  async sendMail(options: Mail.Options) {
    return this.nodemailerTransport.sendMail(options)
  }
}
