import { useNavigate } from 'react-router-dom'
import { Cards, Home, More, Send, Wallet } from '../../components/icons/Icons'
import { Menu } from '../../components/Cards'
import { useGetAuthUserQuery } from '../../appSlices/apiSlice'
import { useEffect } from 'react'
import { Customized } from 'src/components/customizeNav/Tabs'

const AppLayout = ({ child }: { child: React.ReactNode }) => {
  const navigate = useNavigate()
  const { error }: any = useGetAuthUserQuery({})
  useEffect(() => {
    setTimeout(() => {
      if (error?.status === 401) {
        localStorage.removeItem('access_token')
        navigate('/sign-in')
      }
    }, 500)
  })

  return (
    <div className=" bg-gradient-to-br from-orange-50 to-orange-100  w-full min-h-[100vh] flex justify-center">
      {child}

      <section className="w-full max-w-full mx-auto h-[54px] flex justify-between fixed bottom-0 bg-white p-[5px_24px] md:p-[5px_65px]">
        <Customized />
      </section>
    </div>
  )
}

export default AppLayout
