/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,js,ts,vue}'],
  theme: {
    screens: {
      'sm': {'max': '600px'},
      'md': {'max': '1024px'},
      'lg': {'max': '1440px'},
      'xl': {'max': '1920px'},
      '2xl': {'min': '1920px'}
    },
    extend: {},
  },
  plugins: [],
};
