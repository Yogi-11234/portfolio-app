import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { motion } from 'framer-motion';
import {
  Download,
  Sparkles,
  Code2,
  Layers,
  GraduationCap,
  Award,
  BookOpen,
  MessageSquare,
} from 'lucide-react';
import { api } from '@/lib/api';
import { SkillGroup } from '@/types';
import { SanitizedHTML } from '@/components/SanitizedHTML';
import { SkillCard } from '@/components/SkillCard';
import { ProjectCard } from '@/components/ProjectCard';
import { ExperienceItem } from '@/components/ExperienceItem';
import { CertificateCarousel } from '@/components/CertificateCarousel';
import { PostCard } from '@/components/PostCard';
import { ContactForm } from '@/components/ContactForm';
import { AnimatedButton } from '@/components/AnimatedButton';
import { ProfileHeroCard } from '@/components/ProfileHeroCard';
import { TextHoverEffect } from '@/components/TextHoverEffect';
import { CenteredCardContainer } from '@/components/CenteredCardContainer';
import { GlowingEffect } from '@/components/GlowingEffect';
import { WaveMesh } from '@/components/WaveMesh';
import ScrollVelocity from '@/components/ui/ScrollVelocity';
import { getMediaUrl } from '@/lib/utils';

const SKILL_GROUPS: { key: SkillGroup | 'all'; label: string }[] = [
  { key: 'all', label: 'All' },
  { key: 'language', label: 'Languages' },
  { key: 'framework', label: 'Frameworks' },
  { key: 'database', label: 'Databases' },
  { key: 'cloud', label: 'Cloud & DevOps' },
  { key: 'tool', label: 'Tools' },
];

