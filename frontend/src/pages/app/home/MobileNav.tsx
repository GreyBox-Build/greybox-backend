import { useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { Close, Menu } from "@mui/icons-material";
import { LogoTextIcon } from "../../../components/icons/Icons";
import { HomeButton } from "../../../components/buttons/HomeButton";

const MobileNav = () => {
  const [isOpen, setIsOpen] = useState(false);
  const navigate = useNavigate();

  const navigation = [
    {
      name: "Home",
      href: "/",
    },
    {
      name: "About",
      href: "/about-greybox",
    },
    {
      name: "Contact",
      href: "/contact",
    },
  ];

  return (
    <nav className="block md:hidden relative">
      <section className="flex items-center justify-between p-[15px_25px]">
        <LogoTextIcon />
        <div onClick={() => setIsOpen(true)} className="cursor-pointer ">
          <Menu className="text-black-1" />
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
          label="Get Started"
          onClick={() => navigate("/sign-up")}
          extraClass="text-white bg-orange-1 w-[197px]"
        />
      </section>
    </nav>
  );
};

export default MobileNav;
