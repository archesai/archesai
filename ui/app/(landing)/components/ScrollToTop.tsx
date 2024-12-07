'use client'

import { Button } from '@/components/ui/button'
import { ArrowUpToLine } from 'lucide-react'
import { useEffect, useState } from 'react'

export const ScrollToTop = () => {
  const [showTopBtn, setShowTopBtn] = useState(false)

  useEffect(() => {
    window.addEventListener('scroll', () => {
      if (window.scrollY > 400) {
        setShowTopBtn(true)
      } else {
        setShowTopBtn(false)
      }
    })
  }, [])

  const goToTop = () => {
    window.scroll({
      left: 0,
      top: 0
    })
  }

  return (
    <>
      {showTopBtn && (
        <Button className='fixed bottom-4 right-4 opacity-90 shadow-md' onClick={goToTop} size='icon'>
          <ArrowUpToLine className='h-4 w-4' />
        </Button>
      )}
    </>
  )
}
