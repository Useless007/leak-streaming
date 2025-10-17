'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { AlertCircle, Loader2, RefreshCw, Subtitles } from 'lucide-react';
import { Button } from '@/components/ui/button';
import type { Caption, Movie } from '@/lib/api';
import { usePlayback } from '@/lib/hooks/usePlayback';

type PlayerProps = {
  movie: Movie;
};

export function MoviePlayer({ movie }: PlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const { state, retry, activeCaption, selectCaption } = usePlayback({
    slug: movie.slug,
    captions: movie.captions ?? []
  });
  const [playerError, setPlayerError] = useState<string | null>(null);

  useEffect(() => {
    const video = videoRef.current;
    if (!video || !video.textTracks) {
      return;
    }

    Array.from(video.textTracks).forEach((track) => {
      if (activeCaption && track.language === activeCaption) {
        track.mode = 'showing';
      } else {
        track.mode = 'disabled';
      }
    });
  }, [activeCaption, state.status]);

  const captionOptions = useMemo(() => movie.captions ?? [], [movie.captions]);

  useEffect(() => {
    setPlayerError(null);
  }, [state.token]);

  useEffect(() => {
    const video = videoRef.current;
    if (!video || state.status !== 'ready' || !state.src) {
      return;
    }

    let destroyed = false;
    let hlsInstance: { destroy: () => void } | null = null;

    const setupPlayback = async () => {
      if (!video) {
        return;
      }

      const canPlayHLSNatively = video.canPlayType('application/vnd.apple.mpegurl');
      if (canPlayHLSNatively) {
        video.src = state.src;
        video.load();
        try {
          await video.play();
        } catch {
          // ignore autoplay rejection; user can press play manually
        }
        return;
      }

      const { default: Hls } = await import('hls.js');
      if (destroyed) {
        return;
      }

      if (!Hls.isSupported()) {
        setPlayerError('เบราว์เซอร์ไม่รองรับการเล่นสตรีม HLS');
        return;
      }

      const instance = new Hls({
        enableWorker: true,
        lowLatencyMode: true,
        backBufferLength: 90
      });
      hlsInstance = instance;

      instance.attachMedia(video);
      instance.on(Hls.Events.MEDIA_ATTACHED, () => {
        if (!destroyed) {
          instance.loadSource(state.src);
        }
      });
      instance.on(Hls.Events.ERROR, (_event, data) => {
        if (destroyed || !data?.fatal) {
          return;
        }
        instance.destroy();
        hlsInstance = null;
        setPlayerError('ไม่สามารถเล่นสตรีมได้');
      });
    };

    void setupPlayback();

    return () => {
      destroyed = true;
      if (hlsInstance) {
        hlsInstance.destroy();
        hlsInstance = null;
      }
      if (video) {
        video.removeAttribute('src');
        video.load();
      }
    };
  }, [state.status, state.src, state.token]);

  return (
    <div className="space-y-4">
      <div className="relative aspect-video overflow-hidden rounded-3xl border border-border bg-black shadow">
        {state.status === 'loading' && (
          <div className="flex h-full items-center justify-center gap-2 text-muted-foreground">
            <Loader2 className="size-5 animate-spin" aria-hidden />
            <span>กำลังเตรียมสตรีม...</span>
          </div>
        )}
        {state.status === 'error' && (
          <div className="flex h-full flex-col items-center justify-center gap-3 text-center text-muted-foreground">
            <AlertCircle className="size-6 text-destructive" aria-hidden />
            <p>{state.error ?? 'ไม่สามารถเริ่มเล่นได้'}</p>
            <Button
              variant="outline"
              onClick={() => {
                void retry();
              }}
              className="inline-flex items-center gap-2"
            >
              <RefreshCw className="size-4" aria-hidden />
              ลองอีกครั้ง
            </Button>
          </div>
        )}
        {state.status === 'ready' && state.src && (
          <>
            <video
              key={state.token}
              ref={videoRef}
              className="size-full object-cover"
              poster={movie.posterUrl}
              controls
              autoPlay
              playsInline
              preload="auto"
              onError={() => {
                setPlayerError('ไม่สามารถเล่นสตรีมได้');
              }}
            >
              {(movie.captions ?? []).map((caption) => (
                <track
                  key={caption.languageCode}
                  label={caption.label}
                  kind="subtitles"
                  srcLang={caption.languageCode}
                  src={caption.captionUrl}
                  default={caption.languageCode === activeCaption}
                />
              ))}
              เบราว์เซอร์ของคุณไม่รองรับการเล่นวิดีโอ HTML5
            </video>
            {playerError && (
              <div className="absolute inset-0 flex flex-col items-center justify-center gap-3 bg-black/70 text-center text-muted-foreground">
                <AlertCircle className="size-6 text-destructive" aria-hidden />
                <p>{playerError}</p>
                <Button
                  variant="outline"
                  onClick={() => {
                    setPlayerError(null);
                    void retry();
                  }}
                  className="inline-flex items-center gap-2"
                >
                  <RefreshCw className="size-4" aria-hidden />
                  ลองอีกครั้ง
                </Button>
              </div>
            )}
          </>
        )}
        {state.status === 'ready' && !state.src && (
          <div className="flex h-full flex-col items-center justify-center gap-3 text-center text-muted-foreground">
            <AlertCircle className="size-6 text-destructive" aria-hidden />
            <p>ไม่สามารถสร้าง URL สำหรับเล่นสตรีมได้</p>
            <Button
              variant="outline"
              onClick={() => {
                void retry();
              }}
              className="inline-flex items-center gap-2"
            >
              <RefreshCw className="size-4" aria-hidden />
              ลองอีกครั้ง
            </Button>
          </div>
        )}
      </div>

      {captionOptions.length > 0 && (
        <div className="flex flex-wrap items-center gap-3">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Subtitles className="size-4" aria-hidden />
            คำบรรยาย:
          </div>
          <Button
            variant={activeCaption === null ? 'default' : 'outline'}
            size="sm"
            onClick={() => selectCaption(null)}
          >
            ปิด
          </Button>
          {captionOptions.map((caption: Caption) => (
            <Button
              key={caption.languageCode}
              variant={activeCaption === caption.languageCode ? 'default' : 'outline'}
              size="sm"
              onClick={() => selectCaption(caption.languageCode)}
            >
              {caption.label}
            </Button>
          ))}
        </div>
      )}
    </div>
  );
}
