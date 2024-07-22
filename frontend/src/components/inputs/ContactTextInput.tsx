interface FormFieldProps
  extends React.DetailedHTMLProps<
    React.InputHTMLAttributes<HTMLInputElement>,
    HTMLInputElement
  > {}
export const ContactTextInput = ({ ...props }: FormFieldProps) => {
  return (
    <input
      {...props}
      className="w-full h-[54px] rounded-[4px] border-b-[1px] border-b-grey-2 bg-pink-1 text-[0.875rem] text-black-3 font-[300] placeholder:text-[0.875rem] placeholder:text-black-3 placeholder:font-[300] outline-none pl-[11px]"
    />
  );
};
