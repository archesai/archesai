import { addDays } from "date-fns";
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
  from: new Date(2022, 0, 20),
  to: addDays(new Date(2022, 0, 20), 20),
});
