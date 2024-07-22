import { useNavigate } from "react-router-dom";
import { SettingsCard } from "../../components/Cards";
import {
  AboutIcon,
  CancelIconWhite,
  ChangePinIcon,
  CopyWhite,
  PaymentDetailsIcon,
  SignOutIcon,
  UpdateWalletIcon,
  UserPicture,
  WelcomeIcon,
} from "../../components/icons/Icons";
import AppLayout from "./AppLayout";

const Settings = () => {
  const navigate = useNavigate();
  return (
    <AppLayout
      child={
        <div className="w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <section className=" bg-grey-2 pt-[28px] pb-[40px]">
            <div className="flex items-center justify-center relative ">
              <span
                className="absolute left-[24px]"
                onClick={() => navigate(-1)}
              >
                <CancelIconWhite />
              </span>
              <h2 className=" text-white text-[1.5rem] font-[600]">Settings</h2>
            </div>
            <div className="w-full flex items-center pl-[29px] mt-[44px] gap-x-[45px]">
              <UserPicture />
              <div>
                <p className="text-white text-[1rem] leading-[22px]">
                  Lorem Ipsum
                </p>
                <p className="flex items-center gap-x-[17px] text-white text-[0.875rem] leading-[12px]">
                  oifoo4i3049420fjrjnvzn wi3904 <CopyWhite />
                </p>
              </div>
            </div>
          </section>

          <section className="flex flex-col gap-y-[32px] mt-[47px]">
            <SettingsCard
              text="Welcome to Greybox"
              subText="Lorem Ipsum is a dummy text..."
              onClick={() => navigate("")}
              icon={<WelcomeIcon />}
            />
            <SettingsCard
              text="Update Wallet Details"
              subText="Lorem Ipsum is a dummy text..."
              onClick={() => navigate("/update-wallet-details")}
              icon={<UpdateWalletIcon />}
            />
            <SettingsCard
              text="Change Pin"
              subText="Lorem Ipsum is a dummy text..."
              onClick={() => navigate("/change-passcode")}
              icon={<ChangePinIcon />}
            />
            <SettingsCard
              text="Payment Details"
              subText="Lorem Ipsum is a dummy text..."
              onClick={() => navigate("/payment-details")}
              icon={<PaymentDetailsIcon />}
            />
            <SettingsCard
              text="About Greybox"
              subText="Lorem Ipsum is a dummy text..."
              onClick={() => navigate("/about-greybox")}
              icon={<AboutIcon />}
            />
            <SettingsCard
              text="Sign Out"
              subText="Lorem Ipsum is a dummy text..."
              onClick={() => navigate("/")}
              icon={<SignOutIcon />}
            />
          </section>
        </div>
      }
    />
  );
};

export default Settings;
