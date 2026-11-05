import React from 'react';
import { getMediaUrl } from '@/lib/utils';
import {
  Code2,
  Terminal,
  Database,
  Cloud,
  Globe,
} from 'lucide-react';

interface TechIconProps {
  name: string;
  iconKey?: string;
  className?: string;
  brandColor?: string;
}

// Built-in SVG maps for top tech stacks
export const TECH_PRESETS: { key: string; label: string; color: string }[] = [
  { key: 'golang', label: 'Go (Golang)', color: '#00ADD8' },
  { key: 'react', label: 'React / React Native', color: '#61DAFB' },
  { key: 'typescript', label: 'TypeScript', color: '#3178C6' },
  { key: 'javascript', label: 'JavaScript', color: '#F7DF1E' },
  { key: 'python', label: 'Python', color: '#3776AB' },
  { key: 'nodejs', label: 'Node.js', color: '#5FA04E' },
  { key: 'nextjs', label: 'Next.js', color: '#ffffff' },
  { key: 'vue', label: 'Vue.js', color: '#4FC08D' },
  { key: 'angular', label: 'Angular', color: '#DD0031' },
  { key: 'tailwind', label: 'Tailwind CSS', color: '#06B6D4' },
  { key: 'html5', label: 'HTML5', color: '#E34F26' },
  { key: 'css3', label: 'CSS3', color: '#1572B6' },
  { key: 'php', label: 'PHP', color: '#777BB4' },
  { key: 'laravel', label: 'Laravel', color: '#FF2D20' },
  { key: 'java', label: 'Java', color: '#007396' },
  { key: 'rust', label: 'Rust', color: '#DEA584' },
  { key: 'cpp', label: 'C++', color: '#00599C' },
  { key: 'csharp', label: 'C# / .NET', color: '#239120' },
  { key: 'docker', label: 'Docker', color: '#2496ED' },
  { key: 'kubernetes', label: 'Kubernetes', color: '#326CE5' },
  { key: 'aws', label: 'Amazon Web Services', color: '#FF9900' },
  { key: 'gcp', label: 'Google Cloud Platform', color: '#4285F4' },
  { key: 'azure', label: 'Microsoft Azure', color: '#0078D4' },
  { key: 'linux', label: 'Linux', color: '#FCC624' },
  { key: 'postgresql', label: 'PostgreSQL', color: '#4169E1' },
  { key: 'mysql', label: 'MySQL', color: '#4479A1' },
  { key: 'mongodb', label: 'MongoDB', color: '#47A248' },
  { key: 'redis', label: 'Redis', color: '#DC382D' },
  { key: 'graphql', label: 'GraphQL', color: '#E10098' },
  { key: 'git', label: 'Git', color: '#F05032' },
  { key: 'github', label: 'GitHub', color: '#ffffff' },
  { key: 'figma', label: 'Figma', color: '#F24E1E' },
];

