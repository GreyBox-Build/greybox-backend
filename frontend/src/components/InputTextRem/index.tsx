import React, { useState } from "react";
import { FaRegEye } from "react-icons/fa";
import { IoEyeOffOutline } from "react-icons/io5";
import { Controller } from "react-hook-form";
import { useTranslation } from "react-i18next";

interface InputTextRemProps {
  control: any;
  name: string;
}

const InputTextRem: React.FC<InputTextRemProps> = ({ control, name }) => {
  const [showPassword, setShowPassword] = useState(false);
  const { t }: { t: any } = useTranslation();
  return (
    <Controller
      control={control}
      name={name}
      render={({ field: { value, onChange }, fieldState: { error } }) => (
        <div>
          <div className="relative h-[48px] pr-[50px] pl-[20px] rounded-[8px] flex items-center border-[#99999961] border-[1px] gap-x-[5px] shadow-1">
            <input
              type={showPassword ? "text" : "password"}
              value={value}
              onChange={onChange}
              className="flex-grow w-full h-[90%] text-black-3 placeholder:text-black-3 text-[0.875rem] leading-[18px] bg-transparent outline-none px-2"
              placeholder={t("Password")}
            />
            <span
              className="absolute right-5 top-3 cursor-pointer"
              onClick={() => setShowPassword((prev) => !prev)}
            >
              {showPassword ? (
                <FaRegEye size={24} className="text-[#525252]" />
              ) : (
                <IoEyeOffOutline size={24} className="text-[#525252]" />
              )}
            </span>
          </div>
          {error && (
            <p className="text-red-700 text-[0.875rem] mt-[8px]">
              {error.message}
            </p>
          )}
        </div>
      )}
    />
  );
};

export default InputTextRem;
