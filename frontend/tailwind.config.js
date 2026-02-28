/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        // Optional: semantic game colors
        primary: {
          DEFAULT: "#2563eb",
          dark: "#1d4ed8",
        },
        accent: {
          DEFAULT: "#10b981",
          dark: "#059669",
        },
      },
    },
  },
  plugins: [],
};
