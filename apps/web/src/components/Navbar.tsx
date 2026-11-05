import React, { useState, useEffect, useCallback } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Menu, X, Code2 } from 'lucide-react';
import { Profile } from '@/types';
import { AnimatedButton } from './AnimatedButton';

interface NavbarProps {
  profile?: Profile;
}

// Section IDs that live on the homepage (mapped to their hash)
const HOME_SECTIONS = ['about', 'skills', 'projects', 'experience', 'certificates', 'posts', 'contact'];

function scrollToSection(sectionId: string) {
  const el = document.getElementById(sectionId);
  if (!el) return;
  const navHeight = 72; // approximate navbar height
  const y = el.getBoundingClientRect().top + window.scrollY - navHeight;
  window.scrollTo({ top: y, behavior: 'smooth' });
}

export const Navbar: React.FC<NavbarProps> = ({ profile }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const [activeSection, setActiveSection] = useState<string>('');
  const location = useLocation();
  const navigate = useNavigate();

  const isHomePage = location.pathname === '/';

  // Track scroll to make navbar opaque
  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 20);
    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  // Intersection Observer: track which section is visible
  useEffect(() => {
    if (!isHomePage) {
      setActiveSection('');
      return;
    }

    const observers: IntersectionObserver[] = [];

    const sectionVisibility: Record<string, number> = {};

    HOME_SECTIONS.forEach((id) => {
      const el = document.getElementById(id);
      if (!el) return;

      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            sectionVisibility[id] = entry.intersectionRatio;
            // Pick the section with the highest visibility ratio
            const best = Object.entries(sectionVisibility).sort((a, b) => b[1] - a[1])[0];
            if (best && best[1] > 0) {
              setActiveSection(best[0]);
            }
          });
        },
        {
          rootMargin: '-72px 0px -40% 0px', // account for navbar height at top
          threshold: [0, 0.1, 0.25, 0.5, 0.75, 1.0],
        }
      );

      observer.observe(el);
      observers.push(observer);
    });

    return () => observers.forEach((o) => o.disconnect());
  }, [isHomePage, location.pathname]);

  // After navigating to homepage via hash, scroll to target section
  useEffect(() => {
    if (isHomePage && location.hash) {
      const sectionId = location.hash.replace('#', '');
      // Wait for page render then scroll
      const timer = setTimeout(() => {
        scrollToSection(sectionId);
      }, 120);
      return () => clearTimeout(timer);
    }
  }, [isHomePage, location.hash]);

  const handleNavClick = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement | HTMLButtonElement>, sectionId: string | null, path?: string) => {
      e.preventDefault();
      setIsOpen(false);

      // Pure page link (/projects, /posts)
      if (!sectionId && path) {
        navigate(path);
        return;
      }

      if (!sectionId) return;

      if (isHomePage) {
        // Already on homepage → smooth scroll
        scrollToSection(sectionId);
        // Update URL hash without triggering navigation
        window.history.replaceState(null, '', `/#${sectionId}`);
      } else {
        // Navigate to homepage with hash; useEffect above will scroll
        navigate(`/#${sectionId}`);
      }
    },
    [isHomePage, navigate]
  );

  const navLinks: { label: string; section?: string; path?: string }[] = [
    { label: 'About',        section: 'about' },
    { label: 'Skills',       section: 'skills' },
    { label: 'Projects',     path: '/projects' },
    { label: 'Experience',   section: 'experience' },
    { label: 'Certificates', section: 'certificates' },
    { label: 'Articles',     path: '/posts' },
  ];

  const checkIsActive = (link: typeof navLinks[number]): boolean => {
    if (link.path) return location.pathname === link.path;
    if (link.section) return isHomePage && activeSection === link.section;
    return false;
  };

  const linkClass = (isActive: boolean) =>
    `px-3.5 py-1.5 text-xs font-mono rounded-xl transition-all duration-300 cursor-pointer select-none border ${
      isActive
        ? 'bg-[#39FF14]/15 text-[#39FF14] font-bold border-[#39FF14]/60 shadow-[0_0_15px_rgba(57,255,20,0.3)] scale-[1.02]'
        : 'text-text-muted hover:text-[#39FF14] hover:bg-[#39FF14]/10 border-transparent'
    }`;

  const mobileLinkClass = (isActive: boolean) =>
    `block px-4 py-2.5 text-xs font-mono rounded-xl transition-all duration-300 border ${
      isActive
        ? 'bg-[#39FF14]/15 text-[#39FF14] font-bold border-[#39FF14]/60 shadow-[0_0_15px_rgba(57,255,20,0.3)]'
        : 'text-text-muted hover:text-[#39FF14] hover:bg-[#39FF14]/10 border-transparent'
    }`;

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        scrolled
          ? 'bg-[#121212]/80 backdrop-blur-md py-3 border-b border-white/10 shadow-2xl'
          : 'bg-transparent py-5'
      }`}
    >
      <div className="max-w-6xl mx-auto px-4 sm:px-6 flex items-center justify-between">
        {/* Brand / Monogram — always goes to top of homepage */}
        <a
          href="/"
          onClick={(e) => {
            e.preventDefault();
            setIsOpen(false);
            if (isHomePage) {
              window.scrollTo({ top: 0, behavior: 'smooth' });
              window.history.replaceState(null, '', '/');
            } else {
              navigate('/');
            }
          }}
          className="flex items-center gap-2 group text-text font-bold text-lg tracking-tight"
        >
          <div className="w-9 h-9 rounded-xl bg-surface-2/80 border border-border flex items-center justify-center text-[#39FF14] group-hover:border-[#39FF14]/80 group-hover:shadow-[0_0_12px_rgba(57,255,20,0.4)] transition-all">
            <Code2 className="w-5 h-5 text-[#39FF14]" />
          </div>
          <span className="font-mono">
            {profile?.full_name ? profile.full_name.split(' ')[0] : 'Portfolio'}
            <span className="text-[#39FF14] font-bold">/ADITIA</span>
          </span>
        </a>

        {/* Desktop Nav */}
        <nav className="hidden md:flex items-center gap-1.5 p-1.5 rounded-2xl bg-surface-2/60 border border-border/80 backdrop-blur-md shadow-inner">
          {navLinks.map((link) => {
            const isActive = checkIsActive(link);
            return (
              <a
                key={link.label}
                href={link.path ?? `/#${link.section}`}
                onClick={(e) => handleNavClick(e, link.section ?? null, link.path)}
                className={linkClass(isActive)}
              >
                {link.label}
              </a>
            );
          })}
        </nav>

        <div className="hidden md:flex items-center gap-3">
          <AnimatedButton
            as="a"
            size="sm"
            href="/#contact"
            onClick={(e: React.MouseEvent<HTMLElement>) =>
              handleNavClick(e as React.MouseEvent<HTMLAnchorElement>, 'contact')
            }
          >
            <span>Contact</span>
          </AnimatedButton>
        </div>

        {/* Mobile Hamburger */}
        <button
          onClick={() => setIsOpen(!isOpen)}
          aria-label="Toggle Navigation"
          className="md:hidden p-2 text-text-muted hover:text-[#39FF14] rounded-xl bg-surface-2/80 border border-border transition-all"
        >
          {isOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
        </button>
      </div>

      {/* Mobile Drawer */}
      {isOpen && (
        <div className="md:hidden bg-[#121212]/95 backdrop-blur-xl border-b border-border/80 mt-2 px-4 py-4 space-y-2 shadow-2xl">
          {navLinks.map((link) => {
            const isActive = checkIsActive(link);
            return (
              <a
                key={link.label}
                href={link.path ?? `/#${link.section}`}
                onClick={(e) => handleNavClick(e, link.section ?? null, link.path)}
                className={mobileLinkClass(isActive)}
              >
                {link.label}
              </a>
            );
          })}

          <div className="pt-2">
            <AnimatedButton
              as="a"
              fullWidth
              size="sm"
              href="/#contact"
              onClick={(e: React.MouseEvent<HTMLElement>) =>
                handleNavClick(e as React.MouseEvent<HTMLAnchorElement>, 'contact')
              }
            >
              <span>Contact</span>
            </AnimatedButton>
          </div>
        </div>
      )}
    </header>
  );
};