export const TechIcon: React.FC<TechIconProps> = ({
  name,
  iconKey,
  className = 'w-8 h-8',
  brandColor,
}) => {
  const normalizedKey = (iconKey || name || '').toLowerCase().replace(/[^a-z0-9]/g, '');

  // 1. If uploaded media path (starts with skills/ or uploads/ or has image extension)
  if (
    iconKey &&
    (iconKey.startsWith('skills/') ||
      iconKey.startsWith('uploads/') ||
      /\.(png|jpe?g|webp|gif|svg)$/i.test(iconKey))
  ) {
    const src = getMediaUrl(iconKey);
    return (
      <img
        src={src}
        alt={name}
        className={`${className} object-contain transition-transform`}
      />
    );
  }

  // 2. Specific tech SVGs
  switch (normalizedKey) {
    case 'golang':
    case 'go':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#00ADD8'} className={className}>
          <path d="M1.811 10.231c-.047-.16-.058-.328-.031-.494.053-.325.264-.606.568-.758l8.28-4.14a1.8 1.8 0 0 1 1.611 0l8.28 4.14c.304.152.515.433.568.758.027.166.016.334-.031.494l-2.072 6.906a1.8 1.8 0 0 1-1.127 1.189l-5.892 1.964a1.8 1.8 0 0 1-1.139 0l-5.892-1.964a1.8 1.8 0 0 1-1.127-1.189l-2.072-6.906z" opacity="0.3"/>
          <path d="M12 2a1.5 1.5 0 0 0-.67.16l-8.28 4.14A1.5 1.5 0 0 0 2.2 7.64l2.07 6.9a1.5 1.5 0 0 0 .94 1l5.89 1.96a1.5 1.5 0 0 0 .95 0l5.89-1.96a1.5 1.5 0 0 0 .94-1l2.07-6.9a1.5 1.5 0 0 0-.85-1.34L12.67 2.16A1.5 1.5 0 0 0 12 2zm0 2.5 6.6 3.3-1.6 5.3-5 1.66-5-1.66-1.6-5.3L12 4.5z"/>
        </svg>
      );
    case 'react':
    case 'reactnative':
    case 'reactjs':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke={brandColor || '#61DAFB'} strokeWidth="1.5" className={className}>
          <ellipse cx="12" cy="12" rx="4" ry="11" />
          <ellipse cx="12" cy="12" rx="4" ry="11" transform="rotate(60 12 12)" />
          <ellipse cx="12" cy="12" rx="4" ry="11" transform="rotate(120 12 12)" />
          <circle cx="12" cy="12" r="2" fill={brandColor || '#61DAFB'} />
        </svg>
      );
    case 'typescript':
    case 'ts':
      return (
        <svg viewBox="0 0 24 24" className={className}>
          <rect width="24" height="24" rx="4" fill={brandColor || '#3178C6'} />
          <path d="M4 10h6v2.2H7.7v6.6H5.4v-6.6H4V10zm8.3 4.2c.5-.3 1.2-.5 1.9-.5.8 0 1.3.4 1.3 1 0 .6-.4.9-1.3 1.3-1.3.5-2.2 1.2-2.2 2.4 0 1.5 1.2 2.5 3 2.5 1 0 1.8-.3 2.3-.6l-.5-1.7c-.5.3-1.1.5-1.8.5-.7 0-1.1-.3-1.1-.8 0-.6.5-.9 1.4-1.3 1.4-.6 2.2-1.3 2.2-2.5 0-1.5-1.1-2.4-2.8-2.4-.9 0-1.8.2-2.3.6l-.1 1.7z" fill="#ffffff" />
        </svg>
      );
    case 'javascript':
    case 'js':
      return (
        <svg viewBox="0 0 24 24" className={className}>
          <rect width="24" height="24" rx="4" fill={brandColor || '#F7DF1E'} />
          <path d="M7 11.5v6c0 1.5-.7 2-2 2H4v-1.8h.8c.5 0 .6-.2.6-.6v-5.6H7zm6.7 4c.6-.4 1.3-.6 2.1-.6.9 0 1.4.4 1.4 1.1 0 .7-.4 1-1.4 1.4-1.4.6-2.4 1.3-2.4 2.7 0 1.7 1.3 2.8 3.3 2.8 1.1 0 2-.3 2.5-.7l-.5-1.9c-.6.3-1.2.6-2 .6-.8 0-1.2-.3-1.2-.9 0-.7.5-1 1.5-1.5 1.5-.7 2.4-1.4 2.4-2.8 0-1.7-1.2-2.7-3.1-2.7-1 0-2 .2-2.6.7l0 1.8z" fill="#000000" />
        </svg>
      );
    case 'python':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#3776AB'} className={className}>
          <path d="M11.87 2c-3.1 0-4.9.42-4.9 2.5v1.85h5v.62H4.45c-2.07 0-3.45 1.35-3.45 3.75 0 2.36 1.16 3.73 3.45 3.73h1.54v-1.85c0-1.53 1.25-2.77 2.8-2.77h4.92c1.25 0 2.15-.9 2.15-2.15V5.15C15.86 3.1 14.28 2 11.87 2zm-1.85 1.23a.92.92 0 1 1 0 1.85.92.92 0 0 1 0-1.85z"/>
          <path d="M12.13 22c3.1 0 4.9-.42 4.9-2.5v-1.85h-5v-.62h7.52c2.07 0 3.45-1.35 3.45-3.75 0-2.36-1.16-3.73-3.45-3.73H18v1.85c0 1.53-1.25 2.77-2.8 2.77h-4.92c-1.25 0-2.15.9-2.15 2.15v2.53C8.13 20.9 9.72 22 12.13 22zm1.85-1.23a.92.92 0 1 1 0-1.85.92.92 0 0 1 0 1.85z" opacity="0.8"/>
        </svg>
      );
    case 'nodejs':
    case 'node':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#5FA04E'} className={className}>
          <path d="M12 2.5 3.5 7.4v9.8L12 22.1l8.5-4.9V7.4L12 2.5zm0 2.3 6.5 3.7v7.5L12 19.8 5.5 16V8.5L12 4.8z"/>
          <circle cx="12" cy="12" r="3"/>
        </svg>
      );
    case 'nextjs':
    case 'next':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke={brandColor || '#ffffff'} strokeWidth="1.6" className={className}>
          <circle cx="12" cy="12" r="9.5" />
          <path d="M8.5 15.5V8.5L16 16.5" />
          <path d="M15.5 8.5v4" />
        </svg>
      );
    case 'tailwind':
    case 'tailwindcss':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#06B6D4'} className={className}>
          <path d="M12 6c-3 0-4.5 1.5-4.5 4.5 0 2.2 1.5 3.5 3.5 4.5 2 .9 2.5 1.5 2.5 2.5 0 1.5-1 2.5-2.5 2.5-1.8 0-3-1-3-2.5H6c0 2.5 1.8 4.5 4.5 4.5 3 0 4.5-1.5 4.5-4.5 0-2.2-1.5-3.5-3.5-4.5-2-.9-2.5-1.5-2.5-2.5 0-1.5 1-2.5 2.5-2.5 1.8 0 3 1 3 2.5H16C16 7.5 14.5 6 12 6z"/>
        </svg>
      );
    case 'docker':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#2496ED'} className={className}>
          <path d="M3 13h18c-.5 4-4 7-9 7s-8.5-3-9-7zm2-2h2.5V8.5H5V11zm3.5 0H11V8.5H8.5V11zm3.5 0h2.5V8.5H12V11zm3.5 0H18V8.5h-2.5V11zm-7-3.5H11V5H8.5v2.5zm3.5 0h2.5V5H12v2.5z"/>
        </svg>
      );
    case 'postgresql':
    case 'postgres':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#4169E1'} className={className}>
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 16.93V18h-2v.93c-3.95-.49-7-3.85-7-7.93s3.05-7.44 7-7.93V4h2v-.07c3.95.49 7 3.85 7 7.93s-3.05 7.44-7 7.93zM12 6a6 6 0 1 0 0 12 6 6 0 0 0 0-12z"/>
        </svg>
      );
    case 'mongodb':
    case 'mongo':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#47A248'} className={className}>
          <path d="M12 2C11.5 4 6 9.5 6 14.5a6 6 0 0 0 11.5 2.5c.3-.8.5-1.6.5-2.5C18 9.5 12.5 4 12 2zm0 18.5c-.3 0-.5-.2-.5-.5v-4c0-.3.2-.5.5-.5s.5.2.5.5v4c0 .3-.2.5-.5.5z"/>
        </svg>
      );
    case 'redis':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#DC382D'} className={className}>
          <path d="M2 9.5 12 4.5l10 5-10 5L2 9.5zm0 5 10 5 10-5M2 14.5l10 5 10-5" stroke={brandColor || '#DC382D'} strokeWidth="2" fill="none"/>
        </svg>
      );
    case 'aws':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#FF9900'} className={className}>
          <path d="M12 3a9 9 0 0 0-9 9 9 9 0 0 0 9 9 9 9 0 0 0 9-9 9 9 0 0 0-9-9zm4.5 12.5c-2.5 1.5-6.5 1.5-9 0-.3-.2 0-.6.3-.5 2.3 1.2 6.1 1.2 8.4 0 .4-.2.6.2.3.5z"/>
        </svg>
      );
    case 'git':
    case 'github':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#F05032'} className={className}>
          <path d="M21.6 10.9 13.1 2.4a1.8 1.8 0 0 0-2.5 0L8.8 4.2l3.2 3.2a2.1 2.1 0 0 1 2.7 2.7l3.1 3.1a2.1 2.1 0 1 1-1.3 1.3l-2.9-2.9v4.2a2.1 2.1 0 1 1-1.8 0v-4.4a2.1 2.1 0 0 1-1.1-2.7L7.6 5.5 2.4 10.7a1.8 1.8 0 0 0 0 2.5l8.5 8.5a1.8 1.8 0 0 0 2.5 0l8.2-8.2a1.8 1.8 0 0 0 0-2.6z"/>
        </svg>
      );
    case 'linux':
      return (
        <svg viewBox="0 0 24 24" fill={brandColor || '#FCC624'} className={className}>
          <path d="M12 2c-2.2 0-4 1.8-4 4v2c-2 1-3 3-3 5 0 2.5 1.5 4.5 3 5.5v1.5c0 1.1.9 2 2 2h4c1.1 0 2-.9 2-2V18.5c1.5-1 3-3 3-5.5 0-2-1-4-3-5V6c0-2.2-1.8-4-4-4z"/>
        </svg>
      );
    case 'figma':
      return (
        <svg viewBox="0 0 24 24" className={className}>
          <path d="M8 2h4v4H8z" fill="#F24E1E"/>
          <path d="M12 2h4a4 4 0 0 1 0 8h-4z" fill="#FF7262"/>
          <path d="M8 6h4v4H8a4 4 0 0 1 0-8z" fill="#F24E1E"/>
          <path d="M8 10h4v4H8a4 4 0 0 1 0-8z" fill="#A259FF"/>
          <circle cx="14" cy="12" r="2" fill="#1ABCFE"/>
          <path d="M8 14h4v4a4 4 0 0 1-4-4z" fill="#0ACF83"/>
        </svg>
      );
    default:
      // Generic categorized Lucide icon
      if (normalizedKey.includes('db') || normalizedKey.includes('sql') || normalizedKey.includes('data')) {
        return <Database className={className} style={{ color: brandColor || '#39FF14' }} />;
      }
      if (normalizedKey.includes('cloud') || normalizedKey.includes('devops') || normalizedKey.includes('infra')) {
        return <Cloud className={className} style={{ color: brandColor || '#39FF14' }} />;
      }
      if (normalizedKey.includes('tool') || normalizedKey.includes('terminal')) {
        return <Terminal className={className} style={{ color: brandColor || '#39FF14' }} />;
      }
      if (normalizedKey.includes('web') || normalizedKey.includes('api') || normalizedKey.includes('net')) {
        return <Globe className={className} style={{ color: brandColor || '#39FF14' }} />;
      }
      return <Code2 className={className} style={{ color: brandColor || '#39FF14' }} />;
  }
};

export default TechIcon;
