import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuth } from "@/hooks/useAuth";
import { FC } from "react";

interface UserButtonProps {
  size: "lg" | "sm";
}

export const UserButton: FC<UserButtonProps> = ({ size }) => {
  const { defaultOrgname, logout, user } = useAuth();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          aria-label="User menu"
          className={`relative flex items-center justify-between ${
            size === "lg"
              ? "p-2 pl-2 text-sm leading-6 font-semibold"
              : "h-8 w-8 rounded-full"
          }`}
          variant="ghost"
        >
          {size === "lg" ? (
            <>
              <div className="flex items-center gap-x-2.5 max-w-[calc(100%-24px)]">
                <Avatar className="h-9 w-9">
                  <AvatarImage
                    alt={user?.displayName || "User"}
                    src={user?.photoUrl}
                  />
                  <AvatarFallback>
                    {user?.displayName
                      ?.split(" ")
                      .map((x) => x[0])
                      .join("")}
                  </AvatarFallback>
                </Avatar>
                <div className="overflow-hidden flex-grow">
                  <div className="text-start overflow-hidden">
                    <p
                      aria-hidden="true"
                      className="text-sm text-dark font-medium truncate"
                    >
                      {user?.displayName}
                    </p>
                    <p
                      aria-hidden="true"
                      className="text-xs font-normal truncate overflow-ellipsis text-gray-alpha-500"
                    >
                      {defaultOrgname}
                    </p>
                  </div>
                </div>
              </div>
              <svg
                className="text-gray-alpha-950 mr-1 flex-shrink-0 transition-transform duration-200 rotate-90 [button[data-state=open]_&]:rotate-0"
                fill="currentColor"
                height="14"
                stroke="currentColor"
                strokeWidth="0"
                viewBox="0 0 512 512"
                width="14"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path d="M294.1 256L167 129c-9.4-9.4-9.4-24.6 0-33.9s24.6-9.3 34 0L345 239c9.1 9.1 9.3 23.7.7 33.1L201.1 417c-4.7 4.7-10.9 7-17 7s-12.3-2.3-17-7c-9.4-9.4-9.4-24.6 0-33.9l127-127.1z"></path>
              </svg>
            </>
          ) : (
            <Avatar className="h-8 w-8">
              <AvatarImage
                alt={user?.displayName || "User"}
                src={user?.photoUrl}
              />
              <AvatarFallback>
                {user?.displayName
                  ?.split(" ")
                  .map((x) => x[0])
                  .join("")}
              </AvatarFallback>
            </Avatar>
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56" forceMount>
        <DropdownMenuLabel className="font-normal">
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-medium leading-none">
              {user?.displayName}
            </p>
            <p className="text-xs leading-none text-muted-foreground">
              {user?.email}
            </p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem>
            Settings
            <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuItem>
            Support
            <DropdownMenuShortcut>⌘B</DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuItem>New Team</DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={async () => await logout()}>
          Log out
          <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
