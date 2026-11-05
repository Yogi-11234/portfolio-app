import React from 'react';
import { Skill, SkillBrief } from '@/types';
import { TechIcon } from './TechIcon';
import { GlowingEffect } from './GlowingEffect';

interface SkillCardProps {
  skill: Skill | SkillBrief;
}

export const SkillCard: React.FC<SkillCardProps> = ({ skill }) => {
  const brandColor = skill.brand_color || '#39FF14';

  return (
    <div
      className="relative neumorphic-card group select-none flex flex-col items-center justify-center p-3 sm:p-4 rounded-2xl w-full h-[105px] sm:h-[115px] transition-all duration-300 hover:-translate-y-1 cursor-pointer overflow-hidden"
      title={skill.name}
    >
      {/* Interactive Glowing Border & Spotlight tracking pointer */}
      <GlowingEffect proximity={70} spread={75} borderWidth={2} halo={true} />

      {/* Content above the glow */}
      <div className="relative z-10 flex flex-col items-center justify-center w-full">
        {/* 1. Neumorphic Logo / Icon Button */}
        <div
          className="w-10 h-10 sm:w-12 sm:h-12 rounded-xl flex items-center justify-center bg-gradient-to-br from-[#1c1c26] to-[#101016] shadow-[4px_4px_10px_rgba(0,0,0,0.5),-3px_-3px_8px_rgba(255,255,255,0.02)] border border-white/5 group-hover:border-accent/40 group-hover:shadow-[inset_2px_2px_5px_rgba(0,0,0,0.7),0_0_12px_rgba(57,255,20,0.3)] transition-all duration-300"
          style={{
            borderColor: `${brandColor}30`,
          }}
        >
          <TechIcon
            name={skill.name}
            iconKey={skill.icon_key || (skill as any).icon_name}
            brandColor={brandColor}
            className="w-5 h-5 sm:w-6 sm:h-6 transition-transform duration-300 group-hover:scale-95 drop-shadow-[0_0_6px_rgba(57,255,20,0.3)]"
          />
        </div>

        {/* 2. Skill Name */}
        <h4 className="mt-2 text-xs sm:text-[13px] font-semibold text-text group-hover:text-accent transition-colors text-center truncate max-w-full px-1 tracking-tight">
          {skill.name}
        </h4>
      </div>
    </div>
  );
};
