import React from 'react';
import { Link } from 'react-router-dom';

export interface AnimatedButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  as?: 'button' | 'a' | 'link';
  to?: string;
  href?: string;
  target?: string;
  rel?: string;
  size?: 'sm' | 'default' | 'lg';
  fullWidth?: boolean;
  children: React.ReactNode;
}

export const AnimatedButton: React.FC<AnimatedButtonProps> = ({
  as = 'button',
  to,
  href,
  target,
  rel,
  size = 'default',
  fullWidth = false,
  className = '',
  children,
  disabled,
  type = 'button',
  ...rest
}) => {
  const sizeClass = size === 'sm' ? 'animated-button-sm' : size === 'lg' ? 'animated-button-lg' : '';
  const widthClass = fullWidth ? 'w-full flex' : '';
  const combinedClass = `animated-button ${sizeClass} ${widthClass} ${className}`.trim();

  const content = (
    <>
      <span className="text">{children}</span>
      <span className="circle" aria-hidden="true" />
    </>
  );

  if (as === 'link' || to) {
    return (
      <Link to={to || '#'} className={combinedClass} {...(rest as any)}>
        {content}
      </Link>
    );
  }

  if (as === 'a' || href) {
    return (
      <a href={href} target={target} rel={rel || (target === '_blank' ? 'noopener noreferrer' : undefined)} className={combinedClass} {...(rest as any)}>
        {content}
      </a>
    );
  }

  return (
    <button type={type} disabled={disabled} className={combinedClass} {...rest}>
      {content}
    </button>
  );
};

export default AnimatedButton;
