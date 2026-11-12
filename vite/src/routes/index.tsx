// src/routes/index.tsx
import { Routes, Route } from 'react-router-dom';
import { lazy, Suspense } from 'react';

// Auth Components
import { ProtectedRoute, PublicRoute } from '@/components/auth';
import { DashboardRedirect } from '@/components/DashboardRedirect';

// Shared Pages
const LoginPage = lazy(() => import('@/page/globals/login/page'));
const RegisterPage = lazy(() => import('@/page/globals/register/page'));
const ProfilePage = lazy(() => import('@/page/globals/profile/page'));
const HomePage = lazy(() => import('@/page/globals/home/page'));


// Super Pages
const SuperSystemSettingsPage = lazy(() => import('@/page/super/system-settings/page'));

// Admin Pages
const AdminDashboardPage = lazy(() => import('@/page/admin/dashboard/page'));

const AdminBookingsPage = lazy(() => import('@/page/admin/bookings/page'));

const AdminToursPage = lazy(() => import('@/page/admin/tours/page'));
const AdminPeriodsPage = lazy(() => import('@/page/admin/periods/page'));
const AdminCitiesPage = lazy(() => import('@/page/admin/cities/page'));
const AdminUsersPage = lazy(() => import('@/page/admin/users/page'));
const AdminSalesPage = lazy(() => import('@/page/admin/sales/page'));
const AdminProfilesPage = lazy(() => import('@/page/admin/profiles/page'));
const AdminBannersPage = lazy(() => import('@/page/admin/banners/page'));
const AdminNotificationsPage = lazy(() => import('@/page/admin/notifications/page'));

const AdminReportsToursPage = lazy(() => import('@/page/admin/reports-tours/page'));
const AdminReportsDaterangePage = lazy(() => import('@/page/admin/reports-daterange/page'));

// Agent Pages
const AgentDashboardPage = lazy(() => import('@/page/agent/dashboard/page'));
const AgentToursPage = lazy(() => import('@/page/agent/tours/page'));
const AgentToursCityPage = lazy(() => import('@/page/agent/tours-city/page'));
const AgentToursHotPage = lazy(() => import('@/page/agent/tours-hot/page'));
const AgentBookingsPage = lazy(() => import('@/page/agent/bookings/page'));
const AgentProfilePage = lazy(() => import('@/page/agent/profile/page'));
const AgentNotificationsPage = lazy(() => import('@/page/agent/notifications/page'));
const AgentReportsPage = lazy(() => import('@/page/agent/reports/page'));
const AgentApiKeysPage = lazy(() => import('@/page/agent/api-keys/page'));

// Sale Pages
const SaleDashboardPage = lazy(() => import('@/page/sale/dashboard/page'));


const SaleBookingsPage = lazy(() => import('@/page/sale/bookings/page'));


// Layouts
import RootLayout from '@/layouts/root-layout';
import PageLayout from '@/layouts/page';
import AdminToursRegularPage from '@/page/admin/tours-regular/page';
import AdminToursSpecialPage from '@/page/admin/tours-special/page';

// Error Pages
const UnauthorizedPage = lazy(() => import('@/page/globals/unauthorized/page'));

// Loading component
const LoadingFallback = () => (
  <div className="flex items-center justify-center min-h-screen">
    <div className="text-center">
      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
      <p className="mt-2 text-sm text-muted-foreground">กำลังโหลด...</p>
    </div>
  </div>
);

