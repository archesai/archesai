import type { JSX } from 'react'

import { useLocation } from '@tanstack/react-router'

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from '#components/shadcn/breadcrumb'

export const BreadCrumbs = (): JSX.Element => {
  const location = useLocation()

  // Split the pathname into segments and create breadcrumbs
  const pathSegments = location.pathname.split('/').filter(Boolean)

  const breadcrumbs = pathSegments.map((segment, index) => {
    const path = '/' + pathSegments.slice(0, index + 1).join('/')
    const title = segment
      .split('-')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')

    return {
      isLast: index === pathSegments.length - 1,
      path,
      title
    }
  })

  return (
    <>
      <Breadcrumb>
        <BreadcrumbList>
          {breadcrumbs.map((breadcrumb, index) => (
            <div
              className='flex items-center'
              key={breadcrumb.path}
            >
              {index > 0 && <BreadcrumbSeparator className='hidden md:block' />}
              <BreadcrumbItem className={index === 0 ? 'hidden md:block' : ''}>
                {breadcrumb.isLast ?
                  <BreadcrumbPage>{breadcrumb.title}</BreadcrumbPage>
                : <BreadcrumbLink href={breadcrumb.path}>
                    {breadcrumb.title}
                  </BreadcrumbLink>
                }
              </BreadcrumbItem>
            </div>
          ))}
        </BreadcrumbList>
      </Breadcrumb>
    </>
  )
}
