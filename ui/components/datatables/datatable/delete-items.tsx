import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'
import { useToast } from '@/hooks/use-toast'
import * as VisuallyHidden from '@radix-ui/react-visually-hidden'
import { Trash } from 'lucide-react'
import { useState } from 'react'

import { Button } from '../../ui/button'
import { ScrollArea } from '../../ui/scroll-area'
import { Separator } from '../../ui/separator'

export interface DeleteProps<TDeleteVariables> {
  deleteFunction: (params: TDeleteVariables) => Promise<void>
  deleteVariables: TDeleteVariables[]
  items: {
    id: string
    name: string
  }[]
  itemType: string
  variant?: 'lg' | 'md' | 'sm'
}

// create a functional component called DeleteItems
export const DeleteItems = <TDeleteVariables,>({
  deleteFunction,
  deleteVariables,
  items,
  itemType,
  variant = 'sm'
}: DeleteProps<TDeleteVariables>) => {
  const [openConfirmDelete, setOpenConfirmDelete] = useState(false)
  const t = (text: string) => text
  const { toast } = useToast()
  const handleDelete = async () => {
    for (const deleteVars of deleteVariables) {
      try {
        await deleteFunction(deleteVars)
        setOpenConfirmDelete(false)
        toast({ title: t(`The ${itemType} has been removed`) })
      } catch (err) {
        console.error(err)
        toast({ title: t(`Could not remove ${itemType}`) })
      }
    }
  }

  return (
    <Dialog onOpenChange={(open) => setOpenConfirmDelete(open)} open={openConfirmDelete}>
      <VisuallyHidden.Root>
        <DialogTitle />
        <DialogDescription />
      </VisuallyHidden.Root>
      <DialogTrigger asChild>
        {variant === 'sm' ? (
          <div className='cursor-pointer text-destructive' onClick={() => setOpenConfirmDelete(true)}>
            <Trash className='h-5 w-5 text-destructive' />
          </div>
        ) : variant === 'md' ? (
          <div className='w-full' onClick={() => setOpenConfirmDelete(true)}>
            {t('Delete')}
          </div>
        ) : (
          <Button className='h-8' onClick={() => setOpenConfirmDelete(true)} variant='destructive'>
            {t('Delete')}
          </Button>
        )}
      </DialogTrigger>

      <DialogContent className='gap-0 p-0'>
        <div className='flex flex-col items-center justify-center gap-3 p-4'>
          <Trash className='text-destructive' />
          <p className='text-center'>
            {t(`Are you sure you want to permanently delete the following ${itemType}${items?.length > 1 ? 's' : ''}?`)}
          </p>
          {
            <ScrollArea>
              <div className='max-h-72 p-4'>{items?.map((item, i) => <p key={i}>{item.name}</p>)}</div>
            </ScrollArea>
          }
        </div>
        <Separator />
        <DialogFooter className='flex rounded-lg bg-gray-50 p-6 dark:bg-black'>
          <div className='flex w-full items-center gap-4'>
            <Button className='w-full' onClick={() => setOpenConfirmDelete(false)} size='sm'>
              {t('Cancel')}
            </Button>
            <Button className='w-full' onClick={async () => await handleDelete()} size='sm' variant={'destructive'}>
              {t('Delete')}
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
