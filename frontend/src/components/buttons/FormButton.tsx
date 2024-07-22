import { Oval } from "react-loader-spinner";

interface FormButtonProps
  extends React.DetailedHTMLProps<
    React.ButtonHTMLAttributes<HTMLButtonElement>,
    HTMLButtonElement
  > {
  label: React.ReactNode;
  onClick?: React.MouseEventHandler<HTMLButtonElement> | undefined;
  extraClass?: string;
  loading?: boolean;
}

export const FormButton = ({
  label,
  onClick,
  extraClass,
  loading,
  ...props
}: FormButtonProps) => (
  <button
    disabled={loading}
    onClick={onClick}
    {...props}
    className={`w-full h-[52px] bg-grey-2 rounded-[48px] flex items-center justify-center text-white text-[0.875rem] font-[700] ${extraClass}`}
  >
    {loading ? (
      <Oval
        height={20}
        width={20}
        color="#fff"
        wrapperStyle={{}}
        wrapperClass=""
        visible={true}
        ariaLabel="oval-loading"
        secondaryColor="#22262B"
        strokeWidth={2}
        strokeWidthSecondary={2}
      />
    ) : (
      label
    )}
  </button>
);
