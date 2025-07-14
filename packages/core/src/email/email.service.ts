import type { SendMailOptions } from 'nodemailer'

import nodemailer from 'nodemailer'

import type { ConfigService } from '#config/config.service'

export const createEmailService = (configService: ConfigService) => {
  const nodemailerTransport = nodemailer.createTransport({
    auth: {
      pass: configService.get('email.password'),
      user: configService.get('email.user')
    },
    service: configService.get('email.service')
  })

  return {
    sendMail: async (options: SendMailOptions): Promise<void> => {
      await nodemailerTransport.sendMail(options)
    }
  }
}

export type EmailService = ReturnType<typeof createEmailService>
