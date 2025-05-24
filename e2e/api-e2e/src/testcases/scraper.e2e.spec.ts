// import fs from 'fs'
// import path from 'path'
// import type { TestingModule } from '@nestjs/testing'

// import { Test } from '@nestjs/testing'

// describe('ScraperService - Integration Tests', () => {
//   let service: ScraperService

//   beforeEach(async () => {
//     const module: TestingModule = await Test.createTestingModule({
//       providers: [
//         ScraperService,
//         {
//           provide: ConfigService,
//           useValue: {
//             get: jest.fn().mockReturnValue('http://playwright:3000')
//           }
//         }
//       ]
//     }).compile()

//     service = module.get<ScraperService>(ScraperService)
//     await service.onModuleInit()
//   })

//   afterEach(async () => {
//     await service.onModuleDestroy()
//   })

//   it('should be defined', () => {
//     expect(service).toBeDefined()
//   })

//   it('should take a screenshot of a website', async () => {
//     const url = 'https://archesai.com'
//     const mimeType = await service.detectMimeType(url)
//     expect(mimeType).toBe('text/html')

//     const screenshot = await service.generateThumbnail(url, '', mimeType)
//     expect(screenshot).toBeDefined()
//     expect(screenshot.length).toBeGreaterThan(0)
//     expect(screenshot).toBeInstanceOf(Buffer)
//     fs.writeFileSync(path.join(__dirname, 'testdata/website.png'), screenshot)
//   })

//   it('should take a screenshot of a pdf', async () => {
//     const url = 'https://www.ucd.ie/t4cms/Test%20PDF-8mb.pdf'
//     const mimeType = await service.detectMimeType(url)
//     expect(mimeType).toBe('application/pdf')

//     const screenshot = await service.generateThumbnail(url, '', mimeType)
//     expect(screenshot).toBeDefined()
//     expect(screenshot.length).toBeGreaterThan(0)
//     expect(screenshot).toBeInstanceOf(Buffer)
//     fs.writeFileSync(path.join(__dirname, 'testdata/pdf.png'), screenshot)
//   })

//   it('should take a screenshot of text content', async () => {
//     const url =
//       'https://gist.githubusercontent.com/rt2zz/e0a1d6ab2682d2c47746950b84c0b6ee/raw/83b8b4814c3417111b9b9bef86a552608506603e/markdown-sample.md'
//     const mimeType = await service.detectMimeType(url)
//     expect(mimeType).toBe('text/plain')

//     const screenshot = await service.generateThumbnail(url, '', mimeType)
//     expect(screenshot).toBeDefined()
//     expect(screenshot.length).toBeGreaterThan(0)
//     expect(screenshot).toBeInstanceOf(Buffer)
//     fs.writeFileSync(path.join(__dirname, 'testdata/text.png'), screenshot)
//   })

//   it('should throw an error for an invalid URL', async () => {
//     await expect(service.detectMimeType('invalid-url')).rejects.toThrow()
//     await expect(
//       service.generateThumbnail('invalid-url', '', '')
//     ).rejects.toThrow()
//   })
// })
