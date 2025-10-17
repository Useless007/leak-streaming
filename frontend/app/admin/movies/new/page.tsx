import type { Metadata } from 'next';

import { AdminNewMovieForm } from './form';

export const metadata: Metadata = {
	title: 'เพิ่มภาพยนตร์ใหม่'
};

export default function AdminNewMoviePage() {
	return (
		<div className="mx-auto flex max-w-3xl flex-col gap-8">
			<header className="space-y-2">
				<h2 className="text-3xl font-semibold tracking-tight">เพิ่มภาพยนตร์ใหม่</h2>
				<p className="text-sm text-muted-foreground">
					กรอกข้อมูลรายละเอียดภาพยนตร์ กำหนดช่วงเวลาการเผยแพร่ และตั้งค่าลิงก์สตรีมเพื่อให้ผู้ชมเริ่มรับชมได้ทันที
				</p>
			</header>
			<AdminNewMovieForm />
		</div>
	);
}
