import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Search, Filter, Layers, ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';
import { api } from '@/lib/api';
import { ProjectCard } from '@/components/ProjectCard';
import { WaveMesh } from '@/components/WaveMesh';

export const ProjectsPage: React.FC = () => {
  const [selectedCategory, setSelectedCategory] = useState<string>('');
  const [selectedSkill, setSelectedSkill] = useState<string>('');
  const [searchQuery, setSearchQuery] = useState<string>('');

  const { data: categories = [] } = useQuery({
    queryKey: ['categories-project'],
    queryFn: () => api.getCategories('project'),
  });

  const { data: skills = [] } = useQuery({
    queryKey: ['skills'],
    queryFn: api.getSkills,
  });

  const { data: projectsData, isLoading } = useQuery({
    queryKey: ['projects', selectedCategory, selectedSkill],
    queryFn: () =>
      api.getProjects({
        category: selectedCategory || undefined,
        skill: selectedSkill || undefined,
      }),
  });

  const projects = projectsData?.projects || [];

  const filteredProjects = projects.filter((p) => {
    if (!searchQuery.trim()) return true;
    const q = searchQuery.toLowerCase();
    return (
      p.title.toLowerCase().includes(q) ||
      p.summary.toLowerCase().includes(q) ||
      p.skills.some((s) => s.name.toLowerCase().includes(q))
    );
  });

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

      {/* 🚀 MAIN CONTENT (Layer z-10 di atas background) */}
      <div className="relative z-10 max-w-6xl mx-auto px-4 sm:px-6 pt-24 pb-16">
        {/* Header */}
        <div className="mb-8">
          <Link
            to="/"
            className="inline-flex items-center gap-1.5 text-xs text-text-muted hover:text-[#39FF14] font-mono mb-4 transition-colors"
          >
            <ArrowLeft className="w-3.5 h-3.5" />
            <span>Back to Home</span>
          </Link>
          <div className="flex items-center gap-2 text-xs font-mono text-accent uppercase tracking-wider">
            <Layers className="w-4 h-4 text-[#39FF14]" />
            <span>Portfolio Archive</span>
          </div>
          <h1 className="text-2xl sm:text-4xl font-extrabold text-text mt-1">All Projects</h1>
          <p className="mt-2 text-sm text-text-muted max-w-2xl">
            A curated collection of web applications, backend services, APIs, and open-source projects I've built.
          </p>
        </div>

        {/* Filters & Search */}
        <div className="flex flex-col md:flex-row gap-3 mb-8 p-3 rounded-xl bg-surface-2/60 backdrop-blur-md border border-border">
          {/* Search Input */}
          <div className="relative flex-1">
            <Search className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-text-faint" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search by title, description, or skill..."
              className="w-full pl-10 pr-3.5 py-2 text-xs rounded-md bg-surface border border-border text-text placeholder:text-text-faint focus:outline-none focus:border-[#39FF14] transition-colors"
            />
          </div>

          {/* Category Filter */}
          <div className="flex items-center gap-2">
            <Filter className="w-4 h-4 text-text-faint hidden sm:block" />
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="px-3 py-2 text-xs rounded-md bg-surface border border-border text-text focus:outline-none focus:border-[#39FF14] transition-colors"
            >
              <option value="">All Categories</option>
              {categories.map((c) => (
                <option key={c.id} value={c.slug}>
                  {c.name}
                </option>
              ))}
            </select>

            {/* Skill Filter */}
            <select
              value={selectedSkill}
              onChange={(e) => setSelectedSkill(e.target.value)}
              className="px-3 py-2 text-xs rounded-md bg-surface border border-border text-text focus:outline-none focus:border-[#39FF14] transition-colors"
            >
              <option value="">All Skills</option>
              {skills.map((s) => (
                <option key={s.id} value={s.slug}>
                  {s.name}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Projects Grid */}
        {isLoading ? (
          <div className="py-20 text-center text-xs text-text-faint font-mono">
            Loading projects...
          </div>
        ) : filteredProjects.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredProjects.map((proj) => (
              <ProjectCard key={proj.id} project={proj} />
            ))}
          </div>
        ) : (
          <div className="py-16 text-center rounded-xl bg-surface-2/60 backdrop-blur-md border border-border text-xs text-text-faint font-mono">
            No projects match the selected filters or search query.
          </div>
        )}
      </div>
    </div>
  );
};