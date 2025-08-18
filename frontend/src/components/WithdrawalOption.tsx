import AppLayout from '../pages/app/AppLayout'
import { useNavigate } from 'react-router-dom'
import { FlatButton } from './buttons/FlatButton'
import { FormButton } from './buttons/FormButton'

const WithdrawalOption = () => {
  const navigate = useNavigate()
  return (
    <AppLayout
      child={
        <div className="p-[25px] m-auto w-full h-fit md:w-[50.33%] lg:w-[45.33%] bg-white shadow-md  rounded-[32px_32px_0px_0px]">
          <FlatButton
            label="Withdraw Via Bank "
            extraClass="bg-orange-1 text-white mb-[8px]"
            onClick={() => navigate('/withdraw-via-bank')}
          />
          <FlatButton
            label="Withdraw Via Mobile Money"
            extraClass="bg-orange-1 text-white mb-[8px]"
            onClick={() => navigate('/withdraw-fund')}
          />

          <FormButton
            label="Cancel"
            className="bg-orange-1 text-white"
            onClick={() => navigate(-1)}
          />
        </div>
      }
    />
  )
}

export default WithdrawalOption
