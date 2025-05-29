import { useNavigate } from "react-router-dom";
import Navigation from "./home/Navigation";
import Footer from "./home/Footer";
import { AboutLadyFlat, LightIcon } from "../../components/icons/Icons";
import { HomeButton } from "../../components/buttons/HomeButton";
import { useScrollToTop } from "../../utils/ScrollToTop";
import { useTranslation } from "react-i18next";
import LanguageSelector from "../../components/LanguageSelector";

const About = () => {
  const navigate = useNavigate();
  const { t }: { t: any } = useTranslation();
  const AboutCard = ({
    title,
    text,
    other,
    last,
  }: {
    title: string;
    text: string;
    other?: React.ReactNode;
    last?: boolean;
  }) => {
    return (
      <section
        className={`flex gap-x-[16px] pb-[47px] items-start ${
          !last && "border-b-[1px] border-grey-1"
        }`}
      >
        <LightIcon />
        <div>
          <h2 className="text-[1.5rem] text-black-2 font-[600] leading-[24px] mb-[16px]">
            {title}
          </h2>
          <p className="text-[1rem] text-black-3 leading-[22px] max-w-[517px] ">
            {text}
          </p>
          {other}
        </div>
      </section>
    );
  };
  useScrollToTop();
  return (
    <>
      {" "}
      <LanguageSelector />
      <section>
        <section className=" bg-pink-1 pb-[69px] flex flex-col">
          <Navigation />
          <div className="flex flex-col items-center px-[25px]">
            <h2 className="text-center text-[2rem] text-black-2 font-[700] leading-[40.63px] mt-[64px] mb-[16px]">
              {t("aboutGreybox")}
            </h2>
            <p className="text-center text-[1rem] text-black-3 leading-[22px] max-w-[725px]">
              {t("aboutGreyboxTitle")}
            </p>
          </div>
        </section>
        <section className="flex flex-col md:flex-row px-[25px] md:px-[5%] lg:px-[10%] py-[81px] gap-x-[2%] gap-y-[74px]">
          <AboutLadyFlat />
          <div className="flex flex-col gap-y-[47px]">
            <AboutCard
              title={t("ourMission")}
              text={t("ourMissionTitle")}
              other={
                <HomeButton
                  label={t("getStarted")}
                  onClick={() => navigate("/sign-up")}
                  extraClass="bg-orange-1 text-white w-[166px] mt-[32px]"
                />
              }
            />
            <AboutCard title={t("ourVision")} text={t("OurVisionTitle")} />
            <AboutCard title={t("ourGoal")} text={t("ourGoalTitle")} last />
          </div>
        </section>
        <Footer />
      </section>
    </>
  );
};

export default About;
