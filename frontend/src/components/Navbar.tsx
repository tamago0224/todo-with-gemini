import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Navbar: React.FC = () => {
  const { isAuthenticated, logout } = useAuth();

  return (
    <nav style={{ padding: '1rem', background: '#f0f0f0' }}>
      <Link to="/" style={{ marginRight: '1rem' }}>Home</Link>
      {isAuthenticated ? (
        <button onClick={logout}>Logout</button>
      ) : (
        <>
          <Link to="/signup" style={{ marginRight: '1rem' }}>Signup</Link>
          <Link to="/login">Login</Link>
        </>
      )}
    </nav>
  );
};

export default Navbar;
