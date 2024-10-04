import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { SquareX, X } from "lucide-react";
import { useState } from "react";

import { Button } from "../ui/button";
import { useToast } from "../ui/use-toast";

export interface DeleteProps<TMutationVariables> {
  items: {
    id: string;
    name: string;
  }[];
  itemType: string;
  mutationFunction: (params: TMutationVariables) => Promise<void>;
  mutationVariables: TMutationVariables[];
  variant?: "large" | "small";
}

// create a functional component called DeleteItems
export const DeleteItems = <TMutationVariables,>({
  items,
  itemType,
  mutationFunction,
  mutationVariables,
  variant = "small",
}: DeleteProps<TMutationVariables>) => {
  const [openConfirmDelete, setOpenConfirmDelete] = useState(false);
  const t = (text: string) => text;
  const { toast } = useToast();
  const handleDelete = async () => {
    for (const mutationVars of mutationVariables) {
      try {
        await mutationFunction(mutationVars);
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
      <DialogTrigger>
        {variant === "small" ? (
          <div
            className="text-destructive"
            onClick={() => setOpenConfirmDelete(true)}
          >
            <SquareX />
          </div>
        ) : (
          <div onClick={() => setOpenConfirmDelete(true)}>{t("Delete")}</div>
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
            <div className="px-6 py-4">
              {items?.map((item, i) => <p key={i}>{item.name}</p>)}
            </div>
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
