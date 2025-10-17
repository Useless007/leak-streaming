'use client';

import { useState, useTransition, type FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { Plus, Trash } from 'lucide-react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useFieldArray, useForm, type Path } from 'react-hook-form';

import { Button } from '@/components/ui/button';
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
	FormDescription
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

import { createMovieAction } from './actions';
import { createMovieFormSchema, type CreateMovieFormValues } from './schema';

export function AdminNewMovieForm() {
	const router = useRouter();
	const [formError, setFormError] = useState<string | null>(null);
	const [isPending, startTransition] = useTransition();

	const form = useForm<CreateMovieFormValues>({
		resolver: zodResolver(createMovieFormSchema),
		defaultValues: {
			title: '',
			synopsis: '',
			posterUrl: '',
			availabilityStart: '',
			availabilityEnd: '',
			isVisible: true,
			streamUrl: '',
			drmKeyId: '',
			allowedHosts: '',
			captions: []
		}
	});

	const captions = useFieldArray({
		control: form.control,
		name: 'captions'
	});

	const handleValidSubmit = form.handleSubmit((values) => {
		setFormError(null);
		startTransition(() => {
			form.clearErrors();
			void (async () => {
				const result = await createMovieAction(values);
				if (!result.success) {
					if (result.fieldErrors) {
						for (const [field, message] of Object.entries(result.fieldErrors)) {
							form.setError(field as Path<CreateMovieFormValues>, {
								type: 'server',
								message
							});
						}
					}
					if (result.formError) {
						setFormError(result.formError);
					}
					return;
				}

				form.reset();
				router.push(`/movies/${result.slug}`);
				router.refresh();
			})();
		});
	});

	const onSubmit = (event: FormEvent<HTMLFormElement>) => {
		void handleValidSubmit(event);
	};

	return (
		<Form {...form}>
			<form onSubmit={onSubmit} className="space-y-8">
				<div className="grid gap-6 md:grid-cols-2">
					<FormField
						control={form.control}
						name="title"
						render={({ field }) => (
							<FormItem>
								<FormLabel>ชื่อเรื่อง *</FormLabel>
								<FormControl>
									<Input placeholder="เช่น สปริงรีลีส 2025" {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="posterUrl"
						render={({ field }) => (
							<FormItem>
								<FormLabel>โปสเตอร์ *</FormLabel>
								<FormControl>
									<Input placeholder="https://cdn.example.com/posters/movie.jpg" {...field} />
								</FormControl>
								<FormDescription>ต้องเป็น URL ที่เข้าถึงได้ผ่าน http หรือ https</FormDescription>
								<FormMessage />
							</FormItem>
						)}
					/>
				</div>

				<FormField
					control={form.control}
					name="synopsis"
					render={({ field }) => (
						<FormItem>
							<FormLabel>คำบรรยายภาพยนตร์ *</FormLabel>
							<FormControl>
								<Textarea rows={4} placeholder="สรุปเรื่องย่อสั้น ๆ" {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>

				<div className="grid gap-6 md:grid-cols-2">
					<FormField
						control={form.control}
						name="availabilityStart"
						render={({ field }) => (
							<FormItem>
								<FormLabel>เริ่มฉาย *</FormLabel>
								<FormControl>
									<Input type="datetime-local" {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
					<FormField
						control={form.control}
						name="availabilityEnd"
						render={({ field }) => (
							<FormItem>
								<FormLabel>สิ้นสุด *</FormLabel>
								<FormControl>
									<Input type="datetime-local" {...field} />
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>
				</div>

				<FormField
					control={form.control}
					name="streamUrl"
					render={({ field }) => (
						<FormItem>
							<FormLabel>ลิงก์สตรีม (.m3u8) *</FormLabel>
							<FormControl>
								<Input placeholder="https://cdn.example.com/path/master.m3u8" {...field} />
							</FormControl>
							<FormDescription>ระบบจะสร้าง token และ proxy ให้โดยอัตโนมัติ</FormDescription>
							<FormMessage />
						</FormItem>
					)}
				/>

				<FormField
					control={form.control}
					name="allowedHosts"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Allowed hosts *</FormLabel>
							<FormControl>
								<Textarea rows={3} placeholder={['ตัวอย่าง:', 'main.cdn.example.com', 'm42.edge.example.com'].join('\n')} {...field} />
							</FormControl>
							<FormDescription>ใส่ 1 hostname ต่อ 1 บรรทัด (ระบบจะเพิ่ม host หลักจาก URL ให้อัตโนมัติ)</FormDescription>
							<FormMessage />
						</FormItem>
					)}
				/>

				<FormField
					control={form.control}
					name="drmKeyId"
					render={({ field }) => (
						<FormItem>
							<FormLabel>DRM Key ID (ถ้ามี)</FormLabel>
							<FormControl>
								<Input placeholder="เว้นว่างได้หากไม่มีการเข้ารหัส" {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>

				<div className="space-y-4">
					<div className="flex items-center justify-between">
						<div>
							<h3 className="text-base font-medium">คำบรรยาย</h3>
							<p className="text-sm text-muted-foreground">รองรับหลายภาษา กำหนดภาษาละหนึ่งบรรทัด หากไม่ต้องการให้เว้นว่างได้</p>
						</div>
						<Button
							type="button"
							variant="outline"
							onClick={() => captions.append({ languageCode: '', label: '', captionUrl: '' })}
							className="inline-flex items-center gap-2"
						>
							<Plus className="size-4" aria-hidden /> เพิ่มคำบรรยาย
						</Button>
					</div>
					<div className="space-y-4">
						{captions.fields.length === 0 ? (
							<p className="rounded-xl border border-dashed border-border/60 p-4 text-sm text-muted-foreground">
								ยังไม่มีคำบรรยาย หากต้องการเพิ่มให้กดปุ่ม &quot;เพิ่มคำบรรยาย&quot;
							</p>
						) : (
							captions.fields.map((field, index) => (
								<div key={field.id} className="rounded-2xl border border-border/80 p-4">
									<div className="flex items-center justify-between gap-4">
										<h4 className="text-sm font-medium">คำบรรยาย #{index + 1}</h4>
										<Button
											type="button"
											variant="ghost"
											size="icon"
											onClick={() => captions.remove(index)}
										>
											<span className="sr-only">ลบคำบรรยาย</span>
											<Trash className="size-4" aria-hidden />
										</Button>
									</div>
									<div className="mt-4 grid gap-4 md:grid-cols-[160px_1fr]">
										<FormField
											control={form.control}
											name={`captions.${index}.languageCode` as const}
											render={({ field }) => (
												<FormItem>
													<FormLabel>รหัสภาษา</FormLabel>
													<FormControl>
														<Input placeholder="en" {...field} />
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>
										<div className="space-y-4">
											<FormField
												control={form.control}
												name={`captions.${index}.label` as const}
												render={({ field }) => (
													<FormItem>
														<FormLabel>ชื่อที่แสดง</FormLabel>
														<FormControl>
															<Input placeholder="English" {...field} />
														</FormControl>
														<FormMessage />
													</FormItem>
												)}
											/>
											<FormField
												control={form.control}
												name={`captions.${index}.captionUrl` as const}
												render={({ field }) => (
													<FormItem>
														<FormLabel>ไฟล์คำบรรยาย</FormLabel>
														<FormControl>
															<Input placeholder="/captions/movie-en.vtt" {...field} />
														</FormControl>
														<FormMessage />
													</FormItem>
												)}
											/>
										</div>
									</div>
								</div>
							))
						)}
					</div>
				</div>

				<FormField
					control={form.control}
					name="isVisible"
					render={({ field }) => (
						<FormItem className="flex items-center justify-between rounded-2xl border border-border/70 bg-background/60 px-4 py-3">
							<div>
								<FormLabel className="text-sm font-medium">เผยแพร่ทันที</FormLabel>
								<FormDescription>หากปิดไว้ ภาพยนตร์จะไม่ปรากฏในหน้าผู้ชมจนกว่าจะเปิดใช้งาน</FormDescription>
							</div>
							<FormControl>
								<input
									type="checkbox"
									className="size-5 rounded border border-border transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
									checked={field.value}
									onChange={(event) => field.onChange(event.target.checked)}
								/>
							</FormControl>
						</FormItem>
					)}
				/>

				{formError && (
					<div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
						{formError}
					</div>
				)}

				<div className="flex justify-end gap-3">
					<Button type="reset" variant="ghost" onClick={() => form.reset()} disabled={isPending}>
						ล้างฟอร์ม
					</Button>
					<Button type="submit" disabled={isPending}>
						{isPending ? 'กำลังบันทึก...' : 'บันทึกภาพยนตร์'}
					</Button>
				</div>
			</form>
		</Form>
	);
}
