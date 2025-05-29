import React from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { FiSend } from "react-icons/fi";
import { AiOutlineHome } from "react-icons/ai";
import { MdOutlineLogout } from "react-icons/md";

const Sidebar = () => {
  const navigate = useNavigate();
  const links = [
    {
      name: "dashboard",
      to: "/adminDashboard",
      img: <AiOutlineHome size={24} />,
    },
    { name: "Report", to: "summary", img: <FiSend size={24} /> },
    { name: "Confirmed", to: "Set", img: <FiSend size={24} /> },
    { name: "Setting", to: "setting", img: <FiSend size={24} /> },
  ];

  const logoutFunc = () => {
    localStorage.removeItem("adminLogin");
    localStorage.removeItem("access_token");

    navigate("/admin");
  };

  return (
    <div className="w-full ">
      <div className="w-full flex justify-center  ">
        <img src="images/logoText.svg" alt="" className="mt-[56px]" />
      </div>

      <nav className="flex flex-col  justify-between h-[calc(100vh-200px)]">
        <ul className="flex flex-col gap-4 mt-10">
          {links.map((link, index) => (
            <li key={index} className="">
              <NavLink
                className={({ isActive }) =>
                  `flex items-center gap-2 px-8 py-3 rounded-[4px] ${
                    isActive ? "bg-orange-1  text-white" : "text-black"
                  }`
                }
                to={link.to}
                end
              >
                <span>{link.img}</span>
                <span className="capitalize font-medium">{link.name}</span>
              </NavLink>
            </li>
          ))}
        </ul>

        <div
          className="flex items-center gap-2 px-8 py-3 rounded-[4px] cursor-pointer"
          onClick={() => "Logout"}
        >
          <span>
            <MdOutlineLogout size={24} />
          </span>
          <span onClick={logoutFunc}>Logout</span>
        </div>
      </nav>
    </div>
  );
};

export default Sidebar;
