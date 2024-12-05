import LanguageSelector from "../../components/LanguageSelector";

const AuthLayout = ({ child }: { child: React.ReactNode }) => {
  return (
    <div className=" bg-grey-box-bg bg-cover bg-no-repeat w-full min-h-[100vh] flex justify-center">
      <LanguageSelector />
      {child}
    </div>
  );
};

export default AuthLayout;
