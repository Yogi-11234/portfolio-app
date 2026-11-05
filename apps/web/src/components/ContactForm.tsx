import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { toast } from 'sonner';
import { Send, Loader2 } from 'lucide-react';
import { api } from '@/lib/api';
import { AnimatedButton } from './AnimatedButton';

const contactSchema = z.object({
  name: z.string().min(1, 'Full name is required').max(100, 'Maximum 100 characters'),
  email: z.string().email('Invalid email address format').max(160, 'Maximum 160 characters'),
  subject: z.string().max(160, 'Maximum 160 characters').optional(),
  body: z.string().min(10, 'Message must be at least 10 characters').max(5000, 'Message must not exceed 5000 characters'),
  website: z.string().optional(), // Honeypot
});

type ContactFormData = z.infer<typeof contactSchema>;

export const ContactForm: React.FC = () => {
  const [submitting, setSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ContactFormData>({
    resolver: zodResolver(contactSchema),
    defaultValues: {
      name: '',
      email: '',
      subject: '',
      body: '',
      website: '',
    },
  });

  const onSubmit = async (data: ContactFormData) => {
    setSubmitting(true);
    try {
      await api.sendContactMessage(data);
      toast.success('Your message has been sent successfully! I will reply soon.');
      reset();
    } catch (err: any) {
      if (err.code === 'RATE_LIMIT_EXCEEDED') {
        toast.error('Too many requests. Please try again in a few moments.');
      } else {
        toast.error(err.message || 'Failed to send message. Please try again.');
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {/* Honeypot field - visually hidden */}
      <div className="hidden" aria-hidden="true">
        <label htmlFor="website">Website (do not fill this out):</label>
        <input
          id="website"
          type="text"
          tabIndex={-1}
          autoComplete="off"
          {...register('website')}
        />
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label htmlFor="name" className="block text-xs font-medium text-text-muted mb-1.5">
            Full Name <span className="text-danger">*</span>
          </label>
          <input
            id="name"
            type="text"
            placeholder="e.g. John Doe"
            {...register('name')}
            className={`w-full px-3.5 py-2.5 rounded-md bg-surface-2 border text-sm text-text placeholder:text-text-faint focus:outline-none focus:border-accent transition-colors ${
              errors.name ? 'border-danger' : 'border-border'
            }`}
          />
          {errors.name && (
            <p className="mt-1 text-[11px] text-danger">{errors.name.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="email" className="block text-xs font-medium text-text-muted mb-1.5">
            Email Address <span className="text-danger">*</span>
          </label>
          <input
            id="email"
            type="email"
            placeholder="e.g. john@example.com"
            {...register('email')}
            className={`w-full px-3.5 py-2.5 rounded-md bg-surface-2 border text-sm text-text placeholder:text-text-faint focus:outline-none focus:border-accent transition-colors ${
              errors.email ? 'border-danger' : 'border-border'
            }`}
          />
          {errors.email && (
            <p className="mt-1 text-[11px] text-danger">{errors.email.message}</p>
          )}
        </div>
      </div>

      <div>
        <label htmlFor="subject" className="block text-xs font-medium text-text-muted mb-1.5">
          Subject (Optional)
        </label>
        <input
          id="subject"
          type="text"
          placeholder="e.g. Web Development Collaboration / Freelance Project"
          {...register('subject')}
          className="w-full px-3.5 py-2.5 rounded-md bg-surface-2 border border-border text-sm text-text placeholder:text-text-faint focus:outline-none focus:border-accent transition-colors"
        />
      </div>

      <div>
        <label htmlFor="body" className="block text-xs font-medium text-text-muted mb-1.5">
          Message <span className="text-danger">*</span>
        </label>
        <textarea
          id="body"
          rows={5}
          placeholder="Write your message or project details here..."
          {...register('body')}
          className={`w-full px-3.5 py-2.5 rounded-md bg-surface-2 border text-sm text-text placeholder:text-text-faint focus:outline-none focus:border-accent transition-colors resize-y ${
            errors.body ? 'border-danger' : 'border-border'
          }`}
        />
        {errors.body && (
          <p className="mt-1 text-[11px] text-danger">{errors.body.message}</p>
        )}
      </div>

      <div className="pt-2">
        <AnimatedButton
          type="submit"
          disabled={submitting}
          className="w-full sm:w-auto"
        >
          {submitting ? (
            <>
              <Loader2 className="w-4 h-4 animate-spin inline mr-1" />
              <span>Sending...</span>
            </>
          ) : (
            <>
              <Send className="w-4 h-4 inline mr-1" />
              <span>Send Message</span>
            </>
          )}
        </AnimatedButton>
      </div>
    </form>
  );
};