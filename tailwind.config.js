/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{templ,html,js}"],
  theme: {
    extend: {
      transitionProperty: {
        height: "height",
        width: "width",
      },
      colors: {
        text: "#2F2F2F",
        hover: "#FFC000",
        accent: "#FFC000",
        delete: "#FF5F5F",
        border: "#D7DAE2",
        error: "#FF5F5F",
        "submit-disabled": "#D9D9D9",
        "delete-disabled": "#EE6460",
      },
    },
  },
  plugins: [],
};
