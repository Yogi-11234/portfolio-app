import React, { useEffect, useRef } from 'react';
import { cn } from '@/lib/utils';

interface GlowingEffectProps extends React.HTMLAttributes<HTMLDivElement> {
  /** Distance in px outside the card at which the glow starts to fade in */
  proximity?: number;
  /** Length of the lit arc, in degrees of the border */
  spread?: number;
  /** Thickness of the glowing border in px */
  borderWidth?: number;
  /** Add a soft blurred halo behind the lit border */
  halo?: boolean;
}

export function GlowingEffect({
  proximity = 80,
  spread = 80,
  borderWidth = 2,
  halo = true,
  className,
  style,
  ...props
}: GlowingEffectProps) {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;

    let x = -1e5;
    let y = -1e5;
    let angle = 0;
    let lit = false;
    let raf = 0;

    const frame = () => {
      raf = 0;
      const r = el.getBoundingClientRect();
      const dx = Math.max(r.left - x, 0, x - r.right);
      const dy = Math.max(r.top - y, 0, y - r.bottom);
      const dist = Math.hypot(dx, dy);
      const near = Math.max(0, 1 - dist / Math.max(proximity, 1));

      el.style.setProperty('--glow-o', near.toFixed(3));
      el.style.setProperty('--glow-x', `${(x - r.left).toFixed(1)}px`);
      el.style.setProperty('--glow-y', `${(y - r.top).toFixed(1)}px`);

      if (!near) {
        lit = false;
        return;
      }

      // Calculate angle from card center towards cursor
      const target =
        (Math.atan2(x - (r.left + r.width / 2), (r.top + r.height / 2) - y) * 180) / Math.PI;
      const d = ((((target - angle) % 360) + 540) % 360) - 180;
      angle += !lit || Math.abs(d) < 0.5 ? d : d * 0.25;
      lit = true;
      el.style.setProperty('--glow-a', `${angle.toFixed(2)}deg`);

      if (angle !== target && Math.abs(d) >= 0.5) {
        raf = requestAnimationFrame(frame);
      }
    };

    const queue = () => {
      if (!raf) raf = requestAnimationFrame(frame);
    };

    const onMove = (e: PointerEvent | MouseEvent) => {
      x = e.clientX;
      y = e.clientY;
      queue();
    };

    const onOut = (e: MouseEvent) => {
      if (e.relatedTarget) return;
      x = -1e5;
      y = -1e5;
      queue();
    };

    window.addEventListener('pointermove', onMove, { passive: true });
    window.addEventListener('mouseout', onOut);
    window.addEventListener('scroll', queue, { passive: true });

    return () => {
      window.removeEventListener('pointermove', onMove);
      window.removeEventListener('mouseout', onOut);
      window.removeEventListener('scroll', queue);
      cancelAnimationFrame(raf);
    };
  }, [proximity]);

  return (
    <div
      {...props}
      ref={ref}
      aria-hidden="true"
      data-slot="glowing-effect"
      style={
        {
          '--glow-s': `${spread}deg`,
          '--glow-w': `${borderWidth}px`,
          '--brand-from': '#39FF14',
          '--brand-via': '#00ff88',
          '--brand-to': '#10b981',
          opacity: 'var(--glow-o, 0)',
          ...style,
        } as React.CSSProperties
      }
      className={cn('glowing-effect-root', className)}
    >
      {/* 1. Internal interactive spotlight following cursor position */}
      <div className="glow-spotlight" />

      {/* 2. Soft halo around lit border arc */}
      {halo && (
        <div className="glow-arc-container glow-arc-halo">
          <div className="glow-arc-border" style={{ padding: `${borderWidth + 4}px` }} />
        </div>
      )}

      {/* 3. Sharp lit border arc following cursor angle */}
      <div className="glow-arc-container">
        <div className="glow-arc-border" style={{ padding: `${borderWidth}px` }} />
      </div>
    </div>
  );
}
