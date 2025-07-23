import React from 'react';
import { useNavigate } from 'react-router-dom';
import AuthForm from '../components/AuthForm';
import api from '../services/api';
import { useAuth } from '../context/AuthContext';

const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleLogin = async (username: string, password: string) => {
    try {
      const response = await api.login(username, password);
      login(response.token);
      navigate('/'); // Redirect to todo page on successful login
    } catch (error: any) {
      alert(`Login failed: ${error.message}`);
    }
  };

  return (
    <div>
      <h1>Login Page</h1>
      <AuthForm isSignup={false} onSubmit={handleLogin} />
    </div>
  );
};

export default LoginPage;
