import AuthLayout from "./AuthLayout";
import { BackArrow, LockOpen } from "../../components/icons/Icons";
import { TextInput } from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";

const RecoverPassword = () => {
  const navigate = useNavigate();
  const { control } = useForm();
  return (
    <AuthLayout
      child={
        <div className="w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1 p-[51px_25px]">
          <BackArrow />
          <h2 className=" text-[2.5rem] text-black-1 font-[700] mt-[50px] leading-[40px]">
            Forgot Password
          </h2>
          <p className="mt-[13px] text-[0.875rem] text-black-2">
            Fill in the details below, to recover your password.
          </p>
          <form className="mt-[24px]">
            <section className="flex flex-col gap-y-[32px]">
              <TextInput
                name="password"
                control={control}
                placeholder="Password"
                type="password"
                onChange={() => {}}
                img={<LockOpen />}
              />
              <TextInput
                name="confirm"
                control={control}
                placeholder="Confirm Password"
                type="password"
                onChange={() => {}}
                img={<LockOpen />}
              />
            </section>
            <button className="text-[0.875rem] text-black-3 leading-[12px] mt-[12px] border-none bg-transparent">
              Show Password
            </button>

            <FormButton
              label="Recover"
              onClick={() => {
                navigate("/verify-otp");
              }}
              extraClass="mt-[169px]"
            />
          </form>
        </div>
      }
    />
  );
};

export default RecoverPassword;
