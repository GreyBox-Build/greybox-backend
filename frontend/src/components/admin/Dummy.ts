type OrderStatus = "Pending" | "Confirmed" | "Cancelled";

interface Transaction {
  serialNo: number;
  name: string;
  accountNo: string;
  amountNaira: string | number; // Naira symbol ₦, or numeric
  amountDollar: string | number; // Dollar amount $
  date: string;
  status: string;
  ref: string;
  action?: React.ReactNode; // Optional, will only appear for pending transactions
}
export const transactions: Transaction[] = [
  {
    serialNo: 1,
    name: "John Doe",
    accountNo: "1234567890",
    amountNaira: 100000,
    amountDollar: 200,
    date: "2023-10-01",
    status: "Pending",
    ref: "REF123",
  },
  {
    serialNo: 2,
    name: "Jane Smith",
    accountNo: "0987654321",
    amountNaira: 50000,
    amountDollar: 100,
    date: "2023-09-25",
    status: "Confirmed",
    ref: "REF456",
  },
  {
    serialNo: 3,
    name: "Samuel Johnson",
    accountNo: "1122334455",
    amountNaira: 70000,
    amountDollar: 150,
    date: "2023-09-20",
    status: "Cancelled",
    ref: "REF789",
  },
  {
    serialNo: 4,
    name: "Michael Lee",
    accountNo: "3344556677",
    amountNaira: 200000,
    amountDollar: 400,
    date: "2023-09-28",
    status: "Pending",
    ref: "REF101",
  },
  {
    serialNo: 5,
    name: "Emily Davis",
    accountNo: "7766554433",
    amountNaira: 150000,
    amountDollar: 300,
    date: "2023-09-29",
    status: "Confirmed",
    ref: "REF102",
  },
  {
    serialNo: 6,
    name: "David Wilson",
    accountNo: "2233445566",
    amountNaira: 85000,
    amountDollar: 170,
    date: "2023-10-02",
    status: "Pending",
    ref: "REF103",
  },
  {
    serialNo: 7,
    name: "Sophia Brown",
    accountNo: "9988776655",
    amountNaira: 20000,
    amountDollar: 240,
    date: "2023-09-30",
    status: "Cancelled",
    ref: "REF104",
  },
  {
    serialNo: 8,
    name: "Chris Martin",
    accountNo: "6655443322",
    amountNaira: 60000,
    amountDollar: 120,
    date: "2023-09-27",
    status: "Confirmed",
    ref: "REF105",
  },
  {
    serialNo: 9,
    name: "Olivia Clark",
    accountNo: "1111222233",
    amountNaira: 80000,
    amountDollar: 360,
    date: "2023-10-03",
    status: "Pending",
    ref: "REF106",
  },
  {
    serialNo: 10,
    name: "James Walker",
    accountNo: "4444555566",
    amountNaira: 95000,
    amountDollar: 190,
    date: "2023-10-04",
    status: "Cancelled",
    ref: "REF107",
  },
];

// // Sample handler to update the status
