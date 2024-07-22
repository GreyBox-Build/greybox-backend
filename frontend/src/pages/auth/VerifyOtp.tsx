import AuthLayout from "./AuthLayout";
import { BackArrow } from "../../components/icons/Icons";
import { Link, useNavigate } from "react-router-dom";
import { FormButton } from "../../components/buttons/FormButton";
import { useState } from "react";

const VerifyOtp = () => {
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
          <h2 className=" text-[2.5rem] text-black-1 font-[700] mt-[50px] leading-[40px]">
            OTP Sent!
          </h2>
          <p className="mt-[13px] text-[0.875rem] text-black-2">
            Enter the 4-digit code sent to you at +234-8209-2798
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
              label="Resend OTP"
              onClick={() => {
                navigate("/recover-password");
              }}
              extraClass="mt-[74px]"
            />
            <section className="flex flex-col gap-y-[8px] mt-[30px]">
              <Link
                to={""}
                className=" text-[0.875rem] text-orange-1 leading-[18px] font-[700]"
              >
                I didn’t receive code
              </Link>
            </section>
          </form>
        </div>
      }
    />
  );
};

export default VerifyOtp;
