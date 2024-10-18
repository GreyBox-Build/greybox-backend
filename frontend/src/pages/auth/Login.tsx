import AuthLayout from "./AuthLayout";
import { BackArrow, LockOpen, Mail } from "../../components/icons/Icons";
import { TextInput } from "../../components/inputs/TextInput";
import { Link, useNavigate } from "react-router-dom";
import { FormButton } from "../../components/buttons/FormButton";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { obtainTokenSchema } from "../../utils/Validations";
import { useObtainTokenMutation } from "../../appSlices/apiSlice";
import { useSnackbar } from "notistack";

const Login = () => {
  const navigate = useNavigate();
  const [obtainToken, { isLoading }] = useObtainTokenMutation();
  const { enqueueSnackbar } = useSnackbar();

  const { control, handleSubmit } = useForm({
    defaultValues: {
      email: "",
      password: "",
    },
    resolver: zodResolver(obtainTokenSchema),
  });
  const handleObtainToken = async (data: any) => {
    try {
      const response: any = await obtainToken(data).unwrap();
      if (response.status) {
        localStorage.setItem("access_token", response?.data?.access_token);
        navigate("/dashboard");
      }
    } catch (error: any) {
      enqueueSnackbar(
        error?.data?.error ? error?.data?.error : "Connction failed!",
        { variant: "success" }
      ); 
    }
  };
  return (
    <AuthLayout
      child={
        <div className="w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1 p-[51px_25px]">
          <BackArrow />
          <h2 className=" text-[2.5rem] text-black-1 font-[700] mt-[50px] leading-[40px]">
            Welcome Back!
          </h2>
          <p className="mt-[13px] text-[0.875rem] text-black-2">
            Fill in the details below, to create your account.
          </p>
          <form
            className="mt-[24px]"
            onSubmit={handleSubmit(handleObtainToken)}
          >
            <section className="flex flex-col gap-y-[32px]">
              <TextInput
                control={control}
                name="email"
                placeholder="Email Address"
                type="email"
                img={<Mail />}
              />
              <TextInput
                control={control}
                name="password"
                placeholder="Password"
                type="password"
                img={<LockOpen />}
              />
            </section>
            <Link
              to={"/forgot-password"}
              className="text-[0.875rem] text-black-3 leading-[12px] mt-[12px]"
            >
              Forgot Password?
            </Link>

            <FormButton
              label="Continue"
              extraClass="mt-[72px]"
              loading={isLoading}
            />
            <section className="flex flex-col gap-y-[8px] mt-[55px]">
              <p className="text-[0.875rem] text-black-3 leading-[18px]">
                Don't Have an Account?
              </p>
              <Link
                to={"/sign-up"}
                className=" text-[0.875rem] text-orange-1 leading-[18px] font-[700]"
              >
                Signup here &gt;
              </Link>
            </section>
          </form>
        </div>
      }
    />
  );
};

export default Login;
