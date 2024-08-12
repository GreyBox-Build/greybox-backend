import AuthLayout from "./AuthLayout";
import { BackArrow } from "../../components/icons/Icons";
import { Link, useNavigate } from "react-router-dom";
import { FormButton } from "../../components/buttons/FormButton";
import { useState } from "react";

const VerifyOtp = () => {
  const navigate = useNavigate();
  const [otp, setOtp] = useState(new Array(4).fill(""));

  const handleChange = (element: any, index: number) => {
    const value = element.value;

    // When the value is empty and the field was previously filled, move backward
    if (value === "" && otp[index] !== "") {
      if (element.previousSibling) {
        (element.previousSibling as HTMLInputElement).focus();
      }
    } else if (!isNaN(value)) {
      // Handle forward navigation
      setOtp([...otp.map((d, idx) => (idx === index ? value : d))]);
      if (element.nextSibling) {
        element.nextSibling.focus();
      }
    }
  };

  const handleKeyDown = (
    e: React.KeyboardEvent<HTMLInputElement>,
    index: number
  ) => {
    if (e.key === "Backspace") {
      // Prevent default Backspace behavior
      e.preventDefault();

      // Delete the value of the current field
      if (otp[index] !== "") {
        setOtp([...otp.map((d, idx) => (idx === index ? "" : d))]);
      }

      // Move focus to the previous input field if current field is empty
      if (index > 0) {
        (e.currentTarget.previousSibling as HTMLInputElement).focus();
      }
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
            Enter the 4-digit code sent to your email
          </p>
          <form className="mt-[24px]">
            <section className="flex gap-x-[17px] mt-[24px] w-full">
              {otp.map((data, index) => (
                <input
                  name="otp"
                  maxLength={1}
                  key={index}
                  value={data}
                  pattern="[0-9]"
                  required
                  onChange={(e) => handleChange(e.target, index)}
                  onFocus={(e) => e.target.select()}
                  onKeyDown={(e) => handleKeyDown(e, index)}
                  className=" w-[24%] h-[48px] p-[11px_9.5%] text-black-3 placeholder:text-black-3 text-[0.875rem] leading-[18px] border-[#99999961] border-[1px] gap-x-[5px] shadow-shadow-1 rounded-[8px] outline-none"
                  autoFocus={index === 0}
                />
              ))}
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
