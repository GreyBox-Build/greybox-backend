import AuthLayout from "./AuthLayout";
import {
  BackArrow,
  DropDown,
  LockOpen,
  Mail,
  Person,
} from "../../components/icons/Icons";
import { TextInput } from "../../components/inputs/TextInput";
import { Link, useNavigate } from "react-router-dom";
import { FormButton } from "../../components/buttons/FormButton";
import { useState } from "react";
import SelectBox from "../../components/modals/SelectBox";
import { countryData, currencyData } from "../../utils/Dummies";
import { useForm } from "react-hook-form";
import { createUserSchema } from "../../utils/Validations";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  useCreateUserMutation,
  useGetChainsQuery,
} from "../../appSlices/apiSlice";
import { useSnackbar } from "notistack";

const SignUp = () => {
  const navigate = useNavigate();
  const [openCurrency, setOpenCurrency] = useState<boolean>(false);
  const [openCountry, setOpenCountry] = useState<boolean>(false);
  const [openChain, setOpenChain] = useState<boolean>(false);
  const [createUser, { isLoading }] = useCreateUserMutation();
  const { enqueueSnackbar } = useSnackbar();

  const { control, handleSubmit, clearErrors, setValue, getValues } = useForm({
    defaultValues: {
      first_name: "",
      last_name: "",
      email: "",
      password: "",
      currency: "",
      country: "",
      chain: "",
      country_code: "",
    },
    resolver: zodResolver(createUserSchema),
  });
  const { currentData: chains } = useGetChainsQuery({});

  const handleCreateUser = async (data: any) => {
    const updatedData = { ...data, country_code: getValues("country_code") };
    try {
      const response: any = await createUser(updatedData).unwrap();
      console.log(response);
      enqueueSnackbar(response?.message, { variant: "success" });
      setTimeout(() => {
        navigate("/sign-in");
      }, 5000);
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
            Sign up
          </h2>
          <p className="mt-[13px] text-[0.875rem] text-black-2">
            Fill in the details below, to create your account.
          </p>
          <form className="mt-[24px]" onSubmit={handleSubmit(handleCreateUser)}>
            <section className="flex flex-col gap-y-[32px]">
              <TextInput
                control={control}
                name="first_name"
                placeholder="First Name"
                type="text"
                img={<Person />}
              />
              <TextInput
                control={control}
                name="last_name"
                placeholder="Last Name"
                type="text"
                img={<Person />}
              />
              <TextInput
                control={control}
                name="email"
                placeholder="Email Address"
                type="email"
                img={<Mail />}
              />
              <TextInput
                control={control}
                name="currency"
                placeholder="Currency"
                readOnly
                type="text"
                onClick={() => {
                  setOpenCurrency(true);
                }}
                img={<DropDown />}
              />
              <TextInput
                control={control}
                name="country"
                placeholder="Country"
                readOnly
                type="text"
                onClick={() => {
                  setOpenCountry(true);
                }}
                img={<DropDown />}
              />
              <TextInput
                control={control}
                name="chain"
                placeholder="Chain"
                readOnly
                type="text"
                onClick={() => {
                  setOpenChain(true);
                }}
                img={<DropDown />}
              />
              <TextInput
                control={control}
                name="password"
                placeholder="Password"
                type="password"
                img={<LockOpen />}
              />
            </section>

            <FormButton
              label="Continue"
              extraClass="mt-[72px]"
              type="submit"
              loading={isLoading}
            />
            <section className="flex flex-col gap-y-[8px] mt-[55px]">
              <p className="text-[0.875rem] text-black-3 leading-[18px]">
                Already Have an Account?
              </p>
              <Link
                to={"/sign-in"}
                className=" text-[0.875rem] text-orange-1 leading-[18px] font-[700]"
              >
                Login here &gt;
              </Link>
            </section>

            <SelectBox
              state={openCurrency}
              title="Select Currency"
              placeholder="Search Currency"
              childList={currencyData}
              onPickChild={(list) => {
                setValue("currency", list?.code!);
                clearErrors("currency");
              }}
              onClose={() => setOpenCurrency(false)}
            />
            <SelectBox
              state={openCountry}
              title="Select Country"
              placeholder="Search Country"
              childList={countryData}
              onPickChild={(list) => {
                setValue("country", list?.name);
                setValue("country_code", list?.code!);
                clearErrors("country");
              }}
              onClose={() => setOpenCountry(false)}
            />
            <SelectBox
              state={openChain}
              title="Select Chain"
              placeholder="Search Chain"
              childList={chains?.data === undefined ? [] : chains?.data}
              onPickChild={(list: any) => {
                setValue("chain", list?.chain);
                clearErrors("country");
              }}
              onClose={() => setOpenChain(false)}
            />
          </form>
        </div>
      }
    />
  );
};

export default SignUp;
