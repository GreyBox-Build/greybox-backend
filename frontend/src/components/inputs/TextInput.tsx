import { Controller } from "react-hook-form";

interface FormFieldProps
  extends React.DetailedHTMLProps<
    React.InputHTMLAttributes<HTMLInputElement>,
    HTMLInputElement
  > {
  control: any;
  name?: string;
  img?: React.ReactNode;
  imgP?: React.ReactNode;
  pass?: boolean;
  isSmall?: boolean;
  localType?: string;
  onLocalChange?: () => void;
  onClick?: React.MouseEventHandler<HTMLDivElement> | undefined;
}
export const TextInput = ({
  control,
  name,
  img,
  isSmall,
  localType,
  onClick,
  onLocalChange,
  ...props
}: FormFieldProps) => {
  const formatInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (localType === "figure") {
      const value = e.target.value
        .replace(/[^0-9.]/g, "") // Allow digits and decimal points
        .replace(/(\..*?)\..*/g, "$1"); // Allow only one decimal point

      const parts = value.split("."); // Split into whole and decimal parts

      const formattedWholePart = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ","); // Format the whole part with commas

      const decimalPart = parts.length > 1 ? `.${parts[1]}` : ""; // Retain the decimal part if it exists

      return formattedWholePart + decimalPart;
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
            className={`h-[48px] ${
              isSmall ? "p-[11px_9.5%]" : "p-[0px_19px]"
            } rounded-[8px]  flex items-center border-[#99999961] border-[1px] gap-x-[5px] shadow-shadow-1`}
            onClick={onClick}
          >
            <input
              value={value}
              onChange={(e) => {
                onChange(formatInput(e));
                onLocalChange && onLocalChange();
              }}
              {...props}
              className={`flex-grow w-[24%] h-[90%] text-black-3 placeholder:text-black-3 text-[0.875rem] leading-[18px]  bg-transparent outline-none ${
                isSmall ? "flex items-center justify-center px-0" : "px-[5px]"
              }`}
            />
            {!isSmall && <span className=" w-fit">{img}</span>}
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

export const InputLabel = ({ text }: { text: string }) => (
  <label htmlFor={text} className="text-[0.875rem] text-black-2 mb-[8px]">
    {text}
  </label>
);

export const InputInfoLabel = ({
  title,
  value,
}: {
  title: string;
  value: string;
}) => (
  <div className="w-full rounded-[0px_0px_8px_8px] mt-[-5px] bg-orange-2 flex items-center justify-between p-[8px_22px] text-[0.875rem] text-black-2">
    <span>{title}</span> <span>{value}</span>
  </div>
);
