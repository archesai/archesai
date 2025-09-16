import { PureUserButton } from "@archesai/ui";
import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";

export function UserButtonContainer(): JSX.Element {
  const navigate = useNavigate();

  // Fetch user data
  const { data: user } = useQuery({
    queryFn: async () => {
      // TODO: Replace with actual API call
      return {
        email: "john@example.com",
        id: "1",
        image: null,
        name: "John Doe",
      };
    },
    queryKey: ["currentUser"],
  });

  const handleLogout = async () => {
    // TODO: Implement logout logic
    await navigate({ to: "/auth/login" });
  };

  const handleNavigateToProfile = async () => {
    await navigate({ to: "/profile" });
  };

  const handleNavigateToBilling = async () => {
    await navigate({ to: "/organization" });
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
