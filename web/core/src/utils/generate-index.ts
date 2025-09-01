import { readdirSync, statSync, writeFileSync } from 'node:fs'
import { join, relative } from 'node:path'

export const generateExports = (projectSrc: string): void => {
  const INDEX_FILE = join(projectSrc, 'index.ts') // Output file

  const files = getAllFiles(projectSrc)

  const exports = files.map((file) => {
    const relativePath =
      './' + relative(projectSrc, file).replace(/\\/g, '/').replace(/\.ts$/, '')
    return `export * from '${relativePath}';`
  })

  writeFileSync(INDEX_FILE, exports.join('\n') + '\n', 'utf8')
  console.log(
    `✅ Generated ${INDEX_FILE} with ${exports.length.toString()} exports.`
  )
}

export const getAllFiles = (dir: string, fileList: string[] = []): string[] => {
  const files = readdirSync(dir)

  files.forEach((file) => {
    const fullPath = join(dir, file)
    if (statSync(fullPath).isDirectory()) {
      getAllFiles(fullPath, fileList)
    } else if (
      fullPath.endsWith('.ts') &&
      !fullPath.endsWith('index.ts') &&
      !fullPath.endsWith('.spec.ts')
    ) {
      fileList.push(fullPath)
    }
  })

  return fileList
}
