import { useNavigate } from "react-router-dom";
import { HomeButton } from "../../../components/buttons/HomeButton";
import { IntroLady } from "../../../components/icons/Icons";

const Intro = () => {
  const navigate = useNavigate();
  // const Stat = ({
  //   text,
  //   subText,
  //   last,
  // }: {
  //   text: string;
  //   subText: string;
  //   last?: boolean;
  // }) => (
  //   <div
  //     className={`py-[2px] pr-[32px] ${
  //       last ? "" : "border-r-[2px] border-grey-1"
  //     }`}
  //   >
  //     <p className="text-[1.5rem] text-black-2 font-[600] leading-[24px]">
  //       {text}
  //     </p>
  //     <p className="mt-[8px] text-[1rem] text-black-3 leading-[22px]">
  //       {subText}
  //     </p>
  //   </div>
  // );
  return (
    <section className="pb-[100px] bg-white px-[25px] md:px-[5%] lg:px-[10%]">
      <div className="w-full flex flex-col-reverse md:flex-row items-center justify-center gap-x-[2%]">
        <section>
          <h2 className="max-w-[597px] text-[3rem] font-[700] leading-[66px]">
            Simplifying <span className=" text-orange-1">Cross-border </span>
            Payments in Africa
          </h2>

          <div className="mt-[78px] flex flex-col md:flex-row md:items-center gap-[28px]">
            <HomeButton
              label="Learn More"
              onClick={() => {}}
              extraClass="w-full md:w-[225px] bg-[#fff] text-orange-1"
            />
            <HomeButton
              label="Get Started"
              onClick={() => navigate("/sign-up")}
              extraClass="w-full md:w-[225px] bg-orange-1 text-[#fff]"
            />
          </div>
          {/* <div className="mt-[108px] flex items-center gap-x-[32px]">
            <Stat text="100k" subText="Lorem Ipsum" />
            <Stat text="150k" subText="Lorem Ipsum" />
            <Stat text="150k" subText="Lorem Ipsum" last />
          </div> */}
        </section>
        <section>
          <IntroLady />
        </section>
      </div>
    </section>
  );
};

export default Intro;
