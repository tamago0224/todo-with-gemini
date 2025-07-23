import { render, screen, fireEvent } from '@testing-library/react';
import AuthForm from './AuthForm';

describe('AuthForm', () => {
  test('renders signup form correctly', () => {
    render(<AuthForm isSignup={true} onSubmit={() => {}} />);
    expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /signup/i })).toBeInTheDocument();
  });

  test('renders login form correctly', () => {
    render(<AuthForm isSignup={false} onSubmit={() => {}} />);
    expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
  });

  test('calls onSubmit with correct values on signup', () => {
    const handleSubmit = jest.fn();
    render(<AuthForm isSignup={true} onSubmit={handleSubmit} />);

    fireEvent.change(screen.getByLabelText(/username/i), { target: { value: 'testuser' } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /signup/i }));

    expect(handleSubmit).toHaveBeenCalledWith('testuser', 'password123');
  });

  test('calls onSubmit with correct values on login', () => {
    const handleSubmit = jest.fn();
    render(<AuthForm isSignup={false} onSubmit={handleSubmit} />);

    fireEvent.change(screen.getByLabelText(/username/i), { target: { value: 'testuser' } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));

    expect(handleSubmit).toHaveBeenCalledWith('testuser', 'password123');
  });
});
