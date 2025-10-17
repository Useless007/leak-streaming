'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import type { Caption } from '@/lib/api';

export type PlaybackStatus = 'idle' | 'loading' | 'ready' | 'error';

type UsePlaybackArgs = {
  slug: string;
  captions: Caption[];
};

type PlaybackState = {
  status: PlaybackStatus;
  token?: string;
  src?: string;
  error?: string;
};

type UsePlaybackResult = {
  state: PlaybackState;
  retry: () => Promise<void>;
  activeCaption: string | null;
  selectCaption: (languageCode: string | null) => void;
};

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080';

export function usePlayback({ slug, captions }: UsePlaybackArgs): UsePlaybackResult {
  const [state, setState] = useState<PlaybackState>({ status: 'idle' });
  const [activeCaption, setActiveCaption] = useState<string | null>(captions[0]?.languageCode ?? null);

  const sourceBuilder = useCallback(
    (token: string | undefined, url: string | undefined) => {
      if (url) {
        return url;
      }
      if (!token) {
        return '';
      }
      const encodedToken = encodeURIComponent(token);
      return `${API_BASE_URL}/movies/${slug}/manifest.m3u8?token=${encodedToken}`;
    },
    [slug]
  );

  const requestToken = useCallback(async () => {
    setState({ status: 'loading' });
    try {
      const response = await fetch(`${API_BASE_URL}/movies/${slug}/playback-token`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        cache: 'no-store'
      });

      if (!response.ok) {
        const message = await response.text();
        throw new Error(message || 'ไม่สามารถขอโทเคนสำหรับเล่นได้');
      }

      const json = (await response.json()) as { token?: string; url?: string };
      if (!json?.token && !json?.url) {
        throw new Error('ข้อมูลสตรีมไม่ถูกต้อง');
      }

      const src = sourceBuilder(json.token, json.url);
      if (!src) {
        throw new Error('ไม่สามารถสร้าง URL สำหรับสตรีมได้');
      }

      setState({ status: 'ready', token: json.token, src });
    } catch (error) {
      setState({
        status: 'error',
        error: error instanceof Error ? error.message : 'เกิดข้อผิดพลาดไม่คาดคิด'
      });
    }
  }, [slug, sourceBuilder]);

  useEffect(() => {
    void requestToken();
  }, [requestToken]);

  const retry = useCallback(async () => {
    await requestToken();
  }, [requestToken]);

  const selectCaption = useCallback((languageCode: string | null) => {
    setActiveCaption(languageCode);
  }, []);

  return useMemo(
    () => ({
      state,
      retry,
      activeCaption,
      selectCaption
    }),
    [state, retry, activeCaption, selectCaption]
  );
}
