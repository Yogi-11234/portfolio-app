import React from 'react';
import { Link } from 'react-router-dom';
import { Clock, Calendar, ArrowRight } from 'lucide-react';
import { Post } from '@/types';
import { GlowingEffect } from './GlowingEffect';
import { getMediaUrl } from '@/lib/utils';

interface PostCardProps {
  post: Post;
}

export const PostCard: React.FC<PostCardProps> = ({ post }) => {
  const formatDate = (iso?: string) => {
    if (!iso) return '-';
    try {
      const d = new Date(iso);
      return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' });
    } catch {
      return iso;
    }
  };

  return (
    <article className="relative group flex flex-col justify-between rounded-2xl neumorphic-card overflow-hidden p-5 transition-all duration-300 w-full max-w-[360px] mx-auto">
      <GlowingEffect proximity={60} spread={70} borderWidth={2} halo={true} />

      {/* Cover image if any */}
      {post.cover_path && (
        <div className="aspect-video w-full overflow-hidden rounded-xl bg-surface-2/80 border border-white/5 mb-4 shadow-inner">
          <img
            src={getMediaUrl(post.cover_path)}
            alt={post.title}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500 ease-out"
            loading="lazy"
          />
        </div>
      )}

      <div className="p-5 flex-1 flex flex-col justify-between">
        <div>
          {/* Metadata */}
          <div className="flex items-center gap-3 text-xs text-text-faint font-mono">
            {post.category && (
              <span className="text-accent">{post.category.name}</span>
            )}
            <span className="inline-flex items-center gap-1">
              <Calendar className="w-3 h-3" />
              <span>{formatDate(post.published_at || post.created_at)}</span>
            </span>
            <span className="inline-flex items-center gap-1">
              <Clock className="w-3 h-3" />
              <span>{post.reading_minutes} min</span>
            </span>
          </div>

          <h3 className="mt-3 text-base font-bold text-text group-hover:text-accent transition-colors line-clamp-2">
            <Link to={`/posts/${post.slug}`}>{post.title}</Link>
          </h3>

          <p className="mt-2 text-xs text-text-muted line-clamp-3 leading-relaxed">
            {post.excerpt}
          </p>
        </div>

        <div className="mt-4 pt-3 border-t border-border/50">
          <Link
            to={`/posts/${post.slug}`}
            className="neumorphic-mini-btn"
          >
            <span>Read More</span>
            <ArrowRight className="w-3.5 h-3.5 group-hover:translate-x-1 transition-transform" />
          </Link>
        </div>
      </div>
    </article>
  );
};
