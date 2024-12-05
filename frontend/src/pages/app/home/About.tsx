import { HomeButton } from "../../../components/buttons/HomeButton";
import { AboutLady, FiveStars } from "../../../components/icons/Icons";
import { useTranslation } from "react-i18next";

const About = () => {
  const { t }: { t: any } = useTranslation();
  return (
    <>

      <section className="pb-[100px] bg-white px-[25px] md:px-[5%] lg:px-[10%]">
        <div className="w-full flex flex-col md:flex-row items-center gap-x-[2%]">
          <section className="w-full md:w-[50%] flex justify-center md:justify-start ">
            <AboutLady />
          </section>
          <section className="w-full md:w-[50%] flex flex-col md:block">
            <h2 className="text-[1.525rem]  text-orange-1 font-[700] leading-[24.75px]">
              {t("aboutUs")}
            </h2>
            <p className="text-[2rem] text-black-2 font-[700] leading-[40.63px] max-w-[456px] mt-[12px]">
              {t("aboutTitle")}
            </p>
            <div className="my-[26px]">
              <FiveStars />
            </div>
            <p className="text-[1rem] text-justify text-black-3 leading-[22px] max-w-[548px]">
              {t("aboutSubTitle")}
            </p>
            <HomeButton
              label={`${t("learnMore")}`}
              onClick={() => {}}
              extraClass="text-white bg-orange-1 border-white w-[204px] mt-[42px]"
            />
          </section>
        </div>
      </section>
    </>
  );
};

export default About;
