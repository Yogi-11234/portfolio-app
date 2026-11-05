import React, { useState } from 'react';
import { Award, ExternalLink, Eye, Calendar } from 'lucide-react';
import { Certificate } from '@/types';
import { CertificateModal } from './CertificateModal';
import { GlowingEffect } from './GlowingEffect';
import { getMediaUrl } from '@/lib/utils';

interface CertificateCardProps {
  certificate: Certificate;
}

export const CertificateCard: React.FC<CertificateCardProps> = ({ certificate }) => {
  const [modalOpen, setModalOpen] = useState(false);

  const formatDate = (iso: string) => {
    try {
      const d = new Date(iso);
      return d.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
    } catch {
      return iso;
    }
  };

  const certTitle = certificate.name || (certificate as any).title || 'Professional Certificate';
  const hasImage = Boolean(certificate.file_path);

  return (
    <>
      <div className="relative neumorphic-card group flex flex-col justify-between p-4 sm:p-5 transition-all duration-300 w-full max-w-[340px] sm:w-[320px] md:w-[330px] mx-auto min-h-[350px] overflow-hidden">
        {/* Interactive Glowing Border & Spotlight tracking pointer */}
        <GlowingEffect proximity={80} spread={75} borderWidth={2} halo={true} />

        <div className="relative z-10 flex-1 flex flex-col justify-between">
          <div>
            {/* Certificate Image or Emblem Preview Frame */}
            <div
              onClick={() => hasImage && setModalOpen(true)}
              className={`relative aspect-[16/10] rounded-xl bg-surface-2/80 border border-white/5 overflow-hidden shadow-inner mb-4 flex items-center justify-center ${
                hasImage ? 'cursor-pointer' : ''
              }`}
            >
              {hasImage ? (
                <>
                  <img
                    src={getMediaUrl(certificate.file_path)}
                    alt={certTitle}
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500 ease-out"
                    loading="lazy"
                  />
                  {/* <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-end justify-center pb-3">
                    <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-semibold bg-accent text-accent-fg shadow-[0_0_12px_rgba(57,255,20,0.6)]">
                      <Eye className="w-3.5 h-3.5" />
                      <span>Preview</span>
                    </span>
                  </div> */}
                </>
              ) : (
                <div className="flex flex-col items-center justify-center gap-2 p-4 text-center">
                  <div className="w-12 h-12 rounded-xl flex items-center justify-center bg-gradient-to-br from-[#1c1c26] to-[#101016] border border-accent/30 shadow-[0_0_15px_rgba(57,255,20,0.2)]">
                    <Award className="w-6 h-6 text-accent" />
                  </div>
                  <span className="text-[11px] font-mono text-text-faint">
                    {certificate.issuer}
                  </span>
                </div>
              )}

              {/* Issuer pill tag */}
              <span className="absolute top-2 left-2 px-2 py-0.5 rounded-full text-[10px] font-mono bg-black/80 border border-white/10 text-white/90 backdrop-blur-md">
                {certificate.issuer}
              </span>
            </div>

            {/* Meta row: Issue date */}
            <div className="flex items-center justify-between text-[11px] text-text-faint font-mono mb-2">
              <span className="flex items-center gap-1 text-text-muted">
                <Calendar className="w-3 h-3 text-accent" />
                {formatDate(certificate.issued_at)}
              </span>
              {certificate.expires_at && (
                <span className="text-[10px]">
                  Expires: {formatDate(certificate.expires_at)}
                </span>
              )}
            </div>

            {/* Title */}
            <h4 className="text-sm font-bold text-text group-hover:text-accent transition-colors line-clamp-2 leading-snug">
              {certTitle}
            </h4>
          </div>

          {/* Action Buttons */}
          <div className="pt-4 mt-3 border-t border-white/5 flex items-center justify-between gap-2 text-xs">
            {certificate.file_path ? (
              <button
                onClick={() => setModalOpen(true)}
                className="neumorphic-mini-btn flex-1"
              >
                <Eye className="w-3.5 h-3.5 text-accent" />
                <span>View Document</span>
              </button>
            ) : (
              <div className="flex-1" />
            )}

            {certificate.credential_url && (
              <a
                href={certificate.credential_url}
                target="_blank"
                rel="noopener noreferrer"
                className="neumorphic-mini-btn flex-1"
              >
                <span>Verify</span>
                <ExternalLink className="w-3.5 h-3.5 text-accent" />
              </a>
            )}
          </div>
        </div>
      </div>

      {modalOpen && (
        <CertificateModal
          certificate={certificate}
          onClose={() => setModalOpen(false)}
        />
      )}
    </>
  );
};