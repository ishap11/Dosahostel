module.exports = {
    plugins: [
      require('@tailwindcss/postcss')({
        // Optional Tailwind config path if it's not the default
        config: './tailwind.config.js',
      }),
      require('autoprefixer'),
    ],
  }
  