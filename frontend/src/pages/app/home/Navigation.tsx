import { Link, useNavigate } from "react-router-dom";
import { LogoTextIcon } from "../../../components/icons/Icons";
import { HomeButton } from "../../../components/buttons/HomeButton";
import MobileNav from "./MobileNav";

const Navigation = () => {
  const navigate = useNavigate();
  return (
    <>
      <nav className="w-full hidden md:flex items-center justify-between p-[15px_25px] lg:p-[15px_67px]">
        <LogoTextIcon />
        <section className="flex items-center gap-x-[77px] ">
          <div className="flex items-center gap-x-[34px]">
            <Link to={"/"} className="text-[1rem] text-black-2 leading-[22px]">
              Home
            </Link>
            <Link
              to={"/about-greybox"}
              className="text-[1rem] text-black-2 leading-[22px]"
            >
              About
            </Link>
            <Link
              to={"/contact"}
              className="text-[1rem] text-black-2 leading-[22px]"
            >
              Contact
            </Link>
          </div>

          <HomeButton
            label="Get Started"
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
