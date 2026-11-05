import React from 'react';
import { SkillBrief, Skill } from '@/types';

interface SkillBadgeProps {
  skill: Skill | SkillBrief;
  showLevel?: boolean;
}

export const SkillBadge: React.FC<SkillBadgeProps> = ({ skill, showLevel = false }) => {
  const brandColor = skill.brand_color || '#3b82f6';

  return (
    <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-md bg-surface-2 border border-border hover:border-border-strong transition-colors text-xs font-mono">
      <span
        className="w-2 h-2 rounded-full flex-shrink-0"
        style={{ backgroundColor: brandColor }}
      />
      <span className="text-text font-medium">{skill.name}</span>
      {showLevel && 'level' in skill && skill.level && (
        <span className="text-text-faint text-[10px] ml-1">
          {'★'.repeat(skill.level)}
        </span>
      )}
    </div>
  );
};