export const HomePage: React.FC = () => {
  const [selectedSkillGroup, setSelectedSkillGroup] = useState<SkillGroup | 'all'>('all');
  const [experienceTab, setExperienceTab] = useState<'all' | 'work' | 'education'>('all');

  // Queries
  const { data: profile } = useQuery({
    queryKey: ['profile'],
    queryFn: api.getProfile,
  });

  const { data: skills = [] } = useQuery({
    queryKey: ['skills'],
    queryFn: api.getSkills,
  });

  const { data: projectsData } = useQuery({
    queryKey: ['featured-projects'],
    queryFn: () => api.getProjects({ featured: true }),
  });

  const { data: experiences = [] } = useQuery({
    queryKey: ['experiences'],
    queryFn: api.getExperiences,
  });

  const { data: certificates = [] } = useQuery({
    queryKey: ['certificates'],
    queryFn: api.getCertificates,
  });

  const { data: postsData } = useQuery({
    queryKey: ['latest-posts'],
    queryFn: () => api.getPosts({ page: 1 }),
  });

  const filteredSkills = skills.filter((s) =>
    selectedSkillGroup === 'all' ? true : s.group_name === selectedSkillGroup
  );

  const filteredExperiences = experiences.filter((e) =>
    experienceTab === 'all' ? true : e.type === experienceTab
  );

  const featuredProjects = projectsData?.projects || [];
  const latestPosts = postsData?.posts?.slice(0, 3) || [];

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

      {/* 🚀 MAIN CONTENT */}
      <div className="relative z-10 w-full space-y-16 sm:space-y-24">
        {/* 1. HERO SECTION */}
        <motion.section
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.9, ease: [0.16, 1, 0.3, 1] }}
          className="pt-20 pb-10 sm:pt-32 sm:pb-16 max-w-4xl mx-auto text-center px-4"
        >
          {/* Open to Work Badge */}
          {profile?.open_to_work && (
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-950/60 border border-emerald-800/80 text-emerald-400 text-xs font-mono mb-6 shadow-[0_0_12px_rgba(57,255,20,0.25)]">
              <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
              <span>Open to Work / Available for New Projects</span>
            </div>
          )}

          {/* Interactive 3D Reveal Profile Card */}
          <div className="flex justify-center mb-6">
            <ProfileHeroCard profile={profile} />
          </div>

          <div className="max-w-3xl mx-auto">
            <h1 className="text-sm sm:text-lg font-mono text-text-muted tracking-wide">
              Hi, I'm
            </h1>
            <div className="w-full max-w-lg sm:max-w-xl mx-auto my-1 px-1 sm:px-2">
              <TextHoverEffect
                text={(profile?.full_name || 'YOGI ADITIA').toUpperCase()}
                className="font-extrabold tracking-tight"
              />
            </div>
          </div>

          {/* Hero CTAs */}
          <div className="mt-8 flex flex-col sm:flex-row flex-wrap items-center justify-center gap-3 w-full sm:w-auto max-w-xs sm:max-w-none mx-auto">
            {profile?.resume_path && (  
              <AnimatedButton
                as="a"
                href={getMediaUrl(profile.resume_path)}
                target="_blank"
                rel="noopener noreferrer"
                className="w-full sm:w-auto"
              >
                <Download className="w-4 h-4 inline mr-1" />
                <span>Download CV (PDF)</span>
              </AnimatedButton>
            )}
          </div>
        </motion.section>

        {/* ⚡ SCROLL VELOCITY SECTION */}
        <motion.section
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          className="py-8 my-10 overflow-hidden bg-transparent"
        >
          <ScrollVelocity
            texts={[
              'ABOUT ME |• YOGI ADITIA |•',
              'AI & Data Enthusiast |• Backend Enthusiast |• Software Development |• Machine Learning & Computer Vision |• Generative AI |•'
            ]}
            velocity={70}
            skewWithVelocity={6}
            scaleWithVelocity={0.05}
            color="#39FF14"
            fontSize="clamp(3rem, 6vw, 7rem)"
            className="font-mono font-black uppercase tracking-widest"
            fadeEdges={true}
          />
        </motion.section>

        {/* 2. ABOUT / BIO SECTION */}
        {profile?.bio_html && (
          <motion.section
            id="about"
            initial={{ opacity: 0, y: 35 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: '-60px' }}
            transition={{ duration: 0.85, ease: [0.16, 1, 0.3, 1] }}
            className="max-w-4xl mx-auto px-4 scroll-mt-24"
          >
            <div className="relative neumorphic-card p-6 sm:p-8 overflow-hidden">
              <GlowingEffect proximity={60} spread={70} borderWidth={2} halo={true} />
              <div className="flex items-center gap-2 text-xs font-mono text-accent uppercase tracking-wider mb-3">
                <Sparkles className="w-4 h-4" />
                <span>About Me</span>
              </div>
              <SanitizedHTML
                html={profile.bio_html}
                className="text-sm sm:text-base text-text-muted leading-relaxed space-y-4 prose prose-invert max-w-none"
              />
            </div>
          </motion.section>
        )}

        {/* 3. SKILLS SECTION */}
        <motion.section
          id="skills"
          initial={{ opacity: 0, y: 35 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-60px' }}
          transition={{ duration: 0.85, ease: [0.16, 1, 0.3, 1] }}
          className="max-w-5xl mx-auto px-4 scroll-mt-24"
        >
          <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4 mb-6">
            <div>
              <div className="flex items-center gap-2 text-xs font-mono text-accent uppercase tracking-wider">
                <Code2 className="w-4 h-4" />
                <span>Technical Skills</span>
              </div>
              <h2 className="text-xl sm:text-2xl font-bold text-text mt-1">Tech Stack & Tools</h2>
            </div>

            {/* Group Filter Tabs */}
            <div className="flex flex-wrap gap-2 p-1.5 rounded-2xl bg-surface-2/60 backdrop-blur-md border border-border/80 shadow-inner">
              {SKILL_GROUPS.map((g) => {
                const isActive = selectedSkillGroup === g.key;
                return (
                  <button
                    key={g.key}
                    onClick={() => setSelectedSkillGroup(g.key)}
                    className={`px-3.5 py-1.5 text-xs font-mono rounded-xl transition-all duration-300 relative ${
                      isActive
                        ? 'bg-[#39FF14]/15 text-[#39FF14] font-bold border border-[#39FF14]/60 shadow-[0_0_15px_rgba(57,255,20,0.3)] scale-[1.02]'
                        : 'text-text-muted hover:text-text hover:bg-surface-2 border border-transparent'
                    }`}
                  >
                    {g.label}
                  </button>
                );
              })}
            </div>
          </div>

          <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6 gap-2.5 sm:gap-3.5 max-w-4xl mx-auto">
            {filteredSkills.map((s, index) => (
              <motion.div
                key={s.id}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, margin: '-30px' }}
                transition={{
                  duration: 0.65,
                  delay: Math.min(index * 0.035, 0.35),
                  ease: [0.16, 1, 0.3, 1],
                }}
                className="flex justify-center"
              >
                <SkillCard skill={s} />
              </motion.div>
            ))}
            {filteredSkills.length === 0 && (
              <div className="col-span-full py-12 text-center text-xs text-text-faint font-mono rounded-2xl bg-surface border border-border">
                No skills found in this category.
              </div>
            )}
          </div>
        </motion.section>

        {/* 4. FEATURED PROJECTS SECTION */}
        <motion.section
          id="projects"
          initial={{ opacity: 0, y: 35 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-60px' }}
          transition={{ duration: 0.85, ease: [0.16, 1, 0.3, 1] }}
          className="max-w-6xl mx-auto px-4 scroll-mt-24"
        >
          <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4 mb-6">
            <div>
              <div className="flex items-center gap-2 text-xs font-mono text-accent uppercase tracking-wider">
                <Layers className="w-4 h-4" />
                <span>Featured Works</span>
              </div>
              <h2 className="text-xl sm:text-2xl font-bold text-text mt-1">Featured Projects</h2>
            </div>

            <div className="self-start sm:self-auto">
              <AnimatedButton as="link" to="/projects" size="sm">
                <span>All Projects</span>
              </AnimatedButton>
            </div>
          </div>

          <CenteredCardContainer
            items={featuredProjects}
            emptyMessage="No featured projects available."
            renderItem={(proj) => <ProjectCard project={proj} />}
          />
        </motion.section>

        {/* 5. EXPERIENCE TIMELINE SECTION */}
        <motion.section
          id="experience"
          initial={{ opacity: 0, y: 35 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-60px' }}
          transition={{ duration: 0.85, ease: [0.16, 1, 0.3, 1] }}
          className="max-w-4xl mx-auto px-4 scroll-mt-24"
        >
          <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4 mb-6">
            <div>
              <div className="flex items-center gap-2 text-xs font-mono text-accent uppercase tracking-wider">
                <GraduationCap className="w-4 h-4" />
                <span>Career Journey</span>
              </div>
              <h2 className="text-xl sm:text-2xl font-bold text-text mt-1">Experience & Education</h2>
            </div>

            <div className="flex gap-2 p-1.5 rounded-2xl bg-surface-2/60 backdrop-blur-md border border-border/80 shadow-inner text-xs">
              <button
                onClick={() => setExperienceTab('all')}
                className={`px-3.5 py-1.5 font-mono rounded-xl transition-all duration-300 relative ${
                  experienceTab === 'all'
                    ? 'bg-[#39FF14]/15 text-[#39FF14] font-bold border border-[#39FF14]/60 shadow-[0_0_15px_rgba(57,255,20,0.3)] scale-[1.02]'
                    : 'text-text-muted hover:text-text hover:bg-surface-2 border border-transparent'
                }`}
              >
                All
              </button>
              <button
                onClick={() => setExperienceTab('work')}
                className={`px-3.5 py-1.5 font-mono rounded-xl transition-all duration-300 relative ${
                  experienceTab === 'work'
                    ? 'bg-[#39FF14]/15 text-[#39FF14] font-bold border border-[#39FF14]/60 shadow-[0_0_15px_rgba(57,255,20,0.3)] scale-[1.02]'
                    : 'text-text-muted hover:text-text hover:bg-surface-2 border border-transparent'
                }`}
              >
                Work
              </button>
              <button
                onClick={() => setExperienceTab('education')}
                className={`px-3.5 py-1.5 font-mono rounded-xl transition-all duration-300 relative ${
                  experienceTab === 'education'
                    ? 'bg-[#39FF14]/15 text-[#39FF14] font-bold border border-[#39FF14]/60 shadow-[0_0_15px_rgba(57,255,20,0.3)] scale-[1.02]'
                    : 'text-text-muted hover:text-text hover:bg-surface-2 border border-transparent'
                }`}
              >
                Education
              </button>
            </div>
          </div>

          <div className="mt-8 space-y-4">
            {filteredExperiences.map((exp, index) => (
              <motion.div
                key={exp.id}
                initial={{ opacity: 0, y: 25 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, margin: '-40px' }}
                transition={{
                  duration: 0.75,
                  delay: Math.min(index * 0.08, 0.35),
                  ease: [0.16, 1, 0.3, 1],
                }}
              >
                <ExperienceItem experience={exp} />
              </motion.div>
            ))}
            {filteredExperiences.length === 0 && (
              <p className="text-xs text-text-faint font-mono text-center py-6">No experience records found.</p>
            )}
          </div>
        </motion.section>

        {/* 6. CERTIFICATES SECTION */}
        <motion.section
          id="certificates"
          initial={{ opacity: 0, y: 35 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-60px' }}
          transition={{ duration: 0.85, ease: [0.16, 1, 0.3, 1] }}
          className="max-w-6xl mx-auto px-4 scroll-mt-24"
        >
          <div className="mb-6">
            <div className="flex items-center gap-2 text-xs font-mono text-accent uppercase tracking-wider">
              <Award className="w-4 h-4" />
              <span>Professional Credentials</span>
            </div>
            <h2 className="text-xl sm:text-2xl font-bold text-text mt-1">Licenses & Certifications</h2>
          </div>

          <CertificateCarousel certificates={certificates} />
        </motion.section>

        {/* 7. LATEST POSTS SECTION */}
        <motion.section
          id="posts"
          initial={{ opacity: 0, y: 35 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-60px' }}
          transition={{ duration: 0.85, ease: [0.16, 1, 0.3, 1] }}
          className="max-w-6xl mx-auto px-4 scroll-mt-24"
        >
          <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4 mb-6">
            <div>
              <div className="flex items-center gap-2 text-xs font-mono text-accent uppercase tracking-wider">
                <BookOpen className="w-4 h-4" />
                <span>Articles & Writings</span>
              </div>
              <h2 className="text-xl sm:text-2xl font-bold text-text mt-1">Latest Posts</h2>
            </div>

            <div className="self-start sm:self-auto">
              <AnimatedButton as="link" to="/posts" size="sm">
                <span>All Articles</span>
              </AnimatedButton>
            </div>
          </div>

          <div className="flex flex-wrap justify-center gap-6 items-stretch">
            {latestPosts.map((post, index) => (
              <motion.div
                key={post.id}
                initial={{ opacity: 0, y: 30 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, margin: '-50px' }}
                transition={{
                  duration: 0.8,
                  delay: Math.min(index * 0.1, 0.35),
                  ease: [0.16, 1, 0.3, 1],
                }}
                className="flex justify-center w-full max-w-[360px]"
              >
                <PostCard post={post} />
              </motion.div>
            ))}
            {latestPosts.length === 0 && (
              <div className="w-full max-w-lg py-12 text-center text-xs text-text-faint font-mono rounded-2xl bg-surface border border-border">
                No articles published yet.
              </div>
            )}
          </div>
        </motion.section>

        {/* 8. CONTACT FORM SECTION */}
        <motion.section
          id="contact"
          initial={{ opacity: 0, y: 35 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-60px' }}
          transition={{ duration: 0.85, ease: [0.16, 1, 0.3, 1] }}
          className="max-w-4xl mx-auto px-4 scroll-mt-24 pb-12"
        >
          <div className="relative neumorphic-card p-6 sm:p-10 overflow-hidden">
            <GlowingEffect proximity={60} spread={70} borderWidth={2} halo={true} />
            <div className="max-w-xl mx-auto text-center mb-8">
              <div className="inline-flex items-center gap-2 text-xs font-mono text-accent uppercase tracking-wider mb-2">
                <MessageSquare className="w-4 h-4" />
                <span>Get In Touch</span>
              </div>
              <h2 className="text-2xl sm:text-3xl font-bold text-text">Let's Connect & Collaborate</h2>
              <p className="mt-2 text-xs sm:text-sm text-text-muted">
                Have an interesting project idea, job opportunity, or just want to discuss tech? 
                Feel free to send a message through the form below.
              </p>
            </div>

            <div className="max-w-xl mx-auto">
              <ContactForm />
            </div>
          </div>
        </motion.section>
      </div>
    </div>
  );
};