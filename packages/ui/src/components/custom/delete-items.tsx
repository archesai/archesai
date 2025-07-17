'use no memo'

import { useState } from 'react'
import { Trash } from 'lucide-react'
import { toast } from 'sonner'

import type { BaseEntity } from '@archesai/schemas'

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

export interface DeleteItemsProps<TEntity extends BaseEntity> {
  deleteItem: (id: string) => Promise<void>
  entityKey: string
  items: TEntity[]
  variant?: 'lg' | 'md' | 'sm'
}

export const DeleteItems = <TEntity extends BaseEntity>(
  props: DeleteItemsProps<TEntity>
) => {
  const [openConfirmDelete, setOpenConfirmDelete] = useState(false)
  const t = (text: string) => text
  const handleDelete = async () => {
    for (const item of props.items) {
      try {
        await props.deleteItem(item.id)
        setOpenConfirmDelete(false)
        toast(t(`The ${props.entityKey} has been removed`))
      } catch (error: unknown) {
        if (error instanceof Error) {
          toast(t(`Could not remove ${props.entityKey}`), {
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
        {props.variant === 'sm' ?
          <div
            className='cursor-pointer text-destructive'
            onClick={() => {
              setOpenConfirmDelete(true)
            }}
          >
            <Trash className='h-5 w-5 text-destructive' />
          </div>
        : props.variant === 'md' ?
          <div
            className='w-full'
            onClick={() => {
              setOpenConfirmDelete(true)
            }}
          >
            {t('Delete')}
          </div>
        : <Button
            className='h-8'
            onClick={() => {
              setOpenConfirmDelete(true)
            }}
            variant='destructive'
          >
            {t('Delete')}
          </Button>
        }
      </DialogTrigger>

      <DialogContent className='gap-0 p-0'>
        <div className='flex flex-col items-center justify-center gap-3 p-4'>
          <Trash className='text-destructive' />
          <p className='text-center'>
            {t(
              `Are you sure you want to permanently delete the following ${props.entityKey}${props.items.length > 1 ? 's' : ''}?`
            )}
          </p>
          {
            <ScrollArea>
              <div className='max-h-72 p-4'>
                {props.items.map((item, i) => (
                  <p key={i}>{item.id}</p>
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
