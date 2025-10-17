import Image from 'next/image';
import { notFound } from 'next/navigation';
import { Suspense } from 'react';
import { createApiClient, type Movie } from '@/lib/api';
import { MoviePlayer } from '@/components/movie/player';

type PageProps = {
  params: Promise<{ movieId: string }>;
};

const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080';

export default async function MoviePage({ params }: PageProps) {
  const { movieId } = await params;
  const api = createApiClient({ baseUrl: apiBaseUrl });

  let movie: Movie;
  try {
    movie = await api.getMovie(movieId);
  } catch (error) {
    if (error instanceof Error && /not found/i.test(error.message)) {
      notFound();
    }
    throw error;
  }

  return (
    <div className="container grid gap-12 py-12 lg:grid-cols-[420px_1fr]">
      <aside className="space-y-6">
        <div className="relative aspect-[3/4] overflow-hidden rounded-3xl border border-border shadow">
          {movie.posterUrl ? (
            <Image
              src={movie.posterUrl}
              alt={movie.title}
              fill
              priority
              sizes="(min-width: 1024px) 420px, 100vw"
              className="object-cover"
            />
          ) : (
            <div className="flex h-full items-center justify-center bg-muted text-muted-foreground">
              ไม่มีโปสเตอร์
            </div>
          )}
        </div>
        <div className="rounded-3xl border border-border/80 bg-card/60 p-6 shadow-sm backdrop-blur">
          <h1 className="text-3xl font-semibold">{movie.title}</h1>
          {movie.synopsis && <p className="mt-4 text-sm leading-6 text-muted-foreground">{movie.synopsis}</p>}
          <Availability window={movie} />
        </div>
      </aside>
      <section className="space-y-8">
        <Suspense fallback={<div className="aspect-video w-full animate-pulse rounded-3xl bg-muted" />}>
          <MoviePlayer movie={movie} />
        </Suspense>
      </section>
    </div>
  );
}

type AvailabilityProps = {
  window: Pick<Movie, 'availabilityStart' | 'availabilityEnd'>;
};

function Availability({ window }: AvailabilityProps) {
  if (!window.availabilityStart && !window.availabilityEnd) {
    return null;
  }

  const start = window.availabilityStart ? new Date(window.availabilityStart).toLocaleString('th-TH') : 'ไม่ระบุ';
  const end = window.availabilityEnd ? new Date(window.availabilityEnd).toLocaleString('th-TH') : 'ไม่ระบุ';

  return (
    <div className="mt-6 space-y-2 text-sm text-muted-foreground">
      <div>
        <span className="font-medium text-foreground">เริ่มฉาย:</span> {start}
      </div>
      <div>
        <span className="font-medium text-foreground">สิ้นสุด:</span> {end}
      </div>
    </div>
  );
}
