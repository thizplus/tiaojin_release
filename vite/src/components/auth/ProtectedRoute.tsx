import { useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRoles?: string[];
  redirectTo?: string;
}

export function ProtectedRoute({
  children,
  requiredRoles = [],
  redirectTo = "/login"
}: ProtectedRouteProps) {
  const { user, isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    // Don't redirect while still loading auth state
    if (isLoading) return;

    // If not authenticated, redirect to login with return URL
    if (!isAuthenticated) {
      navigate(redirectTo, {
        replace: true,
        state: { from: location.pathname }
      });
      return;
    }

    // If authenticated but no user data yet, wait
    if (!user) return;

    // Check role requirements if specified
    if (requiredRoles.length > 0 && !requiredRoles.includes(user.role)) {
      navigate("/unauthorized", { replace: true });
      return;
    }

  }, [isAuthenticated, user, isLoading, navigate, redirectTo, location.pathname, requiredRoles]);

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

  // Don't render content if not authenticated or role check failed
  if (!isAuthenticated || !user) {
    return null;
  }

  // Check role requirements
  if (requiredRoles.length > 0 && !requiredRoles.includes(user.role)) {
    return null;
  }

  return <>{children}</>;
}