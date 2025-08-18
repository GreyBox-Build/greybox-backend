import React from "react";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "../../components/ui/table";
import { Button } from "../ui/button";
import KYCAction from "./KYCAction";
import { useApproveKycMutation } from "src/appSlices/apiSlice";

const KYCTable = ({ kyc }: { kyc: any }) => {
  const [approveKyc, { isLoading: ApproveLoading }] = useApproveKycMutation();
  const handleApprove = (id: any) => {
  
    approveKyc(id);
  };

  return (
    <Table>
      <TableCaption>A list of your recent invoices.</TableCaption>
      <TableHeader>
        <TableRow>
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
        </TableRow>
      </TableHeader>
      <TableBody>
        {kyc.map((kyc: any) => (
          <TableRow key={kyc.id}>
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
                ApproveLoading={ApproveLoading}
              />
            </TableCell>

            {/* 
"id":10,"id_type":"passport","tax_id":"1111111","id_number":"11111111","date_of_birth":"1999-01-01","issue_date":"2024-01-01","expiry_date":"2025-12-31","front_photo":"","back_photo":"","status":"Pending","phone":"09066112233","street_address":"Somewhere","borderless_identity_id":"","city":"Port Harcourt","state":"Rivers","postal_code":"501002","country":"NG","rejection_reason":"Invalid Backphoto","created_at":"2025-03-26T12:47:07.122373+01:00","approved_at":"0001-01-01T00:13:35+00:13","rejected_at":"2025-03-26T12:49:21.16416+01:00","updated_at":"2025-03-26T12:50:01.006172+01:00"}, */}
          </TableRow>
        ))}
      </TableBody>
      <TableFooter>
        <TableRow>
          <TableCell colSpan={3}>Total</TableCell>
          <TableCell className="text-right">$2,500.00</TableCell>
        </TableRow>
      </TableFooter>
    </Table>
  );
};

export default KYCTable;
