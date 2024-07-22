/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      backgroundImage: {
        "grey-box-bg": "url('./publicImages/grey-box-bg.png')",
        "get-started-bg": "url('./publicImages/getStartedBg.png')",
        "round-grey-bg": "url('./publicImages/roundGreyBg.png')",
      },
      colors: {
        "grey-1": "#D9D9D9",
        "grey-2": "#4D4D4D",
        "grey-3": "#595959",
        "grey-4": "#B3B3B3",
        "grey-5": "#F2F2F2",
        "black-1": "#1B1A1A",
        "black-2": "#101010",
        "black-3": "#404040",
        "orange-1": "#CD5928",
        "orange-2": "#E49776",
        "orange-3": "#F5D8CC",
        "pink-1": "#FDF9F6",
        "success-bg": "#9EEDAA",
        "success-text": "#125814",
        "pending-bg": "#ECED9E",
        "pending-text": "#584012",
        "failed-bg": "#F1BABA",
        "failed-text": "#6D1313",
      },
      boxShadow: {
        "shadow-1": "0px 1px 1px 0px #0000000D",
      },
    },
    screens: {
      sm: "640px",

      md: "768px",

      lg: "1060px",

      xl: "1280px",

      "2xl": "1536px",
    },
  },
  plugins: [
    // ...
    require("tailwind-scrollbar"),
  ],
};
