import { Modal } from "@mui/material";
import { useState } from "react";
import { changeFirstLetterToUpperCase } from "../../utils/ChangeFirstLetterToUpperCase";

type ChildListProps = {
  logo?: React.ReactNode;
  name: string;
  code?: string;
};

type ChildListPropsArray = ChildListProps[];

const SelectBox = ({
  state,
  title,
  placeholder,
  childList,
  type,
  selectedCountry, // Add selectedCountry here
  onPickChild,
  onClose,
}: {
  state: boolean;
  title: string;
  placeholder: string;
  childList: ChildListPropsArray;
  selectedCountry?: any; // Type the selectedCountry prop
  type?: string;
  onPickChild: (list: ChildListProps) => void;
  onClose: () => void;
}) => {
  const [searchValue, setSearchValue] = useState<string>("");
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);

  const isGhana = selectedCountry?.countryName === "Ghana"; // Check if the selected country is Ghana

  return (
    <Modal
      disableAutoFocus
      open={state}
      onClose={onClose}
      className=" bg-grey-box-bg bg-cover bg-no-repeat flex justify-center"
    >
      <div className="p-[28px_26px] m-auto w-full h-fit md:w-[50.33%] lg:w-[45.33%] bg-grey-5 rounded-[32px_32px_0px_0px]">
        <p>{title}</p>
        <input
          type="text"
          placeholder={placeholder}
          className="w-full h-[48px] p-[15px_24px] rounded-[8px] bg-grey-1 text-[#999999] text-[0.875rem] placeholder:text-[#999999] border-[1px] border-[#999999] shadow-sm mt-[24px] mb-[30px] outline-none"
          onChange={(e) => setSearchValue(e.target.value)}
          autoFocus={true}
          onFocus={() => setSearchValue("")}
        />
        <section className="flex flex-col gap-y-[8px]">
          {type === "countryName"
            ? childList
                .filter((child: any) =>
                  child.countryName
                    .toLowerCase()
                    .includes(searchValue.toLowerCase())
                )
                .map((list: any, index) => (
                  <div
                    key={index}
                    className={` flex items-center justify-between cursor-pointer p-[16px] rounded-[4px] border-grey-1 border-[1px] transition-all duration-300 ease-in-out ${
                      selectedIndex === index ? "bg-grey-1" : "bg-grey-5"
                    }`}
                    onClick={() => {
                      onPickChild(list);
                      setSelectedIndex(index);
                      onClose();
                    }}
                  >
                    <div className="flex items-center gap-x-[11px]">
                      {list.logo ? (
                        <div
                          className={`w-[30px] h-[30px] rounded-[50%] flex items-center justify-center ${
                            selectedIndex === index
                              ? "bg-orange-2"
                              : "bg-grey-1"
                          }`}
                        >
                          {list.logo}
                        </div>
                      ) : list.image_url ? (
                        <img
                          src={list?.image_url}
                          alt="logo"
                          className="h-[32px] w-[32px] rounded-[50%]"
                        />
                      ) : (
                        ""
                      )}
                      <div className="flex flex-col gap-y-[3px]">
                        <p className=" text-black-2 text-[0.875rem] leading-[18px] ">
                          {changeFirstLetterToUpperCase(list?.countryName)}
                        </p>
                        {list.code && (
                          <p className="text-black-3 text-[0.875rem] leading-[12px]">
                            {list.code}
                          </p>
                        )}
                        {list.chain && (
                          <p className="text-black-3 text-[0.875rem] leading-[12px]">
                            {list.chain}
                          </p>
                        )}
                      </div>
                    </div>

                    <img
                      src={`/images/${
                        selectedIndex === index ? "orangeBgTick" : "greyBgTick"
                      }.png`}
                      alt=""
                      className="h-[18px] w-[18px]"
                    />
                  </div>
                ))
            : type === "network"
            ? childList
                ?.filter((child: any) =>
                  child?.toLowerCase().includes(searchValue.toLowerCase())
                )
                .map((list: any, index) => (
                  <div
                    key={index}
                    className={` flex items-center justify-between cursor-pointer p-[16px] rounded-[4px] border-grey-1 border-[1px] transition-all duration-300 ease-in-out ${
                      selectedIndex === index ? "bg-grey-1" : "bg-grey-5"
                    }`}
                    onClick={() => {
                      onPickChild(list);
                      setSelectedIndex(index);
                      onClose();
                    }}
                  >
                    <div className="flex items-center gap-x-[11px]">
                      <div className="flex flex-col gap-y-[3px]">
                        <p className=" text-black-2 text-[0.875rem] leading-[18px] ">
                          {isGhana
                            ? list.toUpperCase() === "VODAFONE"
                              ? "TELECEL"
                              : list.toUpperCase() === "TIGO"
                              ? "AT (TIGO)"
                              : list.toUpperCase() === "AIRTEL"
                              ? "AT (AIRTEL)"
                              : list.toUpperCase()
                            : list.toUpperCase()}
                        </p>
                      </div>
                    </div>

                    <img
                      src={`/images/${
                        selectedIndex === index ? "orangeBgTick" : "greyBgTick"
                      }.png`}
                      alt=""
                      className="h-[18px] w-[18px]"
                    />
                  </div>
                ))
            : childList
                .filter((child: any) =>
                  child?.countryName
                    ?.toLowerCase()
                    .includes(searchValue?.toLowerCase())
                )
                .map((list: any, index) => (
                  <div
                    key={index}
                    className={` flex items-center justify-between cursor-pointer p-[16px] rounded-[4px] border-grey-1 border-[1px] transition-all duration-300 ease-in-out ${
                      selectedIndex === index ? "bg-grey-1" : "bg-grey-5"
                    }`}
                    onClick={() => {
                      onPickChild(list);
                      setSelectedIndex(index);
                      onClose();
                    }}
                  >
                    <div className="flex items-center gap-x-[11px]">
                      {list.logo ? (
                        <div
                          className={`w-[30px] h-[30px] rounded-[50%] flex items-center justify-center ${
                            selectedIndex === index
                              ? "bg-orange-2"
                              : "bg-grey-1"
                          }`}
                        >
                          {list.logo}
                        </div>
                      ) : list.image_url ? (
                        <img
                          src={list?.image_url}
                          alt="logo"
                          className="h-[32px] w-[32px] rounded-[50%]"
                        />
                      ) : (
                        ""
                      )}
                      <div className="flex flex-col gap-y-[3px]">
                        <p className=" text-black-2 text-[0.875rem] leading-[18px] ">
                          {changeFirstLetterToUpperCase(list?.name)}
                        </p>
                        {list.code && (
                          <p className="text-black-3 text-[0.875rem] leading-[12px]">
                            {list.code}
                          </p>
                        )}
                        {list.chain && (
                          <p className="text-black-3 text-[0.875rem] leading-[12px]">
                            {list.chain}
                          </p>
                        )}
                      </div>
                    </div>

                    <img
                      src={`/images/${
                        selectedIndex === index ? "orangeBgTick" : "greyBgTick"
                      }.png`}
                      alt=""
                      className="h-[18px] w-[18px]"
                    />
                  </div>
                ))}
        </section>
      </div>
    </Modal>
  );
};

export default SelectBox;
