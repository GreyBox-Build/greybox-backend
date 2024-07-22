import AuthLayout from "./AuthLayout";
import { Success } from "../../components/icons/Icons";
import { FormButton } from "../../components/buttons/FormButton";
import { useNavigate } from "react-router-dom";

const WelcomeScreen = () => {
  const navigate = useNavigate();
  return (
    <AuthLayout
      child={
        <div className="w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1 p-[51px_25px]">
          <div className="flex justify-center">
            <Success />
          </div>
          <section className="flex flex-col gap-y-[8px] mt-[81px]">
            <p className="text-[1.5rem] text-black-2 font-[600] leading-[24px] text-center">
              Account was created
            </p>
            <p className="text-[1rem] text-black-3 leading-[22px] text-center mt-[24px]">
              Lorem Ipsum is a dummy text Lorem Ipsum is a dummy text Lorem
              Ipsum is a dummy text
            </p>
          </section>
          <FormButton
            label="Continue"
            onClick={() => {
              navigate("/dashboard");
            }}
            extraClass="mt-[132px]"
          />
        </div>
      }
    />
  );
};

export default WelcomeScreen;
