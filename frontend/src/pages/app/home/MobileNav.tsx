import { useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { Close, Menu } from "@mui/icons-material";
import { LogoTextIcon } from "../../../components/icons/Icons";
import { HomeButton } from "../../../components/buttons/HomeButton";
import LanguageSelector from "../../../components/LanguageSelector";
import { useTranslation } from "react-i18next";

const MobileNav = () => {
  const [isOpen, setIsOpen] = useState(false);
  const navigate = useNavigate();
  const { t }: { t: any } = useTranslation();

  const navigation = [
    {
      name: t("home"),
      href: "/",
    },
    {
      name: t("about"),
      href: "/about-greybox",
    },
    {
      name: t("contact"),
      href: "/contact",
    },
  ];

  return (
    <nav className="block md:hidden relative">
      <section className="flex items-center justify-between p-[15px_25px]">
        <LogoTextIcon />

        <div className="flex items-center gap-5">
          {" "}
          <LanguageSelector />
          <div onClick={() => setIsOpen(true)} className="cursor-pointer ">
            <Menu className="text-black-1" />
          </div>
        </div>
      </section>

      <section
        className={`${
          isOpen
            ? "right-0 bottom-0 top-0 left-0"
            : "left-full right-0 bottom-0 top-0"
        } fixed w-full flex flex-col justify-center items-center transition-all
           duration-300 overflow-hidden bg-black-3
           `}
      >
        <div className="cursor-pointer absolute top-[15px] right-[25px] text-[#fefefe]">
          <Close onClick={() => setIsOpen(false)} />
        </div>

        {navigation.map((item, index) => {
          return (
            <li key={index} className="mb-8 text-white list-none">
              <NavLink
                onClick={() => {
                  setIsOpen(false);
                }}
                style={({ isActive }) =>
                  isActive
                    ? {
                        textDecoration: "none",
                        color: "#CD5928",
                      }
                    : {}
                }
                to={item.href}
                className=" transition-all duration-300 ring-offset-[-70]"
              >
                {item.name}
              </NavLink>
            </li>
          );
        })}
        <HomeButton
          label={`${t("getStarted")}`}
          onClick={() => navigate("/sign-up")}
          extraClass="text-white bg-orange-1 w-[197px]"
        />
      </section>
    </nav>
  );
};

export default MobileNav;
