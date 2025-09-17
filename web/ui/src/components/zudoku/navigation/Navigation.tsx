import { DrawerContent, DrawerTitle } from "#components/shadcn/drawer";
import type { NavigationItem as NavigationItemType } from "../context/ZudokuContext";
import { Slot } from "../Slot";
import { NavigationItem } from "./NavigationItem";
import { NavigationWrapper } from "./NavigationWrapper";

export const Navigation = ({
  onRequestClose,
  navigation,
}: {
  onRequestClose?: () => void;
  navigation: NavigationItemType[];
}) => (
  <>
    <NavigationWrapper>
      <Slot.Target name="navigation-before" />
      {navigation.map((item) => (
        <NavigationItem
          item={item}
          key={
            item.type +
            (item.label ?? "") +
            ("path" in item ? item.path : "") +
            ("file" in item ? item.file : "") +
            ("to" in item ? item.to : "")
          }
        />
      ))}
      <Slot.Target name="navigation-after" />
    </NavigationWrapper>
    <DrawerContent
      aria-describedby={undefined}
      className="start-0 h-[100dvh] w-[320px] rounded-none lg:hidden"
    >
      <div className="overflow-y-auto overscroll-none p-4">
        <DrawerTitle className="sr-only">Navigation</DrawerTitle>
        {navigation.map((item) => (
          <NavigationItem
            item={item}
            key={item.label}
            {...(onRequestClose && { onRequestClose })}
          />
        ))}
      </div>
    </DrawerContent>
  </>
);
