import {
  ArrowForward,
  DeleteIcon,
  ToolIcon,
  TransferIcon,
} from "./icons/Icons";

export type MenuProps = {
  icon: React.ReactNode;
  label: string;
  onClick: React.MouseEventHandler<HTMLDivElement> | undefined;
};
export const Menu = ({ icon, label, onClick }: MenuProps) => (
  <div onClick={onClick} className="flex flex-col items-center">
    {icon} <p className="text-[0.875rem] text-black">{label}</p>
  </div>
);

export const QuickLink = ({
  label,
  onClick,
  icon,
}: {
  label: string;
  onClick: React.MouseEventHandler<HTMLDivElement> | undefined;
  icon: React.ReactNode;
}) => {
  return (
    <div
      onClick={onClick}
      className="flex flex-col items-center gap-y-[8px] cursor-pointer"
    >
      {icon}
      <p className="text-[0.875rem] text-white font-[700] leading-[18px] text-center">
        {label}
      </p>
    </div>
  );
};

export const DateHead = ({ date }: { date: string }) => (
  <div className="text-[0.875rem] leading-[12px] text-black-3 mb-[8px] mt-[20px] ml-[24px]">
    {date}
  </div>
);

export const DetailsCard = ({
  label,
  time,
  amount,
  onClick,
  index,
  length,
}: {
  label: string;
  time: string;
  amount: string;
  onClick: React.MouseEventHandler<HTMLDivElement> | undefined;
  index: number;
  length: number;
}) => (
  <section
    onClick={onClick}
    className="w-full p-[0px_24px] bg-white cursor-pointer"
  >
    <div
      className={`${
        index !== length - 1 && "border-b-[1px] border-b-grey-1"
      } flex justify-between py-[14px] `}
    >
      <div className="flex items-center gap-x-[16px]">
        <TransferIcon />
        <div className="flex flex-col gap-y-[5px]">
          <p className="text-[0.875rem] leading-[18px] text-black-2">{label}</p>
          <p className="text-[0.875rem] leading-[12px] text-black-3">{time}</p>
        </div>
      </div>
      <div className="text-[0.875rem] leading-[12px] text-black-3">
        {amount}
      </div>
    </div>
  </section>
);

export const PaymentDetailsCard = ({
  text,
  subText,
  deleteFunc,
}: {
  text: string;
  subText: string;
  deleteFunc: React.MouseEventHandler<HTMLImageElement> | undefined;
}) => (
  <section className="w-full p-[24px_16px] bg-[transparent] cursor-pointer rounded-[4px] border-[1px] border-[#999999]">
    <div className="flex items-center justify-between">
      <div className="flex flex-col gap-x-[16px]">
        <p className="text-[0.875rem] leading-[18px] text-black-2">{text}</p>
        <p className="text-[0.875rem] leading-[12px] text-black-3">{subText}</p>
      </div>
      <DeleteIcon onClick={deleteFunc} />
    </div>
  </section>
);

export const NotificationsCard = ({
  text,
  subText,
  date,
  onClick,
}: {
  text: string;
  subText: string;
  date: string;
  onClick?: React.MouseEventHandler<HTMLDivElement> | undefined;
}) => (
  <section onClick={onClick} className="w-full p-[0px_24px]  cursor-pointer">
    <div
      className={` flex justify-between items-center px-[14px] gap-x-[55px] `}
    >
      <div className="flex items-center gap-x-[16px]">
        <ToolIcon />
        <div className="flex flex-col gap-y-[6px]">
          <p className="text-[0.875rem] leading-[18px] text-black-2">{text}</p>
          <p className="text-[0.875rem] leading-[12px] text-black-3">
            {subText}
          </p>
        </div>
      </div>
      <div className="text-[0.875rem] leading-[12px] text-black-3 min-w-fit">
        {date}
      </div>
    </div>
  </section>
);

export const SettingsCard = ({
  text,
  subText,
  onClick,
  icon,
}: {
  text: string;
  subText: string;
  onClick?: React.MouseEventHandler<HTMLDivElement> | undefined;
  icon: React.ReactNode;
}) => (
  <section onClick={onClick} className="w-full p-[0px_24px]  cursor-pointer">
    <div
      className={` flex justify-between items-center px-[14px] gap-x-[55px] `}
    >
      <div className="flex items-center gap-x-[16px]">
        {icon}
        <div className="flex flex-col gap-y-[6px]">
          <p className="text-[0.875rem] leading-[18px] text-black-2">{text}</p>
          <p className="text-[0.875rem] leading-[12px] text-black-3">
            {subText}
          </p>
        </div>
      </div>
      <ArrowForward />
    </div>
  </section>
);

export const TransactionHistoryCard = ({
  label,
  channel,
  time,
  amount,
  status,
  onClick,
  index,
  length,
}: {
  label: string;
  channel: string;
  time: string;
  amount: string;
  status: string;
  onClick: React.MouseEventHandler<HTMLDivElement> | undefined;
  index: number;
  length: number;
}) => (
  <section
    onClick={onClick}
    className="w-full p-[0px_24px] bg-white cursor-pointer"
  >
    <div
      className={`${
        index !== length - 1 && "border-b-[1px] border-b-grey-1"
      } flex justify-between py-[14px] `}
    >
      <div className="flex items-center gap-x-[16px]">
        <TransferIcon />
        <div className="flex flex-col gap-y-[12px]">
          <div className="flex items-center gap-x-[6px]">
            <span className="text-[0.875rem] leading-[18px] text-black-2">
              {label}
            </span>
            <span
              className={` w-[43px] p-[4px_8px] rounded-[48px] flex items-center justify-center text-[0.5rem] ${
                status === "Success"
                  ? "bg-success-bg text-success-text"
                  : status === "Pending"
                  ? "bg-pending-bg text-pending-text"
                  : " bg-failed-bg text-failed-text"
              } `}
            >
              {status}
            </span>
          </div>
          <p className="text-[0.875rem] leading-[12px] text-black-3">
            {channel}
          </p>
        </div>
      </div>
      <div className="flex flex-col gap-y-[12px]">
        <p className="text-[0.875rem] leading-[12px] text-black-3 font-[700]">
          {amount}
        </p>
        <p className="text-[0.625rem] leading-[10px] text-black-3">{time}</p>
      </div>
    </div>
  </section>
);
