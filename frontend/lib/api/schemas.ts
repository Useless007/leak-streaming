import { z } from 'zod';

const urlOrPathSchema = z
  .string()
  .min(1)
  .refine(
    (value) => {
      if (value.startsWith('/')) {
        return true;
      }
      try {
        const parsed = new URL(value);
        return parsed.protocol === 'http:' || parsed.protocol === 'https:';
      } catch {
        return false;
      }
    },
    { message: 'ต้องเป็น URL แบบ http(s) หรือ path ที่ขึ้นต้นด้วย /' }
  );

export const captionSchema = z.object({
  languageCode: z.string().min(2).max(10),
  label: z.string().min(1),
  captionUrl: urlOrPathSchema
});

export const movieSummarySchema = z.object({
  id: z.string(),
  slug: z.string().min(1),
  title: z.string().min(1, 'ต้องกรอกชื่อเรื่อง'),
  synopsis: z.string().optional(),
  posterUrl: z.string().url().optional(),
  availabilityStart: z.string().datetime().optional().nullable(),
  availabilityEnd: z.string().datetime().optional().nullable(),
  isVisible: z.boolean()
});

export const movieSchema = movieSummarySchema.extend({
  captions: captionSchema.array().default([])
});

export const streamSchema = z.object({
  movieId: z.string(),
  streamUrl: z.string().url().endsWith('.m3u8', 'ต้องเป็นลิงก์ .m3u8'),
  drmKeyId: z.string().optional()
});

export const playbackTokenSchema = z.object({
  token: z.string().min(10),
  movieId: z.string(),
  expiresAt: z.string().datetime(),
  issuedAt: z.string().datetime(),
  viewerId: z.string().optional(),
  url: z.string().url().optional()
});

export type MovieSummary = z.infer<typeof movieSummarySchema>;
export type Movie = z.infer<typeof movieSchema>;
export type Caption = z.infer<typeof captionSchema>;
export type Stream = z.infer<typeof streamSchema>;
export type PlaybackToken = z.infer<typeof playbackTokenSchema>;
