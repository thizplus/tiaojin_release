import { Navigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

export function DashboardRedirect() {
  const { user } = useAuth();

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  // Redirect based on user role
  switch (user.role) {
    case 'super':
    case 'admin':
      return <Navigate to="/admin/dashboard" replace />;
    case 'agent':
      return <Navigate to="/agent/dashboard" replace />;
    case 'sales':
      return <Navigate to="/sale/dashboard" replace />;
    default:
      return <Navigate to="/unauthorized" replace />;
  }
}