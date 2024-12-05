import { useNavigate } from "react-router-dom";
import { HomeButton } from "../../../components/buttons/HomeButton";
import { useTranslation } from "react-i18next";

const GetStarted = () => {
  const navigate = useNavigate();
  const { t }: { t: any } = useTranslation();
  return (
    <section className=" bg-get-started-bg bg-cover bg-no-repeat min-h-[346px]  px-[25px] py-[60px] md:px-[5%] lg:px-[10%]">
      <div className="text-[2rem] text-black-3 max-w-[654px] font-[600] leading-[41.14px] mb-[96px]">
        {t("sendfundTitle1")}
        <span className="text-orange-1">{t("africa")}</span>{" "}
        {t("sendfundTitle2")}
      </div>

      <HomeButton
        label={t("getStarted")}
        onClick={() => navigate("/sign-up")}
        extraClass="w-full md:w-[353px] bg-orange-1 text-white"
      />
    </section>
  );
};

export default GetStarted;
