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
        <p className='text-xl font-semibold text-foreground/85'>{title}</p>
        <p className='text-sm text-muted-foreground'>{description}</p>
      </div>
    </div>
  )
}
