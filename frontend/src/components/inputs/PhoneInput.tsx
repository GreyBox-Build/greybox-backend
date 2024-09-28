import { Controller } from "react-hook-form";

interface FormFieldProps
  extends React.DetailedHTMLProps<
    React.InputHTMLAttributes<HTMLInputElement>,
    HTMLInputElement
  > {
  control: any;
  name?: string;
  countryCode:string;
  localType?: string;
  onLocalChange?: () => void;
  onClick?: React.MouseEventHandler<HTMLDivElement> | undefined;
}
export const PhoneInput = ({
  control,
  name,
  countryCode,
  localType,
  onClick,
  onLocalChange,
  ...props
}: FormFieldProps) => {
  const formatInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (localType === "figure") {
      return e.target.value
        .replace(/[^0-9]/g, "")
        .replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    }

    if (localType === "number") {
      return e.target.value.replace(/[^0-9]/g, "");
    }
    return e.target.value;
  };
  return (
    <Controller
      control={control}
      name={name!}
      render={({ field: { value, onChange }, fieldState: { error } }) => (
        <div>
          <div
            className={`h-[48px] p-[0px_19px] rounded-[8px]  flex items-center border-[#99999961] border-[1px] shadow-shadow-1`}
            onClick={onClick}
          >
            <div className=" w-fit h-[90%] border-none flex items-center justify-center text-black-3 text-[0.875rem] leading-[18px]">{countryCode}</div>
            <input
              value={value}
              onChange={(e) => {
                onChange(formatInput(e));
                onLocalChange && onLocalChange();
              }}
              {...props}
              className={`flex-grow w-[24%] h-[90%] text-black-3 placeholder:text-black-3 text-[0.875rem] flex items-center justify-center leading-[18px]  bg-transparent outline-none px-[5px]`}
            />
          
          </div>
          {error && (
            <p className=" text-red-700 text-[0.875rem] mt-[8px]">
              {error.message}
            </p>
          )}
        </div>
      )}
    />
  );
};


