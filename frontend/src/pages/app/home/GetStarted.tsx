import { useNavigate } from "react-router-dom";
import { HomeButton } from "../../../components/buttons/HomeButton";

const GetStarted = () => {
  const navigate = useNavigate();
  return (
    <section className=" bg-get-started-bg bg-cover bg-no-repeat min-h-[346px]  px-[25px] py-[60px] md:px-[5%] lg:px-[10%]">
      <div className="text-[2rem] text-black-3 max-w-[654px] font-[600] leading-[41.14px] mb-[96px]">
        Send funds across <span className="text-orange-1">Africa</span> at a
        fraction of the cost compared to traditional banks
      </div>

      <HomeButton
        label="Get Started"
        onClick={() => navigate("/sign-up")}
        extraClass="w-full md:w-[353px] bg-orange-1 text-white"
      />
    </section>
  );
};

export default GetStarted;
