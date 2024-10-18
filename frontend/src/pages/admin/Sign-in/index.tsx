import { TextInput } from "../../../components/inputs/TextInput";
import { zodResolver } from "@hookform/resolvers/zod";
import { obtainTokenSchema } from "../../../utils/Validations";
import { useForm } from "react-hook-form";
import { useSnackbar } from "notistack";
import { useState } from "react";
import "./index.css";
import { FormButton } from "../../../components/buttons/FormButton";
import { MdOutlineRemoveRedEye } from "react-icons/md";
import { FaRegEyeSlash } from "react-icons/fa6";
import { useObtainTokenMutation } from "../../../appSlices/apiSlice";
import { useNavigate } from "react-router-dom";

const AdminSignIn = () => {
  const handleChange = () => {
    localStorage.setItem("adminLogin", "true");
  };

  const [typeChange, setTypeChange] = useState("password");
  const [passwordChange, setPasswordChange] = useState(true);
  const [obtainToken, { isLoading }] = useObtainTokenMutation();
  const navigate = useNavigate();

  const changeType = () => {
    setTypeChange(typeChange === "password" ? "text" : "password");
    setPasswordChange(!passwordChange);
  };
  const { enqueueSnackbar } = useSnackbar();

  const { register, handleSubmit } = useForm({
    defaultValues: {
      email: "",
      password: "",
    },
    resolver: zodResolver(obtainTokenSchema),
  });

  // Handle form submission and extract input values
  const handleObtainToken = async (data: any) => {
    try {
      const response: any = await obtainToken(data).unwrap();
      if (response.status) {
        localStorage.setItem("access_token", response?.data?.access_token);
        handleChange();
        navigate("/adminDashboard");
      }
    } catch (error: any) {
      enqueueSnackbar(
        error?.data?.error ? error?.data?.error : "Connection failed!",
        { variant: "success" }
      );
    }
  };

  return (
    <div className="w-full  p-[35px] flex items-center h-screen bg-[#f2f2f2]">
      <div className="flex w-full h-full py-10 gap-9 max-w-[1280px] mx-auto flex-col md:flex-row">
        <div className="bg-orange-1 rounded-[22px] relative w-full hidden md:block">
          <img
            src="images/admin-person.png"
            className="absolute bottom-0 max-h-full"
            alt="admin home pics"
          />
        </div>
        <div className="flex flex-col justify-center gap-10 w-full">
          <div>
            <h1 className="text-[40px] font-[700]">Sign In</h1>
            <p>Welcome back admin</p>
          </div>
          <form
            onSubmit={handleSubmit(handleObtainToken)}
            className="w-full gap-5 flex flex-col"
          >
            <div className="flex flex-col w-full">
              <label htmlFor="email">Email Address</label>{" "}
              <input
                {...register("email")}
                type="email"
                placeholder="Your email here"
                name="email"
                id="email"
                className="py-3 px-6 border-[2px] border-[#3333331A] outline-none bg-white focus:bg-white rounded-[8px] shadow-md"
              />
            </div>
            <div className="flex flex-col w-full">
              <label htmlFor="password">Password</label>

              <div className="flex relative w-full border-[2px] border-[#3333331A]  bg-transparent  rounded-[8px] shadow-md">
                <input
                  type={typeChange}
                  placeholder="Your password here"
                  id="password"
                  className="outline-none focus:bg-white py-3 px-6 bg-white w-full rounded-[8px]"
                  {...register("password")}
                />

                <div
                  className="absolute right-4 top-3 cursor-pointer transition-all duration-300 "
                  onClick={changeType}
                >
                  {typeChange === "password" ? (
                    <MdOutlineRemoveRedEye size={24} />
                  ) : (
                    <FaRegEyeSlash size={24} />
                  )}
                </div>
              </div>
            </div>
            <div className="flex flex-row gap-4">
              <input
                type="checkbox"
                id="rememberMe"
                className="appearance-none h-[18px] w-[18px] rounded-[1px] border-[2px] border-orange-1   focus:outline-none relative inline-block"
                {...register("password")}
              />
              <label htmlFor="rememberMe">Remember me</label>
            </div>

            <FormButton
              label="Login"
              extraClass="mt-6 !bg-orange-1 hover:!bg-orange-1/80"
              loading={isLoading}
            />
          </form>
        </div>
      </div>
    </div>
  );
};

export default AdminSignIn;
