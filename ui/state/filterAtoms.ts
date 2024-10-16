import { subDays } from "date-fns";
import { atom } from "jotai";

interface DateRange {
  from: Date;
  to: Date;
}

// Define atoms
export const pageAtom = atom<number>(0);
export const limitAtom = atom<number>(10);
export const queryAtom = atom<string>("");
export const rangeAtom = atom<DateRange>({
  from: subDays(new Date(), 7),
  to: new Date(),
});
export const sortByAtom = atom<string>("createdAt");
export const sortDirectionAtom = atom<"asc" | "desc">("desc");
