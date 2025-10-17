import Link from 'next/link';
import type { Metadata } from 'next';

type AdminLayoutProps = {
	children: React.ReactNode;
};

export const metadata: Metadata = {
	title: 'ศูนย์จัดการคอนเทนต์'
};

export default function AdminLayout({ children }: AdminLayoutProps) {
	return (
		<div className="container grid gap-8 py-12 lg:grid-cols-[280px_1fr]">
			<aside className="rounded-3xl border border-border/80 bg-card/70 p-6 shadow-sm backdrop-blur">
				<h1 className="text-2xl font-semibold">ศูนย์จัดการคอนเทนต์</h1>
				<p className="mt-2 text-sm text-muted-foreground">
					เพิ่มและจัดตารางภาพยนตร์ใหม่ ตรวจสอบสถานะการเผยแพร่ และเตรียมข้อมูลประกอบสตรีมมิ่ง
				</p>
				<nav className="mt-6 space-y-2 text-sm">
					<Link
						href="/admin/movies/new"
						className="block rounded-xl border border-transparent px-3 py-2 font-medium text-foreground transition hover:border-border hover:bg-accent/40"
					>
						เพิ่มภาพยนตร์ใหม่
					</Link>
				</nav>
			</aside>
			<section className="min-h-[520px] rounded-3xl border border-border/70 bg-card/40 p-6 shadow-sm backdrop-blur">
				{children}
			</section>
		</div>
	);
}
