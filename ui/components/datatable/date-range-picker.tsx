"use client";

import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useFilterItems } from "@/hooks/useFilterItems";
import { cn } from "@/lib/utils";
import { format } from "date-fns";
import { Calendar as CalendarIcon } from "lucide-react";
import * as React from "react";

export function DatePickerWithRange({
  className,
}: React.HTMLAttributes<HTMLDivElement>) {
  const { range, setRange } = useFilterItems();

  return (
    <div className={cn("grid gap-2", className)}>
      <Popover>
        <PopoverTrigger asChild>
          <Button
            className={cn(
              "justify-start text-left font-normal h-8",
              !range && "text-muted-foreground"
            )}
            id="range"
            variant={"outline"}
          >
            <CalendarIcon className="mr-2 h-4 w-4" />
            {range?.from ? (
              range.to ? (
                <>
                  {format(range.from, "LLL dd, y")} -{" "}
                  {format(range.to, "LLL dd, y")}
                </>
              ) : (
                format(range.from, "LLL dd, y")
              )
            ) : (
              <span>Pick a date</span>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent align="start" className="w-auto p-0">
          <Calendar
            defaultMonth={range?.from}
            initialFocus
            mode="range"
            numberOfMonths={2}
            onSelect={(range) => setRange(range as any)}
            selected={range}
          />
        </PopoverContent>
      </Popover>
    </div>
  );
}
