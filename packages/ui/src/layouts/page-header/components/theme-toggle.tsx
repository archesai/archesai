import { useTheme } from 'next-themes'

import { MoonIcon, SunIcon } from '#components/custom/icons'
import { Button } from '#components/shadcn/button'

export function ThemeToggle() {
  const { theme, setTheme } = useTheme()

  return (
    <Button
      className='group/toggle extend-touch-target'
      onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
      size='sm'
      title='Toggle theme'
      variant='outline'
    >
      <SunIcon className='size-4 scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90' />
      <MoonIcon className='absolute size-4 scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0' />
      <span className='sr-only'>Toggle theme</span>
    </Button>
  )
}
