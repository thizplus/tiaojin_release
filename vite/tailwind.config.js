/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        'db': ['DB', 'system-ui', 'sans-serif'],
        'sans': ['Google Sans', 'Noto Sans Thai', 'system-ui', 'sans-serif'],
        'google-sans': ['Google Sans', 'sans-serif'],
        'noto-thai': ['Noto Sans Thai', 'sans-serif'],
      },
      fontSize: {
        'base': ['18px', '1.6'], // เปลี่ยน base จาก 16px เป็น 18px
      },
      keyframes: {
        marquee: {
          '0%': { transform: 'translateX(0)' },
          '100%': { transform: 'translateX(-50%)' },
        },
      },
      animation: {
        marquee: 'marquee 30s linear infinite',
      },
    },
  },
  plugins: [],
}