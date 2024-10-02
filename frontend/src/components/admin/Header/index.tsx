import React from "react";
import { useDispatch } from "react-redux";
import { setSearchQuery } from "../../../adminSlices/searchSlice";

const Header = () => {
  const dispatch = useDispatch();

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    dispatch(setSearchQuery(event.target.value));
  };
  return (
    <header className=" w-full py-3 pl-3 ">
      <div className="flex justify-between items-center">
        <div className=" w-1/2 md:w-[30%] flex shrink-0 px-4 py-2 relative items-center gap-3 border-b border-r-[#D9D9D9]">
          <img src="images/search.png" className="w-5 h-5" alt="search icon" />
          <input
            type="text"
            className="w-full outline-none border-none bg-transparent"
            placeholder="Search..."
            onChange={handleChange}
          />
        </div>

        <div className="flex gap-3 items-center">
          <span>
            <img
              src="images/bell.png"
              className="cursor-pointer"
              alt="notification"
            />
          </span>

          <div className="md:flex items-center gap-2 hidden">
            <img
              src="images/avatarp.png"
              className="cursor-pointer"
              alt="profile pics"
            />
            <div>
              <p>Kwekwu Peter</p>
              <p>@kwekwupeter</p>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
