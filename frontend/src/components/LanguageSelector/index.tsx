import { useState } from "react";
import { useTranslation } from "react-i18next";
import { MdLanguage } from "react-icons/md";
import { FaLanguage } from "react-icons/fa";
const LanguageSelector = () => {
  const { i18n }: { t: any; i18n: any } = useTranslation();

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };
  const [showLanguage, setShowLanguage] = useState(false);

  return (
    <>
      <div className="flex gap-3 fixed left-2 top-20 z-30  ">
        <MdLanguage
          // size={30}
          className="text-orange-1 cursor-pointer border border-white text-4xl  bg-white rounded-full "
          onClick={() => setShowLanguage(!showLanguage)}
          title="Translate (English/French"
        />

        {showLanguage && (
          <div className="bg-orange-1 flex flex-col rounded-sm absolute text-white  left-2 top-12">
            {" "}
            <button
              onClick={() => {
                changeLanguage("en");
                setShowLanguage(!showLanguage);
              }}
              title="English"
              className=" py-2 px-6 "
            >
              English
            </button>
            <button
              className=" py-2  px-6 border-t"
              onClick={() => {
                changeLanguage("fr");
                setShowLanguage(!showLanguage);
              }}
            >
              French
            </button>{" "}
          </div>
        )}
      </div>
    </>
  );
};

export default LanguageSelector;
