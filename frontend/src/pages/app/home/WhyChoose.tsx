import {
  ComplianceIcon,
  CostIcon,
  SecurityIcon,
  StabilityIcon,
} from "../../../components/icons/Icons";

const WhyChoose = () => {
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
      <h2 className="text-[1.125rem] text-orange-1 text-center font-[700] leading-[24.75px]">
        QUALITY SERVICES
      </h2>
      <h3 className="text-[2rem] text-black-2 text-center font-[700] leading-[41.14px] mt-[12px]">
        Why You Should Choose Greybox
      </h3>
      <section className="flex flex-col lg:flex-row gap-[35px] mt-[30px] md:justify-center">
        <div className="flex flex-col md:flex-row gap-[35px]">
          <ReasonsCard
            img={<CostIcon />}
            text="Cost-Effective"
            subText="Save on transaction fees"
            extraClass="bg-orange-3"
          />
          <ReasonsCard
            img={<StabilityIcon />}
            text="Stability"
            subText="Enjoy stable transactions"
            extraClass="bg-grey-1"
          />
        </div>
        <div className="flex flex-col md:flex-row gap-[35px]">
          <ReasonsCard
            img={<SecurityIcon />}
            text="Security"
            subText="Your funds are safe with us"
            extraClass="bg-orange-3"
          />
          <ReasonsCard
            img={<ComplianceIcon />}
            text="Speed"
            subText="Enjoy fast transactions wherever you are"
            extraClass="bg-grey-1"
          />
        </div>
      </section>
    </section>
  );
};

export default WhyChoose;
