import { format } from 'date-fns'

export const Timestamp = ({ date }: { date: string }) => {
  const formattedDate = format(new Date(date), 'dd/MM/yyyy HH:mm:ss')
  return (
    <span>{formattedDate}</span>

    //   <span className='font-light'>
    //         {format(new Date(row.original.createdAt), 'M/d/yy h:mm a')}
    //       </span>
  )
}
