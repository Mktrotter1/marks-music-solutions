import type { Config } from 'tailwindcss';

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        mms: {
          bg: '#000000',
          surface: '#111111',
          sidebar: '#0a0a0a',
          card: '#181818',
          hover: '#1a1a1a',
          border: '#222222',
          accent: '#00FFFF',
          'accent-hover': '#00E5E5',
          'accent-dim': '#00999980',
          'text-primary': '#ffffff',
          'text-secondary': '#999999',
          'text-tertiary': '#666666',
        },
      },
      fontFamily: {
        sans: [
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'Helvetica Neue',
          'Arial',
          'sans-serif',
        ],
      },
      spacing: {
        sidebar: '240px',
        player: '80px',
      },
    },
  },
  plugins: [],
} satisfies Config;
