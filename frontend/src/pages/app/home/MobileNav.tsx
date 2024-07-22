import { useState } from "react";

import { motion } from "framer-motion";
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

  const circleVariants = {
    hidden: {
      scale: 0,
    },
    visible: {
      scale: 180,
      transition: {
        type: "spring",
        stiffness: 160,
        damping: 60,
      },
    },
  };

  const ulVariants = {
    hidden: {
      opacity: 0,
    },
    visible: {
      opacity: 1,
      transition: {
        delay: 0.1,
      },
    },
  };

  return (
    <nav className="block md:hidden relative">
      <section className="flex items-center justify-between p-[15px_25px]">
        <LogoTextIcon />
        <div onClick={() => setIsOpen(true)} className="cursor-pointer ">
          <Menu className="text-black-1" />
        </div>
      </section>

      <motion.div
        variants={circleVariants}
        initial="hidden"
        animate={isOpen ? "visible" : "hidden"}
        className="w-4 h-4 rounded-full bg-accent fixed top-0 right-0"
      ></motion.div>

      <motion.ul
        variants={ulVariants}
        initial="hidden"
        animate={isOpen ? "visible" : ""}
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
            <li key={index} className="mb-8 text-white">
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
      </motion.ul>
    </nav>
  );
};

export default MobileNav;
