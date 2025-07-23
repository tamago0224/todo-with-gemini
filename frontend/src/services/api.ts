const API_BASE_URL = 'http://localhost:8080'; // Replace with your backend API URL

interface AuthResponse {
  token: string;
}

interface Task {
  id?: number;
  title: string;
  completed: boolean;
}

const api = {
  signup: async (username: string, password: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/signup`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Signup failed');
    }
  },

  login: async (username: string, password: string): Promise<AuthResponse> => {
    const response = await fetch(`${API_BASE_URL}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Login failed');
    }
    return response.json();
  },

  getTasks: async (token: string): Promise<Task[]> => {
    const response = await fetch(`${API_BASE_URL}/api/tasks`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Failed to fetch tasks');
    }
    return response.json();
  },

  createTask: async (token: string, title: string): Promise<Task> => {
    const response = await fetch(`${API_BASE_URL}/api/tasks`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({ title, completed: false }),
    });
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Failed to create task');
    }
    return response.json();
  },

  updateTask: async (token: string, taskId: number, completed: boolean): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/api/tasks/${taskId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({ completed }),
    });
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Failed to update task');
    }
  },

  deleteTask: async (token: string, taskId: number): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/api/tasks/${taskId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Failed to delete task');
    }
  },
};

export default api;
