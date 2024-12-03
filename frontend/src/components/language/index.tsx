import { useTranslation } from "react-i18next";

const Example: React.FC = () => {
  const { t, i18n } = useTranslation();

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };

  return (
    <div>
      <h1>{t("welcome")}</h1>
      <p>{t("description")}</p>

      <button onClick={() => changeLanguage("en")}>English</button>
      <button onClick={() => changeLanguage("fr")}>French</button>
    </div>
  );
};

export default Example;
