import { z } from 'zod';

const urlOrPathSchema = z
	.string({ required_error: 'จำเป็นต้องระบุ URL ของไฟล์คำบรรยาย' })
	.trim()
	.refine(
		(value) => {
			if (!value) {
				return false;
			}
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

export const captionInputSchema = z.object({
	languageCode: z
		.string({ required_error: 'กรุณาระบุรหัสภาษา' })
		.trim()
		.min(2, 'รหัสภาษาต้องมีอย่างน้อย 2 ตัวอักษร')
		.max(10, 'รหัสภาษาต้องไม่เกิน 10 ตัวอักษร')
		.transform((value) => value.toLowerCase()),
	label: z.string({ required_error: 'กรุณาระบุชื่อคำบรรยาย' }).trim().min(1, 'ต้องระบุชื่อคำบรรยาย'),
	captionUrl: urlOrPathSchema
});

export const createMovieFormSchema = z
	.object({
		title: z.string({ required_error: 'กรุณาระบุชื่อเรื่อง' }).trim().min(1, 'กรุณาระบุชื่อเรื่อง'),
		synopsis: z.string().trim().optional(),
		posterUrl: z
			.string()
			.trim()
			.optional()
			.refine((value) => !value || /^https?:\/\//.test(value), {
				message: 'โปสเตอร์ต้องเป็น URL แบบ http(s)'
			}),
		availabilityStart: z.string().optional(),
		availabilityEnd: z.string().optional(),
		isVisible: z.boolean().default(true),
		streamUrl: z
			.string({ required_error: 'กรุณาระบุลิงก์สตรีม' })
			.trim()
			.url('ต้องเป็น URL แบบ http(s)')
			.refine((value) => value.toLowerCase().includes('.m3u8'), {
				message: 'ต้องเป็นลิงก์ไฟล์ .m3u8'
			}),
		drmKeyId: z.string().trim().optional(),
		allowedHosts: z.string().trim().optional(),
		captions: z.array(captionInputSchema).default([])
	})
	.superRefine((data, ctx) => {
		if (data.availabilityStart && Number.isNaN(Date.parse(data.availabilityStart))) {
			ctx.addIssue({
				code: z.ZodIssueCode.custom,
				path: ['availabilityStart'],
				message: 'รูปแบบวันที่ไม่ถูกต้อง'
			});
		}
		if (data.availabilityEnd && Number.isNaN(Date.parse(data.availabilityEnd))) {
			ctx.addIssue({
				code: z.ZodIssueCode.custom,
				path: ['availabilityEnd'],
				message: 'รูปแบบวันที่ไม่ถูกต้อง'
			});
		}
		if (data.availabilityStart && data.availabilityEnd) {
			const start = Date.parse(data.availabilityStart);
			const end = Date.parse(data.availabilityEnd);
			if (!Number.isNaN(start) && !Number.isNaN(end) && end < start) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					path: ['availabilityEnd'],
					message: 'วันที่สิ้นสุดต้องอยู่หลังวันที่เริ่มฉาย'
				});
			}
		}
	});

export type CreateMovieFormValues = z.infer<typeof createMovieFormSchema>;
