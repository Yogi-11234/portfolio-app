import React, { useState } from 'react';
import { Code2 } from 'lucide-react';
import { Profile } from '@/types';
import { getMediaUrl } from '@/lib/utils';

interface ProfileHeroCardProps {
  profile?: Profile;
}

export const ProfileHeroCard: React.FC<ProfileHeroCardProps> = ({ profile }) => {
  // Toggle reveal on mobile tap
  const [isRevealed, setIsRevealed] = useState(false);

  // Get social media URLs from CMS profile admin
  const githubUrl =
    profile?.github_url ||
    profile?.social_links?.github ||
    'https://github.com';

  const linkedinUrl =
    profile?.linkedin_url ||
    profile?.social_links?.linkedin ||
    'https://linkedin.com';

  const instagramUrl =
    profile?.social_links?.instagram ||
    profile?.social_links?.insta ||
    'https://instagram.com';

  const fullName = profile?.full_name || 'Software Engineer';
  const roleName = profile?.headline || 'Senior Fullstack Engineer';

  return (
    <div className="py-2 px-6 sm:px-0">
      <div
        onClick={() => setIsRevealed(!isRevealed)}
        className={`profile-card-container mx-auto my-3 group select-none ${
          isRevealed ? 'card-revealed' : ''
        }`}
      >
        {/* Front Card (User Profile Photo / Avatar) */}
        <div className="profile-card-front">
          {profile?.avatar_path ? (
            <img
              src={getMediaUrl(profile.avatar_path)}
              alt={fullName}
              className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
          ) : (
            <div className="w-full h-full flex flex-col items-center justify-center bg-surface-2 text-text-muted gap-2">
              <Code2 className="w-14 h-14 sm:w-16 sm:h-16 text-accent animate-pulse" />
              <span className="text-xs font-mono tracking-wider text-accent">PORTFOLIO</span>
            </div>
          )}

          {/* Hover / Tap Hint Overlay at bottom of front card */}
          <div
            className={`absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/85 via-black/40 to-transparent p-2 text-center transition-opacity duration-300 ${
              isRevealed ? 'opacity-0' : 'group-hover:opacity-0'
            }`}
          >
            <p className="text-[10px] sm:text-[11px] font-mono text-accent tracking-wide">
              {isRevealed ? 'Close' : ' Hover'}
            </p>
          </div>
        </div>

        {/* Back Card (Reveal Info & Social Media Icons) */}
        <div className="profile-card-back">
          {/* Name and Role on the left side, aligned to the bottom */}
          <div className="w-[72%] h-full flex flex-col justify-end p-3 text-left border-r border-[#39FF14]/20 bg-gradient-to-t from-black/95 via-[#111116]/80 to-transparent">
            <p className="text-[11px] sm:text-xs font-extrabold text-white leading-tight line-clamp-1 drop-shadow-md">
              {fullName}
            </p>
            <p className="text-[9px] sm:text-[10px] font-mono text-accent mt-0.5 line-clamp-2 leading-tight">
              {roleName}
            </p>
          </div>

          {/* Social Media Column on the right side */}
          <div className="w-[28%] h-full flex flex-col items-center justify-center gap-4 sm:gap-5 bg-black/60 py-3 z-10">
            {/* 1. GitHub Icon */}
            <a
              href={githubUrl}
              target="_blank"
              rel="noopener noreferrer"
              aria-label="GitHub Profile"
              title="Visit GitHub"
              onClick={(e) => e.stopPropagation()}
              className="p-1 text-[#39FF14] hover:text-white hover:scale-125 transition-all duration-300 flex items-center justify-center drop-shadow-[0_0_8px_rgba(57,255,20,0.8)]"
            >
              <svg
                viewBox="0 0 24 24"
                className="w-4 h-4 sm:w-5 sm:h-5 fill-current"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
              </svg>
            </a>

            {/* 2. LinkedIn Icon */}
            <a
              href={linkedinUrl}
              target="_blank"
              rel="noopener noreferrer"
              aria-label="LinkedIn Profile"
              title="Visit LinkedIn"
              onClick={(e) => e.stopPropagation()}
              className="p-1 text-[#39FF14] hover:text-white hover:scale-125 transition-all duration-300 flex items-center justify-center drop-shadow-[0_0_8px_rgba(57,255,20,0.8)]"
            >
              <svg
                viewBox="0 0 24 24"
                className="w-4 h-4 sm:w-5 sm:h-5 fill-current"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451c.979 0 1.778-.773 1.778-1.729V1.73C24 .774 23.205 0 22.225 0z" />
              </svg>
            </a>

            {/* 3. Instagram Icon */}
            <a
              href={instagramUrl}
              target="_blank"
              rel="noopener noreferrer"
              aria-label="Instagram Profile"
              title="Visit Instagram"
              onClick={(e) => e.stopPropagation()}
              className="p-1 text-[#39FF14] hover:text-white hover:scale-125 transition-all duration-300 flex items-center justify-center drop-shadow-[0_0_8px_rgba(57,255,20,0.8)]"
            >
              <svg
                viewBox="0 0 24 24"
                className="w-4 h-4 sm:w-5 sm:h-5 fill-current"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z" />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProfileHeroCard;
