import type { JSX } from "react";

import { useTransition } from "react";
import { toast } from "sonner";
import { Loader2Icon, TrashIcon } from "#components/custom/icons";
import { Button } from "#components/shadcn/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "#components/shadcn/dialog";
import { ScrollArea } from "#components/shadcn/scroll-area";
import { Separator } from "#components/shadcn/separator";
import type { BaseEntity } from "#types/entities";

interface DeleteItemsProps<TEntity extends BaseEntity>
  extends React.ComponentPropsWithoutRef<typeof Dialog> {
  deleteItem: (id: string) => Promise<void>;
  entityKey: string;
  items: TEntity[];
  showTrigger?: boolean;
}

export const DeleteItems = <TEntity extends BaseEntity>(
  props: DeleteItemsProps<TEntity>,
): JSX.Element => {
  const [isDeletePending, startDeleteTransition] = useTransition();

  function onDelete() {
    startDeleteTransition(async () => {
      for (const item of props.items) {
        try {
          await props.deleteItem(item.id);
          toast(t(`The ${props.entityKey} has been removed`));
        } catch (error: unknown) {
          if (error instanceof Error) {
            toast(t(`Could not remove ${props.entityKey}`), {
              description: error.message,
            });
            props.onOpenChange?.(false);
            toast.success("Tasks deleted");
          } else {
            console.error(error);
          }
        }
      }
    });
  }

  const t = (text: string) => text;

  return (
    <Dialog {...props}>
      {props.showTrigger ? (
        <DialogTrigger asChild>
          <Button
            size="sm"
            variant="outline"
          >
            <TrashIcon
              aria-hidden="true"
              className="mr-2 size-4"
            />
            Delete ({props.items.length})
          </Button>
        </DialogTrigger>
      ) : null}

      <DialogContent className="gap-0 p-0">
        <DialogHeader>
          <DialogTitle>Are you absolutely sure?</DialogTitle>
          <DialogDescription>
            This action cannot be undone. This will permanently delete your{" "}
            <span className="font-medium">{props.items.length}</span>
            {props.items.length === 1 ? " item" : " items"} from our servers.
          </DialogDescription>
        </DialogHeader>
        <div className="flex flex-col items-center justify-center gap-3 p-4">
          <TrashIcon className="text-destructive" />
          <p className="text-center">
            {t(
              `Are you sure you want to permanently delete the following ${props.entityKey}${props.items.length > 1 ? "s" : ""}?`,
            )}
          </p>
          <ScrollArea>
            <div className="max-h-72 p-4">
              {props.items.map((item) => (
                <p key={item.id}>{item.id}</p>
              ))}
            </div>
          </ScrollArea>
        </div>
        <Separator />
        <DialogFooter className="gap-2 sm:gap-x-0">
          <DialogClose asChild>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button
            aria-label="Delete selected rows"
            disabled={isDeletePending}
            onClick={onDelete}
            variant="destructive"
          >
            {isDeletePending && (
              <Loader2Icon
                aria-hidden="true"
                className="mr-2 size-4 animate-spin"
              />
            )}
            Delete
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
