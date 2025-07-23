import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import SignupPage from './pages/SignupPage';
import LoginPage from './pages/LoginPage';
import TodoPage from './pages/TodoPage';
import Navbar from './components/Navbar';
import PrivateRoute from './components/PrivateRoute'; // Add this line

function App() {
  return (
    <Router>
      <Navbar />
      <Routes>
        <Route path="/signup" element={<SignupPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/" element={<PrivateRoute><TodoPage /></PrivateRoute>} />
      </Routes>
    </Router>
  );
}

export default App;