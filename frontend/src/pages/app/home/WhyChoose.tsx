import {
  ComplianceIcon,
  CostIcon,
  SecurityIcon,
  StabilityIcon,
} from "../../../components/icons/Icons";
import { useTranslation } from "react-i18next";

const WhyChoose = () => {
  const { t }: { t: any } = useTranslation();

  const ReasonsCard = ({
    img,
    text,
    subText,
    extraClass,
  }: {
    img: React.ReactNode;
    text: string;
    subText: string;
    extraClass: string;
  }) => (
    <div
      className={`p-[53px_15px] w-full md:w-[50%] lg:w-[230px] xl:w-[260px] h-[222px] rounded-[8px] flex flex-col items-center ${extraClass}`}
    >
      {img}
      <p className="text-[1rem] text-black-2 text-center font-[700] leading-[22px] mt-[26px]">
        {text}
      </p>
      <p className="text-[1rem] text-black-3 text-center leading-[22px] mt-[4px]">
        {subText}
      </p>
    </div>
  );
  return (
    <section className=" pb-[100px] px-[25px] md:px-[66px] bg-white">
      <h2 className="text-[1.125rem] md:text-[2.225rem] text-orange-1 text-center font-[700] leading-[24.75px]">
        {t("whyGreyboxTitle")}
      </h2>

      <section className="flex flex-col lg:flex-row gap-[35px] mt-[30px] md:justify-center">
        <div className="flex flex-col md:flex-row gap-[35px]">
          <ReasonsCard
            img={<CostIcon />}
            text={`${t("costEffective")}`}
            subText={`${t("costEffectiveTit")}`}
            extraClass="bg-orange-3"
          />
          <ReasonsCard
            img={<StabilityIcon />}
            text={t("stability")}
            subText={t("stabilityTitle")}
            extraClass="bg-grey-1"
          />
        </div>
        <div className="flex flex-col md:flex-row gap-[35px]">
          <ReasonsCard
            img={<SecurityIcon />}
            text={t("security")}
            subText={t("securityTitle")}
            extraClass="bg-orange-3"
          />
          <ReasonsCard
            img={<ComplianceIcon />}
            text={t("speed")}
            subText={t("speedTitle")}
            extraClass="bg-grey-1"
          />
        </div>
      </section>
    </section>
  );
};

export default WhyChoose;
