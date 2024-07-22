import { HomeButton } from "../../../components/buttons/HomeButton";

const Services = () => {
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
      <h2 className="text-orange-1 text-center text-[1.125rem] font-[700] leading-[24.75px]">
        OUR SERVICES
      </h2>
      <h2 className=" text-center text-[2rem] text-black-2 font-[700] leading-[40.63px] mt-[2px]">
        What Solutions We Offer You
      </h2>
      <section className="flex flex-col md:flex-row items-center md:justify-center md:items-end mt-[58px] gap-y-[25px]">
        <SolCard
          sn="01."
          title="Low cost cross-border solutions"
          subText="Send cash to your loved ones or business partners across the world"
          extraClass="bg-[#F5D8CC] text-black-2 min-h-[289px] rounded-[8px] md:rounded-[8px_0px_0px_8px]"
        />
        <SolCard
          sn="02."
          title="Currency Stability"
          subText="We leverage stablecoins, pegged to major fiat USD, to ensure your transactions are stable and predictable."
          extraClass="bg-orange-1 text-white min-h-[349px] rounded-[8px] md:rounded-[8px_8px_0px_0px]"
          btn={
            <HomeButton
              label="Learn More"
              onClick={() => {}}
              extraClass="text-white border-white w-[204px] mt-[42px]"
            />
          }
        />
        <SolCard
          sn="03."
          title="Remittances"
          subText="Send money to your loved ones in Africa"
          extraClass="bg-[#F5D8CC] text-black-2 min-h-[289px] rounded-[8px] md:rounded-[0px_8px_8px_0px] "
        />
      </section>
    </section>
  );
};

export default Services;
