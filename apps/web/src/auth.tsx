import { createContext, useContext } from "react";
import { User } from "@/types/user";
import { useQuery } from "@tanstack/react-query";
import { api } from "./lib/ky";

export interface AuthContext {
  user: User | undefined;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContext | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { data: user, isLoading } = useQuery({
    queryKey: ["user"],
    queryFn: async () => {
      const result = await api.get("user", {}).json();
      return User.parse(result);
    },
    retry: false,
  });

  return (
    <AuthContext.Provider value={{ user, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}