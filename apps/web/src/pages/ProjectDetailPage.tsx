import React from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, ExternalLink, Github, Calendar, Tag } from 'lucide-react';
import { api } from '@/lib/api';
import { SanitizedHTML } from '@/components/SanitizedHTML';
import { SkillBadge } from '@/components/SkillBadge';
import { AnimatedButton } from '@/components/AnimatedButton';
import { getMediaUrl } from '@/lib/utils';
import { WaveMesh } from '@/components/WaveMesh';

export const ProjectDetailPage: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();

  const { data: project, isLoading, error } = useQuery({
    queryKey: ['project', slug],
    queryFn: () => api.getProjectBySlug(slug!),
    enabled: Boolean(slug),
  });

  const formatDate = (iso?: string) => {
    if (!iso) return '-';
    try {
      const d = new Date(iso);
      return d.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
    } catch {
      return iso;
    }
  };

  if (isLoading) {
    return (
      <div className="relative w-full min-h-screen">
        <div className="fixed inset-0 z-0 pointer-events-auto w-full h-full">
          <WaveMesh
            colors={["#06070a", "#39FF14", "#00ff88"]}
            crests={true}
            className="w-full h-full"
          />
        </div>
        <div className="relative z-10 max-w-4xl mx-auto px-4 pt-32 pb-20 text-center font-mono text-xs text-text-faint">
          Loading project details...
        </div>
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="relative w-full min-h-screen">
        <div className="fixed inset-0 z-0 pointer-events-auto w-full h-full">
          <WaveMesh
            colors={["#06070a", "#39FF14", "#00ff88"]}
            crests={true}
            className="w-full h-full"
          />
        </div>
        <div className="relative z-10 max-w-4xl mx-auto px-4 pt-32 pb-20 text-center">
          <h2 className="text-xl font-bold text-text">Project Not Found</h2>
          <p className="mt-2 text-xs text-text-muted">
            The project you are looking for might not be published yet or has been removed.
          </p>
          <div className="mt-6 flex justify-center">
            <AnimatedButton as="link" to="/projects">
              <span>Back to Projects</span>
            </AnimatedButton>
          </div>
        </div>
      </div>
    );
  }

  const dateSpan = `${formatDate(project.started_at)} — ${
    project.finished_at ? formatDate(project.finished_at) : 'Completed'
  }`;

  return (
    <div className="relative w-full min-h-screen">
      {/* 🌊 FULLSCREEN FIXED WAVE MESH BACKGROUND */}
      <div className="fixed inset-0 z-0 pointer-events-auto w-full h-full">
        <WaveMesh
          colors={["#06070a", "#39FF14", "#00ff88"]}
          crests={true}
          className="w-full h-full"
        />
      </div>

      {/* 🚀 PROJECT CONTENT (Layer z-10 di atas background) */}
      <article className="relative z-10 max-w-4xl mx-auto px-4 sm:px-6 pt-24 pb-20">
        {/* Back button */}
        <Link
          to="/projects"
          className="inline-flex items-center gap-1.5 text-xs text-text-muted hover:text-[#39FF14] font-mono mb-6 transition-colors"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          <span>Back to All Projects</span>
        </Link>

        {/* Header */}
        <div className="space-y-3">
          <div className="flex flex-wrap items-center gap-2 text-xs text-text-faint font-mono">
            {project.category && (
              <span className="inline-flex items-center gap-1 text-accent">
                <Tag className="w-3 h-3" />
                <span>{project.category.name}</span>
              </span>
            )}
            <span>•</span>
            <span className="inline-flex items-center gap-1">
              <Calendar className="w-3 h-3" />
              <span>{dateSpan}</span>
            </span>
          </div>

          <h1 className="text-2xl sm:text-4xl font-extrabold text-text tracking-tight">
            {project.title}
          </h1>

          <p className="text-sm sm:text-base text-text-muted leading-relaxed">
            {project.summary}
          </p>
        </div>

        {/* Action Links */}
        <div className="mt-6 flex flex-wrap items-center gap-4">
          {project.demo_url && (
            <AnimatedButton
              as="a"
              href={project.demo_url}
              target="_blank"
              rel="noopener noreferrer"
            >
              <ExternalLink className="w-4 h-4 inline mr-1" />
              <span>Visit Live Demo</span>
            </AnimatedButton>
          )}
          {project.repo_url && (
            <AnimatedButton
              as="a"
              href={project.repo_url}
              target="_blank"
              rel="noopener noreferrer"
            >
              <Github className="w-4 h-4 inline mr-1" />
              <span>View Source Code</span>
            </AnimatedButton>
          )}
        </div>

        {/* Hero Thumbnail */}
        {project.thumbnail_path && (
          <div className="mt-8 rounded-lg overflow-hidden border border-border bg-surface-2 shadow-2xl">
            <img
              src={getMediaUrl(project.thumbnail_path)}
              alt={project.title}
              className="w-full h-auto object-contain block"
            />
          </div>
        )}

        {/* Skills & Tech stack used */}
        {project.skills && project.skills.length > 0 && (
          <div className="mt-8 p-5 rounded-lg bg-surface-2/60 backdrop-blur-md border border-border">
            <h3 className="text-xs font-mono uppercase tracking-wider text-text-faint mb-3">
              Technologies & Tools Used
            </h3>
            <div className="flex flex-wrap gap-2">
              {project.skills.map((s) => (
                <SkillBadge key={s.id} skill={s} showLevel />
              ))}
            </div>
          </div>
        )}

        {/* Markdown Description */}
        <div className="mt-10 rounded-lg bg-surface-2/60 backdrop-blur-md border border-border p-6 sm:p-8">
          <h2 className="text-lg font-bold text-text mb-4 pb-2 border-b border-border">
            Overview & Technical Details
          </h2>
          <SanitizedHTML
            html={project.description_html}
            className="text-sm text-text-muted leading-relaxed prose prose-invert max-w-none space-y-4"
          />
        </div>
      </article>
    </div>
  );
};