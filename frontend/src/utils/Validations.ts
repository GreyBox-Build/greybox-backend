import { z } from "zod";

export const createUserSchema = z.object({
  first_name: z.string().min(1, { message: "First name is required" }),
  last_name: z.string().min(1, { message: "Last name is required" }),
  email: z
    .string()
    .min(1, { message: "Email cannot be empty" })
    .email("This is not a valid email."),
  password: z
    .string()
    .refine(
      (value) => /^(?=.*[A-Z])(?=.*\d).{8,}$/.test(value ?? ""),
      "Password must have at least 8 characters, have at least a digit and at least an Upper case letter"
    ),
  currency: z.string().min(1, { message: "Currency is required" }),
  country: z.string().min(1, { message: "Country name is required" }),
  chain: z.string().min(1, { message: "Select chain" }),
});

export const mobileDepositSchema = z.object({
  phoneNumber: z.string(), // At least 10 digits
  countryCode: z
    .string()
    .length(2, "Country code must be exactly 2 characters"),
  network: z.enum(["MPESA", "AIRTEL", "TIGO", "VODAFONE", "MTN"]),
  amount: z.string().min(1, "Amount must be at least 1"),
});

export const obtainTokenSchema = z.object({
  email: z
    .string()
    .min(1, { message: "Email cannot be empty" })
    .email("This is not a valid email."),
  password: z
    .string()
    .refine(
      (value) => /^(?=.*[A-Z])(?=.*\d).{8,}$/.test(value ?? ""),
      "Password must have at least 8 characters, have at least a digit and at least an Upper case letter"
    ),
});

export const forgetPasswordSchema = z.object({
  email: z
    .string()
    .min(1, { message: "Email cannot be empty" })
    .email("This is not a valid email."),
});

export const resetPasswordSchema = z.object({
  password: z
    .string()
    .refine(
      (value) => /^(?=.*[A-Z])(?=.*\d).{8,}$/.test(value ?? ""),
      "Password must have at least 8 characters, have at least a digit and at least an Upper case letter"
    ),
});
export const sendBankSchema = z.object({
  bank: z.string().min(1, { message: "Select bank" }),
  currency: z.string().min(1, { message: "Currency is required" }),
  amount_to_send: z.string().min(1, { message: "Enter amount to send" }),
});

export const depositViaMobileSchema = z.object({
  amount: z
    .string()
    .refine((amount) => parseFloat(amount) !== 0, "Zero amount not allowed"),
});

export const depositViaBankTransferSchema = ({
  currency,
}: {
  currency: string | null;
}) => {
  const check = currency === "NGN" ? 1500 : 50;
  return z.object({
    amount: z
      .string()
      .refine(
        (amount) => parseFloat(amount?.replace(/,/g, "")) >= check,
        `Amount must not be less than ${check}`
      ),
  });
};

export const withdrawViaBankSchema = z.object({
  cryptoAmount: z
    .string()
    .refine(
      (cryptoAmount) => parseFloat(cryptoAmount?.replace(/,/g, "")) >= 1,
      "Amount must not be less than 1"
    ),
  bankName: z.string().min(1, { message: "Bank name is required" }),
  accountNumber: z.string().min(1, { message: "Account number is required" }),
  accountName: z.string().min(1, { message: "Account name is required" }),
});

export const withdrawViaMobileSchema = z.object({
  amount: z
    .string()
    .refine(
      (amount) => parseFloat(amount) >= 1,
      "Amount must not be less than 1406"
    ),
});

export const withdrawPaymentSchema = z.object({
  amount: z.string().min(1, { message: "Enter amount" }),
  account_address: z.string().min(1, { message: "Address is required" }),
  chain: z.string().min(1, { message: "Chain is required" }),
});
