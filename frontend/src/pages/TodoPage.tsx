import React, { useState, useEffect, useMemo } from 'react';
import TodoList from '../components/TodoList';
import AddTodoForm from '../components/AddTodoForm';
import FilterButtons from '../components/FilterButtons';
import api from '../services/api';
import { useAuth } from '../context/AuthContext';

interface Task {
  id: number;
  title: string;
  completed: boolean;
}

const TodoPage: React.FC = () => {
  const { token } = useAuth();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [filter, setFilter] = useState<'all' | 'active' | 'completed'>('all');

  useEffect(() => {
    const fetchTasks = async () => {
      if (token) {
        try {
          const fetchedTasks = await api.getTasks(token);
          setTasks(fetchedTasks || []);
        } catch (error: any) {
          alert(`Failed to fetch tasks: ${error.message}`);
        }
      }
    };
    fetchTasks();
  }, [token]);

  const handleAddTodo = async (title: string) => {
    if (token) {
      try {
        const newTask = await api.createTask(token, title);
        setTasks([...tasks, newTask]);
      } catch (error: any) {
        alert(`Failed to add task: ${error.message}`);
      }
    }
  };

  const handleToggleComplete = async (id: number, completed: boolean) => {
    if (token) {
      try {
        await api.updateTask(token, id, completed);
        setTasks(tasks.map((task) => (task.id === id ? { ...task, completed } : task)));
      } catch (error: any) {
        alert(`Failed to update task: ${error.message}`);
      }
    }
  };

  const handleDeleteTodo = async (id: number) => {
    if (token) {
      try {
        await api.deleteTask(token, id);
        setTasks(tasks.filter((task) => task.id !== id));
      } catch (error: any) {
        alert(`Failed to delete task: ${error.message}`);
      }
    }
  };

  const filteredTasks = useMemo(() => {
    switch (filter) {
      case 'active':
        return tasks.filter((task) => !task.completed);
      case 'completed':
        return tasks.filter((task) => task.completed);
      default:
        return tasks;
    }
  }, [tasks, filter]);

  return (
    <div>
      <h1>Your Todos</h1>
      <AddTodoForm onAdd={handleAddTodo} />
      <FilterButtons currentFilter={filter} onFilterChange={setFilter} />
      <TodoList
        tasks={filteredTasks}
        onToggleComplete={handleToggleComplete}
        onDelete={handleDeleteTodo}
      />
    </div>
  );
};

export default TodoPage;
