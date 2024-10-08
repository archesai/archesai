import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { SquareX, X } from "lucide-react";
import { useState } from "react";

import { Button } from "../ui/button";
import { ScrollArea } from "../ui/scroll-area";
import { useToast } from "../ui/use-toast";

export interface DeleteProps<TDeleteVariables> {
  deleteFunction: (params: TDeleteVariables) => Promise<void>;
  deleteVariables: TDeleteVariables[];
  items: {
    id: string;
    name: string;
  }[];
  itemType: string;
  variant?: "lg" | "md" | "sm";
}

// create a functional component called DeleteItems
export const DeleteItems = <TDeleteVariables,>({
  deleteFunction,
  deleteVariables,
  items,
  itemType,
  variant = "sm",
}: DeleteProps<TDeleteVariables>) => {
  const [openConfirmDelete, setOpenConfirmDelete] = useState(false);
  const t = (text: string) => text;
  const { toast } = useToast();
  const handleDelete = async () => {
    for (const deleteVars of deleteVariables) {
      try {
        await deleteFunction(deleteVars);
        setOpenConfirmDelete(false);
        toast({ title: t(`The ${itemType} has been removed`) });
      } catch (err) {
        console.error(err);
        toast({ title: t(`Could not remove ${itemType}`) });
      }
    }
  };

  return (
    <Dialog
      onOpenChange={(open) => setOpenConfirmDelete(open)}
      open={openConfirmDelete}
    >
      <DialogTrigger asChild>
        {variant === "sm" ? (
          <div
            className="text-destructive cursor-pointer"
            onClick={() => setOpenConfirmDelete(true)}
          >
            <SquareX className="h-5 w-5" />
          </div>
        ) : variant === "md" ? (
          <div className="w-full" onClick={() => setOpenConfirmDelete(true)}>
            {t("Delete")}
          </div>
        ) : (
          <Button
            className="h-8"
            onClick={() => setOpenConfirmDelete(true)}
            variant="destructive"
          >
            {t("Delete")}
          </Button>
        )}
      </DialogTrigger>

      <DialogContent>
        <div className="flex flex-col items-center justify-center p-5 gap-3">
          <X />

          <p className="text-center">
            {t(
              `Are you sure you want to permanently delete the following ${itemType}${
                items?.length > 1 ? "s" : ""
              }?`
            )}
          </p>
          {
            <ScrollArea>
              <div className="max-h-72 px-6 py-4">
                {items?.map((item, i) => <p key={i}>{item.name}</p>)}
              </div>
            </ScrollArea>
          }

          <div className="flex-1 flex gap-4">
            <Button onClick={() => setOpenConfirmDelete(false)} size="sm">
              {t("Cancel")}
            </Button>
            <Button
              onClick={async () => await handleDelete()}
              size="sm"
              variant={"destructive"}
            >
              {t("Delete")}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};
