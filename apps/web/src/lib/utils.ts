import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function getMediaUrl(path?: string | null): string {
  if (!path) return '';
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path;
  }
  const clean = path.replace(/^\/?uploads\//, '').replace(/^\//, '');
  const supabaseUrl = import.meta.env.VITE_SUPABASE_URL || 'https://kpbuhrpgnrqmwrzzxols.supabase.co';
  return `${supabaseUrl}/storage/v1/object/public/portfolio_assets/${clean}`;
}
