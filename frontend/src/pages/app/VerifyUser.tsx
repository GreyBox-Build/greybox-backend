// import React, { useState, useCallback } from 'react'
// import { useKycMutation } from '../../appSlices/apiSlice'
// import { Notification, UserLogo } from '../../components/icons/Icons'
// import AppLayout from './AppLayout'
// import { useNavigate } from 'react-router-dom'
// import { useGetAuthUserQuery } from '../../appSlices/apiSlice'
// import { FormButton } from 'src/components/buttons/FormButton'

// import {
//   Select,
//   SelectContent,
//   SelectGroup,
//   SelectItem,
//   SelectTrigger,
//   SelectValue,
// } from '../../components/ui/select'
// import { BorderlessCountry } from 'src/utils/Dummies'

// interface KycData {
//   email: string
//   id_type: string
//   issue_date: string
//   expiry_date: string
//   front_photo: File | null
//   back_photo: File | null
//   tax_id: string
//   id_number: string
//   date_of_birth: string
//   phone: string
//   street_address: string
//   city: string
//   state: string
//   postal_code: string
//   country: string
// }

// const VerifyUSer: React.FC = () => {
//   const navigate = useNavigate()

//   const { currentData: userData, isFetching } = useGetAuthUserQuery({})

//   const personInfo = userData?.data?.personal_details
//   const [kycData, setKycData] = useState<KycData>({
//     email: '',
//     id_type: '',
//     issue_date: '',
//     expiry_date: '',
//     front_photo: null,
//     back_photo: null,
//     tax_id: '',
//     id_number: '',
//     date_of_birth: '',
//     phone: '',
//     street_address: '',
//     city: '',
//     state: '',
//     postal_code: '',
//     country: '',
//   })
//   const [isSubmitting, setIsSubmitting] = useState(false)
//   const [submissionError, setSubmissionError] = useState<string | null>(null)
//   const [submissionSuccess, setSubmissionSuccess] = useState(false)
//   const [kyc, { isLoading }] = useKycMutation()

//   const handleInputChange = useCallback(
//     (field: keyof KycData, value: string) => {
//       setKycData((prevData) => ({
//         ...prevData,
//         [field]: value,
//       }))
//     },
//     []
//   )

//   const handleFileChange = useCallback(
//     (field: 'front_photo' | 'back_photo', file: File | null) => {
//       if (file && file.size > 5 * 1024 * 1024) {
//         alert('File size exceeds 5MB. Please upload a smaller file.')
//         return
//       }
//       setKycData((prevData) => ({
//         ...prevData,
//         [field]: file,
//       }))
//     },
//     []
//   )

//   const handleSubmit = useCallback(
//     async (e: React.FormEvent) => {
//       e.preventDefault()
//       setIsSubmitting(true)
//       setSubmissionError(null)
//       setSubmissionSuccess(false)

//       // Basic validation
//       const requiredFields: (keyof KycData)[] = [
//         'email',
//         'id_type',
//         'issue_date',
//         'expiry_date',
//         'front_photo',
//         'back_photo',
//         'id_number',
//         'date_of_birth',
//         'phone',
//         'street_address',
//         'city',
//         'state',
//         'postal_code',
//         'country',
//       ]
//       if (requiredFields.some((field) => !kycData[field])) {
//         setSubmissionError('Please fill in all required fields.')
//         setIsSubmitting(false)
//         return
//       }

//       try {
//         const formData = new FormData()
//         for (const key in kycData) {
//           const value = kycData[key as keyof KycData]
//           if (value instanceof File) {
//             formData.append(key, value)
//           } else if (value) {
//             formData.append(key, value)
//           }
//         }

//         const result = await kyc(formData).unwrap() // RTK Query unwraps to data

//         setSubmissionSuccess(true)
//       } catch (error: any) {
//         console.error('Submission error:', error)
//         setSubmissionError(
//           error?.data?.message ||
//             'An error occurred during submission. Please try again.'
//         )
//       } finally {
//         setIsSubmitting(false)
//       }
//     },
//     [kycData, kyc]
//   )

