interface FormFieldProps
  extends React.DetailedHTMLProps<
    React.TextareaHTMLAttributes<HTMLTextAreaElement>,
    HTMLTextAreaElement
  > {}
export const ContactTextArea = ({ ...props }: FormFieldProps) => {
  return (
    <textarea
      {...props}
      className="w-full h-[164px] rounded-[4px] border-b-[1px] border-b-grey-2 bg-pink-1 text-[0.875rem] text-black-3 font-[300] placeholder:text-[0.875rem] placeholder:text-black-3 placeholder:font-[300] resize-none outline-none pl-[11px] pt-[16px]"
    />
  );
};
