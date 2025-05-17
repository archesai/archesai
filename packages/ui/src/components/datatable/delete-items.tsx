import { useState } from 'react'
import { Trash } from 'lucide-react'
import { toast } from 'sonner'

import type { BaseEntity } from '@archesai/domain'

import { Button } from '#components/shadcn/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
  DialogTrigger
} from '#components/shadcn/dialog'
import { ScrollArea } from '#components/shadcn/scroll-area'
import { Separator } from '#components/shadcn/separator'

export interface DeleteProps<TEntity extends BaseEntity> {
  deleteItem: (id: string) => Promise<void>
  entityType: string
  items: TEntity[]
  variant?: 'lg' | 'md' | 'sm'
}

export const DeleteItems = <TEntity extends BaseEntity>({
  deleteItem,
  entityType,
  items,
  variant = 'sm'
}: DeleteProps<TEntity>) => {
  const [openConfirmDelete, setOpenConfirmDelete] = useState(false)
  const t = (text: string) => text
  const handleDelete = async () => {
    for (const item of items) {
      try {
        await deleteItem(item.id)
        setOpenConfirmDelete(false)
        toast(t(`The ${entityType} has been removed`))
      } catch (error: unknown) {
        if (error instanceof Error) {
          toast(t(`Could not remove ${entityType}`), {
            description: error.message
          })
        } else {
          console.error(error)
        }
      }
    }
  }

  return (
    <Dialog
      onOpenChange={(open) => {
        setOpenConfirmDelete(open)
      }}
      open={openConfirmDelete}
    >
      <div className='hidden'>
        <DialogTitle />
        <DialogDescription />
      </div>
      <DialogTrigger asChild>
        {variant === 'sm' ? (
          <div
            className='text-destructive cursor-pointer'
            onClick={() => {
              setOpenConfirmDelete(true)
            }}
          >
            <Trash className='text-destructive h-5 w-5' />
          </div>
        ) : variant === 'md' ? (
          <div
            className='w-full'
            onClick={() => {
              setOpenConfirmDelete(true)
            }}
          >
            {t('Delete')}
          </div>
        ) : (
          <Button
            className='h-8'
            onClick={() => {
              setOpenConfirmDelete(true)
            }}
            variant='destructive'
          >
            {t('Delete')}
          </Button>
        )}
      </DialogTrigger>

      <DialogContent className='gap-0 p-0'>
        <div className='flex flex-col items-center justify-center gap-3 p-4'>
          <Trash className='text-destructive' />
          <p className='text-center'>
            {t(
              `Are you sure you want to permanently delete the following ${entityType}${items.length > 1 ? 's' : ''}?`
            )}
          </p>
          {
            <ScrollArea>
              <div className='max-h-72 p-4'>
                {items.map((item, i) => (
                  <p key={i}>{item.name}</p>
                ))}
              </div>
            </ScrollArea>
          }
        </div>
        <Separator />
        <DialogFooter className='flex rounded-lg bg-gray-50 p-6 dark:bg-black'>
          <div className='flex w-full items-center gap-4'>
            <Button
              className='w-full'
              onClick={() => {
                setOpenConfirmDelete(false)
              }}
              size='sm'
            >
              {t('Cancel')}
            </Button>
            <Button
              className='w-full'
              onClick={async () => {
                await handleDelete()
              }}
              size='sm'
              variant={'destructive'}
            >
              {t('Delete')}
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
