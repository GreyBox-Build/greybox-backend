import { BsCaretUpFill } from "react-icons/bs";

interface CardProps {
  name: string;
  amount: number;
  interest: number;
}

const Card: React.FC<CardProps> = ({ name, amount, interest }) => {
  return (
    <div className="w-full bg-white rounded-md px-10 py-8">
      <div className="flex items-center gap-3">
        {" "}
        <h1 className="text-[28px] font-medium">
          ${amount > 999 ? `${Math.floor(amount / 1000)}K` : amount}
        </h1>{" "}
        <div className="bg-[#E1F4E4] px-1 flex gap-2 items-center rounded-[2px]">
          <p>
            {" "}
            {interest > 999
              ? `${(interest / 1000).toFixed(2)}`
              : interest.toFixed(2)}
            %
          </p>

          <span className="text-green-500">
            <BsCaretUpFill />
          </span>
        </div>
      </div>

      <div className="text-sm font-bold flex items-center gap-3 ml-5 mt-2">
        <div
          className={`w-[18px] h-[18px] rounded-full ${
            name === "Total Revenew"
              ? " bg-orange-1 "
              : name === "All orders"
              ? " bg-[#4D4D4D]"
              : name === "Pending orders"
              ? " bg-orange-1"
              : " bg-[#4D4D4D]"
          }`}
        ></div>
        <p className="text-nowrap">{name}</p>
      </div>
    </div>
  );
};

export default Card;
