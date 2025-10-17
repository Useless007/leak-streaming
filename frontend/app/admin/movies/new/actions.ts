'use server';

import { revalidatePath } from 'next/cache';

import { createMovieFormSchema, type CreateMovieFormValues } from './schema';

type ActionErrorMap = Record<string, string>;

export type CreateMovieActionResult =
	| { success: true; slug: string }
	| { success: false; fieldErrors?: ActionErrorMap; formError?: string };

export async function createMovieAction(values: CreateMovieFormValues): Promise<CreateMovieActionResult> {
	const parsed = createMovieFormSchema.safeParse(values);
	if (!parsed.success) {
		const fieldErrors: ActionErrorMap = {};
		for (const [field, messages] of Object.entries(parsed.error.flatten().fieldErrors)) {
			if (messages && messages.length > 0) {
				fieldErrors[field] = messages[0];
			}
		}
		return { success: false, fieldErrors };
	}

	const data = parsed.data;

	const allowedHosts = (data.allowedHosts ?? '')
		.split(/\r?\n/)
		.map((value) => value.trim())
		.filter((value) => value.length > 0);

	const availabilityStartISO = data.availabilityStart ? new Date(data.availabilityStart).toISOString() : undefined;
	const availabilityEndISO = data.availabilityEnd ? new Date(data.availabilityEnd).toISOString() : undefined;

	const payload = {
		title: data.title,
		synopsis: data.synopsis?.trim() ? data.synopsis.trim() : undefined,
		posterUrl: data.posterUrl?.trim() ? data.posterUrl.trim() : undefined,
		availabilityStart: availabilityStartISO,
		availabilityEnd: availabilityEndISO,
		isVisible: data.isVisible,
		streamUrl: data.streamUrl,
		drmKeyId: data.drmKeyId?.trim() ? data.drmKeyId.trim() : undefined,
		allowedHosts,
		captions: data.captions.map((caption) => ({
			languageCode: caption.languageCode,
			label: caption.label.trim(),
			captionUrl: caption.captionUrl
		}))
	};

	const sanitizedPayload = Object.fromEntries(
		Object.entries(payload).filter(([, value]) => value !== undefined)
	) as typeof payload;

	const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080';
	let response: Response;
	try {
		response = await fetch(`${baseUrl.replace(/\/$/, '')}/movies`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(sanitizedPayload)
		});
	} catch (error) {
		return {
			success: false,
			formError: error instanceof Error ? error.message : 'ไม่สามารถเชื่อมต่อ API ได้'
		};
	}

	if (!response.ok) {
		try {
			const raw: unknown = await response.json();
			if (raw && typeof raw === 'object') {
				const data = raw as { error?: string; details?: Record<string, string> };
				if (response.status === 422 && data.details) {
					return { success: false, fieldErrors: data.details };
				}
				return { success: false, formError: data.error ?? 'ไม่สามารถบันทึกข้อมูลได้' };
			}
		} catch {
			// ignore and fall through
		}
		return { success: false, formError: `คำขอถูกปฏิเสธ (${response.status})` };
	}

	let result: { slug?: string } | undefined;
	try {
		result = (await response.json()) as { slug?: string };
	} catch {
		return { success: false, formError: 'ไม่สามารถอ่านผลลัพธ์จาก API ได้' };
	}

	if (!result?.slug) {
		return { success: false, formError: 'API ไม่ได้ส่ง slug ของภาพยนตร์กลับมา' };
	}

	revalidatePath('/movies');

	return { success: true, slug: result.slug };
}
