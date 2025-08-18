import { useState, useCallback, useMemo } from 'react'
import {
  Eye,
  EyeOff,
  Plus,
  Send,
  Minus,
  ArrowUpRight,
  ArrowDownLeft,
  Wallet,
  X,
  CreditCard,
  Smartphone,
  Building2,
  Bitcoin,
  UserPlus,
  Lock,
  ArrowRight,
  Banknote,
  Bell,
} from 'lucide-react'
import { Link, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { enqueueSnackbar } from 'notistack'
import moment from 'moment'
import { ASSETS, FIAT } from 'src/utils/Dummies'
// API Hooks
import {
  useGetAuthUserQuery,
  useGetTransactionQuery,
  useCreateVirtualAccountMutation,
  useGetMyKycQuery,
  useGetUserAccountdetailsQuery,
} from '../../appSlices/apiSlice'
import AppLayout from './AppLayout'
import SelectBoxT from 'src/components/modals/SelectBoxT'
import { FormButton } from 'src/components/buttons/FormButton'
import { TextInput } from 'src/components/inputs/TextInput'
import { DropDown } from 'src/components/icons/Icons'
import { Modal } from 'src/components/Modal'

// Types
type TabType = 'deposits' | 'withdrawals' | 'history'
type TransactionType = 'deposit' | 'transfer' | 'withdraw'

interface Transaction {
  id: number
  type: TransactionType
  amount: number
  description: string
  date: string
  time: string
  transaction_sub_type?: string
  asset?: string
  CreatedAt?: string
}

interface VirtualAccountFormData {
  fiat: string
  asset: string
}

// Schema
const createVirtualAccountSchema = z.object({
  fiat: z.string().min(1, { message: 'Fiat is required' }),
  asset: z.string().min(1, { message: 'Asset is required' }),
})

const PaymentWebApp = () => {
  const navigate = useNavigate()
  const [openAcceptedTos, setOpenAcceptedTos] = useState(false)
  // State
  const [showBalance, setShowBalance] = useState(true)
  const [openTos, setOpenTos] = useState(false)
  const [activeDialog, setActiveDialog] = useState<null | string>(null)
  const [tab, setTab] = useState<TabType>('deposits')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  // API Queries
  const { currentData: userData, isFetching: isFetchingUser } =
    useGetAuthUserQuery({})
  const { data: virtualAccount, isLoading: virtualIsLoading } =
    useGetMyKycQuery({})
  const { data: userAccountDetails, isLoading: userAccountLoading } =
    useGetUserAccountdetailsQuery({})

  console.log(userAccountDetails)
  const { currentData: transactionsData, isFetching: isFetchingTransactions } =
    useGetTransactionQuery(
      userData?.data?.personal_details?.crypto_currency?.toLowerCase(),
      {
        refetchOnMountOrArgChange: true,
      }
    )

  // Mutations
  const [createVirtualAccountMutation, { isLoading: isCreatingAccount }] =
    useCreateVirtualAccountMutation()

  // Form handling
  const { control, handleSubmit, reset, setValue, clearErrors } =
    useForm<VirtualAccountFormData>({
      resolver: zodResolver(createVirtualAccountSchema),
      defaultValues: { fiat: '', asset: '' },
    })

  // Memoized values
  const personInfo = useMemo(() => userData?.data?.personal_details, [userData])
  const walletInfo = useMemo(() => userData?.data?.wallet_details, [userData])
  const isVerified = useMemo(() => personInfo?.is_verified, [personInfo])
  const balance = useMemo(
    () => parseFloat(walletInfo?.balance || '0'),
    [walletInfo]
  )

  const { transactionArray } = useMemo(() => {
    const deposits =
      transactionsData?.data?.filter(
        (t: Transaction) => t?.transaction_sub_type === 'Deposit'
      ) || []

    const withdrawals =
      transactionsData?.data?.filter(
        (t: Transaction) => t?.transaction_sub_type === 'Withdrawal'
      ) || []

    return {
      depositArray: deposits,
      withdrawalArray: withdrawals,
      transactionArray: tab === 'deposits' ? deposits : withdrawals,
    }
  }, [transactionsData, tab])

  // Handlers
  const handleCreateVirtualAccount = useCallback(
    async (data: VirtualAccountFormData) => {
      try {
        const response: any = await createVirtualAccountMutation(data).unwrap()
        enqueueSnackbar(response?.message, { variant: 'success' })
        setIsModalOpen(false)
      } catch (error: any) {
        enqueueSnackbar(error?.data?.error || 'Connection failed!', {
          variant: 'error',
        })
        setIsModalOpen(false)
      }
    },
    [createVirtualAccountMutation]
  )

  const handleTabChange = useCallback((newTab: TabType) => {
    setTab(newTab)
  }, [])

  const handleDeposit = (type: string) => {
    const input = prompt(`Enter deposit amount via ${type}:`)
    const amount = input !== null ? parseFloat(input) : NaN

    if (amount && amount > 0) {
      enqueueSnackbar(`Deposit of $${amount} via ${type} initiated`, {
        variant: 'success',
      })
      setActiveDialog(null)
    }
  }

  const handleTransfer = () => {
    const input = prompt('Enter transfer amount to bank:')
    const amount = input !== null ? parseFloat(input) : NaN

    if (amount && amount > 0 && amount <= balance) {
      enqueueSnackbar(`Transfer of $${amount} initiated`, {
        variant: 'success',
      })
      setActiveDialog(null)
    } else if (amount > balance) {
      enqueueSnackbar('Insufficient balance', { variant: 'error' })
    }
  }

  const handleWithdraw = (type: string) => {
    const input = prompt(`Enter withdrawal amount via ${type}:`)
    const amount = input !== null ? parseFloat(input) : NaN

    if (amount && amount > 0 && amount <= balance) {
      enqueueSnackbar(`Withdrawal of $${amount} via ${type} initiated`, {
        variant: 'success',
      })
      setActiveDialog(null)
    } else if (amount > balance) {
      enqueueSnackbar('Insufficient balance', { variant: 'error' })
    }
  }

  // UI Components
  const VerificationBanner = () => (
    <section className="py-2 px-4 rounded-md bg-orange-100 mt-4 border border-orange-300">
      <div className="flex flex-col gap-4">
        <div className="flex justify-between items-center">
          <div className="h-8 w-8 bg-orange-500 rounded-full flex items-center justify-center text-white">
            <Lock className="w-4 h-4" />
          </div>
          <button
            onClick={() => navigate('/verify-user')}
            className="px-3 text-sm font-normal py-1 bg-orange-500 rounded-full text-white hover:bg-orange-600 transition-colors"
          >
            Start Verification
          </button>
        </div>
        <div className="flex justify-between items-center">
          <p className="text-base text-gray-800">
            Your Identity needs to be verified
          </p>
          <ArrowRight className="text-orange-500" />
        </div>
      </div>
    </section>
  )

  const VirtualAccountButton = () => {
    if (!virtualAccount?.data?.[0]?.virtual_account_id) {
      return (
        <button
          onClick={() => setIsModalOpen(true)}
          disabled={!isVerified}
          className="bg-gradient-to-r from-white to-gray-50 text-orange-1 px-6 py-3 rounded-full shadow-md hover:shadow-lg transition-all duration-200 hover:scale-105 flex items-center space-x-2 disabled:opacity-70"
        >
          <UserPlus className="w-5 h-5" />
          <span className="font-semibold">Create Virtual Account</span>
        </button>
      )
    }

    return (
      <button className="bg-gradient-to-r from-white to-gray-50 text-orange-1 px-6 py-2 rounded-md shadow shadow-gray-500 hover:shadow-xl transition-all duration-200 hover:scale-105 flex items-center space-x-2">
        <Banknote className="w-5 h-5" />
        <Link to="/account-details" className="font-semibold">
          {`${
            userAccountDetails?.data[0]?.instructions?.bank_account_number.slice(
              0,
              10
            ) || 'Account Rejected'
          }`}
        </Link>
        <ArrowRight className="w-4 h-4" />
      </button>
    )
  }

  const getTransactionIcon = (type: TransactionType) => {
    switch (type) {
      case 'deposit':
        return <ArrowDownLeft className="w-5 h-5 text-green-500" />
      case 'transfer':
        return <Send className="w-5 h-5 text-orange-500" />
      case 'withdraw':
        return <ArrowUpRight className="w-5 h-5 text-red-500" />
      default:
        return <Wallet className="w-5 h-5 text-gray-500" />
    }
  }

  const formatAmount = (amount: number) => {
    const sign = amount >= 0 ? '+' : ''
    return `${sign}$${Math.abs(amount).toFixed(2)}`
  }

  const getAmountColor = (amount: number) => {
    if (amount > 0) return 'text-green-600'
    return 'text-red-600'
  }

  const DialogOverlay: React.FC<{ children: React.ReactNode }> = ({
    children,
  }) => (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-2xl max-w-md w-full shadow-2xl">
        {children}
      </div>
    </div>
  )

  const VirtualAccountDialog = () => (
    <DialogOverlay>
      <div className="p-6">
        <div className="flex justify-end items-center mb-6">
          <button
            onClick={() => setActiveDialog(null)}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X
              className="w-5 h-5 text-gray-500"
              onClick={() => setOpenTos(!openTos)}
            />
          </button>
        </div>

        <div className="space-y-4 text-center">
          {!openAcceptedTos && (
            <div className="bg-orange-1 text-white p-4 rounded-lg">
              <button
                onClick={() => setOpenAcceptedTos(true)}
                className="text-sm text-white mt-2 break-all"
              >
                Accept Terms of Service
              </button>
            </div>
          )}

          {openAcceptedTos && (
            <div className="relative">
              {isLoading && (
                <div className="absolute inset-0 flex items-center justify-center bg-gray-100 z-10">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-1"></div>
                </div>
              )}
              <iframe
                className="w-full h-[50vh]"
                src={userAccountDetails?.data[0]?.tos_link}
                title="Terms of Service"
                onLoad={() => setIsLoading(false)}
              ></iframe>
            </div>
          )}
        </div>
      </div>
    </DialogOverlay>
  )

  return (
    <AppLayout
      child={
        <div className="min-h-screen w-full bg-gradient-to-br from-orange-50 to-orange-100 p-4">
          <div className="w-full max-w-2xl rounded-md p-3 bg-white mx-auto">
            <div className="bg-white shadow-sm">
              <div className="max-w-7xl mx-auto px-1 py-4">
                <div className="flex justify-between items-center">
                  <div className="relative">
                    <Bell className="h-6 w-6 text-gray-600" />
                    <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
                      0
                    </span>
                  </div>

                  <div className="text-center">
                    <p className="text-sm text-gray-600">Welcome,</p>
                    <p className="text-lg font-semibold text-gray-900">
                      Welcome, {personInfo?.first_name || 'User'}
                    </p>
                  </div>

                  <Link
                    to="/settings"
                    className="w-10 h-10 bg-orange-500 rounded-full flex items-center justify-center"
                  >
                    <span className="text-white font-semibold">J</span>
                  </Link>
                </div>
              </div>
            </div>

            {/* Verification Banner */}
            {!isVerified && <VerificationBanner />}

            {/* Balance Card */}
            <div className="bg-gradient-to-r from-orange-1  to-orange-600 rounded-2xl p-6 mb-6 shadow-xl">
              <div className="flex justify-between items-center mb-2">
                <span className="text-orange-100 text-sm font-medium">
                  Current Balance
                </span>
                <button
                  onClick={() => setShowBalance(!showBalance)}
                  className="p-2 hover:bg-orange-500 rounded-lg transition-colors"
                >
                  {showBalance ? (
                    <Eye className="w-5 h-5 text-white" />
                  ) : (
                    <EyeOff className="w-5 h-5 text-white" />
                  )}
                </button>
              </div>
              <div className="text-3xl font-bold text-white">
                {showBalance ? `$${balance.toFixed(2)}` : '••••••'}
              </div>
              <p className="text-orange-100 text-xs mt-2">
                Last updated {moment(personInfo?.UpdatedAt).fromNow()}
              </p>
              <div className="flex justify-center mt-4">
                <VirtualAccountButton />
              </div>
            </div>

            {/* Action Buttons */}
            <div className="grid grid-cols-3 gap-4 mb-8">
              <Link
                to="/deposit-options"
                className="bg-white flex justify-center flex-col items-center rounded-xl p-4 shadow-lg hover:shadow-xl transition-all duration-200 hover:scale-105 border border-orange-100"
              >
                <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-2">
                  <Plus className="w-6 h-6 text-green-600" />
                </div>
                <span className="text-sm font-semibold text-gray-700">
                  Deposit
                </span>
              </Link>

              <button
                onClick={() => setActiveDialog('transfer')}
                className="bg-white rounded-xl p-4 shadow-lg hover:shadow-xl transition-all duration-200 hover:scale-105 border border-orange-100"
              >
                <div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center mx-auto mb-2">
                  <Send className="w-6 h-6 text-orange-600" />
                </div>
                <span className="text-sm font-semibold text-gray-700">
                  Transfer
                </span>
              </button>

              <Link
                to="/withdrawal-options"
                className="bg-white flex items-center justify-center flex-col rounded-xl p-4 shadow-lg hover:shadow-xl transition-all duration-200 hover:scale-105 border border-orange-100"
              >
                <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-2">
                  <Minus className="w-6 h-6 text-orange-1" />
                </div>
                <span className="text-sm inline-block text-center font-semibold text-gray-700">
                  Withdraw
                </span>
              </Link>
            </div>

            {/* Transaction History */}
            <div className="bg-white rounded-2xl shadow-xl overflow-hidden">
              <div className="bg-gradient-to-r from-orange-1 to-orange-600 p-4">
                <div className="flex flex-col md:flex-row justify-between items-center">
                  <h2 className="text-lg font-semibold text-white">
                    Recent Transactions
                  </h2>
                  <div className="flex space-x-2">
                    <button
                      onClick={() => handleTabChange('deposits')}
                      className={`px-2 py-1 flex-1 rounded text-xs font-medium ${
                        tab === 'deposits'
                          ? 'bg-white text-orange-1'
                          : 'text-white'
                      }`}
                    >
                      Deposits
                    </button>
                    <button
                      onClick={() => handleTabChange('withdrawals')}
                      className={`px-2 flex-1 py-1 rounded text-xs font-medium ${
                        tab === 'withdrawals'
                          ? 'bg-white text-orange-1'
                          : 'text-white'
                      }`}
                    >
                      Withdrawals
                    </button>
                  </div>
                </div>
              </div>

              <div className="divide-y divide-gray-100 max-h-96 overflow-y-auto">
                {isFetchingTransactions ? (
                  <div className="p-4 text-center text-gray-500">
                    Loading transactions...
                  </div>
                ) : transactionArray?.length === 0 ? (
                  <div className="p-4 text-center text-gray-500">
                    No transactions yet
                  </div>
                ) : (
                  transactionArray?.map((transaction: any) => (
                    <div
                      key={transaction.id}
                      className="p-4 hover:bg-orange-50 transition-colors"
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <div className="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center">
                            {getTransactionIcon(
                              transaction.transaction_sub_type.toLowerCase() as TransactionType
                            )}
                          </div>
                          <div>
                            <p className="font-medium text-gray-800 capitalize">
                              {transaction.transaction_sub_type}
                            </p>
                            <p className="text-sm text-gray-600">
                              {transaction.asset || 'Transaction'}
                            </p>
                            <p className="text-xs text-gray-500">
                              {moment(transaction.CreatedAt).format(
                                'YYYY-MM-DD • hh:mm A'
                              )}
                            </p>
                          </div>
                        </div>
                        <div className="text-right">
                          <p
                            className={`font-semibold ${getAmountColor(
                              parseFloat(transaction.amount)
                            )}`}
                          >
                            {formatAmount(parseFloat(transaction.amount))}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>

            {/* Footer */}
            <div className="text-center mt-8 pb-8">
              <p className="text-gray-500 text-sm">
                © {new Date().getFullYear()} GreyBox Wallet. Secure & Simple.
              </p>
            </div>

            {openTos &&
            !userAccountLoading &&
            userAccountDetails?.data[0]?.has_accepted_tos === false ? (
              <VirtualAccountDialog />
            ) : (
              ''
            )}

            {/* {openTos && <VirtualAccountDialog />} */}
            <CreateVirtualAccountModal
              isOpen={isModalOpen}
              onClose={() => setIsModalOpen(false)}
              onSubmit={handleCreateVirtualAccount}
              isLoading={isCreatingAccount}
            />
          </div>
        </div>
      }
    />
  )
}

export default PaymentWebApp

const CreateVirtualAccountModal: React.FC<{
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: VirtualAccountFormData) => void
  isLoading: boolean
}> = ({ isOpen, onClose, onSubmit, isLoading }) => {
  const [openChain, setOpenChain] = useState(false)
  const [openAsset, setOpenAssets] = useState(false)

  const { control, handleSubmit, clearErrors, setValue, reset } =
    useForm<VirtualAccountFormData>({
      defaultValues: { fiat: '', asset: '' },
      resolver: zodResolver(createVirtualAccountSchema),
    })

  const handleClose = useCallback(() => {
    reset()
    onClose()
  }, [reset, onClose])

  const handleFormSubmit = useCallback(
    (data: VirtualAccountFormData) => {
      onSubmit(data)
    },
    [onSubmit]
  )

  return (
    <Modal isOpen={isOpen} onClose={handleClose}>
      <form
        className="p-6 max-w-lg mx-auto"
        onSubmit={handleSubmit(handleFormSubmit)}
      >
        <h2 className="text-lg font-semibold mb-4 text-white">
          Create Virtual Account
        </h2>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">
              FIAT
            </label>
            <TextInput
              control={control}
              name="fiat"
              placeholder="SELECT FIAT"
              readOnly
              type="text"
              onClick={() => setOpenChain(true)}
              img={<DropDown />}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">
              ASSET
            </label>
            <TextInput
              control={control}
              name="asset"
              placeholder="SELECT ASSET"
              readOnly
              type="text"
              onClick={() => setOpenAssets(true)}
              img={<DropDown />}
            />
          </div>

          <FormButton
            label="Create Account"
            extraClass="mt-[72px]"
            type="submit"
            loading={isLoading}
          />
        </div>

        <SelectBoxT
          state={openChain}
          title="Select FIAT"
          placeholder="Search FIAT"
          childList={FIAT}
          onPickChild={(list) => {
            setValue('fiat', list?.name!)
            clearErrors('fiat')
          }}
          onClose={() => setOpenChain(false)}
        />

        <SelectBoxT
          state={openAsset}
          title="Select Asset"
          placeholder="Search Asset"
          childList={ASSETS}
          onPickChild={(list) => {
            setValue('asset', list?.name!)
            clearErrors('asset')
          }}
          onClose={() => setOpenAssets(false)}
        />
      </form>
    </Modal>
  )
}
