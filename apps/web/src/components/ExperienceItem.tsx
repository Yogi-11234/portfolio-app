import React from 'react';
import { Briefcase, GraduationCap, Users, HeartHandshake } from 'lucide-react';
import { Experience, ExperienceType } from '@/types';
import { SanitizedHTML } from './SanitizedHTML';

interface ExperienceItemProps {
  experience: Experience;
}

export const ExperienceItem: React.FC<ExperienceItemProps> = ({ experience }) => {
  const getIcon = (type: ExperienceType) => {
    switch (type) {
      case 'work':
        return <Briefcase className="w-4 h-4 text-sky-400" />;
      case 'education':
        return <GraduationCap className="w-4 h-4 text-amber-400" />;
      case 'organization':
        return <Users className="w-4 h-4 text-emerald-400" />;
      case 'volunteer':
        return <HeartHandshake className="w-4 h-4 text-rose-400" />;
    }
  };

  const formatDate = (iso: string) => {
    try {
      const d = new Date(iso);
      return d.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
    } catch {
      return iso;
    }
  };

  const startedAt = experience.started_at || (experience as any).start_date || '';
  const endedAt = experience.ended_at || (experience as any).end_date || '';
  const isCurrent = experience.current ?? (experience as any).is_current ?? false;
  const roleTitle = experience.title || (experience as any).role || '';

  const dateSpan = `${formatDate(startedAt)} — ${
    isCurrent ? 'Present' : endedAt ? formatDate(endedAt) : '-'
  }`;

  return (
    <div className="relative pl-8 pb-8 last:pb-0 group">
      {/* Timeline line */}
      <div className="absolute left-[11px] top-3 bottom-0 w-[2px] bg-border group-last:hidden" />

      {/* Timeline bullet */}
      <div className="absolute left-0 top-1 w-6 h-6 rounded-full bg-surface-2 border border-border flex items-center justify-center">
        {getIcon(experience.type)}
      </div>

      <div className="rounded-lg bg-surface border border-border p-5 hover:border-border-strong transition-colors">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-1 text-xs text-text-faint font-mono">
          <span>{experience.organization} {experience.location ? `• ${experience.location}` : ''}</span>
          <span className="text-text-muted">{dateSpan}</span>
        </div>

        <h4 className="mt-1 text-sm sm:text-base font-bold text-text">
          {roleTitle}
        </h4>

        {experience.description_html && (
          <SanitizedHTML
            html={experience.description_html}
            className="mt-3 text-xs text-text-muted leading-relaxed prose prose-invert max-w-none"
          />
        )}
      </div>
    </div>
  );
};
