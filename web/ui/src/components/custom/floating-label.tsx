import { forwardRef } from 'react'

import { Input } from '#components/shadcn/input'
import { Label } from '#components/shadcn/label'
import { cn } from '#lib/utils'

export type InputProps = React.InputHTMLAttributes<HTMLInputElement>

interface FloatingLabelInputProps extends InputProps {
  id: string
  label: string
}

const FloatingLabelInput = forwardRef<
  HTMLInputElement,
  FloatingLabelInputProps
>(({ className, id, label, value, ...props }, ref) => {
  return (
    <div className='relative'>
      <Input
        className={cn(
          'peer',
          'focus:border-ring focus:ring-0',
          'focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50',
          className
        )}
        id={id}
        placeholder=' '
        ref={ref}
        value={value}
        {...props}
      />
      <Label
        className={cn(
          'pointer-events-none absolute start-2 top-2 z-10 origin-[0] -translate-y-4 scale-75 transform cursor-text bg-background px-1 text-muted-foreground duration-300',
          'peer-placeholder-shown:top-1/2 peer-placeholder-shown:-translate-y-1/2 peer-placeholder-shown:scale-100 peer-placeholder-shown:bg-transparent',
          'peer-focus-visible:top-2 peer-focus-visible:-translate-y-4 peer-focus-visible:scale-75 peer-focus-visible:bg-background'
        )}
        htmlFor={id}
      >
        {label}
      </Label>
    </div>
  )
})

FloatingLabelInput.displayName = 'FloatingLabelInput'

export { FloatingLabelInput }
