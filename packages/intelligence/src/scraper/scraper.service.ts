import type { Browser } from 'playwright-core'

import { chromium } from 'playwright-core'
import sharp from 'sharp'

import type { ConfigService, Logger } from '@archesai/core'

import { BadRequestException, catchErrorAsync, isString } from '@archesai/core'

export const createScraperService = (
  configService: ConfigService,
  logger: Logger
): {
  detectMimeTypeFromUrl(url: string): Promise<string>
  generateThumbnail(
    url: null | string,
    text: null | string,
    mimeType: null | string
  ): Promise<Buffer>
  onDestroy(): Promise<void>
  onInit(): Promise<void>
  screenshot(url: string): Promise<Buffer>
} => {
  let browser: Browser | undefined

  const escapeXml = (unsafe: string): string => {
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

  const getThumbnailFromFirstPagePdf = async (
    _buffer: Buffer
  ): Promise<Buffer> => {
    return new Promise((_resolve, reject) => {
      reject(new Error('PDF thumbnail generation is not implemented'))
    })
  }

  const getThumbnailFromText = async (
    text: string,
    width = 800,
    height = 600
  ): Promise<Buffer> => {
    const maxLength = 1000
    const displayText =
      text.length > maxLength ? text.substring(0, maxLength) + '...' : text
    const escapedText = escapeXml(displayText)

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

  const closeBrowser = async () => {
    if (browser) {
      await browser.close()
      logger.log('browser connection closed')
    }
  }

  const initializeBrowser = async () => {
    if (!configService.get('scraper.enabled')) {
      throw new Error('scraper is not enabled')
    }
    const scraperEndpoint = configService.get('scraper.endpoint')
    browser = await chromium.connect(scraperEndpoint)
    logger.log('browser connection open', {
      scraperEndpoint
    })
  }

  const scraperService = {
    async detectMimeTypeFromUrl(url: string): Promise<string> {
      logger.debug('detecting mime type', { url })
      const [err, response] = await catchErrorAsync(
        fetch(url, {
          method: 'HEAD'
        })
      )
      if (err) {
        logger.error('failed to detect mime type', { error: err })
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
    },

    async generateThumbnail(
      url: null | string,
      text: null | string,
      mimeType: null | string
    ): Promise<Buffer> {
      switch (mimeType) {
        case 'application/pdf':
          if (!url) {
            throw new BadRequestException('PDF URL is required')
          }

          return getThumbnailFromFirstPagePdf(
            Buffer.from((await fetch(url)) as unknown as ArrayBuffer)
          )
        case 'text/plain':
          if (!text) {
            throw new BadRequestException('PDF URL is required')
          }
          return getThumbnailFromText(text)
        default:
          if (url) {
            return scraperService.screenshot(url)
          } else {
            throw new BadRequestException('URL is required')
          }
      }
    },

    async onDestroy() {
      if (configService.get('scraper.enabled')) {
        await closeBrowser()
      }
    },

    async onInit() {
      if (configService.get('scraper.enabled')) {
        await initializeBrowser()
      }
    },

    async screenshot(url: string): Promise<Buffer> {
      logger.debug('taking screenshot', { url })
      if (!browser) {
        logger.error('browser is not initialized')
        throw new Error('scraper is not initialized')
      }
      const context = await browser.newContext({
        viewport: { height: 1080, width: 1920 }
      })
      const page = await context.newPage()
      await page.goto(url, { timeout: 60000, waitUntil: 'load' })
      await page.waitForLoadState('networkidle')
      const screenshot = await page.screenshot({
        fullPage: true
      })
      await context.close()
      logger.debug('screenshot taken', { url })
      return screenshot
    }
  }

  return scraperService
}

export type ScraperService = ReturnType<typeof createScraperService>
