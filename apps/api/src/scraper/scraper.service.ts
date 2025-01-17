import {
  Injectable,
  OnModuleInit,
  OnModuleDestroy,
  Logger,
  BadRequestException
} from '@nestjs/common'
import { chromium, Browser } from 'playwright-core'
import * as mime from 'mime-types'
import sharp from 'sharp'
import gm from 'gm'
import { ConfigService } from '../config/config.service'

@Injectable()
export class ScraperService implements OnModuleInit, OnModuleDestroy {
  private readonly logger = new Logger(ScraperService.name)
  private browser: Browser

  constructor(readonly configService: ConfigService) {}

  async onModuleInit() {
    await this.initializeBrowser()
  }

  async onModuleDestroy() {
    await this.closeBrowser()
  }

  private async initializeBrowser() {
    const scraperEndpoint = this.configService.get('scraper.endpoint')!
    this.browser = await chromium.connect(scraperEndpoint)
    this.logger.log(
      {
        scraperEndpoint
      },
      'browser connection open'
    )
  }

  private async closeBrowser() {
    if (this.browser) {
      await this.browser.close()
      this.logger.log('browser connection closed')
    }
  }

  async takeScreenshot(url: string): Promise<Buffer> {
    this.logger.debug({ url }, 'taking screenshot')
    const context = await this.browser.newContext({
      viewport: { width: 1920, height: 1080 } // Optional: Set a consistent viewport size
    })
    const page = await context.newPage()
    await page.goto(url, { waitUntil: 'load', timeout: 60000 })
    await page.waitForLoadState('networkidle')
    const screenshot = await page.screenshot({
      fullPage: true
    })
    await context.close()
    this.logger.debug({ url }, 'screenshot taken')
    return screenshot
  }

  async detectMimeType(url: string): Promise<string> {
    this.logger.debug({ url }, 'detecting mime type')
    let mimeType: string
    try {
      const response = await fetch(url, { method: 'HEAD' })
      mimeType = response.headers.get('content-type')?.split(';')[0] || ''
      if (mimeType) {
        return mimeType
      }
      this.logger.warn('failed to detect mime type from content-type header')
      const urlObj = new URL(url)
      const pathname = urlObj.pathname
      const fileName = pathname.split('/').pop()
      if (!fileName) {
        throw new BadRequestException('failed to detect mime type')
      }
      mimeType = mime.lookup(fileName) || ''
      if (mimeType === '') {
        throw new BadRequestException('failed to detect mime type')
      }
    } catch (error) {
      throw new BadRequestException({
        message: 'failed to detect mime type',
        cause: error
      })
    }
    return mimeType
  }

  async generateThumbnail(
    url: string | null,
    text: string | null,
    mimeType: string | null
  ): Promise<Buffer> {
    switch (mimeType) {
      case 'application/pdf':
        if (!url) {
          throw new BadRequestException('PDF URL is required')
        }
        return this.getThumbnailFromFirstPagePdf(
          Buffer.from(await (await fetch(url)).arrayBuffer())
        )
      case 'text/plain':
        if (!text) {
          throw new BadRequestException('PDF URL is required')
        }
        return this.getThumbnailFromText(text)
      default:
        if (url) {
          return this.takeScreenshot(url)
        } else {
          throw new BadRequestException('URL is required')
        }
    }
  }

  private resizeImage(
    imageBuffer: Buffer,
    width: number,
    height: number
  ): Promise<Buffer> {
    return sharp(imageBuffer)
      .resize(width, height, { fit: 'inside' })
      .png()
      .toBuffer()
  }

  private getThumbnailFromFirstPagePdf(buffer: Buffer): Promise<Buffer> {
    return new Promise((resolve, reject) => {
      gm(buffer, 'pdf[0]')
        .density(150, 150)
        .quality(90)
        .toBuffer('PNG', (err, buffer) => {
          if (err) {
            reject(err)
          } else {
            resolve(buffer)
          }
        })
    })
  }

  private async getThumbnailFromText(
    text: string,
    width: number = 800,
    height: number = 600
  ): Promise<Buffer> {
    // Truncate text if it's too long to prevent overflowing
    const maxLength = 1000
    const displayText =
      text.length > maxLength ? text.substring(0, maxLength) + '...' : text

    // Escape XML characters to prevent SVG injection
    const escapedText = this.escapeXml(displayText)

    // Create an SVG image with the text
    const svgImage = `
      <svg width="${width}" height="${height}">
        <style>
          .title { fill: black; font-size: 24px; font-family: Arial, sans-serif; }
        </style>
        <rect width="100%" height="100%" fill="white"/>
        <text x="50%" y="50%" text-anchor="middle" dominant-baseline="middle" class="title">${escapedText}</text>
      </svg>
    `

    // Convert SVG to PNG using sharp
    return sharp(Buffer.from(svgImage)).png().toBuffer()
  }

  private async getThumbnailFromYoutubeUrl(url: string): Promise<Buffer> {
    const match = url.match(
      /(?:youtube\.com\/(?:[^/n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})/
    )
    const videoId = match ? match[1] : null
    if (!videoId) {
      throw new Error('Invalid YouTube URL')
    }
    const response = await fetch(
      `https://img.youtube.com/vi/${videoId}/maxresdefault.jpg`
    )
    if (!response.ok) {
      throw new Error('Failed to fetch YouTube thumbnail.')
    }
    return Buffer.from(await response.arrayBuffer())
  }

  private escapeXml(unsafe: string): string {
    return unsafe.replace(/[<>&'"]/g, function (c) {
      switch (c) {
        case '<':
          return '&lt;'
        case '>':
          return '&gt;'
        case '&':
          return '&amp;'
        case "'":
          return '&apos;'
        case '"':
          return '&quot;'
        default:
          return c
      }
    })
  }
}
