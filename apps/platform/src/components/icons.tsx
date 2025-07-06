import iconsHref from './icons.svg?url'

export function Icon({
  name,
  size = 'md',
  spin = false
}: {
  name: string
  size?: 'md' | 'xl'
  spin?: boolean
}) {
  const classNames = {
    md: 'w-4 h-4',
    xl: 'w-8 h-8'
  }
  return (
    <svg
      className={`${classNames[size]} inline self-center ${
        spin ? 'animate-spin' : ''
      }`}
    >
      <use href={`${iconsHref}#${name}`} />
    </svg>
  )
}

export function LoginIcon() {
  return (
    <svg className='inline h-8 w-8 scale-x-[-1] transform self-center text-white'>
      <use href={`${iconsHref}#login`} />
    </svg>
  )
}

export function LogoutIcon() {
  return (
    <svg className='inline h-8 w-8 self-center text-white'>
      <use href={`${iconsHref}#logout`} />
    </svg>
  )
}
