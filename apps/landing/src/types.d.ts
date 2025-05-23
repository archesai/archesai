declare module '*.svg' {
  import type React from 'react'

  const content: string
  export const ReactComponent: React.FC<React.SVGProps<SVGSVGElement>>
  export default content
}

declare module '*.css' {}
