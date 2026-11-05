import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Search, Filter, BookOpen, ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';
import { api } from '@/lib/api';
import { PostCard } from '@/components/PostCard';
import { WaveMesh } from '@/components/WaveMesh';

export const PostsPage: React.FC = () => {
  const [selectedCategory, setSelectedCategory] = useState<string>('');
  const [searchQuery, setSearchQuery] = useState<string>('');

  const { data: categories = [] } = useQuery({
    queryKey: ['categories-post'],
    queryFn: () => api.getCategories('post'),
  });

  const { data: postsData, isLoading } = useQuery({
    queryKey: ['posts', selectedCategory],
    queryFn: () =>
      api.getPosts({
        category: selectedCategory || undefined,
      }),
  });

  const posts = postsData?.posts || [];

  const filteredPosts = posts.filter((p) => {
    if (!searchQuery.trim()) return true;
    const q = searchQuery.toLowerCase();
    return (
      p.title.toLowerCase().includes(q) ||
      p.excerpt.toLowerCase().includes(q)
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
            <BookOpen className="w-4 h-4 text-[#39FF14]" />
            <span>Blog & Publications</span>
          </div>
          <h1 className="text-2xl sm:text-4xl font-extrabold text-text mt-1">
            Articles & Technical Notes
          </h1>
          <p className="mt-2 text-sm text-text-muted max-w-2xl">
            Writings, technology explorations, architecture case studies, and software engineering guides.
          </p>
        </div>

        {/* Filters & Search */}
        <div className="flex flex-col sm:flex-row gap-3 mb-8 p-3 rounded-xl bg-surface-2/60 backdrop-blur-md border border-border">
          {/* Search */}
          <div className="relative flex-1">
            <Search className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-text-faint" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search articles by title or excerpt..."
              className="w-full pl-10 pr-3.5 py-2 text-xs rounded-md bg-surface border border-border text-text placeholder:text-text-faint focus:outline-none focus:border-[#39FF14] transition-colors"
            />
          </div>

          {/* Category */}
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
          </div>
        </div>

        {/* Posts Grid */}
        {isLoading ? (
          <div className="py-20 text-center text-xs text-text-faint font-mono">
            Loading articles...
          </div>
        ) : filteredPosts.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredPosts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        ) : (
          <div className="py-16 text-center rounded-xl bg-surface-2/60 backdrop-blur-md border border-border text-xs text-text-faint font-mono">
            No articles match the selected filters or search query.
          </div>
        )}
      </div>
    </div>
  );
};