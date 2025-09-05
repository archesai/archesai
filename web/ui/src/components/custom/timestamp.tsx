import type { JSX } from 'react'
import { format } from 'date-fns'

export const Timestamp = ({ date }: { date: string }): JSX.Element => {
  const formattedDate = format(new Date(date), 'dd/MM/yyyy HH:mm:ss')
  return <span>{formattedDate}</span>
}
