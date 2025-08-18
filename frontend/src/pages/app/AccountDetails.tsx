import React, { useState } from 'react'
import { Copy, Check, ArrowLeft } from 'lucide-react'
import { useGetUserAccountdetailsQuery } from 'src/appSlices/apiSlice'
import { Link, useNavigate } from 'react-router-dom'

const BankAccountDetails = () => {
  const [copiedField, setCopiedField] = useState<string | null>(null)
  const { data: userAccountDetails, isLoading: userAccountLoading } =
    useGetUserAccountdetailsQuery({})
  const navigate = useNavigate()
  // Mock data based on your API response
  const virtualAccount = {
    data: [
      {
        ID: 14,
        user_id: 86,
        account_id: 'cmda40zhl1bd1pd24litu0x7c',
        asset: 'USDC_POLYGON',
        fiat: 'USD',
        country: 'Nigeria',
        virtual_account_id: 'cmda40zuf1bd5pd24ruolj89s',
        status: 'Active',
        instructions: {
          currency: 'usd',
          bank_name: 'Bank of Nowhere',
          bank_address: '1800 North Pole St., Orlando, FL 32801',
          payment_rail: 'ach_push',
          payment_rails: ['ach_push', 'wire'],
          bank_account_number: '900458412690',
          bank_routing_number: '101019644',
          bank_beneficiary_name: 'Daniel Hobbs',
          bank_beneficiary_address:
            'Jesuloluwa Street, Green mast opposite Tresbol Filling station,, Ado Ekiti, EKITI 331106, NG',
        },
        failure_reason: null,
        has_accepted_tos: true,
        created_at: '2025-07-19T11:33:27.362066+01:00',
        updated_at: '2025-07-19T11:33:57.934709+01:00',
      },
    ],
    status: 'Fetched user accounts successfully',
  }

  const copyToClipboard = async (text: string, fieldName: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopiedField(fieldName)
      setTimeout(() => setCopiedField(null), 2000)
    } catch (err) {
      console.error('Failed to copy: ', err)
    }
  }
  if (userAccountLoading) {
    return <div className="text-center p-4">Loading...</div>
  }
  const accountData = virtualAccount.data[0]
  const instructions = accountData.instructions

  const CopyableField = ({
    label,
    value,
    fieldKey,
    isLong = false,
  }: {
    label: string
    value: string
    fieldKey: string
    isLong?: boolean
  }) => (
    <div className="mb-4">
      <div className="text-sm text-gray-600 mb-1">{label}</div>
      <div className="flex items-center justify-between bg-gray-50 p-3 rounded-lg border">
        <div
          className={`text-gray-900 font-medium ${
            isLong ? 'text-sm leading-5' : ''
          }`}
        >
          {value}
        </div>
        <button
          onClick={() => copyToClipboard(value, fieldKey)}
          className="ml-3 p-2 text-gray-500 hover:text-orange-600 hover:bg-orange-50 rounded-full transition-colors"
          title="Copy to clipboard"
        >
          {copiedField === fieldKey ? (
            <Check className="w-4 h-4 text-green-600" />
          ) : (
            <Copy className="w-4 h-4" />
          )}
        </button>
      </div>
    </div>
  )

  return (
    <div className="max-w-xl mx-auto pt-10 bg-white shadow-md min-h-screen">
      {/* Header */}
      <div className="flex items-center p-4 border-b">
        <ArrowLeft
          className="w-6 h-6 text-gray-700 mr-3"
          onClick={() => navigate(-1)}
        />

        <div className="text-lg font-semibold text-gray-900">
          Bank Account Details
        </div>
      </div>

      {/* Account Balance Section */}
      <div className="p-6 border-b">
        <div className="text-sm text-gray-500 mb-1">AVAILABLE BALANCE</div>
        <div className="text-4xl font-light text-gray-900 mb-2">$0</div>
        {/* <div className="text-sm text-gray-500 mb-4">Pending Balance: $0</div> */}

        <div className="flex space-x-3">
          <button className="flex-1 bg-orange-600 text-white py-2 px-4 rounded-lg text-sm font-medium">
            Account Statements
          </button>
          <button className="flex-1 border border-orange-600 text-orange-600 py-2 px-4 rounded-lg text-sm font-medium">
            Proof of account
          </button>
        </div>
      </div>

      {/* Bank Details Section */}
      <div className="p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          Your virtual account details
        </h3>

        <CopyableField
          label="Account Holder"
          value={
            userAccountDetails?.data[0]?.instructions?.bank_beneficiary_name
          }
          fieldKey="account_holder"
        />

        <CopyableField
          label="Account Number"
          value={userAccountDetails?.data[0]?.instructions?.bank_account_number}
          fieldKey="account_number"
        />

        <CopyableField
          label="Sort Code"
          value={userAccountDetails?.data[0]?.instructions?.bank_routing_number}
          fieldKey="sort_code"
        />

        {/* <CopyableField
          label="Swift Code"
          value="CLJU8521XXX"
          fieldKey="swift_code"
        /> */}

        <CopyableField
          label="IBAN"
          value={`${userAccountDetails?.data[0]?.instructions?.bank_routing_number}${userAccountDetails?.data[0]?.instructions?.bank_account_number}`}
          fieldKey="iban"
        />

        <CopyableField
          label="Bank name"
          value={userAccountDetails?.data[0]?.instructions?.bank_name}
          fieldKey="bank_name"
        />

        <CopyableField
          label="Address"
          value={userAccountDetails?.data[0]?.instructions?.bank_address}
          fieldKey="address"
          isLong={true}
        />
      </div>

      {/* Success Message */}
      {copiedField && (
        <div className="fixed bottom-4 left-1/2 transform -translate-x-1/2 bg-orange-600 text-white px-4 py-2 rounded-lg shadow-lg">
          Copied to clipboard!
        </div>
      )}
    </div>
  )
}

export default BankAccountDetails
