import { useNavigate } from 'react-router-dom'
import { GHFlag, KENFlag, TANZFlag, NGFlag } from '../components/icons/Icons'

export const depositData = [
  {
    id: 1,
    date: '1/12/24',
    details: [
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
      },
    ],
  },
  {
    id: 2,
    date: '1/12/24',
    details: [
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
      },
    ],
  },
]

export const withdrawalData = [
  {
    id: 1,
    date: '1/12/24',
    details: [
      {
        label: 'Local Withdrawal',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Withdrawal',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Withdrawal',
        time: '2:03 PM',
        amount: '$2400',
      },
    ],
  },
  {
    id: 2,
    date: '1/12/24',
    details: [
      {
        label: 'Local Withdrawal',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Withdrawal',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Withdrawal',
        time: '2:03 PM',
        amount: '$2400',
      },
    ],
  },
]

export const historyData = [
  {
    id: 1,
    date: '1/12/24',
    details: [
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Withdrawal',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
      },
    ],
  },
  {
    id: 2,
    date: '1/12/24',
    details: [
      {
        label: 'Local Withdrawal',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Withdrawal',
        time: '2:03 PM',
        amount: '$2400',
      },
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
      },
    ],
  },
]

export const allTransactionData = [
  {
    id: 1,
    date: '1/12/24',
    details: [
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
        channel: 'Deposit Via Paychant',
        status: 'Success',
      },
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
        channel: 'Deposit Via Paychant',
        status: 'Pending',
      },
    ],
  },
  {
    id: 2,
    date: '1/12/24',
    details: [
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
        channel: 'Deposit Via Paychant',
        status: 'Failed',
      },
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
        channel: 'Deposit Via Paychant',
        status: 'Success',
      },
    ],
  },
  {
    id: 3,
    date: '1/12/24',
    details: [
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
        channel: 'Deposit Via Paychant',
        status: 'Success',
      },
      {
        label: 'Local Transfer',
        time: '2:03 PM',
        amount: '$2400',
        channel: 'Deposit Via Paychant',
        status: 'Success',
      },
    ],
  },
]

export const notificationsData = [
  {
    text: 'Nothing for now maintenance | between  now and may 2020',
    subText: 'Lorem Ipsum is a dummy text...',
    date: '18 Jan',
  },
  {
    text: 'Scheduled maintenance | between  now and may 2020',
    subText: 'Lorem Ipsum is a dummy text...',
    date: '18 Jan',
  },
  {
    text: 'Scheduled maintenance | between  now and may 2020',
    subText: 'Lorem Ipsum is a dummy text...',
    date: '18 Jan',
  },
  {
    text: 'Scheduled maintenance | between  now and may 2020',
    subText: 'Lorem Ipsum is a dummy text...',
    date: '18 Jan',
  },
]

export const SettingsData = () => {
  const navigate = useNavigate()
  return [
    {
      text: 'Welcome to Greybox',
      subText: 'Lorem Ipsum is a dummy text...',
      onClick: () => navigate(''),
    },
    {
      text: 'Update Wallet Details',
      subText: 'Lorem Ipsum is a dummy text...',
      onClick: () => navigate('/update-wallet-details'),
    },
    {
      text: 'Change Pin',
      subText: 'Lorem Ipsum is a dummy text...',
      onClick: () => navigate(''),
    },
    {
      text: 'Export Private Key',
      subText: 'Lorem Ipsum is a dummy text...',
      onClick: () => navigate(''),
    },
    {
      text: 'About Greybox',
      subText: 'Lorem Ipsum is a dummy text...',
      onClick: () => navigate('/about-greybox'),
    },
    {
      text: 'Sign Out',
      subText: 'Lorem Ipsum is a dummy text...',
      onClick: () => navigate('/'),
    },
  ]
}

export const banksData = [
  {
    name: 'First Bank of Nigeria',
  },
  {
    name: 'United Bank for Africa',
  },
  {
    name: 'Zenith Bank PLC',
  },
  {
    name: 'Kuda Bank',
  },
]
export const momoData = [
  {
    name: 'Momo Agent 1',
  },
  {
    name: 'Momo Agent 2',
  },
  {
    name: 'Momo Agent 3',
  },
  {
    name: 'Momo Agent 4',
  },
]

export const countryData = [
  {
    name: 'Ghana',
    logo: GHFlag(),
    code: 'GH',
  },
  {
    name: 'Nigeria',
    logo: NGFlag(),
    code: 'NG',
  },
  {
    name: 'Kenya',
    logo: KENFlag(),
    code: 'KEN',
  },
  {
    name: 'Tanzania',
    logo: TANZFlag(),
    code: 'TZ',
  },
]

