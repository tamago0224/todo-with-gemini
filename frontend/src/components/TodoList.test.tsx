import { render, screen, fireEvent } from '@testing-library/react';
import TodoList from './TodoList';

describe('TodoList', () => {
  const tasks = [
    { id: 1, title: 'Task 1', completed: false },
    { id: 2, title: 'Task 2', completed: true },
  ];

  test('renders tasks correctly', () => {
    render(
      <TodoList
        tasks={tasks}
        onToggleComplete={() => {}}
        onDelete={() => {}}
      />
    );

    expect(screen.getByText('Task 1')).toBeInTheDocument();
    expect(screen.getByText('Task 2')).toBeInTheDocument();
    expect(screen.getByLabelText('Task 1')).not.toBeChecked();
    expect(screen.getByLabelText('Task 2')).toBeChecked();
  });

  test('calls onToggleComplete when checkbox is clicked', () => {
    const handleToggleComplete = jest.fn();
    render(
      <TodoList
        tasks={tasks}
        onToggleComplete={handleToggleComplete}
        onDelete={() => {}}
      />
    );

    fireEvent.click(screen.getByLabelText('Task 1'));
    expect(handleToggleComplete).toHaveBeenCalledWith(1, true);
  });

  test('calls onDelete when delete button is clicked', () => {
    const handleDelete = jest.fn();
    render(
      <TodoList
        tasks={tasks}
        onToggleComplete={() => {}}
        onDelete={handleDelete}
      />
    );

    fireEvent.click(screen.getAllByText('Delete')[0]); // Click delete button for Task 1
    expect(handleDelete).toHaveBeenCalledWith(1);
  });
});
