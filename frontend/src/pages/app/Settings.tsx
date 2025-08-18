// import { useNavigate } from 'react-router-dom'
// import { SettingsCard } from '../../components/Cards'
// import {
//   AboutIcon,
//   CancelIconWhite,
//   ChangePinIcon,
//   CopyWhite,
//   SignOutIcon,
//   UpdateWalletIcon,
//   UserPicture,
// } from '../../components/icons/Icons'
// import { useGetAuthUserQuery } from '../../appSlices/apiSlice'
// import AppLayout from './AppLayout'
// import { useSnackbar } from 'notistack'

// const Settings = () => {
//   const navigate = useNavigate()
//   const { enqueueSnackbar } = useSnackbar()
//   const { currentData: userData, isFetching } = useGetAuthUserQuery({})
//   const personInfo = userData?.data?.personal_details

//   // Function to copy text to clipboard
//   const copyToClipboard = (text: string) => {
//     navigator.clipboard
//       .writeText(text)
//       .then(() => {
//         enqueueSnackbar('Account address copied to clipboard!', {
//           variant: 'success',
//         })
//       })
//       .catch((err) => {
//         enqueueSnackbar('Failed to copy address', { variant: 'success' })
//       })
//   }

//   return (
//     <AppLayout
//       child={
//         <div className="w-full md:w-[50.33%] lg:w-[45.33%] min-h-[100vh] bg-grey-1">
//           <section className="bg-orange-1 pt-[28px] pb-[40px]">
//             <div className="flex items-center justify-center relative">
//               <span
//                 className="absolute left-[24px]"
//                 onClick={() => navigate(-1)}
//               >
//                 <CancelIconWhite />
//               </span>
//               <h2 className="text-white text-[1.5rem] font-[600]">Settings</h2>
//             </div>
//             <div className="w-full flex items-center pl-[29px] mt-[44px] gap-x-[45px]">
//               <UserPicture />
//               <div>
//                 <p className="text-white text-[1rem] leading-[22px]">
//                   {isFetching
//                     ? null
//                     : `${personInfo?.first_name} ${personInfo?.last_name}`}
//                 </p>
//                 <p className="flex items-center gap-x-[17px] text-white text-[0.875rem] leading-[12px]">
//                   {isFetching
//                     ? 'Fetching...'
//                     : personInfo?.account_address.substring(0, 15) + '...'}

//                   <div
//                     onClick={() => copyToClipboard(personInfo?.account_address)}
//                   >
//                     <CopyWhite />
//                   </div>
//                 </p>
//               </div>
//             </div>
//           </section>

//           <section className="flex bg-white flex-col gap-y-[32px] mt-[47px]">
//             <SettingsCard
//               text="Update Wallet Details"
//               subText="Coming soon..."
//               onClick={() => navigate('/seupdate-wallet-details')}
//               icon={<UpdateWalletIcon />}
//             />
//             <SettingsCard
//               text="Change Pin"
//               subText="Coming soon"
//               onClick={() => navigate('/change-passcode')}
//               icon={<ChangePinIcon />}
//             />
//             <SettingsCard
//               text="About Greybox"
//               subText="Get to know more about us"
//               onClick={() => navigate('/about-greybox')}
//               icon={<AboutIcon />}
//             />
//             <SettingsCard
//               text="Sign Out"
//               subText="We'll miss you..."
//               onClick={() => {
//                 localStorage.removeItem('access_token')
//                 navigate('/')
//               }}
//               icon={<SignOutIcon />}
//             />
//           </section>
//         </div>
//       }
//     />
//   )
// }

// export default Settings

import {
  User,
  Shield,
  Smartphone,
  Lock,
  Globe,
  Users,
  MessageCircle,
  CreditCard,
  Bell,
  LogOut,
} from 'lucide-react'
import AppLayout from './AppLayout'

export const Settings = () => (
  <AppLayout
    child={
      <div className="p-6 bg-white shadow-md w-full max-w-xl mx-auto min-h-screen">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900">More</h1>
          <button className="bg-orange-600 text-white px-4 py-2 rounded-full text-sm font-medium flex items-center space-x-2">
            <Users className="w-4 h-4" />
            <span>Our community</span>
          </button>
        </div>

        {/* Account Section */}
        <div className="mb-8">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4">
            ACCOUNT
          </h3>
          <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <User className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">Your Profile</span>
              </div>
            </button>

            <div className="border-t border-gray-100"></div>

            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <CreditCard className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">
                  Account Verification
                </span>
              </div>
              <span className="text-sm text-green-600 font-medium">
                Verified
              </span>
            </button>

            <div className="border-t border-gray-100"></div>

            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <Bell className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">Notifications</span>
              </div>
            </button>
          </div>
        </div>

        {/* Finances Section */}
        <div className="mb-8">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4">
            FINANCES
          </h3>
          <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <CreditCard className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">
                  Transaction Limits
                </span>
              </div>
            </button>
          </div>
        </div>

        {/* Security Section */}
        <div className="mb-8">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4">
            SECURITY
          </h3>
          <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <Lock className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">
                  Change Password
                </span>
              </div>
            </button>

            <div className="border-t border-gray-100"></div>

            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <Shield className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">
                  Two Factor Authentication
                </span>
              </div>
            </button>

            <div className="border-t border-gray-100"></div>

            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <Smartphone className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">
                  Devices & Sessions
                </span>
              </div>
            </button>

            <div className="border-t border-gray-100"></div>

            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <Lock className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">
                  Change your PIN
                </span>
              </div>
            </button>
          </div>
        </div>

        {/* Others Section */}
        <div className="mb-8">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4">
            OTHERS
          </h3>
          <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <Globe className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">App Language</span>
              </div>
            </button>

            <div className="border-t border-gray-100"></div>

            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <Users className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">
                  Affiliates & Referrals
                </span>
              </div>
            </button>

            <div className="border-t border-gray-100"></div>

            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <MessageCircle className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">
                  Talk to Support
                </span>
              </div>
            </button>

            <div className="border-t border-gray-100"></div>

            <button className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                  <CreditCard className="w-5 h-5 text-orange-600" />
                </div>
                <span className="font-medium text-gray-900">See our rates</span>
              </div>
            </button>
          </div>
        </div>

        {/* Logout */}
        <div className="flex justify-center mb-8">
          <button className="flex items-center space-x-2 text-red-600 font-medium">
            <span>Log out</span>
            <LogOut className="w-4 h-4" />
          </button>
        </div>

        {/* Version */}
        <div className="text-center text-gray-400 text-sm">v2.5.8</div>
      </div>
    }
  />
)
