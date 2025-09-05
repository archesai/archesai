import type { JSX } from "react"

import { useEffect, useState } from "react"

import { ArrowUpToLineIcon } from "#components/custom/icons"
import { Button } from "#components/shadcn/button"

export const ScrollButton = (): JSX.Element => {
  const [showTopBtn, setShowTopBtn] = useState(false)

  useEffect(() => {
    window.addEventListener("scroll", () => {
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
        <Button
          className="fixed right-4 bottom-4 opacity-90 shadow-md"
          onClick={goToTop}
          size="icon"
        >
          <ArrowUpToLineIcon className="h-4 w-4" />
        </Button>
      )}
    </>
  )
}
