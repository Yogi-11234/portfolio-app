import React, { useEffect, useId, useRef, useState } from "react";
import {
  motion,
  useMotionTemplate,
  useMotionValue,
  useReducedMotion,
  useSpring,
} from "framer-motion";
import { cn } from "@/lib/utils";

interface TextHoverEffectProps
  extends Omit<React.SVGAttributes<SVGSVGElement>, "children"> {
  /** The word or short phrase to draw */
  text: string;
  /** Outline width in SVG units */
  strokeWidth?: number;
  /** Radius of the pointer reveal spotlight in % */
  revealRadius?: number;
  /** Seconds the outline takes to draw itself in on mount */
  drawDuration?: number;
}

const PAD = 4;

export function TextHoverEffect({
  text,
  strokeWidth = 0.8,
  revealRadius = 32,
  drawDuration = 6.5,
  className,
  onPointerMove,
  ...props
}: TextHoverEffectProps) {
  const id = useId();
  const outlineRef = useRef<SVGTextElement>(null);
  const reducedMotion = useReducedMotion();

  const [box, setBox] = useState({
    w: Math.max(text.length * 56 + PAD * 2, 200),
    h: 90,
    x: PAD,
    y: 75,
  });

  const px = useMotionValue(50);
  const py = useMotionValue(50);
  const spring = { stiffness: 260, damping: 32, mass: 0.6 };
  const cx = useMotionTemplate`${useSpring(px, spring)}%`;
  const cy = useMotionTemplate`${useSpring(py, spring)}%`;

  useEffect(() => {
    let live = true;
    document.fonts.ready.then(() => {
      const el = outlineRef.current;
      if (!live || !el) return;
      try {
        const bb = el.getBBox();
        const css = getComputedStyle(el);
        const ctx = document.createElement("canvas").getContext("2d")!;
        ctx.font = `${css.fontStyle} ${css.fontWeight} 100px ${css.fontFamily}`;
        const ink = ctx.measureText(text);
        const top = ink.actualBoundingBoxAscent || 70;
        const bottom = ink.actualBoundingBoxDescent || 20;
        setBox({
          w: Math.max(bb.width + PAD * 2, 100),
          h: Math.max(top + bottom + PAD * 2, 80),
          x: PAD - (bb.x || 0),
          y: PAD + top,
        });
      } catch {
        setBox({
          w: text.length * 58 + PAD * 2,
          h: 90,
          x: PAD,
          y: 75,
        });
      }
    });
    return () => {
      live = false;
    };
  }, [text]);

  // Initial draw-in outline effect
  useEffect(() => {
    const el = outlineRef.current;
    if (!el) return;
    void getComputedStyle(el).strokeDashoffset;
    el.style.strokeDashoffset = "0";
  }, []);

  const frame = { maskUnits: "userSpaceOnUse", x: 0, y: 0, width: "100%", height: "100%" } as const;
  const word = (attrs: React.SVGProps<SVGTextElement>) => (
    <text
      transform={`translate(${box.x} ${box.y})`}
      fontSize={100}
      fontFamily="inherit"
      fontWeight="inherit"
      textAnchor="start"
      dominantBaseline="alphabetic"
      {...attrs}
    >
      {text}
    </text>
  );

  return (
    <svg
      role="img"
      aria-label={text}
      viewBox={`0 0 ${box.w} ${box.h}`}
      data-slot="text-hover-effect"
      className={cn(
        "group/text-hover h-auto w-full select-none font-extrabold cursor-pointer transition-all duration-300",
        className
      )}
      onPointerMove={(event) => {
        onPointerMove?.(event);
        if (reducedMotion) return;
        const rect = event.currentTarget.getBoundingClientRect();
        px.set(((event.clientX - rect.left) / rect.width) * 100);
        py.set(((event.clientY - rect.top) / rect.height) * 100);
      }}
      {...props}
    >
      <defs>
        {/* Neon green gradient matching the theme (#39FF14) */}
        <linearGradient id={`${id}-ink`} x1="0" x2="1" y1="0" y2="0">
          <stop offset="0%" stopColor="#10b981" />
          <stop offset="50%" stopColor="#39FF14" />
          <stop offset="100%" stopColor="#00ff88" />
        </linearGradient>

        {/* Pointer Spotlight Radial Gradient */}
        <motion.radialGradient
          id={`${id}-spot`}
          gradientUnits="userSpaceOnUse"
          cx={cx}
          cy={cy}
          r={`${revealRadius}%`}
        >
          <stop offset="0%" stopColor="white" />
          <stop offset="60%" stopColor="white" stopOpacity={0.7} />
          <stop offset="100%" stopColor="white" stopOpacity={0} />
        </motion.radialGradient>

        <mask id={`${id}-spot-mask`} {...frame}>
          <rect width="100%" height="100%" fill={`url(#${id}-spot)`} />
        </mask>

        {/* Hides inner fill so base outline remains crisp */}
        <mask id={`${id}-outside`} {...frame}>
          <rect width="100%" height="100%" fill="white" />
          {word({ fill: "black" })}
        </mask>
      </defs>

      {/* 1. Base Outline: Clean crisp neon outline without background blur */}
      <g mask={`url(#${id}-outside)`}>
        {word({
          ref: outlineRef,
          strokeWidth: strokeWidth * 2,
          style: { transitionDuration: `${drawDuration}s` },
          className:
            "fill-none stroke-[#39FF14]/40 group-hover/text-hover:stroke-[#39FF14]/70 transition-[stroke-dashoffset,stroke] ease-in-out [stroke-dasharray:1000] motion-safe:[stroke-dashoffset:1000] motion-reduce:transition-none",
          onTransitionEnd: (event) => {
            event.currentTarget.style.strokeDasharray = "none";
          },
        })}
      </g>

      {/* 2. Revealed Text: Specifically ONLY fills the letters without bleeding to background */}
      <g
        mask={`url(#${id}-spot-mask)`}
        className="opacity-0 transition-opacity duration-300 group-hover/text-hover:opacity-100 motion-reduce:opacity-100"
      >
        {/* Crisp solid fill inside glyphs */}
        {word({
          fill: `url(#${id}-ink)`,
          fillOpacity: 0.95,
        })}
        {/* Crisp accent stroke along letter borders */}
        {word({
          fill: "none",
          stroke: `url(#${id}-ink)`,
          strokeWidth: strokeWidth * 2.2,
        })}
      </g>
    </svg>
  );
}

export default TextHoverEffect;
