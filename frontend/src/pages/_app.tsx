import type { AppProps } from 'next/app';
import { Space_Mono, Inter, Geist_Mono } from 'next/font/google';
import '@/app/globals.css';

const spaceMono = Space_Mono({
  variable: '--font-space-mono',
  subsets: ['latin'],
  weight: ['400', '700'],
});

const inter = Inter({
  variable: '--font-inter',
  subsets: ['latin'],
});

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
});

export default function App({ Component, pageProps }: AppProps) {
  const fontClasses = `${spaceMono.variable} ${inter.variable} ${geistMono.variable}`;
  return (
    <div className={fontClasses}>
      <Component {...pageProps} />
    </div>
  );
}
