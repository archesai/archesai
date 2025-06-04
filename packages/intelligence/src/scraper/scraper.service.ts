import type { Browser } from 'playwright-core'

import { chromium } from 'playwright-core'
import sharp from 'sharp'

import type { ConfigService, FetcherService } from '@archesai/core'

import {
  BadRequestException,
  catchErrorAsync,
  isString,
  Logger
} from '@archesai/core'

/**
 * Service for scraping content from URLs.
 */
export class ScraperService {
  private browser?: Browser
  private readonly configService: ConfigService
  private readonly fetcherService: FetcherService
  private readonly logger = new Logger(ScraperService.name)

  constructor(configService: ConfigService, fetcherService: FetcherService) {
    this.configService = configService
    this.fetcherService = fetcherService
  }

  public async detectMimeTypeFromUrl(url: string): Promise<string> {
    this.logger.debug('detecting mime type', { url })
    const [err, response] = await catchErrorAsync(
      this.fetcherService.head<Response>(url, {
        method: 'HEAD'
      })
    )
    if (err) {
      this.logger.error('failed to detect mime type', { error: err })
      throw new BadRequestException(`Failed to detect mime type`)
    }

    const contentType = response.headers.get('Content-Type')
    if (!isString(contentType)) {
      throw new BadRequestException('Failed to detect mime type')
    }
    const mimeType = contentType.split(';')[0]
    if (!isString(mimeType)) {
      throw new BadRequestException('Failed to detect mime type')
    }
    return mimeType
  }

  // public async detectMimeTypeFromFileName(
  //   fileName: string
  // ): Promise<string> {
  //   this.logger.debug('detecting mime type', { fileName })
  //   const mimeType = mime.lookup(fileName)
  //   if (!mimeType) {
  //     this.logger.warn('failed to detect mime type from file name')
  //     throw new BadRequestException('failed to detect mime type')
  //   }
  //   return mimeType
  // }

  public async generateThumbnail(
    url: null | string,
    text: null | string,
    mimeType: null | string
  ): Promise<Buffer> {
    switch (mimeType) {
      case 'application/pdf':
        if (!url) {
          throw new BadRequestException('PDF URL is required')
        }

        return this.getThumbnailFromFirstPagePdf(
          Buffer.from(await this.fetcherService.get<ArrayBuffer>(url))
        )
      case 'text/plain':
        if (!text) {
          throw new BadRequestException('PDF URL is required')
        }
        return this.getThumbnailFromText(text)
      default:
        if (url) {
          return this.screenshot(url)
        } else {
          throw new BadRequestException('URL is required')
        }
    }
  }

  public async onDestroy() {
    if (this.configService.get('scraper.enabled')) {
      await this.closeBrowser()
    }
  }

  public async onInit() {
    if (this.configService.get('scraper.enabled')) {
      await this.initializeBrowser()
    }
  }

  public async screenshot(url: string): Promise<Buffer> {
    this.logger.debug('taking screenshot', { url })
    if (!this.browser) {
      this.logger.error('browser is not initialized')
      throw new Error('scraper is not initialized')
    }
    const context = await this.browser.newContext({
      viewport: { height: 1080, width: 1920 } // Optional: Set a consistent viewport size
    })
    const page = await context.newPage()
    await page.goto(url, { timeout: 60000, waitUntil: 'load' })
    await page.waitForLoadState('networkidle')
    const screenshot = await page.screenshot({
      fullPage: true
    })
    await context.close()
    this.logger.debug('screenshot taken', { url })
    return screenshot
  }

  private async closeBrowser() {
    if (this.browser) {
      await this.browser.close()
      this.logger.log('browser connection closed')
    }
  }

  private escapeXml(unsafe: string): string {
    return unsafe.replace(/[<>&'"]/g, function (c) {
      switch (c) {
        case '"':
          return '&quot;'
        case '&':
          return '&amp;'
        case "'":
          return '&apos;'
        case '<':
          return '&lt;'
        case '>':
          return '&gt;'
        default:
          return c
      }
    })
  }

  private async getThumbnailFromFirstPagePdf(_buffer: Buffer): Promise<Buffer> {
    return new Promise((_resolve, reject) => {
      // gm(buffer, 'pdf[0]')
      //   .density(150, 150)
      //   .quality(90)
      //   .toBuffer('PNG', (err, buffer) => {
      //     if (err) {
      //       reject(err)
      //     } else {
      //       resolve(buffer)
      //     }
      //   })
      reject(new Error('PDF thumbnail generation is not implemented'))
    })
  }

  private async getThumbnailFromText(
    text: string,
    width = 800,
    height = 600
  ): Promise<Buffer> {
    const maxLength = 1000
    const displayText =
      text.length > maxLength ? text.substring(0, maxLength) + '...' : text
    const escapedText = this.escapeXml(displayText)

    const svgImage = `
      <svg width="${width.toString()}" height="${height.toString()}">
        <style>
          .title { fill: black; font-size: 24px; font-family: Arial, sans-serif; }
        </style>
        <rect width="100%" height="100%" fill="white"/>
        <text x="50%" y="50%" text-anchor="middle" dominant-baseline="middle" class="title">${escapedText}</text>
      </svg>
    `

    return sharp(Buffer.from(svgImage)).png().toBuffer()
  }

  private async initializeBrowser() {
    if (!this.configService.get('scraper.enabled')) {
      throw new Error('scraper is not enabled')
    }
    const scraperEndpoint = this.configService.get('scraper.endpoint')
    this.browser = await chromium.connect(scraperEndpoint)
    this.logger.log('browser connection open', {
      scraperEndpoint
    })
  }

  // private async getThumbnailFromYoutubeUrl(url: string): Promise<Buffer> {
  //   const match = url.match(
  //     /(?:youtube\.com\/(?:[^/n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})/
  //   )
  //   const videoId = match ? match[1] : null
  //   if (!videoId) {
  //     throw new Error('Invalid YouTube URL')
  //   }
  //   const response = await fetch(
  //     `https://img.youtube.com/vi/${videoId}/maxresdefault.jpg`
  //   )
  //   if (!response.ok) {
  //     throw new Error('Failed to fetch YouTube thumbnail.')
  //   }
  //   return Buffer.from(await response.arrayBuffer())
  // }

  // private resizeImage(
  //   imageBuffer: Buffer,
  //   width: number,
  //   height: number
  // ): Promise<Buffer> {
  //   return sharp(imageBuffer)
  //     .resize(width, height, { fit: 'inside' })
  //     .png()
  //     .toBuffer()
  // }
}
