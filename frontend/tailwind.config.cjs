/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        ink: '#151312',
        paper: '#f3ece3',
        sand: '#d8c6b1',
        clay: '#b5623f',
        ember: '#d87b4f',
        olive: '#4b5945',
        bronze: '#8f6a47',
      },
      boxShadow: {
        card: '0 24px 40px rgba(21, 19, 18, 0.14)',
      },
      backgroundImage: {
        grain:
          "radial-gradient(circle at 20% 20%, rgba(255,255,255,0.18) 0, rgba(255,255,255,0) 38%), radial-gradient(circle at 80% 0%, rgba(216, 123, 79, 0.16) 0, rgba(216, 123, 79, 0) 30%), radial-gradient(circle at 10% 80%, rgba(75, 89, 69, 0.12) 0, rgba(75, 89, 69, 0) 32%)"
      }
    },
  },
  plugins: [],
};
