import type { ReactNode } from "react";
import { createContext, useContext } from "react";

export interface NavigationItem {
  type?: string;
  label?: string;
  path?: string;
  file?: string;
  to?: string;
  href?: string;
  icon?: string;
  children?: NavigationItem[];
  category?: string;
}

export interface ZudokuConfig {
  basePath?: string;
  site?: {
    title?: string;
    dir?: "ltr" | "rtl";
    logo?: {
      src: {
        light: string;
        dark: string;
      };
      alt?: string;
      width?: number;
      href?: string;
    };
  };
  options?: {
    basePath: string;
    site?: {
      dir?: "ltr" | "rtl";
    };
  };
  navigation?: NavigationItem[];
  plugins?: Record<string, unknown>;
}

export interface AuthContextType {
  isAuthenticated: boolean;
  isAuthEnabled: boolean;
  profile?:
    | {
        name?: string;
        email?: string;
      }
    | undefined;
  login: () => void;
  logout: () => void;
}

const ZudokuContext = createContext<ZudokuConfig | undefined>(undefined);
const NavigationContext = createContext<
  { navigation: NavigationItem[] } | undefined
>(undefined);
const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const ZudokuProvider = ({
  children,
  config,
}: {
  children: ReactNode;
  config?: ZudokuConfig;
}) => {
  const defaultConfig: ZudokuConfig = {
    basePath: "/",
    navigation: [],
    options: {
      basePath: "/",
    },
    plugins: {},
    site: {
      dir: "ltr",
      title: "Documentation",
    },
    ...config,
  };

  return (
    <ZudokuContext.Provider value={defaultConfig}>
      {children}
    </ZudokuContext.Provider>
  );
};

export const NavigationProvider = ({
  children,
  navigation = [],
}: {
  children: ReactNode;
  navigation?: NavigationItem[];
}) => {
  return (
    <NavigationContext.Provider value={{ navigation }}>
      {children}
    </NavigationContext.Provider>
  );
};

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const auth: AuthContextType = {
    isAuthEnabled: false,
    isAuthenticated: false,
    login: () => console.log("Login not implemented"),
    logout: () => console.log("Logout not implemented"),
  };

  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
};

export const useZudoku = () => {
  const context = useContext(ZudokuContext);
  if (!context) {
    return {
      basePath: "/",
      navigation: [],
      options: { basePath: "/", site: { dir: "ltr" as const } },
      plugins: {},
      site: { dir: "ltr" as const, title: "Documentation" },
    };
  }
  return context;
};

export const useCurrentNavigation = () => {
  const context = useContext(NavigationContext);
  if (!context) {
    return { navigation: [] };
  }
  return context;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    return {
      isAuthEnabled: false,
      isAuthenticated: false,
      login: () => console.log("Login not implemented"),
      logout: () => console.log("Logout not implemented"),
      profile: undefined,
    };
  }
  return context;
};