//   return (
//     <AppLayout
//       child={
//         <div className="w-full max-w-xl bg-white px-[25px] pt-[51px] ">
//           <section className="w- flex justify-between items-center">
//             <Notification
//               current={0}
//               onClick={() => navigate('/notifications')}
//             />
//             <div>
//               <p className="text-center text-[0.875rem] text-gray-800 leading-[18px]">
//                 Welcome,
//               </p>
//               <p className="text-center text-gray-800 text-[1.25rem] font-[600]">
//                 {isFetching
//                   ? null
//                   : `${personInfo?.first_name} ${personInfo?.last_name}`}
//               </p>
//             </div>
//             <UserLogo
//               firstCharacterOfName={personInfo?.first_name?.substring(0, 1)}
//               onClick={() => navigate('/settings')}
//             />
//           </section>

//           {/* Form */}
//           <div className="min-h-screen  py-10 mt-4 rounded-lg ">
//             <div className="max-w-4xl mx-auto px-4 bg-white shadow-md p-6 rounded-[20px]">
//               <h2 className="text-2xl font-semibold mb-6 text-orange-1 text-center">
//                 KYC Verification Form
//               </h2>
//               <form onSubmit={handleSubmit} className="space-y-4">
//                 {/* Email */}
//                 <div>
//                   <label
//                     htmlFor="email"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Email <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="email"
//                     id="email"
//                     value={kycData.email}
//                     onChange={(e) => handleInputChange('email', e.target.value)}
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                   />
//                 </div>

//                 {/* ID Type */}
//                 <div>
//                   <label
//                     htmlFor="id_type"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     ID Type <span className="text-red-500">*</span>
//                   </label>

//                   <Select
//                     onValueChange={(value) =>
//                       handleInputChange('id_type', value)
//                     }
//                   >
//                     <SelectTrigger className="w-full border-gray-200 !py-5">
//                       <SelectValue placeholder="Select ID Type" />
//                     </SelectTrigger>
//                     <SelectContent className="bg-white">
//                       <SelectGroup>
//                         <SelectItem
//                           value="Passport"
//                           className="hover:bg-orange-3 border-b"
//                         >
//                           Passport
//                         </SelectItem>
//                         <SelectItem
//                           value="DriversLicense"
//                           className="hover:bg-orange-3 border-b"
//                         >
//                           Drivers License
//                         </SelectItem>
//                         <SelectItem
//                           value="NationalId"
//                           className="hover:bg-orange-3 "
//                         >
//                           National ID
//                         </SelectItem>
//                         <SelectItem
//                           value="ResidencePermit"
//                           className="hover:bg-orange-3 "
//                         >
//                           Residence Permit
//                         </SelectItem>
//                         <SelectItem
//                           value="ProofOfAddress"
//                           className="hover:bg-orange-3 "
//                         >
//                           Proof Of Address
//                         </SelectItem>
//                         <SelectItem
//                           value="Ownership"
//                           className="hover:bg-orange-3 "
//                         >
//                           Ownership
//                         </SelectItem>
//                         <SelectItem
//                           value="Formation"
//                           className="hover:bg-orange-3 "
//                         >
//                           Formation
//                         </SelectItem>
//                       </SelectGroup>
//                     </SelectContent>
//                   </Select>
//                 </div>

//                 {/* Issue Date */}
//                 <div>
//                   <label
//                     htmlFor="issue_date"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Issue Date <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="date"
//                     id="issue_date"
//                     value={kycData.issue_date}
//                     onChange={(e) =>
//                       handleInputChange('issue_date', e.target.value)
//                     }
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                   />
//                 </div>

//                 {/* Expiry Date */}
//                 <div>
//                   <label
//                     htmlFor="expiry_date"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Expiry Date <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="date"
//                     id="expiry_date"
//                     value={kycData.expiry_date}
//                     onChange={(e) =>
//                       handleInputChange('expiry_date', e.target.value)
//                     }
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                   />
//                 </div>

