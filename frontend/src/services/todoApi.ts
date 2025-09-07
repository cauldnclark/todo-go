import type {
  Todo,
  TodosPaginated,
  CreateTodoRequest,
  UpdateTodoRequest,
  TodoFilters,
  ApiError,
} from "@/types/api";

const API_URL = import.meta.env.VITE_API_URL;

class TodoApiService {
  private getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem("todo-auth-token");
    return {
      "Content-Type": "application/json",
      ...(token && { Authorization: `Bearer ${token}` }),
    };
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      const errorText = await response.text();
      const error: ApiError = {
        message: errorText || `HTTP error! status: ${response.status}`,
        status: response.status,
      };
      throw error;
    }

    // Handle empty responses (like DELETE)
    if (response.status === 204) {
      return {} as T;
    }

    return response.json();
  }

  async getTodos(filters: TodoFilters = {}): Promise<TodosPaginated> {
    const params = new URLSearchParams();

    if (filters.completed !== undefined) {
      params.append("completed", filters.completed.toString());
    }
    if (filters.page) {
      params.append("page", filters.page.toString());
    }
    if (filters.limit) {
      params.append("limit", filters.limit.toString());
    }

    const url = `${API_URL}/api/todos${
      params.toString() ? `?${params.toString()}` : ""
    }`;

    const response = await fetch(url, {
      method: "GET",
      headers: this.getAuthHeaders(),
    });

    return this.handleResponse<TodosPaginated>(response);
  }

  async createTodo(todoData: CreateTodoRequest): Promise<Todo> {
    const response = await fetch(`${API_URL}/api/todos`, {
      method: "POST",
      headers: this.getAuthHeaders(),
      body: JSON.stringify(todoData),
    });

    return this.handleResponse<Todo>(response);
  }

  async updateTodo(id: number, todoData: UpdateTodoRequest): Promise<Todo> {
    const response = await fetch(`${API_URL}/api/todos/${id}`, {
      method: "PUT",
      headers: this.getAuthHeaders(),
      body: JSON.stringify(todoData),
    });

    return this.handleResponse<Todo>(response);
  }

  async deleteTodo(id: number): Promise<void> {
    const response = await fetch(`${API_URL}/api/todos/${id}`, {
      method: "DELETE",
      headers: this.getAuthHeaders(),
    });

    await this.handleResponse<void>(response);
  }

  async toggleTodoCompletion(id: number, completed: boolean): Promise<Todo> {
    return this.updateTodo(id, { completed });
  }
}

export const todoApi = new TodoApiService();
export default todoApi;
