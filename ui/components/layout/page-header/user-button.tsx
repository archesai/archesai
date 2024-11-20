import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  useUserControllerFindOne,
  useUserControllerUpdate,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/use-auth";
import { useToast } from "@/hooks/use-toast";
import { useRouter } from "next/navigation";
import { FC } from "react";

import { Badge } from "../../ui/badge";

interface UserButtonProps {
  size: "lg" | "sm";
}

export const UserButton: FC<UserButtonProps> = ({ size }) => {
  const { defaultOrgname, logout } = useAuth();
  const { toast } = useToast();
  const router = useRouter();
  const { data: user } = useUserControllerFindOne({});

  const { mutateAsync: updateDefaultOrg } = useUserControllerUpdate({
    onError: (error) => {
      toast({
        description: error?.stack.message,
        title: "Error updating default organization",
        variant: "destructive",
      });
    },
    onSuccess: () => {
      toast({
        description: "Your default organization has been updated.",
        title: "Default organization updated",
      });
    },
  });

  const memberships = user?.memberships;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          className={`flex h-auto items-center justify-between ${
            size === "lg"
              ? "w-full p-2 font-semibold leading-6"
              : "h-8 w-8 rounded-full p-0"
          }`}
          variant="outline"
        >
          {size === "lg" ? (
            <>
              <div className="flex min-w-0 items-center gap-2">
                <Avatar className="h-10 w-10">
                  <AvatarImage
                    alt={user?.displayName || "User"}
                    src={
                      user?.photoUrl ||
                      `https://ui-avatars.com/api/?name=${user?.displayName
                        ?.split(" ")
                        .map((x) => x[0])
                        .join("+")}&background=3D61FF&color=fff`
                    }
                  />
                  <AvatarFallback>
                    {user?.displayName
                      ?.split(" ")
                      .map((x) => x[0])
                      .join("")}
                  </AvatarFallback>
                </Avatar>
                <div className="overflow-hidden text-start">
                  <p
                    aria-hidden="true"
                    className="text-dark truncate text-sm font-medium"
                  >
                    {user?.displayName}
                  </p>
                  <p
                    aria-hidden="true"
                    className="text-gray-alpha-500 truncate text-xs font-normal"
                  >
                    {defaultOrgname}
                  </p>
                </div>
              </div>
              <svg
                className="text-gray-alpha-950 flex-shrink-0 rotate-90 transition-transform duration-200 [button[data-state=open]_&]:rotate-0"
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
                src={
                  user?.photoUrl ||
                  `https://ui-avatars.com/api/?name=${user?.displayName
                    ?.split(" ")
                    .map((x) => x[0])
                    .join("+")}&background=3D61FF&color=fff`
                }
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
          <div className="flex flex-col gap-1">
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
          <DropdownMenuSub>
            <DropdownMenuSubTrigger>Organizations</DropdownMenuSubTrigger>
            <DropdownMenuPortal>
              <DropdownMenuSubContent>
                {memberships?.map((membership) => (
                  <DropdownMenuItem
                    className="flex justify-between gap-2"
                    key={membership.id}
                    onClick={() => {
                      updateDefaultOrg({
                        body: {
                          defaultOrgname: membership.orgname,
                        },
                      });
                    }}
                  >
                    {membership.orgname}
                    {defaultOrgname === membership.orgname && (
                      <Badge>Current</Badge>
                    )}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuSubContent>
            </DropdownMenuPortal>
          </DropdownMenuSub>
          <DropdownMenuItem
            onClick={() => router.push("/organization/general")}
          >
            Settings
            <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
          </DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem onClick={() => router.push("/profile/general")}>
            Profile
            <DropdownMenuShortcut>⌘P</DropdownMenuShortcut>
          </DropdownMenuItem>
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
