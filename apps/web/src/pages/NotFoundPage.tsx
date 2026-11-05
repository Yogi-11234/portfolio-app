import React from 'react';
import { Home } from 'lucide-react';
import { AnimatedButton } from '@/components/AnimatedButton';

export const NotFoundPage: React.FC = () => {
  return (
    <div className="min-h-[70vh] flex flex-col items-center justify-center text-center px-4">
      <span className="font-mono text-sm text-accent tracking-widest uppercase">404 Error</span>
      <h1 className="mt-2 text-3xl sm:text-5xl font-extrabold text-text">Page Not Found</h1>
      <p className="mt-3 text-sm text-text-muted max-w-md">
        The page you are looking for might have been moved, deleted, or does not exist.
      </p>
      <div className="mt-6">
        <AnimatedButton as="link" to="/">
          <Home className="w-4 h-4 inline mr-1" />
          <span>Back to Home</span>
        </AnimatedButton>
      </div>
    </div>
  );
};
