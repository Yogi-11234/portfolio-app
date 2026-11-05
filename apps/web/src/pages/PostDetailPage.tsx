import React from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, Calendar, Clock, User, Tag } from 'lucide-react';
import { api } from '@/lib/api';
import { SanitizedHTML } from '@/components/SanitizedHTML';
import { AnimatedButton } from '@/components/AnimatedButton';
import { getMediaUrl } from '@/lib/utils';
import { WaveMesh } from '@/components/WaveMesh';

export const PostDetailPage: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();

  const { data: post, isLoading, error } = useQuery({
    queryKey: ['post', slug],
    queryFn: () => api.getPostBySlug(slug!),
    enabled: Boolean(slug),
  });

  const formatDate = (iso?: string) => {
    if (!iso) return '-';
    try {
      const d = new Date(iso);
      return d.toLocaleDateString('en-US', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
      });
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
        <div className="relative z-10 max-w-3xl mx-auto px-4 pt-32 pb-20 text-center font-mono text-xs text-text-faint">
          Loading article...
        </div>
      </div>
    );
  }

  if (error || !post) {
    return (
      <div className="relative w-full min-h-screen">
        <div className="fixed inset-0 z-0 pointer-events-auto w-full h-full">
          <WaveMesh
            colors={["#06070a", "#39FF14", "#00ff88"]}
            crests={true}
            className="w-full h-full"
          />
        </div>
        <div className="relative z-10 max-w-3xl mx-auto px-4 pt-32 pb-20 text-center">
          <h2 className="text-xl font-bold text-text">Article Not Found</h2>
          <p className="mt-2 text-xs text-text-muted">
            The article you are looking for might not be published yet or has been moved.
          </p>
          <div className="mt-6 flex justify-center">
            <AnimatedButton as="link" to="/posts">
              <span>Back to Articles</span>
            </AnimatedButton>
          </div>
        </div>
      </div>
    );
  }

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

      {/* 🚀 ARTICLE CONTENT (Layer z-10 di atas background) */}
      <article className="relative z-10 max-w-3xl mx-auto px-4 sm:px-6 pt-24 pb-20">
        {/* Back button */}
        <Link
          to="/posts"
          className="inline-flex items-center gap-1.5 text-xs text-text-muted hover:text-[#39FF14] font-mono mb-6 transition-colors"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          <span>Back to All Articles</span>
        </Link>

        {/* Header */}
        <header className="space-y-4">
          {post.category && (
            <div className="inline-flex items-center gap-1.5 text-xs font-mono text-accent bg-surface-2 border border-border px-2.5 py-1 rounded">
              <Tag className="w-3 h-3" />
              <span>{post.category.name}</span>
            </div>
          )}

          <h1 className="text-2xl sm:text-4xl font-extrabold text-text tracking-tight leading-tight">
            {post.title}
          </h1>

          <div className="flex flex-wrap items-center gap-4 text-xs text-text-faint font-mono border-y border-border/60 py-3">
            <span className="inline-flex items-center gap-1 text-text-muted">
              <User className="w-3.5 h-3.5" />
              <span>{post.author_name || 'Author'}</span>
            </span>
            <span>•</span>
            <span className="inline-flex items-center gap-1">
              <Calendar className="w-3.5 h-3.5" />
              <span>{formatDate(post.published_at || post.created_at)}</span>
            </span>
            <span>•</span>
            <span className="inline-flex items-center gap-1">
              <Clock className="w-3.5 h-3.5" />
              <span>{post.reading_minutes} min read</span>
            </span>
          </div>
        </header>

        {/* Cover Image */}
        {post.cover_path && (
          <div className="mt-8 rounded-lg overflow-hidden border border-border bg-surface-2 shadow-xl">
            <img
              src={getMediaUrl(post.cover_path)}
              alt={post.title}
              className="w-full h-auto object-contain block"
            />
          </div>
        )}

        {/* Content */}
        <div className="mt-10">
          <SanitizedHTML
            html={post.content_html}
            className="text-base text-text leading-relaxed prose prose-invert max-w-none space-y-6"
          />
        </div>
      </article>
    </div>
  );
};