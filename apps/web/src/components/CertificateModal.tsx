import React from 'react';
import { X, ExternalLink, Download } from 'lucide-react';
import { Certificate } from '@/types';
import { getMediaUrl } from '@/lib/utils';

interface CertificateModalProps {
  certificate: Certificate;
  onClose: () => void;
}

export const CertificateModal: React.FC<CertificateModalProps> = ({ certificate, onClose }) => {
  if (!certificate.file_path) return null;

  const fileUrl = getMediaUrl(certificate.file_path);
  const isPdf = certificate.file_mime === 'application/pdf' || certificate.file_path.endsWith('.pdf');

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
      <div className="relative w-full max-w-4xl max-h-[90vh] bg-surface border border-border rounded-lg shadow-2xl flex flex-col overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-border bg-surface-2">
          <div>
            <h3 className="text-sm font-bold text-text line-clamp-1">
              {certificate.name || (certificate as any).title}
            </h3>
            <p className="text-xs text-text-faint">{certificate.issuer}</p>
          </div>
          <div className="flex items-center gap-2">
            <a
              href={fileUrl}
              download
              className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface transition-colors"
              title="Unduh"
            >
              <Download className="w-4 h-4" />
            </a>
            <a
              href={fileUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface transition-colors"
              title="Buka Tab Baru"
            >
              <ExternalLink className="w-4 h-4" />
            </a>
            <button
              onClick={onClose}
              className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface transition-colors"
              title="Tutup"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        </div>

        {/* Content Preview */}
        <div className="flex-1 overflow-auto p-4 flex items-center justify-center bg-bg/50">
          {isPdf ? (
            <iframe
              src={fileUrl}
              title={certificate.name}
              className="w-full h-[70vh] rounded border border-border"
            />
          ) : (
            <img
              src={fileUrl}
              alt={certificate.name}
              className="max-h-[70vh] max-w-full object-contain rounded"
            />
          )}
        </div>
      </div>
    </div>
  );
};
