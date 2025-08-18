import { useEffect, useState } from 'react'
import { Modal } from '../Modal' // Assuming this is your generic Modal component
import { Button } from '../ui/button'
import { useRejectKycMutation } from 'src/appSlices/apiSlice'
import { AlertCircle, CheckCircle, Loader2, X } from 'lucide-react' // Added X for modal close button
import { twMerge } from 'tailwind-merge'
import { enqueueSnackbar } from 'notistack'

// Define a type for the RTK Query error structure, if known
interface ApiError {
  status: number
  data: {
    message: string
    // ... other error fields
  }
}

const KYCAction = ({
  status,
  id,
  handleApprove,
  ApproveLoading,
}: {
  status: string
  id?: string | number
  handleApprove: (id?: string | number) => void
  ApproveLoading: boolean
}) => {
  const [isOpenModal, setIsModalOpen] = useState<boolean>(false)
  const [rejectKyc, { isLoading, isSuccess, isError, data, error, reset }] =
    useRejectKycMutation()

  const [email, setEmail] = useState('')
  const [rejectionReason, setRejectionReason] = useState('')

  // Handle successful rejection and close modal
  useEffect(() => {
    if (isSuccess) {
      enqueueSnackbar(data?.message || 'Rejection was successful!', {
        variant: 'success',
      })
      setIsModalOpen(false) // Close the modal on success
      setEmail('') // Clear form fields
      setRejectionReason('')
      reset() // Reset RTK Query state for next use
    }
  }, [isSuccess, data, reset]) // Add `data` to dependencies for `data?.message`

  // Handle error notification
  useEffect(() => {
    if (isError) {
      const errorMessage =
        (error as ApiError)?.data?.message ||
        (error as Error)?.message || // Fallback for generic JS error
        'Rejection was not successful!'
      enqueueSnackbar(errorMessage, { variant: 'error' })
    }
  }, [isError, error]) // Only run when isError or error object changes

  // Reset form fields when modal opens/closes or ID changes
  useEffect(() => {
    if (isOpenModal) {
      setEmail('')
      setRejectionReason('')
      reset() // Ensure RTK Query state is clean when opening the modal
    }
  }, [isOpenModal, id, reset]) // Added `id` as a dependency for clear state per KYC item

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!email || !rejectionReason) {
      enqueueSnackbar('Please fill in all fields.', { variant: 'warning' })
      return
    }

    try {
      // The .unwrap() method will throw the error if the mutation fails,
      // allowing the catch block to handle it. The success useEffect will handle the success case.
      await rejectKyc({
        id: id,
        email,
        rejection_reason: rejectionReason,
      }).unwrap()
    } catch (err) {
      // Error is handled by the useEffect(isError) and notistack directly.
      // We don't need another enqueueSnackbar here.
      console.error('Failed to reject KYC:', err)
    }
  }

  return (
    <>
      {status === 'Pending' ? (
        <p className="flex gap-2">
          <Button
            className="bg-orange-1 text-white"
            onClick={() => handleApprove(id)}
            disabled={ApproveLoading || isLoading} // Also disable if rejection is processing
          >
            {ApproveLoading ? 'Wait a bit' : 'Approve'}
          </Button>
          <Button
            className="bg-red-600 text-white"
            onClick={() => setIsModalOpen(true)}
            disabled={ApproveLoading || isLoading} // Also disable if rejection is processing
          >
            Reject
          </Button>
        </p>
      ) : (
        <p className="text-gray-500 dark:text-gray-400">No Action</p>
      )}

      <Modal isOpen={isOpenModal} onClose={() => setIsModalOpen(false)}>
        {/* Modal Header */}
        <div className="flex items-center justify-between pb-4 border-b border-gray-200 dark:border-gray-700 mb-4">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            Reject KYC for ID: {id}
          </h3>
          <button
            type="button"
            className="text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300 transition-colors"
            onClick={() => setIsModalOpen(false)}
            aria-label="Close dialog"
            disabled={isLoading} // Prevent closing while mutation is in progress
          >
            <X className="h-6 w-6" />
          </button>
        </div>

        {/* The form content */}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label
              htmlFor="email"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              Email:
            </label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-red-500 focus:border-red-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              placeholder="e.g., user@example.com"
            />
          </div>

          <div>
            <label
              htmlFor="rejectionReason"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              Rejection Reason:
            </label>
            <textarea
              id="rejectionReason"
              value={rejectionReason}
              onChange={(e) => setRejectionReason(e.target.value)}
              required
              rows={4}
              className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-red-500 focus:border-red-500 sm:text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white"
              placeholder="e.g., Invalid photo, Missing document, etc."
            ></textarea>
          </div>

          {/* Loading Indicator inside the form */}
          {isLoading && (
            <div className="flex items-center text-blue-600 dark:text-blue-400">
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              <p>Processing rejection...</p>
            </div>
          )}

          {/* The success/error alerts are now handled by Notistack,
              so you might not need to render them inside the modal,
              or keep them if you want immediate visual feedback within the form before close.
              I've removed the `isSuccess` and `isError` blocks from direct rendering
              since `enqueueSnackbar` is now used for this.
              If you want visual feedback *inside* the modal *before* it closes,
              you could put them back, but often the snackbar is enough.
          */}
          {/*
          {isError && (
            <div
              className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative dark:bg-red-800 dark:border-red-700 dark:text-red-100"
              role="alert"
            >
              <div className="flex items-center">
                <AlertCircle className="mr-2 h-5 w-5" />
                <strong className="font-bold">Error!</strong>
              </div>
              <span className="block sm:inline ml-1">
                {(error as any)?.data?.message ||
                  (error as any)?.message ||
                  'An unknown error occurred'}
              </span>
            </div>
          )}
          {isSuccess && (
            <div
              className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative dark:bg-green-800 dark:border-green-700 dark:text-green-100"
              role="alert"
            >
              <div className="flex items-center">
                <CheckCircle className="mr-2 h-5 w-5" />
                <strong className="font-bold">Success!</strong>
              </div>
              <span className="block sm:inline ml-1">
                {data?.message || 'KYC rejected successfully.'}
              </span>
            </div>
          )}
          */}

          {/* Action Buttons */}
          <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700 mt-6">
            <button
              type="button"
              onClick={() => setIsModalOpen(false)} // Directly close modal
              disabled={isLoading}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md border border-gray-300 hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 dark:bg-gray-700 dark:text-gray-200 dark:border-gray-600 dark:hover:bg-gray-600"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading || isSuccess} // Disable submit while loading or after success
              className={twMerge(
                'px-4 py-2 text-sm font-medium text-white rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2',
                isLoading
                  ? 'bg-red-400 cursor-not-allowed'
                  : 'bg-red-600 hover:bg-red-700 focus:ring-red-500',
                isSuccess && 'opacity-70 cursor-not-allowed'
              )}
            >
              {isLoading ? 'Rejecting...' : 'Reject KYC'}
            </button>
          </div>
        </form>
      </Modal>
    </>
  )
}

export default KYCAction
