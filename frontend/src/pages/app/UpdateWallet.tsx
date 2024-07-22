import { useNavigate } from "react-router-dom";
import { CancelIcon, DropDown } from "../../components/icons/Icons";
import AppLayout from "./AppLayout";
import { TextInput } from "../../components/inputs/TextInput";
import { FormButton } from "../../components/buttons/FormButton";
import { useForm } from "react-hook-form";

const UpdateWallet = () => {
  const navigate = useNavigate();

  const { control } = useForm();
  return (
    <AppLayout
      child={
        <div className="pt-[51px] w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <div className="flex items-center justify-center relative ">
            <span className="absolute left-[24px]" onClick={() => navigate(-1)}>
              <CancelIcon />
            </span>
            <h2 className=" text-black text-[1.5rem] font-[600]">
              Update Wallet
            </h2>
          </div>
          <form className="mt-[29px] px-[24px] pb-[80px]">
            <section className="flex flex-col gap-y-[32px]">
              <div>
                <TextInput
                  name="currency"
                  control={control}
                  placeholder="Select currency"
                  readOnly
                  type="text"
                  onClick={() => {}}
                  img={<DropDown />}
                />
              </div>
            </section>
            <FormButton
              label="Update"
              extraClass="mt-[80px]"
              onClick={() => {}}
            />
          </form>
        </div>
      }
    />
  );
};

export default UpdateWallet;