//                 {/* Front Photo */}
//                 <div>
//                   <label
//                     htmlFor="front_photo"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Front Photo <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="file"
//                     id="front_photo"
//                     onChange={(e) =>
//                       handleFileChange(
//                         'front_photo',
//                         e.target.files ? e.target.files[0] : null
//                       )
//                     }
//                     className="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-orange-50 file:text-orange-700 hover:file:bg-orange-100"
//                     accept="image/*"
//                     required
//                   />
//                   {kycData.front_photo && (
//                     <p className="text-xs text-gray-500">
//                       Selected: {kycData.front_photo.name}
//                     </p>
//                   )}
//                 </div>

//                 {/* Back Photo */}
//                 <div>
//                   <label
//                     htmlFor="back_photo"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Back Photo <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="file"
//                     id="back_photo"
//                     onChange={(e) =>
//                       handleFileChange(
//                         'back_photo',
//                         e.target.files ? e.target.files[0] : null
//                       )
//                     }
//                     className="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-orange-50 file:text-orange-700 hover:file:bg-orange-100"
//                     accept="image/*"
//                     required
//                   />
//                   {kycData.back_photo && (
//                     <p className="text-xs text-gray-500">
//                       Selected: {kycData.back_photo.name}
//                     </p>
//                   )}
//                 </div>

//                 {/* Tax ID */}
//                 <div>
//                   <label
//                     htmlFor="tax_id"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Tax ID
//                   </label>
//                   <input
//                     type="text"
//                     id="tax_id"
//                     value={kycData.tax_id}
//                     onChange={(e) =>
//                       handleInputChange('tax_id', e.target.value)
//                     }
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                   />
//                 </div>

//                 {/* ID Number */}
//                 <div>
//                   <label
//                     htmlFor="id_number"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     ID Number <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="text"
//                     id="id_number"
//                     value={kycData.id_number}
//                     onChange={(e) =>
//                       handleInputChange('id_number', e.target.value)
//                     }
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                   />
//                 </div>

//                 {/* Date of Birth */}
//                 <div>
//                   <label
//                     htmlFor="date_of_birth"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Date of Birth <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="date"
//                     id="date_of_birth"
//                     value={kycData.date_of_birth}
//                     onChange={(e) =>
//                       handleInputChange('date_of_birth', e.target.value)
//                     }
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                   />
//                 </div>

//                 {/* Phone */}
//                 <div>
//                   <label
//                     htmlFor="phone"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Phone <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="tel"
//                     id="phone"
//                     value={kycData.phone}
//                     onChange={(e) => handleInputChange('phone', e.target.value)}
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                   />
//                 </div>

//                 {/* Street Address */}
//                 <div>
//                   <label
//                     htmlFor="street_address"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Street Address <span className="text-red-500">*</span>
//                   </label>
//                   <textarea
//                     id="street_address"
//                     value={kycData.street_address}
//                     onChange={(e) =>
//                       handleInputChange('street_address', e.target.value)
//                     }
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                     rows={3}
//                   />
//                 </div>

//                 {/* City */}
//                 <div>
//                   <label
//                     htmlFor="city"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     City <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="text"
//                     id="city"
//                     value={kycData.city}
//                     onChange={(e) => handleInputChange('city', e.target.value)}
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                   />
//                 </div>

//                 {/* State */}
//                 <div>
//                   <label
//                     htmlFor="state"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     State <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="text"
//                     id="state"
//                     value={kycData.state}
//                     onChange={(e) => handleInputChange('state', e.target.value)}
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                   />
//                 </div>

//                 {/* Postal Code */}
//                 <div>
//                   <label
//                     htmlFor="postal_code"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Postal Code <span className="text-red-500">*</span>
//                   </label>
//                   <input
//                     type="text"
//                     id="postal_code"
//                     value={kycData.postal_code}
//                     onChange={(e) =>
//                       handleInputChange('postal_code', e.target.value)
//                     }
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                   />
//                 </div>

