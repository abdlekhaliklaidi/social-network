// import LogoutButton from '../components/LogoutButton';
import LoginForm from '../components/LoginForm';
import RegisterForm from '../components/RegisterForm';
// import AuthForm from '../components/AuthForm';

const HomePage = () => {
  return (
    <div>
      <h1>Welcome</h1>
      {/* <LogoutButton /> */}
      <LoginForm />
      <RegisterForm />
      {/* <AuthForm /> */}
    </div>
  );
};

export default HomePage;
