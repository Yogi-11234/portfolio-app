import React from 'react';
import { Link } from 'react-router-dom';
import { ExternalLink, Github, ArrowRight, Layers } from 'lucide-react';
import { Project } from '@/types';
import { SkillBadge } from './SkillBadge';
import { GlowingEffect } from './GlowingEffect';
import { getMediaUrl } from '@/lib/utils';

interface ProjectCardProps {
  project: Project;
}

export const ProjectCard: React.FC<ProjectCardProps> = ({ project }) => {
  return (
    <div className="relative group neumorphic-card flex flex-col justify-between overflow-hidden p-5 transition-all duration-300 w-full max-w-[360px] mx-auto">
      {/* Interactive Glowing Border & Spotlight tracking pointer */}
      <GlowingEffect proximity={80} spread={75} borderWidth={2} halo={true} />

      <div className="relative z-10 flex-1 flex flex-col justify-between">
        <div>
          {/* Thumbnail Frame */}
          <div className="relative aspect-video rounded-2xl bg-surface-2/60 border border-white/5 overflow-hidden shadow-inner mb-4">
            {project.thumbnail_path ? (
              <img
                src={getMediaUrl(project.thumbnail_path)}
                alt={project.title}
                className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500 ease-out"
                loading="lazy"
              />
            ) : (
              <div className="w-full h-full flex flex-col items-center justify-center gap-2 text-text-faint font-mono text-xs">
                <Layers className="w-6 h-6 opacity-40 text-accent" />
                <span>No preview available</span>
              </div>
            )}

            {project.featured && (
              <span className="absolute top-2.5 left-2.5 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-accent text-accent-fg uppercase tracking-wider shadow-[0_0_12px_rgba(57,255,20,0.6)]">
                Featured
              </span>
            )}

            {project.category && (
              <span className="absolute top-2.5 right-2.5 px-2.5 py-0.5 rounded-full text-[10px] font-mono bg-black/70 border border-white/10 text-white/80 backdrop-blur-md">
                {project.category.name}
              </span>
            )}
          </div>

          {/* Content */}
          <div>
            <h3 className="text-base font-bold text-text group-hover:text-accent transition-colors line-clamp-1">
              <Link to={`/projects/${project.slug}`}>{project.title}</Link>
            </h3>

            <p className="mt-2 text-xs text-text-muted line-clamp-2 leading-relaxed">
              {project.summary}
            </p>

            {/* Skills Badges */}
            {project.skills && project.skills.length > 0 && (
              <div className="mt-4 flex flex-wrap gap-1.5">
                {project.skills.slice(0, 3).map((s) => (
                  <SkillBadge key={s.id} skill={s} />
                ))}
                {project.skills.length > 3 && (
                  <span className="px-2 py-0.5 rounded-full bg-surface-2 text-[10px] text-text-faint font-mono border border-white/5">
                    +{project.skills.length - 3}
                  </span>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Card Footer Actions */}
        <div className="mt-5 pt-4 border-t border-white/5 flex items-center justify-between text-xs">
          <Link
            to={`/projects/${project.slug}`}
            className="neumorphic-mini-btn"
          >
            <span>Details</span>
            <ArrowRight className="w-3.5 h-3.5" />
          </Link>

          <div className="flex items-center gap-2">
            {project.repo_url && (
              <a
                href={project.repo_url}
                target="_blank"
                rel="noopener noreferrer"
                aria-label="GitHub Repository"
                className="neumorphic-mini-btn p-2"
                title="View Source Code on GitHub"
              >
                <Github className="w-3.5 h-3.5" />
              </a>
            )}
            {project.demo_url && (
              <a
                href={project.demo_url}
                target="_blank"
                rel="noopener noreferrer"
                aria-label="Live Demo"
                className="neumorphic-mini-btn p-2"
                title="Visit Live Demo"
              >
                <ExternalLink className="w-3.5 h-3.5" />
              </a>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
