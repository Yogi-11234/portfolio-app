import React, { useRef, useState, useEffect } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { motion } from 'framer-motion';
import { Certificate } from '@/types';
import { CertificateCard } from './CertificateCard';

interface CertificateCarouselProps {
  certificates: Certificate[];
}

export const CertificateCarousel: React.FC<CertificateCarouselProps> = ({ certificates }) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(false);
  const [activeIndex, setActiveIndex] = useState(0);

  const checkScroll = () => {
    const el = containerRef.current;
    if (!el) return;
    const { scrollLeft, scrollWidth, clientWidth } = el;
    setCanScrollLeft(scrollLeft > 10);
    setCanScrollRight(scrollLeft < scrollWidth - clientWidth - 10);

    // Calculate active slide index
    const cardWidth = el.querySelector<HTMLElement>('.cert-slide')?.offsetWidth || clientWidth;
    const idx = Math.round(scrollLeft / (cardWidth + 24));
    setActiveIndex(Math.max(0, Math.min(idx, certificates.length - 1)));
  };

  useEffect(() => {
    checkScroll();
    window.addEventListener('resize', checkScroll);
    return () => window.removeEventListener('resize', checkScroll);
  }, [certificates]);

  const scrollByCard = (direction: 'left' | 'right') => {
    const el = containerRef.current;
    if (!el) return;
    const card = el.querySelector<HTMLElement>('.cert-slide');
    const shift = card ? card.offsetWidth + 24 : 330;
    el.scrollBy({
      left: direction === 'left' ? -shift : shift,
      behavior: 'smooth',
    });
  };

  if (certificates.length === 0) {
    return (
      <div className="py-12 text-center text-xs text-text-faint font-mono rounded-2xl bg-surface border border-border w-full max-w-xl mx-auto">
        No certificates to display.
      </div>
    );
  }

  return (
    <div className="relative w-full max-w-5xl mx-auto">
      {/* Header Controls */}
      {certificates.length > 1 && (
        <div className="flex items-center justify-between gap-3 mb-4 px-2">
          <div className="flex items-center gap-2">
            <span className="text-xs font-mono text-accent font-semibold">
              {activeIndex + 1}
            </span>
            <span className="text-xs font-mono text-text-faint">
              / {certificates.length} Certificates
            </span>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => scrollByCard('left')}
              disabled={!canScrollLeft}
              className="neumorphic-mini-btn p-2 rounded-xl disabled:opacity-30 disabled:pointer-events-none"
              aria-label="Previous certificate"
              title="Previous"
            >
              <ChevronLeft className="w-4 h-4 text-accent" />
            </button>
            <button
              onClick={() => scrollByCard('right')}
              disabled={!canScrollRight}
              className="neumorphic-mini-btn p-2 rounded-xl disabled:opacity-30 disabled:pointer-events-none"
              aria-label="Next certificate"
              title="Next"
            >
              <ChevronRight className="w-4 h-4 text-accent" />
            </button>
          </div>
        </div>
      )}

      {/* Horizontal Slide Track */}
      {/* 
        Perbaikan: 
        1. overflow-x-auto overflow-y-hidden -> Mengizinkan scroll horizontal & mematikan scroll vertikal kanan.
        2. [scrollbar-height:6px] -> Menjaga scrollbar bawah tetap tipis & pas.
      */}
      <div
        ref={containerRef}
        onScroll={checkScroll}
        className={`flex gap-6 overflow-x-auto overflow-y-hidden pb-6 pt-2 px-2 snap-x snap-mandatory neon-scrollbar scroll-smooth items-stretch ${
          certificates.length <= 3 ? 'sm:justify-center' : 'justify-start'
        }`}
      >
        {certificates.map((cert, index) => (
          <motion.div
            key={cert.id}
            initial={{ opacity: 0, y: 35 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: '-60px' }}
            transition={{
              duration: 0.85,
              delay: Math.min(index * 0.12, 0.5),
              ease: [0.16, 1, 0.3, 1],
            }}
            className="cert-slide flex-shrink-0 snap-center w-[85vw] max-w-[320px] sm:w-[calc((100%-3rem)/2)] md:w-[calc((100%-3rem)/3)] md:max-w-[340px] flex justify-center"
          >
            <CertificateCard certificate={cert} />
          </motion.div>
        ))}
      </div>

      {/* Mobile Indicator Dots */}
      {certificates.length > 1 && (
        <div className="flex sm:hidden justify-center items-center gap-1.5 mt-2">
          {certificates.map((_, idx) => (
            <span
              key={idx}
              className={`h-1.5 rounded-full transition-all duration-300 ${
                idx === activeIndex
                  ? 'w-6 bg-accent shadow-[0_0_8px_#39FF14]'
                  : 'w-1.5 bg-white/20'
              }`}
            />
          ))}
        </div>
      )}
    </div>
  );
};