const AppRoutes = () => {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <Routes>
        <Route element={<RootLayout />}>
          {/* Public routes - accessible without authentication */}
          <Route
            path="/login"
            element={
              <PublicRoute restricted>
                <LoginPage />
              </PublicRoute>
            }
          />
          <Route
            path="/register"
            element={
              <PublicRoute restricted>
                <RegisterPage />
              </PublicRoute>
            }
          />

          {/* Error pages */}
          <Route
            path="/unauthorized"
            element={
              <PublicRoute>
                <UnauthorizedPage />
              </PublicRoute>
            }
          />

          {/* Public Home Page - แสดงรายการทัวร์ */}
          <Route
            path="/"
            element={
              <PublicRoute>
                <HomePage />
              </PublicRoute>
            }
          />

          {/* Protected routes with layout */}
          <Route
            element={
              <ProtectedRoute>
                <PageLayout />
              </ProtectedRoute>
            }
          >
            {/* Global dashboard route - redirects to role-specific dashboard */}
            <Route
              path="/dashboard"
              element={
                <ProtectedRoute>
                  <DashboardRedirect />
                </ProtectedRoute>
              }
            />

            {/* Shared/Global routes */}
            <Route
              path="/profile"
              element={
                <ProtectedRoute>
                  <ProfilePage />
                </ProtectedRoute>
              }
            />

            {/* Super routes */}
            <Route
              path="/super/system-settings"
              element={
                <ProtectedRoute requiredRoles={['super']}>
                  <SuperSystemSettingsPage />
                </ProtectedRoute>
              }
            />

            {/* Admin routes */}
            <Route
              path="/admin/dashboard"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminDashboardPage />
                </ProtectedRoute>
              }
            />
            
            <Route
              path="/admin/tours-regular"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminToursRegularPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/tours-special"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminToursSpecialPage />
                </ProtectedRoute>
              }
            />


            <Route
              path="/admin/bookings"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminBookingsPage />
                </ProtectedRoute>
              }
            />
            
            <Route
              path="/admin/tours"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminToursPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/tours/:tourId/periods"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminPeriodsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/cities"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminCitiesPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/users"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminUsersPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/sales"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminSalesPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/profiles"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminProfilesPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/banners"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminBannersPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/notifications"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminNotificationsPage />
                </ProtectedRoute>
              }
            />
         
            <Route
              path="/admin/reports-tours"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminReportsToursPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/reports-daterange"
              element={
                <ProtectedRoute requiredRoles={['admin', 'super']}>
                  <AdminReportsDaterangePage />
                </ProtectedRoute>
              }
            />

            {/* Agent routes */}
            <Route
              path="/agent/dashboard"
              element={
                <ProtectedRoute requiredRoles={['agent']}>
                  <AgentDashboardPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/agent/tours"
              element={
                <ProtectedRoute requiredRoles={['agent']}>
                  <AgentToursPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/agent/cities/:cityId"
              element={
                <ProtectedRoute requiredRoles={['agent']}>
                  <AgentToursCityPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/agent/tours-hot"
              element={
                <ProtectedRoute requiredRoles={['agent']}>
                  <AgentToursHotPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/agent/bookings"
              element={
                <ProtectedRoute requiredRoles={['agent']}>
                  <AgentBookingsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/agent/profile"
              element={
                <ProtectedRoute requiredRoles={['agent']}>
                  <AgentProfilePage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/agent/notifications"
              element={
                <ProtectedRoute requiredRoles={['agent']}>
                  <AgentNotificationsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/agent/reports"
              element={
                <ProtectedRoute requiredRoles={['agent']}>
                  <AgentReportsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/agent/api-keys"
              element={
                <ProtectedRoute requiredRoles={['agent']}>
                  <AgentApiKeysPage />
                </ProtectedRoute>
              }
            />

            {/* Sale routes */}
            <Route
              path="/sale/dashboard"
              element={
                <ProtectedRoute requiredRoles={['sales']}>
                  <SaleDashboardPage />
                </ProtectedRoute>
              }
            />
          
         
            <Route
              path="/sale/bookings"
              element={
                <ProtectedRoute requiredRoles={['sales']}>
                  <SaleBookingsPage />
                </ProtectedRoute>
              }
            />
           
          </Route>

          {/* 404 - redirect to home page */}
          <Route
            path="*"
            element={
              <PublicRoute>
                <HomePage />
              </PublicRoute>
            }
          />
        </Route>
      </Routes>
    </Suspense>
  );
};

export default AppRoutes;