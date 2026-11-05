import React from 'react';
import { Link } from 'react-router-dom';
import { Github, Linkedin, Twitter, Instagram, Mail, Globe, ArrowUp, MapPin, Code2 } from 'lucide-react';
import { Profile } from '@/types';

interface FooterProps {
  profile?: Profile;
}

const footerMenus = [
  {
    title: 'Navigation',
    links: [
      { label: 'Home', href: '/' },
      { label: 'About Me', href: '/#about' },
      { label: 'Skills', href: '/#skills' },
      { label: 'Experience', href: '/#experience' },
    ],
  },
  {
    title: 'Content',
    links: [
      { label: 'Projects', href: '/projects' },
      { label: 'Certificates', href: '/#certificates' },
      { label: 'Articles & Blog', href: '/posts' },
      { label: 'Contact Me', href: '/#contact' },
    ],
  },
];

export const Footer: React.FC<FooterProps> = ({ profile }) => {
  const currentYear = new Date().getFullYear();

  const scrollToTop = () => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  // Helper function to get social icon
  const getSocialIcon = (key: string) => {
    switch (key.toLowerCase()) {
      case 'github':
        return <Github className="w-4 h-4" />;
      case 'linkedin':
        return <Linkedin className="w-4 h-4" />;
      case 'twitter':
      case 'x':
        return <Twitter className="w-4 h-4" />;
      case 'instagram':
        return <Instagram className="w-4 h-4" />;
      case 'email':
      case 'mail':
        return <Mail className="w-4 h-4" />;
      default:
        return <Globe className="w-4 h-4" />;
    }
  };

  const socialItems: { name: string; url: string; icon: React.ReactNode }[] = [];

  // 1. Get from social_links map (API Response)
  if (profile?.social_links && Object.keys(profile.social_links).length > 0) {
    Object.entries(profile.social_links).forEach(([key, url]) => {
      if (url) {
        socialItems.push({
          name: key,
          url,
          icon: getSocialIcon(key),
        });
      }
    });
  }

  // 2. Direct property fallback
  if (profile?.github_url && !socialItems.some(item => item.name.toLowerCase() === 'github')) {
    socialItems.push({ name: 'GitHub', url: profile.github_url, icon: <Github className="w-4 h-4" /> });
  }
  if (profile?.linkedin_url && !socialItems.some(item => item.name.toLowerCase() === 'linkedin')) {
    socialItems.push({ name: 'LinkedIn', url: profile.linkedin_url, icon: <Linkedin className="w-4 h-4" /> });
  }
  if (profile?.public_email && !socialItems.some(item => item.name.toLowerCase() === 'email')) {
    socialItems.push({ name: 'Email', url: `mailto:${profile.public_email}`, icon: <Mail className="w-4 h-4" /> });
  }

  // 3. Static fallback
  const defaultSocials = [
    { name: 'GitHub', url: 'https://github.com/Yogi-11234', icon: <Github className="w-4 h-4" /> },
    { name: 'LinkedIn', url: 'https://www.linkedin.com/in/yogi-aditia-01a6ab439', icon: <Linkedin className="w-4 h-4" /> },
    { name: 'Email', url: 'mailto:y.aditia1125@gmail.com', icon: <Mail className="w-4 h-4" /> },
  ];

  const finalSocials = socialItems.length > 0 ? socialItems : defaultSocials;

  return (
    <footer className="relative z-10 border-t border-border/80 bg-surface/60 backdrop-blur-md mt-20">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 py-12 sm:py-16">

        {/* Top Section: Brand + Menu Columns + Social */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-10 lg:gap-8">

          {/* Brand & Description */}
          <div className="sm:col-span-2 lg:col-span-1">
            <Link to="/" className="inline-flex items-center gap-2 group">
              <Code2 className="w-5 h-5 text-accent" />
              <span className="text-base font-bold text-text group-hover:text-accent transition-colors">
                {profile?.full_name || 'YOGI ADITIA'}
              </span>
            </Link>
            <p className="mt-3 text-xs text-text-muted leading-relaxed max-w-xs">
              {profile?.headline || 'Fullstack Software Engineer focused on building modern, fast, and scalable web applications.'}
            </p>
            {profile?.location && (
              <div className="mt-3 flex items-center gap-1.5 text-xs text-text-faint">
                <MapPin className="w-3.5 h-3.5 text-accent" />
                <span>{profile.location}</span>
              </div>
            )}
          </div>

          {/* Menu Columns */}
          {footerMenus.map((menu) => (
            <div key={menu.title}>
              <h4 className="text-xs font-bold text-text uppercase tracking-wider mb-4">
                {menu.title}
              </h4>
              <ul className="space-y-2.5">
                {menu.links.map((link) => (
                  <li key={link.href}>
                    {link.href.startsWith('/#') ? (
                      <a
                        href={link.href}
                        className="text-xs text-text-muted hover:text-accent transition-colors duration-200"
                      >
                        {link.label}
                      </a>
                    ) : (
                      <Link
                        to={link.href}
                        className="text-xs text-text-muted hover:text-accent transition-colors duration-200"
                      >
                        {link.label}
                      </Link>
                    )}
                  </li>
                ))}
              </ul>
            </div>
          ))}

          {/* Social Media */}
          <div>
            <h4 className="text-xs font-bold text-text uppercase tracking-wider mb-4">
              Social Media
            </h4>
            <div className="flex flex-wrap gap-2.5">
              {finalSocials.map((item) => (
                <a
                  key={item.name}
                  href={item.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  aria-label={item.name}
                  title={item.name}
                  className="w-9 h-9 rounded-xl bg-gradient-to-br from-[#1a1a22] to-[#111116] border border-white/5 flex items-center justify-center text-text-muted hover:text-[#39FF14] hover:border-[#39FF14]/50 hover:shadow-[0_0_14px_rgba(57,255,20,0.3)] transition-all duration-300"
                >
                  {item.icon}
                </a>
              ))}
            </div>

            {profile?.public_email && (
              <a
                href={`mailto:${profile.public_email}`}
                className="mt-4 inline-flex items-center gap-1.5 text-xs text-text-muted hover:text-accent transition-colors"
              >
                <Mail className="w-3.5 h-3.5" />
                <span>{profile.public_email}</span>
              </a>
            )}
          </div>
        </div>

        {/* Bottom Bar */}
        <div className="mt-10 pt-6 border-t border-border/50 flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <p className="text-xs text-text-faint font-mono">
              © {currentYear} {profile?.full_name || 'YOGI ADITIA'}. All rights reserved.
            </p>
            {profile?.open_to_work && (
              <div className="flex items-center gap-1.5 text-emerald-400 text-xs">
                <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
                <span className="font-mono">Open to Work</span>
              </div>
            )}
          </div>

          <button
            onClick={scrollToTop}
            aria-label="Back to top"
            className="w-9 h-9 rounded-xl bg-gradient-to-br from-[#1a1a22] to-[#111116] border border-white/5 flex items-center justify-center text-text-muted hover:text-[#39FF14] hover:border-[#39FF14]/50 hover:shadow-[0_0_14px_rgba(57,255,20,0.3)] transition-all duration-300"
          >
            <ArrowUp className="w-4 h-4" />
          </button>
        </div>
      </div>
    </footer>
  );
};