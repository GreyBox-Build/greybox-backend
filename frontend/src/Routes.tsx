import { createBrowserRouter } from 'react-router-dom'
import Login from './pages/auth/Login'
import SignUp from './pages/auth/SignUp'
import VerifyOtp from './pages/auth/VerifyOtp'
import ForgotPassword from './pages/auth/ForgotPassword'
import RecoverPassword from './pages/auth/RecoverPassword'
import LockScreen from './pages/auth/LockScreen'
import WelcomeScreen from './pages/auth/WelcomeScreen'
import Dashboard from './pages/app/Dashboard'
import AllTransaction from './pages/app/AllTransaction'

import Notifications from './pages/app/Notifications'

import UpdateWallet from './pages/app/UpdateWallet'
import About from './pages/app/About'
import WithdrawalOption from './components/WithdrawalOption'
import DepositOption from './components/DepositOption'
import SendBank from './pages/app/SendBank'
import SendOption from './components/SendOption'
import CryptoRampInterface from './pages/app/SendMobileMoney'
import ChangePassCode from './pages/app/ChangePassCode'
import NewPassCode from './pages/app/NewPassCode'
import PaymentDetails from './pages/app/PaymentDetails'
import Home from './pages/app/home/Home'
import Contact from './pages/app/Contact'
import PrivacyPolicy from './pages/app/PrivacyPolicy'
import AMLPolicy from './pages/app/AMLPolicy'
import AppLayout from './pages/app/AppLayout'
import BankTransferDeposit from './pages/app/BankTransferDeposit'
import BankWithdrawal from './pages/app/BankWithdrawal'
import MobileDeposit from './pages/app/MobileDeposit'
import AdminSignIn from './pages/admin/Sign-in'
import AdminDashboard from './pages/admin/Dashboard'
import DashboardInner from './pages/admin/DashboardInner'
import Summary from './pages/admin/Summary'
import BorderlessOnramp from './pages/app/BorderlessOnramp'
import VerifyUSer from './pages/app/VerifyUser'
import KYC from './pages/admin/KYC'
import BankAccountDetails from './pages/app/AccountDetails'
import VirtualCard from './pages/app/Cards'
import { TransactionHistory } from './pages/app/History'
import { Settings } from './pages/app/Settings'
import OnRampDepositForm from './components/BorderlessDeposit'
import OffRampWithdrawalForm from './components/BorderlessWithdrawalForm'

export const routes = createBrowserRouter([
  {
    path: '/',
    element: <Home />,
  },
  {
    path: '/sign-in',
    element: <Login />,
  },
  {
    path: '/sign-up',
    element: <SignUp />,
  },
  {
    path: '/welcome',
    element: <WelcomeScreen />,
  },
  {
    path: '/forgot-password',
    element: <ForgotPassword />,
  },
  {
    path: 'verify-user',
    element: <VerifyUSer />,
  },
  {
    path: '/recover-password',
    element: <RecoverPassword />,
  },

  {
    path: '/borderless',
    element: <BorderlessOnramp />,
  },
  {
    path: '/verify-otp',
    element: <VerifyOtp />,
  },
  {
    path: '/lock-screen',
    element: <LockScreen />,
  },
  {
    path: '/dashboard/*',
    element: <Dashboard />,
  },
  {
    path: '/all-transactions',
    element: <AllTransaction />,
  },
  {
    path: '/send-options',
    element: <SendOption />,
  },
  {
    path: '/send-via-bank',
    element: <SendBank />,
  },
  {
    path: '/send-via-mobile-money',
    element: <CryptoRampInterface />,
  },
  {
    path: '/deposit-options',
    element: <DepositOption />,
  },
  {
    path: '/deposit-via-bank-transfer',
    element: <BankTransferDeposit />,
  },

  {
    path: '/deposit-mobile-money',
    element: <MobileDeposit />,
  },
  {
    path: '/withdraw-fund',
    element: <OffRampWithdrawalForm />,
  },
  {
    path: '/deposit-funds',
    element: <OnRampDepositForm />,
  },
  // {
  //   path: "/deposit-via-exchange",
  //   element: <ExchangeDeposit />,
  // },
  {
    path: '/withdrawal-options',
    element: <WithdrawalOption />,
  },
  {
    path: '/withdraw-via-bank',
    element: <BankWithdrawal />,
  },

  {
    path: '/notifications',
    element: <Notifications />,
  },
  {
    path: '/settings',
    element: <Settings />,
  },
  {
    path: '/virtual-card',
    element: <VirtualCard />,
  },
  {
    path: '/transaction-history',
    element: <TransactionHistory />,
  },
  {
    path: '/account-details',
    element: <BankAccountDetails />,
  },
  {
    path: '/update-wallet-details',
    element: <UpdateWallet />,
  },
  {
    path: '/change-passcode',
    element: <ChangePassCode />,
  },
  {
    path: '/new-passcode',
    element: <NewPassCode />,
  },
  {
    path: '/about-greybox',
    element: <About />,
  },
  {
    path: '/contact',
    element: <Contact />,
  },
  {
    path: '/privacy-policy',
    element: <PrivacyPolicy />,
  },
  {
    path: '/aml-policy',
    element: <AMLPolicy />,
  },
  {
    path: '/payment-details',
    element: <PaymentDetails />,
  },
  {
    path: '/admin',
    element: <AdminSignIn />,
  },
  {
    path: '/adminDashboard',
    element: <AdminDashboard />,
    children: [
      {
        path: '',
        element: <DashboardInner />,
      },
      { path: 'summary', element: <Summary /> },
      { path: 'kyc', element: <KYC /> },
    ],
  },
  {
    path: '*',
    element: (
      <AppLayout
        child={
          <div className="m-[auto] text-white text-[1.5rem]">
            Page not found!
          </div>
        }
      />
    ),
  },
])
