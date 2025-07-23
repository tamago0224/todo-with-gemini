import { render, screen, fireEvent } from '@testing-library/react';
import AddTodoForm from './AddTodoForm';

describe('AddTodoForm', () => {
  test('calls onAdd with the new todo title when submitted', () => {
    const handleAdd = jest.fn();
    render(<AddTodoForm onAdd={handleAdd} />);

    fireEvent.change(screen.getByPlaceholderText(/add new todo/i), { target: { value: 'New Test Todo' } });
    fireEvent.click(screen.getByRole('button', { name: /add todo/i }));

    expect(handleAdd).toHaveBeenCalledWith('New Test Todo');
    expect(screen.getByPlaceholderText(/add new todo/i)).toHaveValue(''); // Input should be cleared
  });

  test('does not call onAdd if input is empty', () => {
    const handleAdd = jest.fn();
    render(<AddTodoForm onAdd={handleAdd} />);

    fireEvent.click(screen.getByRole('button', { name: /add todo/i }));

    expect(handleAdd).not.toHaveBeenCalled();
  });
});
