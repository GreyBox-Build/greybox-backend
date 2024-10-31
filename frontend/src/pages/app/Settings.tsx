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
} from "../../components/icons/Icons";
import { useGetAuthUserQuery } from "../../appSlices/apiSlice";
import AppLayout from "./AppLayout";
import { useSnackbar } from "notistack";

const Settings = () => {
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();
  const { currentData: userData, isFetching } = useGetAuthUserQuery({});
  const personInfo = userData?.data?.personal_details;

  // Function to copy text to clipboard
  const copyToClipboard = (text: string) => {
    navigator.clipboard
      .writeText(text)
      .then(() => {
        enqueueSnackbar("Account address copied to clipboard!", {
          variant: "success",
        });
      })
      .catch((err) => {
        enqueueSnackbar("Failed to copy address", { variant: "success" });
      });
  };

  return (
    <AppLayout
      child={
        <div className="w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
          <section className="bg-grey-2 pt-[28px] pb-[40px]">
            <div className="flex items-center justify-center relative">
              <span
                className="absolute left-[24px]"
                onClick={() => navigate(-1)}
              >
                <CancelIconWhite />
              </span>
              <h2 className="text-white text-[1.5rem] font-[600]">Settings</h2>
            </div>
            <div className="w-full flex items-center pl-[29px] mt-[44px] gap-x-[45px]">
              <UserPicture />
              <div>
                <p className="text-white text-[1rem] leading-[22px]">
                  {isFetching
                    ? null
                    : `${personInfo?.first_name} ${personInfo?.last_name}`}
                </p>
                <p className="flex items-center gap-x-[17px] text-white text-[0.875rem] leading-[12px]">
                  {isFetching
                    ? "Fetching..."
                    : personInfo?.account_address.substring(0, 15) + "..."}

                  <div
                    onClick={() => copyToClipboard(personInfo?.account_address)}
                  >
                    <CopyWhite />
                  </div>
                </p>
              </div>
            </div>
          </section>

          <section className="flex flex-col gap-y-[32px] mt-[47px]">
            <SettingsCard
              text="Update Wallet Details"
              subText="Coming soon..."
              onClick={() => navigate("/seupdate-wallet-details")}
              icon={<UpdateWalletIcon />}
            />
            <SettingsCard
              text="Change Pin"
              subText="Coming soon"
              onClick={() => navigate("/change-passcode")}
              icon={<ChangePinIcon />}
            />
            <SettingsCard
              text="About Greybox"
              subText="Get to know more about us"
              onClick={() => navigate("/about-greybox")}
              icon={<AboutIcon />}
            />
            <SettingsCard
              text="Sign Out"
              subText="We'll miss you..."
              onClick={() => {
                localStorage.removeItem("access_token");
                navigate("/");
              }}
              icon={<SignOutIcon />}
            />
          </section>
        </div>
      }
    />
  );
};

export default Settings;
