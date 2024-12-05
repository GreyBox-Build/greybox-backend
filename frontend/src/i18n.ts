// src/i18n.ts
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import HttpApi from "i18next-http-backend";
import LanguageDetector from "i18next-browser-languagedetector";

// i18Next initialization
i18n
  .use(HttpApi) // Load translations from external files
  .use(LanguageDetector) // Detect user's language
  .use(initReactI18next) // Bind i18Next to React
  .init({
    fallbackLng: "en", // Fallback language if detection fails
    supportedLngs: ["en", "fr", "es"], // List of supported languages
    debug: true, // Enable console debugging
    interpolation: {
      escapeValue: false, // React already handles XSS
    },
    backend: {
      loadPath: "/locales/{{lng}}/{{ns}}.json", // Translation files location
    },
  });

export default i18n;
