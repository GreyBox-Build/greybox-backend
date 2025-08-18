import { Home, Settings, IdCard, ChartArea } from 'lucide-react'
import { LimelightNav } from './CustomizeNav'

const customNavItems = [
  {
    id: 'home',
    icon: <Home />,
    label: 'Home',
    onClick: () => console.log('Home Clicked!'),
    navigate: '/dashboard',
  },

  {
    id: 'card',
    icon: <IdCard />,
    label: 'Add New',
    onClick: () => console.log('Add Clicked!'),
    navigate: '/virtual-card',
  },
  {
    id: 'history',
    icon: <ChartArea />,
    label: 'History',
    onClick: () => console.log('Profile Clicked!'),
    navigate: '/transaction-history',
  },
  {
    id: 'settings',
    icon: <Settings />,
    label: 'Settings',
    onClick: () => console.log('Settings Clicked!'),
    navigate: '/settings',
  },
]

const Customized = () => {
  return (
    <LimelightNav
      className="bg-secondary dark:bg-card/50 dark:border-accent/50 rounded-xl"
      items={customNavItems}
    />
  )
}

export { Customized }
