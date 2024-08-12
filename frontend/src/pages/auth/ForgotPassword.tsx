import AuthLayout from "./AuthLayout";
import { BackArrow, Mail } from "../../components/icons/Icons";
import { TextInput } from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { forgetPasswordSchema } from "../../utils/Validations";
import { useForgetPasswordMutation } from "../../appSlices/apiSlice";
import { useSnackbar } from "notistack";

const ForgotPassword = () => {
  const { control, handleSubmit } = useForm({
    defaultValues: {
      email: "",
    },
    resolver: zodResolver(forgetPasswordSchema),
  });

  const { enqueueSnackbar } = useSnackbar();

  const [forgetPassword, { isLoading }] = useForgetPasswordMutation();

  const handleForgetPassword = async (data: { email: string }) => {
    try {
      const response = await forgetPassword(data).unwrap();
      enqueueSnackbar(response?.message, { variant: "success" });
    } catch (error) {
      console.log(error);
    }
  };
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
          <form
            className="mt-[24px]"
            onSubmit={handleSubmit(handleForgetPassword)}
          >
            <section className="flex flex-col">
              <TextInput
                name="email"
                control={control}
                placeholder="Email Address"
                type="email"
                img={<Mail />}
              />
            </section>

            <FormButton
              label="Recover"
              extraClass="mt-[169px]"
              loading={isLoading}
            />
          </form>
        </div>
      }
    />
  );
};

export default ForgotPassword;
