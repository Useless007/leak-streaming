import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { ThemeProvider } from '@/components/providers/theme-provider';
import { SiteHeader } from '@/components/layout/site-header';
import './globals.css';

const fontSans = Inter({
  subsets: ['latin'],
  variable: '--font-sans',
  display: 'swap'
});

export const metadata: Metadata = {
  title: {
    default: 'Leak Streaming Portal',
    template: '%s · Leak Streaming Portal'
  },
  description:
    'แพลตฟอร์มสตรีมมิ่งที่ออกแบบมาสำหรับการเผยแพร่และจัดการภาพยนตร์ผ่าน Next.js 15 และบริการ Go backend.'
};

type RootLayoutProps = {
  children: React.ReactNode;
};

export default function RootLayout({ children }: RootLayoutProps) {
	return (
		<html lang="th" suppressHydrationWarning>
			<body className={`${fontSans.variable} min-h-screen bg-background font-sans text-foreground antialiased`}>
				<ThemeProvider>
					<div className="flex min-h-screen flex-col">
						<SiteHeader />
						<main className="flex-1 pb-16 pt-10">{children}</main>
					</div>
				</ThemeProvider>
			</body>
		</html>
	);
}
