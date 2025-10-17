import { z } from 'zod';
import {
  movieSchema,
  movieSummarySchema,
  playbackTokenSchema,
  streamSchema,
  captionSchema,
  type Movie,
  type MovieSummary,
  type Caption,
  type Stream,
  type PlaybackToken
} from './schemas';

const errorResponseSchema = z.object({
  error: z.string(),
  code: z.string().optional()
});

type FetcherOptions = {
  baseUrl?: string;
  headers?: HeadersInit;
};

type RequestOptions = FetcherOptions & {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  body?: unknown;
  cache?: RequestCache;
  next?: NextFetchRequestConfig;
};

type NextFetchRequestConfig = {
  revalidate?: number | false;
  tags?: string[];
};

async function resolveJSON<ResponseType>(response: Response, schema: z.ZodSchema<ResponseType>): Promise<ResponseType> {
  const data: unknown = await response.json();
  const parsed = schema.safeParse(data);
  if (!parsed.success) {
    throw new Error('ไม่สามารถตีความข้อมูลที่ได้จาก API ได้');
  }
  return parsed.data;
}

async function handleError(response: Response): Promise<never> {
  try {
    const payload: unknown = await response.json();
    const parsed = errorResponseSchema.parse(payload);
    const error = new Error(parsed.error);
    throw error;
  } catch {
    throw new Error(`คำขอ API ล้มเหลวด้วยสถานะ ${response.status}`);
  }
}

async function request<ResponseType>(
  endpoint: string,
  schema: z.ZodSchema<ResponseType>,
  options: RequestOptions = {}
): Promise<ResponseType> {
  const { baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? '', method = 'GET', body, headers, cache, next } = options;

  const requestInit: RequestInit = {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...headers
    },
    cache
  };
  if (next) {
    Object.assign(requestInit, { next });
  }

  if (body !== undefined) {
    requestInit.body = JSON.stringify(body);
  }

  const response = await fetch(new URL(endpoint, baseUrl).toString(), requestInit);
  if (!response.ok) {
    await handleError(response);
  }

  return resolveJSON(response, schema);
}

export function createApiClient(options: FetcherOptions = {}) {
  return {
    async listMovies(): Promise<MovieSummary[]> {
      return request('/movies', movieSummarySchema.array(), {
        ...options,
        cache: 'no-store'
      });
    },
    async getMovie(movieId: string): Promise<Movie> {
      return request(`/movies/${movieId}`, movieSchema, {
        ...options,
        cache: 'no-store'
      });
    },
    async listCaptions(movieId: string): Promise<Caption[]> {
      return request(`/movies/${movieId}/captions`, captionSchema.array(), {
        ...options,
        cache: 'no-store'
      });
    },
    async getStream(movieId: string): Promise<Stream> {
      return request(`/movies/${movieId}/stream`, streamSchema, {
        ...options,
        cache: 'no-store'
      });
    },
    async createPlaybackToken(movieId: string): Promise<PlaybackToken> {
      return request(
        `/movies/${movieId}/playback-token`,
        playbackTokenSchema,
        {
          ...options,
          method: 'POST'
        }
      );
    }
  };
}

export type ApiClient = ReturnType<typeof createApiClient>;
