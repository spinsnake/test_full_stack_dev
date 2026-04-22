/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        ink: '#1d1d1f',
        paper: '#f5f5f7',
        sand: '#d2d2d7',
        clay: '#0071e3',
        ember: '#2997ff',
        olive: '#6e6e73',
        bronze: '#ffffff',
      },
      boxShadow: {
        card: '0 18px 44px rgba(0, 0, 0, 0.08)',
      },
      backgroundImage: {
        grain:
          "radial-gradient(circle at top, rgba(0, 113, 227, 0.12) 0, rgba(0, 113, 227, 0) 26%), radial-gradient(circle at 82% 12%, rgba(255,255,255,0.96) 0, rgba(255,255,255,0) 32%), linear-gradient(180deg, rgba(255,255,255,0.85) 0%, rgba(245,245,247,0) 58%)"
      }
    },
  },
  plugins: [],
};
