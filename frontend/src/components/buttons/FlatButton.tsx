export const FlatButton = ({
  label,
  extraClass,
  onClick,
}: {
  label: string;
  extraClass: string;
  onClick: React.MouseEventHandler<HTMLButtonElement> | undefined;
}) => (
  <button
    className={`w-full rounded-[4px] border-[1px] border-grey-1 h-[50px] p-[16px] text-black-2 text-[0.875rem] font-[700] ${extraClass}`}
    onClick={onClick}
  >
    {label}
  </button>
);
