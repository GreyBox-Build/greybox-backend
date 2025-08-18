const AuthLayout = ({ child }: { child: React.ReactNode }) => {
  return (
    <div className=" bg-gradient-to-br from-orange-50 to-orange-100 w-full min-h-[100vh] flex justify-center">
      {child}
    </div>
  )
}

export default AuthLayout