//                 {/* Country */}
//                 <div>
//                   <label
//                     htmlFor="country"
//                     className="block text-sm font-medium text-gray-700"
//                   >
//                     Country <span className="text-red-500">*</span>
//                   </label>
//                   <select
//                     id="country"
//                     value={kycData.country}
//                     onChange={(e) =>
//                       handleInputChange('country', e.target.value)
//                     }
//                     className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-orange-2 focus:border-orange-2 sm:text-sm"
//                     required
//                   >
//                     <option value="">Select a country</option>
//                     {BorderlessCountry.map((country) => (
//                       <option key={country.code} value={country.code}>
//                         {country.name}
//                       </option>
//                     ))}
//                   </select>
//                 </div>

//                 {/* Submit Button */}

//                 <FormButton
//                   label={isLoading || isSubmitting ? 'Submitting' : 'Submit'}
//                   extraClass="mt-[72px]"
//                   loading={isLoading || isSubmitting}
//                 />
//                 {/* Error and Success Messages */}
//                 {submissionError && (
//                   <p className="text-red-500">{submissionError}</p>
//                 )}
//                 {submissionSuccess && (
//                   <p className="text-green-500">
//                     KYC data submitted successfully!
//                   </p>
//                 )}
//               </form>
//             </div>
//           </div>
//         </div>
//       }
//     />
//   )
// }

// export default VerifyUSer
import React, { useState, useCallback } from 'react'
import { ChevronDown, Calendar, Upload, X, Bell, User } from 'lucide-react'
import DatePicker from 'react-datepicker'
import 'react-datepicker/dist/react-datepicker.css'

// Mock data for countries
import { BorderlessCountry } from 'src/utils/Dummies'
import AppLayout from './AppLayout'
import { useKycMutation } from '../../appSlices/apiSlice'
import { useNavigate } from 'react-router-dom'

interface KycData {
  email: string
  id_type: string
  issue_date: string
  expiry_date: string
  front_photo: File | null
  back_photo: File | null
  tax_id: string
  id_number: string
  date_of_birth: string
  phone: string
  street_address: string
  city: string
  state: string
  postal_code: string
  country: string
}

