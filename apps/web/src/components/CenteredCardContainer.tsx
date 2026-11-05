import React, { useRef } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { motion } from 'framer-motion';

interface CenteredCardContainerProps<T> {
  items: T[];
  renderItem: (item: T, index: number) => React.ReactNode;
  emptyMessage?: string;
}

export function CenteredCardContainer<T extends { id: string | number }>({
  items,
  renderItem,
  emptyMessage = 'No data to display.',
}: CenteredCardContainerProps<T>) {
  const scrollRef = useRef<HTMLDivElement>(null);

  const scroll = (direction: 'left' | 'right') => {
    if (scrollRef.current) {
      const scrollAmount = 350;
      scrollRef.current.scrollBy({
        left: direction === 'left' ? -scrollAmount : scrollAmount,
        behavior: 'smooth',
      });
    }
  };

  if (items.length === 0) {
    return (
      <div className="py-12 text-center text-xs text-text-faint font-mono rounded-2xl bg-surface border border-border w-full max-w-xl mx-auto">
        {emptyMessage}
      </div>
    );
  }

  const isScrollable = items.length > 4;

  return (
    <div className="relative w-full">
      {/* Scroll controls when > 4 items */}
      {isScrollable && (
        <div className="flex items-center justify-between sm:justify-end gap-2 mb-3 px-2">
          <span className="text-[11px] font-mono text-text-faint sm:mr-auto">
            Scroll to see all ({items.length} items)
          </span>
          <button
            onClick={() => scroll('left')}
            className="neumorphic-mini-btn p-2 rounded-xl"
            aria-label="Scroll left"
            title="Scroll left"
          >
            <ChevronLeft className="w-4 h-4" />
          </button>
          <button
            onClick={() => scroll('right')}
            className="neumorphic-mini-btn p-2 rounded-xl"
            aria-label="Scroll right"
            title="Scroll right"
          >
            <ChevronRight className="w-4 h-4" />
          </button>
        </div>
      )}

      {/* Container */}
      <div
        ref={scrollRef}
        className={
          isScrollable
            ? 'flex gap-6 overflow-x-auto pb-6 pt-2 px-2 snap-x snap-mandatory neon-scrollbar scroll-smooth items-stretch'
            : 'flex flex-wrap justify-center items-stretch gap-6'
        }
      >
        {items.map((item, index) => (
          <motion.div
            key={item.id}
            initial={{ opacity: 0, y: 35 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: '-60px' }}
            transition={{
              duration: 0.85,
              delay: Math.min(index * 0.12, 0.45),
              ease: [0.16, 1, 0.3, 1],
            }}
            className={`flex justify-center ${isScrollable ? 'flex-shrink-0 snap-center' : ''}`}
          >
            {renderItem(item, index)}
          </motion.div>
        ))}
      </div>
    </div>
  );
}
