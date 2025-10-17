import Image from 'next/image';
import Link from 'next/link';

import { createApiClient, type MovieSummary } from '@/lib/api';

const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080';

function formatDate(value?: string | null) {
  if (!value) {
    return null;
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return null;
  }
  return new Intl.DateTimeFormat('th-TH', {
    dateStyle: 'medium',
    timeStyle: 'short'
  }).format(date);
}

export default async function MoviesPage() {
  const api = createApiClient({ baseUrl: apiBaseUrl });
  const movies = await api.listMovies();

  return (
    <main className="container space-y-10 py-12">
      <header className="space-y-2">
        <p className="text-sm uppercase tracking-wide text-muted-foreground">กำลังฉาย</p>
        <h1 className="text-3xl font-semibold md:text-4xl">เลือกภาพยนตร์ที่คุณอยากรับชม</h1>
        <p className="max-w-2xl text-sm text-muted-foreground">
          ระบบจะสตรีมผ่าน backend proxy พร้อมปกป้องลิงก์ .m3u8—คลิกที่ภาพยนตร์เพื่อเริ่มเล่นทันที
        </p>
      </header>

      {movies.length === 0 ? (
        <div className="rounded-3xl border border-dashed border-border/60 bg-muted/20 p-10 text-center text-muted-foreground">
          ยังไม่มีภาพยนตร์ที่พร้อมสตรีมในขณะนี้
        </div>
      ) : (
        <section className="grid gap-6 sm:grid-cols-2 xl:grid-cols-3">
          {movies.map((movie) => (
            <MovieCard key={movie.id} movie={movie} />
          ))}
        </section>
      )}
    </main>
  );
}

function MovieCard({ movie }: { movie: MovieSummary }) {
  const availabilityLabel = (() => {
    const starts = formatDate(movie.availabilityStart ?? undefined);
    const ends = formatDate(movie.availabilityEnd ?? undefined);

    if (starts && ends) {
      return `ฉายตั้งแต่ ${starts} ถึง ${ends}`;
    }
    if (starts) {
      return `พร้อมฉายตั้งแต่ ${starts}`;
    }
    if (ends) {
      return `ฉายถึง ${ends}`;
    }
    return 'พร้อมให้รับชม';
  })();

  return (
    <article
      className="group flex h-full flex-col overflow-hidden rounded-3xl border border-border bg-card/60 shadow-sm transition hover:-translate-y-1 hover:shadow-lg"
      data-testid="movie-card"
      data-movie-slug={movie.slug}
    >
      <Link href={`/movies/${movie.slug}`} className="flex h-full flex-col">
        <div className="relative aspect-[3/4] w-full overflow-hidden bg-muted">
          {movie.posterUrl ? (
            <Image
              src={movie.posterUrl}
              alt={movie.title}
              fill
              sizes="(max-width: 768px) 100vw, (max-width: 1280px) 50vw, 33vw"
              className="object-cover transition duration-500 group-hover:scale-105"
            />
          ) : (
            <div className="flex h-full items-center justify-center text-muted-foreground">ไม่มีโปสเตอร์</div>
          )}
        </div>
        <div className="flex flex-1 flex-col gap-3 p-5">
          <div className="space-y-1">
            <h2 className="text-lg font-semibold text-foreground transition group-hover:text-primary">
              {movie.title}
            </h2>
            {movie.synopsis ? (
              <p className="line-clamp-2 text-sm text-muted-foreground">{movie.synopsis}</p>
            ) : (
              <p className="text-sm text-muted-foreground">ไม่มีคำบรรยายประกอบ</p>
            )}
          </div>
          <div className="mt-auto text-xs text-muted-foreground">{availabilityLabel}</div>
        </div>
      </Link>
    </article>
  );
}