// Reusable CustomSelect component
const CustomSelect: React.FC<{
  value: string | null
  onChange: (value: string) => void
  options: { value: string; label: string }[]
  placeholder: string
  label: string
  required?: boolean
}> = ({ value, onChange, options, placeholder, label, required = false }) => {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <div className="relative">
      <label className="block text-sm font-medium text-gray-700 mb-2">
        {label} {required && <span className="text-red-500">*</span>}
      </label>
      <div
        className="relative w-full border border-gray-300 rounded-lg shadow-sm py-3 px-4 bg-white cursor-pointer hover:border-orange-400 focus-within:ring-2 focus-within:ring-orange-500 focus-within:border-orange-500 transition-all duration-200"
        onClick={() => setIsOpen(!isOpen)}
      >
        <span
          className={`block truncate ${
            !value ? 'text-gray-400' : 'text-gray-900'
          }`}
        >
          {value
            ? options.find((opt) => opt.value === value)?.label
            : placeholder}
        </span>
        <ChevronDown
          className={`absolute right-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400 transition-transform duration-200 ${
            isOpen ? 'rotate-180' : ''
          }`}
        />
      </div>

      {isOpen && (
        <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-auto">
          {options.map((option) => (
            <div
              key={option.value}
              className="px-4 py-3 hover:bg-orange-50 cursor-pointer transition-colors duration-150 border-b border-gray-100 last:border-b-0"
              onClick={() => {
                onChange(option.value)
                setIsOpen(false)
              }}
            >
              {option.label}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

// React-DatePicker wrapper component
const CustomDatePicker: React.FC<{
  value: string
  onChange: (value: string) => void
  label: string
  required?: boolean
  placeholder?: string
}> = ({
  value,
  onChange,
  label,
  required = false,
  placeholder = 'Select date',
}) => {
  // Convert string to Date for DatePicker
  const selectedDate = value ? new Date(value) : null

  // Handle date change and convert back to string
  const handleDateChange = (date: Date | null) => {
    if (date) {
      const formattedDate = date.toISOString().split('T')[0]
      onChange(formattedDate)
    } else {
      onChange('')
    }
  }

  return (
    <div className="relative">
      <label className="block text-sm font-medium text-gray-700 mb-2">
        {label} {required && <span className="text-red-500">*</span>}
      </label>
      <div className="relative">
        <DatePicker
          selected={selectedDate}
          onChange={handleDateChange}
          showYearDropdown
          showMonthDropdown
          dropdownMode="select"
          yearDropdownItemNumber={15}
          dateFormat="MMMM d, yyyy"
          placeholderText={placeholder}
          className="w-full border border-gray-300 rounded-lg shadow-sm py-3 px-4 pr-12 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all duration-200"
          wrapperClassName="w-full"
          calendarClassName="react-datepicker-custom"
        />
        <Calendar className="absolute right-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400 pointer-events-none" />
      </div>
    </div>
  )
}

// Reusable FileUpload component
const FileUpload: React.FC<{
  value: File | null
  onChange: (file: File | null) => void
  label: string
  required?: boolean
}> = ({ value, onChange, label, required = false }) => {
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null
    if (file && file.size > 5 * 1024 * 1024) {
      alert('File size exceeds 5MB. Please upload a smaller file.')
      return
    }
    onChange(file)
  }

  const removeFile = () => {
    onChange(null)
  }

  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-2">
        {label} {required && <span className="text-red-500">*</span>}
      </label>

      <div className="relative">
        <input
          type="file"
          onChange={handleFileChange}
          accept="image/*"
          className="hidden"
          id={`file-${label.replace(/\s+/g, '-').toLowerCase()}`}
        />

        <label
          htmlFor={`file-${label.replace(/\s+/g, '-').toLowerCase()}`}
          className="flex items-center justify-center w-full border-2 border-dashed border-gray-300 rounded-lg py-6 px-4 hover:border-orange-400 hover:bg-orange-50 transition-all duration-200 cursor-pointer"
        >
          <div className="text-center">
            <Upload className="mx-auto h-8 w-8 text-gray-400 mb-2" />
            <span className="text-sm text-gray-600">
              {value ? 'Change file' : 'Click to upload or drag and drop'}
            </span>
            <p className="text-xs text-gray-500 mt-1">
              PNG, JPG, GIF up to 5MB
            </p>
          </div>
        </label>

        {value && (
          <div className="mt-2 flex items-center justify-between bg-orange-50 border border-orange-200 rounded-lg px-3 py-2">
            <span className="text-sm text-gray-700 truncate">{value.name}</span>
            <button
              type="button"
              onClick={removeFile}
              className="ml-2 text-red-500 hover:text-red-700"
            >
              <X className="h-4 w-4" />
            </button>
          </div>
        )}
      </div>
    </div>
  )
}

const VerifyUser = () => {
  const [kycData, setKycData] = useState<KycData>({
    email: '',
    id_type: '',
    issue_date: '',
    expiry_date: '',
    front_photo: null,
    back_photo: null,
    tax_id: '',
    id_number: '',
    date_of_birth: '',
    phone: '',
    street_address: '',
    city: '',
    state: '',
    postal_code: '',
    country: '',
  })

  const [kyc, { isLoading }] = useKycMutation()
  const [submissionError, setSubmissionError] = useState<string | null>(null)
  const [submissionSuccess, setSubmissionSuccess] = useState(false)
  const navigate = useNavigate()

  const handleInputChange = useCallback(
    (field: keyof KycData, value: string | File | null) => {
      setKycData((prevData) => ({
        ...prevData,
        [field]: value,
      }))
    },
    []
  )

  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault()
      setSubmissionError(null)
      setSubmissionSuccess(false)

      // Basic validation
      const requiredFields: (keyof KycData)[] = [
        'email',
        'id_type',
        'issue_date',
        'expiry_date',
        'front_photo',
        'back_photo',
        'id_number',
        'date_of_birth',
        'phone',
        'street_address',
        'city',
        'state',
        'postal_code',
        'country',
      ]

      if (requiredFields.some((field) => !kycData[field])) {
        setSubmissionError('Please fill in all required fields.')
        return
      }

      try {
        const formData = new FormData()
        for (const key in kycData) {
          const value = kycData[key as keyof KycData]
          if (value instanceof File) {
            formData.append(key, value)
          } else if (value) {
            formData.append(key, value as string)
          }
        }

        const result = await kyc(formData).unwrap() // RTK Query unwraps to data

        setSubmissionSuccess(true)
        setTimeout(() => navigate('/dashboard'), 2000)
      } catch (error: any) {
        console.error('Submission error:', error)
        setSubmissionError(
          error?.data?.message ||
            'An error occurred during submission. Please try again.'
        )
      }
    },
    [kycData, kyc]
  )

  const idTypeOptions = [
    { value: 'Passport', label: 'Passport' },
    { value: 'DriversLicense', label: "Driver's License" },
    { value: 'NationalId', label: 'National ID' },
    { value: 'ResidencePermit', label: 'Residence Permit' },
    { value: 'ProofOfAddress', label: 'Proof of Address' },
  ]

  const countryOptions = BorderlessCountry.map((country) => ({
    value: country.code,
    label: country.name,
  }))

  return (
    <AppLayout
      child={
        <div className="min-h-screen bg-gradient-to-br from-orange-50 to-orange-100">
          {/* Custom CSS for React DatePicker */}

          {/* Header */}
          <div className="bg-white shadow-sm">
            <div className="max-w-7xl mx-auto px-6 py-4">
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
                    John Doe
                  </p>
                </div>

                <div className="w-10 h-10 bg-orange-500 rounded-full flex items-center justify-center">
                  <span className="text-white font-semibold">J</span>
                </div>
              </div>
            </div>
          </div>

          {/* Form */}
          <div className="max-w-4xl mx-auto px-4 py-8">
            <div className="bg-white rounded-2xl shadow-xl overflow-hidden">
              {/* Form Header */}
              <div className="bg-gradient-to-r from-orange-500 to-orange-600 px-8 py-6">
                <h2 className="text-2xl font-bold text-white text-center">
                  KYC Verification Form
                </h2>
                <p className="text-orange-100 text-center mt-2">
                  Please provide accurate information to verify your identity
                </p>
              </div>

              <form onSubmit={handleSubmit} className="p-8 space-y-6">
                {/* Personal Information Section */}
                <div className="space-y-6">
                  <h3 className="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">
                    Personal Information
                  </h3>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Email <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="email"
                        value={kycData.email}
                        onChange={(e) =>
                          handleInputChange('email', e.target.value)
                        }
                        className="w-full border border-gray-300 rounded-lg shadow-sm py-3 px-4 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all duration-200"
                        placeholder="Enter your email address"
                        required
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Phone <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="tel"
                        value={kycData.phone}
                        onChange={(e) =>
                          handleInputChange('phone', e.target.value)
                        }
                        className="w-full border border-gray-300 rounded-lg shadow-sm py-3 px-4 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all duration-200"
                        placeholder="Enter your phone number"
                        required
                      />
                    </div>
                  </div>

                  <CustomDatePicker
                    value={kycData.date_of_birth}
                    onChange={(date) =>
                      handleInputChange('date_of_birth', date)
                    }
                    label="Date of Birth"
                    required
                    placeholder="Select your date of birth"
                  />
                </div>

                {/* ID Information Section */}
                <div className="space-y-6">
                  <h3 className="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">
                    ID Information
                  </h3>

                  <CustomSelect
                    value={kycData.id_type}
                    onChange={(value) => handleInputChange('id_type', value)}
                    options={idTypeOptions}
                    placeholder="Select ID Type"
                    label="ID Type"
                    required
                  />

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      ID Number <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      value={kycData.id_number}
                      onChange={(e) =>
                        handleInputChange('id_number', e.target.value)
                      }
                      className="w-full border border-gray-300 rounded-lg shadow-sm py-3 px-4 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all duration-200"
                      placeholder="Enter your ID number"
                      required
                    />
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <CustomDatePicker
                      value={kycData.issue_date}
                      onChange={(date) => handleInputChange('issue_date', date)}
                      label="Issue Date"
                      required
                      placeholder="Select issue date"
                    />

                    <CustomDatePicker
                      value={kycData.expiry_date}
                      onChange={(date) =>
                        handleInputChange('expiry_date', date)
                      }
                      label="Expiry Date"
                      required
                      placeholder="Select expiry date"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Tax ID
                    </label>
                    <input
                      type="text"
                      value={kycData.tax_id}
                      onChange={(e) =>
                        handleInputChange('tax_id', e.target.value)
                      }
                      className="w-full border border-gray-300 rounded-lg shadow-sm py-3 px-4 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all duration-200"
                      placeholder="Enter your tax ID (optional)"
                    />
                  </div>
                </div>

                {/* Document Upload Section */}
                <div className="space-y-6">
                  <h3 className="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">
                    Document Upload
                  </h3>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <FileUpload
                      value={kycData.front_photo}
                      onChange={(file) =>
                        handleInputChange('front_photo', file)
                      }
                      label="Front Photo"
                      required
                    />

                    <FileUpload
                      value={kycData.back_photo}
                      onChange={(file) => handleInputChange('back_photo', file)}
                      label="Back Photo"
                      required
                    />
                  </div>
                </div>

                {/* Address Information Section */}
                <div className="space-y-6">
                  <h3 className="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">
                    Address Information
                  </h3>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Street Address <span className="text-red-500">*</span>
                    </label>
                    <textarea
                      value={kycData.street_address}
                      onChange={(e) =>
                        handleInputChange('street_address', e.target.value)
                      }
                      className="w-full border border-gray-300 rounded-lg shadow-sm py-3 px-4 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all duration-200"
                      placeholder="Enter your full street address"
                      rows={3}
                      required
                    />
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        City <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="text"
                        value={kycData.city}
                        onChange={(e) =>
                          handleInputChange('city', e.target.value)
                        }
                        className="w-full border border-gray-300 rounded-lg shadow-sm py-3 px-4 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all duration-200"
                        placeholder="Enter your city"
                        required
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        State <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="text"
                        value={kycData.state}
                        onChange={(e) =>
                          handleInputChange('state', e.target.value)
                        }
                        className="w-full border border-gray-300 rounded-lg shadow-sm py-3 px-4 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all duration-200"
                        placeholder="Enter your state"
                        required
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Postal Code <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="text"
                        value={kycData.postal_code}
                        onChange={(e) =>
                          handleInputChange('postal_code', e.target.value)
                        }
                        className="w-full border border-gray-300 rounded-lg shadow-sm py-3 px-4 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all duration-200"
                        placeholder="Enter your postal code"
                        required
                      />
                    </div>

                    <CustomSelect
                      value={kycData.country}
                      onChange={(value) => handleInputChange('country', value)}
                      options={countryOptions}
                      placeholder="Select your country"
                      label="Country"
                      required
                    />
                  </div>
                </div>

                {/* Submit Button */}
                <div className="pt-6">
                  <button
                    type="submit"
                    disabled={isLoading}
                    className="w-full bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold py-4 px-6 rounded-lg shadow-lg transition-all duration-200 transform hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                  >
                    {isLoading ? (
                      <div className="flex items-center justify-center">
                        <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                        Submitting...
                      </div>
                    ) : (
                      'Submit KYC Application'
                    )}
                  </button>
                </div>

                {/* Status Messages */}
                {submissionError && (
                  <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                    <p className="text-red-600 text-sm">{submissionError}</p>
                  </div>
                )}

                {submissionSuccess && (
                  <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                    <p className="text-green-600 text-sm font-medium">
                      🎉 KYC data submitted successfully! We'll review your
                      application and get back to you soon.
                    </p>
                  </div>
                )}
              </form>
            </div>
          </div>
        </div>
      }
    />
  )
}

export default VerifyUser
