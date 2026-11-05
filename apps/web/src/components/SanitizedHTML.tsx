import React from 'react';
import DOMPurify from 'dompurify';

interface SanitizedHTMLProps {
  html: string;
  className?: string;
}

export const SanitizedHTML: React.FC<SanitizedHTMLProps> = ({ html, className }) => {
  const clean = DOMPurify.sanitize(html || '', {
    USE_PROFILES: { html: true },
    ADD_ATTR: ['target', 'rel'],
  });

  return (
    <div
      className={className}
      dangerouslySetInnerHTML={{ __html: clean }}
    />
  );
};
