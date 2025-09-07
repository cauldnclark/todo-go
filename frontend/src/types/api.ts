// API Types based on backend models

export interface User {
  id: number;
  google_id: string;
  email: string;
  name: string;
  picture_url: string;
  created_at: string;
  updated_at: string;
}

export interface Todo {
  id: number;
  user_id: number;
  title: string;
  description: string;
  completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface MetaPagination {
  total: number;
  page: number;
  limit: number;
}

export interface TodosPaginated {
  todos: Todo[];
  meta: MetaPagination;
}

export interface CreateTodoRequest {
  title: string;
  description: string;
  completed?: boolean;
}

export interface UpdateTodoRequest {
  title?: string;
  description?: string;
  completed?: boolean;
}

export interface ApiError {
  error?: string;
  message: string;
  status?: number;
}

export interface AuthResponse {
  user: User;
  token: string;
}

// Frontend-specific types
export interface TodoFormData {
  title: string;
  description: string;
}

export interface TodoFilters {
  completed?: boolean;
  page?: number;
  limit?: number;
  search?: string;
}
