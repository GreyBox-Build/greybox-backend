import { Link, useNavigate } from "react-router-dom";
import { HomeButton } from "../../../components/buttons/HomeButton";
import {
  FacebookIcon,
  InstagramIcon,
  LocationIcon,
  LogoTextIcon,
  MailFIcon,
  PhoneIcon,
  TwitterIcon,
  WhatsappIcon,
} from "../../../components/icons/Icons";
import { FaLinkedin, FaFacebookSquare, FaWhatsapp } from "react-icons/fa";
import { RiInstagramFill } from "react-icons/ri";
import { FaSquareXTwitter } from "react-icons/fa6";
import { useTranslation } from "react-i18next";

const Footer = () => {
  const { t }: { t: any } = useTranslation();
  const navigate = useNavigate();
  return (
    <section className="min-h-[500px] bg-pink-1 pt-[100px] px-[25px] md:px-[5%] lg:px-[10%]">
      <section className="w-full p-[54px_25px] bg-orange-1 text-white rounded-[24px] flex flex-col items-center">
        <h2 className="text-center text-[2.5rem] font-[600] leading-[51.43px] mb-8">
          {t("getStartedWithGreyBox")}
        </h2>

        <HomeButton
          label={t("getStarted")}
          onClick={() => navigate("/sign-up")}
          extraClass="border-white w-[204px]"
        />
      </section>

      <section className="w-full flex justify-between flex-wrap gap-[50px] mt-[104px]">
        <div>
          <LogoTextIcon />
          {/* <p className="text-[1rem] text-black-3 leading-[22px] max-w-[218px]">
            Lorem Ipsum is a dummy text,Lorem Ipsum
          </p> */}
          <div className="flex items-center gap-x-[24px] mt-[14px]">
            <a href="https://instagram.com/greybox_">
              <RiInstagramFill size={30} />
            </a>
            <a href="https://linkedin.com/company/greyboxpay">
              <FaLinkedin size={30} />
            </a>
            <a href="https://facebook.com/profile.php?id=100089993932799">
              <FaFacebookSquare size={30} />
            </a>
            <a href="https://x.com/@greyboxpay">
              <FaSquareXTwitter size={30} />
            </a>
          </div>
        </div>

        <div className="flex flex-col">
          <h3 className="text-[1.25rem] text-black-2 font-[600] leading-[20px] mb-[10px]">
            Useful Links
          </h3>
          <Link to={"/"} className="text-[1rem] text-black-3 leading-[22px]">
            Home
          </Link>
          <Link
            to={"/about-greybox"}
            className="text-[1rem] text-black-3 leading-[22px]"
          >
            About
          </Link>
          <Link
            to={"/contact"}
            className="text-[1rem] text-black-3 leading-[22px]"
          >
            Contact
          </Link>
        </div>

        <div className="flex flex-col">
          <h3 className="text-[1.25rem] text-black-2 font-[600] leading-[20px] mb-[10px]">
            Explore
          </h3>
          <Link to={""} className="text-[1rem] text-black-3 leading-[22px]">
            Terms & Condition
          </Link>
          <Link
            to={"/privacy-policy"}
            className="text-[1rem] text-black-3 leading-[22px]"
          >
            Privacy Policy
          </Link>
          <Link
            to={"/aml-policy"}
            className="text-[1rem] text-black-3 leading-[22px]"
          >
            AML Policy
          </Link>
        </div>

        <div>
          <h3 className="text-[1.25rem] text-black-2 font-[600] leading-[20px] mb-[10px]">
            Contact
          </h3>

          <div className="flex text-[1rem] text-black-3 leading-[22px] gap-x-[8px]">
            <MailFIcon /> info@greyboxpay.com
          </div>
          <div className="flex text-[1rem] text-black-3 leading-[22px] gap-x-[8px]">
            <PhoneIcon /> +233 2022680388
          </div>
        </div>
      </section>

      <section className="text-center text-black-3 text-[0.875rem] leading-[18px] p-[36px_25px]">
        Copyright © {new Date().getFullYear()}. All rights reserved.
      </section>
    </section>
  );
};

export default Footer;
