'use client'

import { useState } from 'react'
import Link from 'next/link'

import { ArchesLogo } from '@archesai/ui/components/custom/arches-logo'
import { Github, Menu } from '@archesai/ui/components/custom/icons'
import { Button, buttonVariants } from '@archesai/ui/components/shadcn/button'
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuList
} from '@archesai/ui/components/shadcn/navigation-menu'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTrigger
} from '@archesai/ui/components/shadcn/sheet'
import { useIsTop } from '@archesai/ui/hooks/use-is-top'
import { ModeToggle } from '@archesai/ui/layouts/page-header/mode-toggle'

interface RouteProps {
  href: string
  label: string
}

const routeList: RouteProps[] = [
  {
    href: '#features',
    label: 'Features'
  },
  {
    href: '#testimonials',
    label: 'Testimonials'
  },
  {
    href: '#pricing',
    label: 'Pricing'
  },
  {
    href: '#faq',
    label: 'FAQ'
  }
]

export const Navbar = () => {
  const [isOpen, setIsOpen] = useState<boolean>(false)
  const isTop = useIsTop()

  return (
    <header
      className={`sticky top-0 z-40 w-full ${
        isTop
          ? 'bg-transparent'
          : 'border-b bg-white shadow-xs transition-all dark:bg-background'
      }`}
    >
      <NavigationMenu>
        <NavigationMenuList className='flex h-[56px] w-screen justify-between px-4'>
          <div className='flex items-center justify-center gap-3'>
            <NavigationMenuItem className='flex font-bold'>
              <ArchesLogo />
            </NavigationMenuItem>
            {/* mobile */}
            <span className='flex md:hidden'>
              <Sheet
                onOpenChange={setIsOpen}
                open={isOpen}
              >
                <SheetTrigger className='px-2'>
                  <Menu
                    className='flex h-5 w-5 md:hidden'
                    onClick={() => {
                      setIsOpen(true)
                    }}
                  ></Menu>
                </SheetTrigger>

                <SheetContent side={'left'}>
                  <SheetHeader>
                    <ArchesLogo />
                  </SheetHeader>
                  <nav className='mt-4 flex flex-col items-center justify-center gap-2'>
                    {routeList.map(({ href, label }: RouteProps) => (
                      <a
                        className={buttonVariants({ variant: 'ghost' })}
                        href={href}
                        key={label}
                        onClick={() => {
                          setIsOpen(false)
                        }}
                        rel='noreferrer noopener'
                      >
                        {label}
                      </a>
                    ))}
                    <a
                      className={`w-[110px] border ${buttonVariants({
                        variant: 'secondary'
                      })}`}
                      href='https://github.com/leoMirandaa/shadcn-landing-page.git'
                      rel='noreferrer noopener'
                      target='_blank'
                    >
                      <Github className='mr-2 h-5 w-5' />
                      Github
                    </a>
                  </nav>
                </SheetContent>
              </Sheet>
            </span>
            {/* desktop */}
            <nav className='hidden gap-2 md:flex'>
              {routeList.map((route: RouteProps, i) => (
                <a
                  className={`text-[17px] ${buttonVariants({
                    variant: 'ghost'
                  })}`}
                  href={route.href}
                  key={i}
                  rel='noreferrer noopener'
                >
                  {route.label}
                </a>
              ))}
            </nav>
          </div>
          <div className='hidden items-center gap-2 md:flex'>
            <ModeToggle h={'h-10'} />

            <Link href='/auth/login'>
              <Button variant={'outline'}>Log in</Button>
            </Link>

            <Link href='/auth/register'>
              <Button>Sign up for free </Button>
            </Link>
          </div>
        </NavigationMenuList>
      </NavigationMenu>
    </header>
  )
}