export const currencyData = [
  {
    name: 'Ghana Cedi',
    logo: GHFlag(),
    code: 'GHS',
  },
  {
    name: 'Nigerian Naira',
    logo: NGFlag(),
    code: 'NGN',
  },
  {
    name: 'Kenyan Shilling',
    logo: KENFlag(),
    code: 'KES',
  },
  {
    name: 'Tanzania Shilling',
    logo: TANZFlag(),
    code: 'TZS',
  },
  {
    name: 'US Dollar',
    logo: 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a4/Flag_of_the_United_States.svg/512px-Flag_of_the_United_States.svg.png',
    code: 'USD',
  },
]

export const currencyDataT = [
  {
    name: 'Ghana Cedi',
    logo: GHFlag(),
    code: 'GHS',
  },
  {
    name: 'Nigerian Naira',
    logo: NGFlag(),
    code: 'NGN',
  },
  {
    name: 'Kenyan Shilling',
    logo: KENFlag(),
    code: 'KES',
  },
  {
    name: 'Tanzania Shilling',
    logo: TANZFlag(),
    code: 'TZS',
  },
  {
    name: 'US Dollar',
    logo: '$',
    code: 'USD',
  },
]

export const FIAT = [
  {
    name: 'US',
  },
  {
    name: 'USD',
  },
]

export const ASSETS = [
  {
    name: 'ETH',
  },
  {
    name: 'USDC_BASE',
  },
  {
    name: 'USDC_ETHERIUM',
  },
  {
    name: 'USDC_OPTIMISM',
  },
  {
    name: 'USDC_POLYGON',
  },
  {
    name: 'USDC_SOLANA',
  },
  {
    name: 'USDT_ETHERIUM',
  },
]

export const BorderlessCountry = [
  { code: 'US', name: 'United States' },
  { code: 'NG', name: 'Nigeria' },
  { code: 'AD', name: 'Andorra' },
  { code: 'AL', name: 'Albania' },
  { code: 'CH', name: 'Switzerland' },
  { code: 'GB', name: 'United Kingdom' },
  { code: 'IS', name: 'Iceland' },
  { code: 'LI', name: 'Liechtenstein' },
  { code: 'MC', name: 'Monaco' },
  { code: 'MD', name: 'Moldova' },
  { code: 'ME', name: 'Montenegro' },
  { code: 'MK', name: 'North Macedonia' },
  { code: 'NO', name: 'Norway' },
  { code: 'SM', name: 'San Marino' },
  { code: 'VA', name: 'Vatican City' },
  { code: 'FR', name: 'France' },
  { code: 'DE', name: 'Germany' },
  { code: 'IT', name: 'Italy' },
  { code: 'DK', name: 'Denmark' },
  { code: 'PT', name: 'Portugal' },
  { code: 'AT', name: 'Austria' },
  { code: 'BE', name: 'Belgium' },
  { code: 'BG', name: 'Bulgaria' },
  { code: 'CY', name: 'Cyprus' },

  { code: 'CZ', name: 'Czech Republic' },
  { code: 'EE', name: 'Estonia' },
  { code: 'ES', name: 'Spain' },
  { code: 'FI', name: 'Finland' },
  { code: 'GR', name: 'Greece' },
  { code: 'HR', name: 'Croatia' },
  { code: 'HU', name: 'Hungary' },
  { code: 'IE', name: 'Ireland' },
  { code: 'LT', name: 'Lithuania' },
  { code: 'LU', name: 'Luxembourg' },
  { code: 'LV', name: 'Latvia' },
  { code: 'MT', name: 'Malta' },
  { code: 'NL', name: 'Netherlands' },
  { code: 'PL', name: 'Poland' },
  { code: 'RO', name: 'Romania' },
  { code: 'SE', name: 'Sweden' },
  { code: 'SI', name: 'Slovenia' },
  { code: 'SK', name: 'Slovakia' },
  { code: 'CL', name: 'Chile' },
  { code: 'CO', name: 'Colombia' },
  { code: 'MX', name: 'Mexico' },
  { code: 'PE', name: 'Peru' },
  { code: 'AR', name: 'Argentina' },
  { code: 'BJ', name: 'Benin' },
  { code: 'UG', name: 'Uganda' },
  { code: 'ZA', name: 'South Africa' },
  { code: 'ZM', name: 'Zambia' },
  { code: 'CD', name: 'Democratic Republic of the Congo' },
  { code: 'CI', name: 'Ivory Coast' },
  { code: 'CM', name: 'Cameroon' },
  { code: 'GH', name: 'Ghana' },
  { code: 'KE', name: 'Kenya' },
  { code: 'MW', name: 'Malawi' },
  { code: 'RW', name: 'Rwanda' },
  { code: 'SN', name: 'Senegal' },
  { code: 'TZ', name: 'Tanzania' },
]
