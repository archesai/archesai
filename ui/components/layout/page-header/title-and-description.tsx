export const TitleAndDescription = ({
  description,
  Icon,
  title
}: {
  description?: string
  Icon: any
  title?: string
}) => {
  if (!title) return null
  return (
    <div className='flex items-center gap-3 border-b bg-muted/30 px-3 py-3'>
      {Icon && <Icon className='h-8 w-8 text-slate-700 dark:text-slate-500' />}
      <div>
        <p className='text-xl font-semibold text-foreground'>{title}</p>
        <p className='text-sm text-muted-foreground'>{description}</p>
      </div>
    </div>
  )
}
