import { PureUserButton } from "@archesai/ui";
import type { JSX } from "react";
import { useCallback } from "react";

// import { useGetUser } from "#lib/index";

type UserButtonContainerProps = {
  navigate: (to: { to: string }) => Promise<void>;
  session?: {
    data: {
      id: string;
      userID: string;
    };
  };
  userData: {
    data: {
      email: string;
      id: string;
      name?: string | null;
    };
  };
  deleteSession: (args: { id: string }) => Promise<void>;
};

export function UserButtonContainer({
  navigate,
  session,
  userData,
  deleteSession,
}: UserButtonContainerProps): JSX.Element {
  // const navigate = useNavigate();
  // const queryClient = useQueryClient();
  // const { session } = useRouteContext({ from: "__root__" });
  //   const { mutateAsync: deleteSession } = useDeleteSession();

  const sessionData = session?.data;

  //   const { data: userData } = useGetUser(sessionData?.userID, {
  //     query: {
  //       enabled: !!sessionData?.userID,
  //     },
  //   });

  const user = userData?.data
    ? {
        email: userData.data.email,
        id: userData.data.id,
        name: userData.data.name || userData.data.email,
        picture: null,
      }
    : null;

  const handleLogout = useCallback(async () => {
    try {
      if (sessionData?.id) {
        await deleteSession({ id: sessionData.id });
      }
      // queryClient.clear();
      await navigate({
        to: "/auth/login",
      });
    } catch (error) {
      console.error("Logout error:", error);
      await navigate({
        to: "/auth/login",
      });
    }
  }, [
    sessionData?.id,
    deleteSession,
    // queryClient,
    navigate,
  ]);

  const handleNavigateToProfile = async () => {
    await navigate({
      to: "/profile",
    });
  };

  const handleNavigateToBilling = async () => {
    await navigate({
      to: "/organization",
    });
  };

  if (!user) {
    return <PureUserButton />;
  }

  return (
    <PureUserButton
      onLogout={handleLogout}
      onNavigateToBilling={handleNavigateToBilling}
      onNavigateToProfile={handleNavigateToProfile}
      user={user}
    />
  );
}
