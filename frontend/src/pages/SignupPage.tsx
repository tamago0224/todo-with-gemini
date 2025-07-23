import React from 'react';
import { useNavigate } from 'react-router-dom';
import AuthForm from '../components/AuthForm';
import api from '../services/api';

const SignupPage: React.FC = () => {
  const navigate = useNavigate();

  const handleSignup = async (username: string, password: string) => {
    try {
      await api.signup(username, password);
      alert('Signup successful! Please login.');
      navigate('/login');
    } catch (error: any) {
      alert(`Signup failed: ${error.message}`);
    }
  };

  return (
    <div>
      <h1>Signup Page</h1>
      <AuthForm isSignup={true} onSubmit={handleSignup} />
    </div>
  );
};

export default SignupPage;
