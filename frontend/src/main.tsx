import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { RouterProvider } from "react-router/dom";
import router from './app/router';
import { QueryClientProvider } from '@tanstack/react-query';
import queryClient from './shared/lib/query-client';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <div className="typeset">
        <RouterProvider router={router} />
      </div>
    </QueryClientProvider>
  </StrictMode>,
)
