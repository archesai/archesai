export const TitleAndDescription = ({
  description,
  // Icon,
  title
}: {
  description?: string
  Icon: any
  title?: string
}) => {
  if (!title) return null
  return (
    <div className='container flex items-center gap-3 border-b px-4 py-3'>
      {/* {Icon && <Icon className='h-8 w-8' />} */}
      <div>
        <p className='text-foreground/85 text-xl font-semibold'>{title}</p>
        <p className='text-muted-foreground text-sm'>{description}</p>
      </div>
    </div>
  )
}
