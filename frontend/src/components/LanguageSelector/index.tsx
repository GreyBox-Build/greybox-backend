import { useTranslation } from "react-i18next";

const LanguageSelector = () => {
  const { t, i18n }: { t: any; i18n: any } = useTranslation();

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };

  return (
    <div>
      <h1>{t("welcome")}</h1>
      <p>{t("description")}</p>
      <button onClick={() => changeLanguage("en")}>English</button>
      <button onClick={() => changeLanguage("fr")}>French</button>
      <button onClick={() => changeLanguage("es")}>Spanish</button>
    </div>
  );
};

export default LanguageSelector;
