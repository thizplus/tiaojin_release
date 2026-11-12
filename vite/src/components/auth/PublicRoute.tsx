import { useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";

interface PublicRouteProps {
  children: React.ReactNode;
  redirectTo?: string;
  restricted?: boolean; // If true, authenticated users will be redirected
}

export function PublicRoute({
  children,
  redirectTo = "/dashboard",
  restricted = false
}: PublicRouteProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    // Don't redirect while still loading auth state
    if (isLoading) return;

    // If this is a restricted public route (like login/register) and user is authenticated
    if (restricted && isAuthenticated) {
      // Check if there's a return URL from login redirect
      const from = location.state?.from || redirectTo;
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, isLoading, navigate, redirectTo, restricted, location.state]);

  // Show loading while checking auth
  if (isLoading) {
    return (
      <div className="flex min-h-svh items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
          <p className="mt-2 text-sm text-muted-foreground">กำลังตรวจสอบสิทธิ์...</p>
        </div>
      </div>
    );
  }

  // For restricted routes, don't show content if authenticated
  if (restricted && isAuthenticated) {
    return null;
  }

  return <>{children}</>;
}