import React, { useState } from 'react';

interface AddTodoFormProps {
  onAdd: (title: string) => void;
}

const AddTodoForm: React.FC<AddTodoFormProps> = ({ onAdd }) => {
  const [title, setTitle] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (title.trim()) {
      onAdd(title);
      setTitle('');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Add new todo"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />
      <button type="submit">Add Todo</button>
    </form>
  );
};

export default AddTodoForm;
