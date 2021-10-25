const defaultTheme = require('tailwindcss/defaultTheme');

module.exports = {
  purge: ['./src/**/*.{js,jsx,ts,tsx}', './public/index.html'],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        primary: '#2F1160',
        primary_neon: '#6200FF',
        dunkelgrau: '#B9B7BD', 
        grau: '#EBE9F0',
        hellgrau: '#FBFAFC',
      },
      fontFamily:{
          sans: ['Poppins', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
