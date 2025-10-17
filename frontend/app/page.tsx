import Link from 'next/link';
import { ArrowRight } from 'lucide-react';
import { cn } from '@/lib/utils';

const heroStyles =
  'rounded-3xl border border-border bg-card/60 p-10 shadow-sm backdrop-blur supports-[backdrop-filter]:bg-card/40';

export default function HomePage() {
  return (
    <main className="container flex min-h-[calc(100vh-4rem)] flex-col justify-center gap-10 py-12">
      <section className={heroStyles}>
        <div className="flex flex-col gap-6">
          <p className="text-sm uppercase tracking-widest text-muted-foreground">Leak Streaming Studio</p>
          <h1 className="text-balance text-4xl font-semibold leading-tight md:text-5xl">
            สร้างประสบการณ์ดูหนังระดับสตูดิโอด้วย Next.js 15 และระบบหลังบ้าน Go
          </h1>
          <p className="max-w-2xl text-balance text-lg text-muted-foreground">
            โปรเจ็กต์นี้เป็นฐานสำหรับแพลตฟอร์มภาพยนตร์ที่สตรีมด้วย HLS, รองรับ dark mode, มีระบบจัดการคอนเทนต์,
            และออกแบบมาเพื่อรองรับการขยายตัวในระดับโปรดักชันตั้งแต่วินาทีแรก
          </p>
          <div className="flex flex-wrap gap-4">
            <Link
              href="/movies/sample-movie"
              className={cn(
                'inline-flex items-center justify-center rounded-full bg-primary px-6 py-2 text-sm font-medium text-primary-foreground shadow transition hover:bg-primary/90'
              )}
            >
              เริ่มสำรวจหนัง <ArrowRight className="ml-2 size-4" aria-hidden="true" />
            </Link>
            <Link
              href="/admin"
              className="inline-flex items-center justify-center rounded-full border border-border px-6 py-2 text-sm font-medium text-foreground transition hover:bg-secondary"
            >
              สำหรับผู้ดูแลระบบ
            </Link>
          </div>
        </div>
      </section>
      <section className="grid gap-6 md:grid-cols-3">
        {[
          {
            title: 'สตรีมคุณภาพสูง',
            description: 'พร้อมรองรับลิ้งค์ .m3u8 แบบเซ็นลายเซ็นตามสเปก.'
          },
          {
            title: 'ระบบจัดการภาพยนตร์',
            description: 'ออกแบบ workflow รองรับการเพิ่ม/แก้ไข/ซ่อนภาพยนตร์สำหรับทีมคอนเทนต์.'
          },
          {
            title: 'ขยายสเกลง่าย',
            description: 'ระบบ telemetry, rate-limit และ caching blueprint อยู่ในดีไซน์ตั้งแต่แรก.'
          }
        ].map((item) => (
          <article
            key={item.title}
            className="rounded-3xl border border-dashed border-border/70 bg-card/50 p-6 shadow-sm backdrop-blur"
          >
            <h2 className="text-lg font-semibold">{item.title}</h2>
            <p className="mt-2 text-sm text-muted-foreground">{item.description}</p>
          </article>
        ))}
      </section>
    </main>
  );
}
