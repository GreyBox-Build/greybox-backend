// KYC.tsx
import React, { useState } from 'react'
import { useApproveKycMutation, useGetKycQuery } from 'src/appSlices/apiSlice'
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table'

import KYCAction from 'src/components/admin/KYCAction'

const KYC = () => {
  const { data, isLoading } = useGetKycQuery({})
  const [approveKyc] = useApproveKycMutation()
  const [approveLoadingStates, setApproveLoadingStates] = useState<
    Record<string | number, boolean>
  >({})

  const handleApprove = async (id: any) => {
    setApproveLoadingStates((prev) => ({ ...prev, [id]: true }))
    try {
      await approveKyc(id)
    } finally {
      setApproveLoadingStates((prev) => ({ ...prev, [id]: false }))
    }
  }


  return (
    <div>
      <div>
        {isLoading ? (
          <p> "Loading"</p>
        ) : (
          <Table>
            <TableCaption>A list of your recent invoices.</TableCaption>
            <TableHeader>
              <TableHead className="w-[100px]">ID TYPE</TableHead>

              <TableHead>ID NO</TableHead>

              <TableHead>DOB</TableHead>

              <TableHead className="">ISSUER DATE</TableHead>

              <TableHead className="">EXPIRY DATE</TableHead>

              <TableHead className="">FRONT PIC</TableHead>

              <TableHead className="">BACK PIC</TableHead>

              <TableHead className="">PHone NO</TableHead>

              <TableHead className="">STATUS</TableHead>

              <TableHead className="">ACTION</TableHead>
            </TableHeader>
            <TableBody>
              {data?.kycs.map((kyc: any) => (
                <TableRow key={kyc.id}>
                  {/* ... table cells ... */}
                  <TableCell className="font-medium">{kyc.id_type}</TableCell>

                  <TableCell>{kyc.tax_id}</TableCell>

                  <TableCell>{kyc.id_number}</TableCell>

                  <TableCell>{kyc.date_of_birth}</TableCell>

                  <TableCell>{kyc.issue_date}</TableCell>

                  <TableCell>{kyc.expiry_date}</TableCell>

                  <TableCell>{kyc.country}</TableCell>

                  <TableCell>{kyc.phone}</TableCell>

                  <TableCell>{kyc.status}</TableCell>
                  <TableCell className="cursor-pointer">
                    <KYCAction
                      status={kyc.status}
                      id={kyc.id}
                      handleApprove={handleApprove}
                      ApproveLoading={approveLoadingStates[kyc.id] || false}
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </div>
    </div>
  )
}

export default KYC
