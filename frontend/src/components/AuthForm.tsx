import React, { useState } from 'react';

interface AuthFormProps {
  isSignup: boolean;
  onSubmit: (username: string, password: string) => void;
}

const AuthForm: React.FC<AuthFormProps> = ({ isSignup, onSubmit }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(username, password);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label htmlFor="username">Username:</label>
        <input
          type="text"
          id="username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
      </div>
      <div>
        <label htmlFor="password">Password:</label>
        <input
          type="password"
          id="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
      </div>
      <button type="submit">{isSignup ? 'Signup' : 'Login'}</button>
    </form>
  );
};

export default AuthForm;
