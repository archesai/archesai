import { useCallback, useEffect, useState } from 'react'
import { useTheme } from 'next-themes'

import { MoonIcon, SunIcon } from '#components/custom/icons'
import { Button } from '#components/shadcn/button'

export function ModeToggle() {
  const { resolvedTheme, setTheme } = useTheme()

  const toggleTheme = useCallback(() => {
    setTheme(resolvedTheme === 'dark' ? 'light' : 'dark')
  }, [resolvedTheme, setTheme])

  const [isMounted, setIsMounted] = useState(false)
  useEffect(() => {
    setIsMounted(true)
  }, [])

  if (!isMounted) {
    return <div />
  }

  return (
    <Button
      className='group/toggle extend-touch-target'
      onClick={toggleTheme}
      size='sm'
      title='Toggle theme'
      variant='outline'
    >
      {resolvedTheme === 'dark' ?
        <MoonIcon />
      : <SunIcon />}
      <span className='sr-only'>Toggle theme</span>
    </Button>
  )
}
