import type { Metadata } from 'next';
import { createApiClient } from '@/lib/api';

type MetadataProps = {
  params: Promise<{ movieId: string }>;
};

const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080';

export async function generateMetadata({ params }: MetadataProps): Promise<Metadata> {
  const { movieId } = await params;
  const api = createApiClient({ baseUrl: apiBaseUrl });

  try {
    const movie = await api.getMovie(movieId);
    return {
      title: movie.title,
      description: movie.synopsis ?? 'ภาพยนตร์พร้อมสตรีมจาก Leak Streaming Portal',
      openGraph: {
        title: movie.title,
        description: movie.synopsis ?? undefined,
        images: movie.posterUrl ? [{ url: movie.posterUrl }] : undefined
      }
    };
  } catch {
    return {
      title: 'ไม่พบภาพยนตร์',
      description: 'ไม่สามารถดึงข้อมูลภาพยนตร์ที่ต้องการได้'
    };
  }
}
