import { Link, useNavigate } from "react-router-dom";
import { LogoTextIcon } from "../../../components/icons/Icons";
import { HomeButton } from "../../../components/buttons/HomeButton";
import MobileNav from "./MobileNav";
import { useTranslation } from "react-i18next";
import LanguageSelector from "../../../components/LanguageSelector";
const Navigation = () => {
  const navigate = useNavigate();
  const { t }: { t: any } = useTranslation();
  return (
    <>
      <nav className="w-full hidden md:flex items-center justify-between p-[15px_25px] lg:p-[15px_67px]">
        <LogoTextIcon />
        <section className="flex items-center gap-x-[77px] ">
          <div className="flex items-center gap-x-[34px]">
            <LanguageSelector />
            <Link to={"/"} className="text-[1rem] text-black-2 leading-[22px]">
              {t("home")}
            </Link>
            <Link
              to={"/about-greybox"}
              className="text-[1rem] text-black-2 leading-[22px]"
            >
              {t("about")}
            </Link>
            <Link
              to={"/contact"}
              className="text-[1rem] text-black-2 leading-[22px]"
            >
              {t("contact")}
            </Link>
          </div>

          <HomeButton
            label={`${t("getStarted")}`}
            onClick={() => navigate("/sign-up")}
            extraClass="text-white bg-orange-1 w-[197px]"
          />
        </section>
      </nav>
      <MobileNav />
    </>
  );
};

export default Navigation;
