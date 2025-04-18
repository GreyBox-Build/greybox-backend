import { HomeButton } from "../../../components/buttons/HomeButton";
import { useTranslation } from "react-i18next";
const Services = () => {
  const { t }: { t: any } = useTranslation();
  const SolCard = ({
    sn,
    title,
    subText,
    btn,
    extraClass,
  }: {
    sn: string;
    title: string;
    subText?: string;
    btn?: React.ReactNode;
    extraClass: string;
  }) => {
    return (
      <div
        className={`p-[43px_24px]  border-[1px] w-full md:w-[383px]  ${extraClass}`}
      >
        <p className="font-[700] text-[2rem] leading-[40.63px]">{sn}</p>
        <h3 className="font-[600] text-[1.25rem] leading-[27.5px] mt-[10px]">
          {title}
        </h3>
        <p className="text-[1rem] leading-[22px] mt-[15px] max-w-[328px]">
          {subText}
        </p>
        {btn}
      </div>
    );
  };
  return (
    <section className=" pb-[100px] bg-white px-[25px] md:px-[5%] lg:px-[10%]">
      <h2 className="text-orange-1 text-center text-[1.225rem] md:text-[2.225rem] font-[700] leading-[24.75px]">
        {t("ourServices")}
      </h2>

      <section className="flex flex-col md:flex-row items-center md:justify-center md:items-end mt-[58px] gap-y-[25px]">
        <SolCard
          sn="01."
          title={`${t("sendCash")}`}
          subText={`${t("sendCashTitle")}`}
          extraClass="bg-[#F5D8CC] text-black-2 min-h-[289px] rounded-[8px] md:rounded-[8px_0px_0px_8px]"
        />
        <SolCard
          sn="02."
          title={`${t("currencyStability")}`}
          subText={`${t("currencyStabilityTitle")}`}
          extraClass="bg-orange-1 text-white min-h-[349px] rounded-[8px] md:rounded-[8px_8px_0px_0px]"
          btn={
            <HomeButton
              label={t("learnMore")}
              onClick={() => {}}
              extraClass="text-white border-white w-[204px] mt-[42px]"
            />
          }
        />
        <SolCard
          sn="03."
          title={`${t("remittances")}`}
          subText={`${t("remittancesTitle")}`}
          extraClass="bg-[#F5D8CC] text-black-2 min-h-[289px] rounded-[8px] md:rounded-[0px_8px_8px_0px] "
        />
      </section>
    </section>
  );
};

export default Services;
