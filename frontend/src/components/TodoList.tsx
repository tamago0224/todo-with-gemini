import React from 'react';

interface Task {
  id: number;
  title: string;
  completed: boolean;
}

interface TodoListProps {
  tasks: Task[];
  onToggleComplete: (id: number, completed: boolean) => void;
  onDelete: (id: number) => void;
}

const TodoList: React.FC<TodoListProps> = ({ tasks, onToggleComplete, onDelete }) => {
  return (
    <ul>
      {tasks.map((task) => (
        <li key={task.id} style={{ textDecoration: task.completed ? 'line-through' : 'none' }}>
          <input
            type="checkbox"
            checked={task.completed}
            onChange={() => onToggleComplete(task.id, !task.completed)}
          />
          {task.title}
          <button onClick={() => onDelete(task.id)} style={{ marginLeft: '10px' }}>Delete</button>
        </li>
      ))}
    </ul>
  );
};

export default TodoList;
