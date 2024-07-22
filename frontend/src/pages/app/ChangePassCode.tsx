import { BackArrow, LockIcon } from "../../components/icons/Icons";
import { Link, useNavigate } from "react-router-dom";
import { FormButton } from "../../components/buttons/FormButton";
import { useState } from "react";
import AuthLayout from "../auth/AuthLayout";

const ChangePassCode = () => {
  const navigate = useNavigate();
  const [otp, setOtp] = useState(new Array(4).fill(""));

  const handleChange = (element: any, index: number) => {
    if (isNaN(element.value)) return false;
    setOtp([...otp.map((d, idx) => (idx === index ? element.value : d))]);
    console.log(otp);
    if (element.nextSibling) {
      element.nextSibling.focus();
    }
  };

  return (
    <AuthLayout
      child={
        <div className="w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1 p-[51px_25px]">
          <BackArrow />
          <div className="flex justify-center mb-[36px]">
            <LockIcon />
          </div>
          <p className="mt-[13px] text-[0.875rem] text-center text-black-2">
            Enter Your current Pin
          </p>
          <form className="mt-[24px]">
            <section className="flex gap-x-[17px] mt-[24px] w-full">
              {otp.map((data, index) => {
                return (
                  <input
                    name="otp"
                    maxLength={1}
                    key={index}
                    value={data}
                    pattern="[0-9]"
                    required
                    onChange={(e) => handleChange(e.target, index)}
                    onFocus={(e) => e.target.select()}
                    className=" w-[24%] h-[48px] p-[11px_9.5%] text-black-3 placeholder:text-black-3 text-[0.875rem] leading-[18px] border-[#99999961] border-[1px] gap-x-[5px] shadow-shadow-1 rounded-[8px] outline-none"
                    autoFocus={index === 0}
                  />
                );
              })}
            </section>

            <FormButton
              label="Continue"
              onClick={() => {
                navigate("/new-passcode");
              }}
              extraClass="mt-[74px]"
            />
            <section className="flex flex-col gap-y-[8px] mt-[30px]">
              <Link
                to={""}
                className=" text-[0.875rem] text-center text-orange-1 leading-[18px] font-[700]"
              >
                Forgot Pin
              </Link>
            </section>
          </form>
        </div>
      }
    />
  );
};

export default ChangePassCode;
