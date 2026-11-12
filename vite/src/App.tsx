// src/App.tsx
import { BrowserRouter } from 'react-router-dom'

import { ThemeProvider } from '@/theme/theme-provider';
import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'

import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { ConditionalWebSocketProvider } from '@/components/providers/ConditionalWebSocketProvider';
import AppRoutes from './routes';
import { clearOldStorage, setGlobalQueryClient } from '@/stores/authStore';

// Clear old storage format ก่อน app start
clearOldStorage();

const queryClient = new QueryClient()

// Set global query client สำหรับ forceLogout
setGlobalQueryClient(queryClient);

function App() {
  return (
    <BrowserRouter>
      <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <QueryClientProvider client={queryClient}>
          <ConditionalWebSocketProvider>
            <AppRoutes />
            <ReactQueryDevtools initialIsOpen={false} />
          </ConditionalWebSocketProvider>
        </QueryClientProvider>
      </ThemeProvider>
    </BrowserRouter>
  )
}

export default App;