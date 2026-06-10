/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f5f3ff',
          100: '#edd8ff',
          200: '#ddbaff',
          300: '#c28cff',
          400: '#a254ff',
          500: '#851cff',
          600: '#6d07eb',
          700: '#5700cf',
          800: '#4800b0',
          900: '#3c0091',
          950: '#230061',
        }
      },
      fontFamily: {
        sans: ['Outfit', 'Inter', 'system-ui', 'sans-serif'],
      },
      animation: {
        'spin-slow': 'spin 20s linear infinite',
      }
    },
  },
  plugins: [],
}
