export const HomeButton = ({
  label,
  extraClass,
  onClick,
}: {
  label: string;
  extraClass: string;
  onClick: React.MouseEventHandler<HTMLButtonElement> | undefined;
}) => {
  return (
    <button
      onClick={onClick}
      className={`p-[18px] rounded-[48px] border-[2px] border-orange-1 text-[0.875rem] font-[700] ${extraClass}`}
    >
      {label}
    </button>
  );
};
