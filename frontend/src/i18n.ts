import i18n from "i18next";
import { initReactI18next } from "react-i18next";

export const resources = {
  en: {
    translation: {
      welcome: "Welcome to the app!",
      description: "This is an example description.",
    },
  },
  fr: {
    translation: {
      welcome: "Bienvenue dans l'application!",
      description: "Ceci est un exemple de description.",
    },
  },
};

i18n.use(initReactI18next).init({
  missingKeyHandler: (lng, ns, key) => {
    console.warn(
      `Missing key "${key}" in namespace "${ns}" for language "${lng}"`
    );
  },
  resources,
  lng: "en",
  debug: true, // Enable debug mode
  fallbackLng: "en",
  interpolation: {
    escapeValue: false,
  },
});

export default i18n;
