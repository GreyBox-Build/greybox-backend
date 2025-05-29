import { useState } from "react";
import { useTranslation } from "react-i18next";
import { RxCaretDown } from "react-icons/rx";
const LanguageSelector = () => {
  const { i18n }: { t: any; i18n: any } = useTranslation();

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };
  const [showLanguage, setShowLanguage] = useState(false);

  return (
    <>
      <div className="flex gap-3   relative">
        <p
          className="flex gap-1 items-center cursor-pointer"
          onClick={() => setShowLanguage(!showLanguage)}
          title="Translate (English/French"
        >
          <span>{i18n.language.toUpperCase()} </span> <RxCaretDown />
        </p>

        {showLanguage && (
          <div className="bg-orange-1 flex flex-col rounded-sm absolute text-white  -left-5 top-8">
            {" "}
            <button
              onClick={() => {
                changeLanguage("en");
                setShowLanguage(!showLanguage);
              }}
              title="English"
              className=" py-2 px-6 "
            >
              EN
            </button>
            <button
              className=" py-2  px-6 border-t"
              onClick={() => {
                changeLanguage("fr");
                setShowLanguage(!showLanguage);
              }}
            >
              FR
            </button>{" "}
          </div>
        )}
      </div>
    </>
  );
};

export default LanguageSelector;
