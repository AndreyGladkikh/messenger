import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from "react-router/dom";
import { ThemeProvider } from "@/shared/ui/theme-provider"
import './index.css'
import router from '@/app/router';
import {
  QueryClientProvider,
} from '@tanstack/react-query'
import queryClient from './shared/lib/query-client';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </ThemeProvider>
  </StrictMode>,
)
