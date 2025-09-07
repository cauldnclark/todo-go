import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Trash2,
  Plus,
  Calendar,
  CheckCircle2,
  Circle,
  Search,
  Loader2,
} from "lucide-react";
import { todoApi } from "@/services/todoApi";
import type { Todo, ApiError } from "@/types/api";

export default function TodoListPage() {
  const todoAuthToken = localStorage.getItem("todo-auth-token");

  useEffect(() => {
    if (!todoAuthToken) {
      window.location.href = "/login";
    }
  }, [todoAuthToken]);

  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTodo, setNewTodo] = useState("");
  const [newTodoDescription, setNewTodoDescription] = useState("");
  const [filter, setFilter] = useState<"all" | "completed" | "pending">("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState({
    total: 0,
    page: 1,
    limit: 50,
  });

  const fetchTodos = async () => {
    try {
      setLoading(true);
      setError(null);
      const completedFilter =
        filter === "all" ? undefined : filter === "completed";
      const response = await todoApi.getTodos({
        completed: completedFilter,
        page: pagination.page,
        limit: pagination.limit,
      });
      setTodos(response.todos || []);
      setPagination(response.meta);
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.message || "Failed to fetch todos");
    } finally {
      setLoading(false);
    }
  };

  const addTodo = async () => {
    if (!newTodo.trim()) return;

    try {
      setLoading(true);
      setError(null);
      await todoApi.createTodo({
        title: newTodo.trim(),
        description: newTodoDescription.trim() || "",
        completed: false,
      });
      setNewTodo("");
      setNewTodoDescription("");
      await fetchTodos();
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.message || "Failed to create todo");
    } finally {
      setLoading(false);
    }
  };

  const toggleTodo = async (id: number) => {
    const todo = (todos || []).find((t) => t.id === id);
    if (!todo) return;

    try {
      setError(null);
      await todoApi.toggleTodoCompletion(id, !todo.completed);
      await fetchTodos();
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.message || "Failed to update todo");
    }
  };

  const deleteTodo = async (id: number) => {
    try {
      setError(null);
      await todoApi.deleteTodo(id);
      await fetchTodos();
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.message || "Failed to delete todo");
    }
  };

  // Load todos on component mount and when filter changes
  useEffect(() => {
    if (todoAuthToken) {
      fetchTodos();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [todoAuthToken, filter]);

  const filteredTodos = (todos || []).filter((todo) => {
    // Filter by completion status
    let matchesFilter = true;
    if (filter === "completed") matchesFilter = todo.completed;
    if (filter === "pending") matchesFilter = !todo.completed;

    // Filter by search query
    let matchesSearch = true;
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      matchesSearch =
        todo.title.toLowerCase().includes(query) ||
        todo.description.toLowerCase().includes(query);
    }

    return matchesFilter && matchesSearch;
  });

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  const completedCount = (todos || []).filter((todo) => todo.completed).length;
  const totalCount = (todos || []).length;

  return (
    <div className="min-h-screen bg-gray-50 p-4">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">My Tasks</h1>
          <p className="text-gray-600">
            {completedCount} of {totalCount} tasks completed
          </p>
        </div>

        {error && (
          <Card className="mb-6 border-red-200 bg-red-50">
            <CardContent className="p-4">
              <p className="text-red-800 text-sm">{error}</p>
            </CardContent>
          </Card>
        )}

        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="text-lg">Add New Task</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <Input
                placeholder="Task title"
                value={newTodo}
                onChange={(e) => setNewTodo(e.target.value)}
                onKeyPress={(e) =>
                  e.key === "Enter" && !e.shiftKey && addTodo()
                }
                disabled={loading}
              />
              <Input
                placeholder="Description (optional)"
                value={newTodoDescription}
                onChange={(e) => setNewTodoDescription(e.target.value)}
                disabled={loading}
              />
              <Button
                onClick={addTodo}
                disabled={loading || !newTodo.trim()}
                className="bg-blue-600 hover:bg-blue-700 w-full"
              >
                {loading ? (
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                ) : (
                  <Plus className="h-4 w-4 mr-2" />
                )}
                Add Task
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card className="mb-6">
          <CardContent className="p-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
              <Input
                placeholder="Search todos..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>
          </CardContent>
        </Card>

        <div className="flex gap-2 mb-6">
          <Button
            variant={filter === "all" ? "default" : "outline"}
            onClick={() => setFilter("all")}
            className={filter === "all" ? "bg-blue-600 hover:bg-blue-700" : ""}
            disabled={loading}
          >
            All ({pagination.total})
          </Button>
          <Button
            variant={filter === "pending" ? "default" : "outline"}
            onClick={() => setFilter("pending")}
            className={
              filter === "pending" ? "bg-blue-600 hover:bg-blue-700" : ""
            }
            disabled={loading}
          >
            Pending ({pagination.total - completedCount})
          </Button>
          <Button
            variant={filter === "completed" ? "default" : "outline"}
            onClick={() => setFilter("completed")}
            className={
              filter === "completed" ? "bg-blue-600 hover:bg-blue-700" : ""
            }
            disabled={loading}
          >
            Completed ({completedCount})
          </Button>
        </div>

        <div className="space-y-3">
          {loading && todos.length === 0 ? (
            <Card>
              <CardContent className="flex items-center justify-center py-8">
                <div className="text-center">
                  <Loader2 className="h-8 w-8 text-gray-400 mx-auto mb-4 animate-spin" />
                  <p className="text-gray-500">Loading todos...</p>
                </div>
              </CardContent>
            </Card>
          ) : filteredTodos.length === 0 ? (
            <Card>
              <CardContent className="flex items-center justify-center py-8">
                <div className="text-center">
                  <CheckCircle2 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                  <p className="text-gray-500">
                    {searchQuery.trim()
                      ? "No todos match your search"
                      : filter === "completed"
                      ? "No completed tasks yet"
                      : filter === "pending"
                      ? "All tasks completed! Great job!"
                      : "No tasks yet. Add one above!"}
                  </p>
                </div>
              </CardContent>
            </Card>
          ) : (
            filteredTodos.map((todo) => (
              <Card
                key={todo.id}
                className={`transition-all hover:shadow-md ${
                  todo.completed ? "bg-gray-50" : "bg-white"
                }`}
              >
                <CardContent className="p-4">
                  <div className="flex items-start gap-3">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => toggleTodo(todo.id)}
                      className="p-1 h-6 w-6 mt-1"
                      disabled={loading}
                    >
                      {todo.completed ? (
                        <CheckCircle2 className="h-5 w-5 text-green-600" />
                      ) : (
                        <Circle className="h-5 w-5 text-gray-400" />
                      )}
                    </Button>

                    <div className="flex-1">
                      <h3
                        className={`font-medium ${
                          todo.completed
                            ? "line-through text-gray-500"
                            : "text-gray-900"
                        }`}
                      >
                        {todo.title}
                      </h3>
                      {todo.description && (
                        <p
                          className={`text-sm mt-1 ${
                            todo.completed
                              ? "line-through text-gray-400"
                              : "text-gray-600"
                          }`}
                        >
                          {todo.description}
                        </p>
                      )}
                      <div className="flex items-center gap-2 mt-2">
                        <div className="flex items-center gap-1 text-xs text-gray-500">
                          <Calendar className="h-3 w-3" />
                          Created: {formatDate(todo.created_at)}
                        </div>
                        {todo.updated_at !== todo.created_at && (
                          <div className="flex items-center gap-1 text-xs text-gray-500">
                            Updated: {formatDate(todo.updated_at)}
                          </div>
                        )}
                      </div>
                    </div>

                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => deleteTodo(todo.id)}
                      className="text-red-600 hover:text-red-700 hover:bg-red-50"
                      disabled={loading}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>

        {pagination.total > 0 && (
          <div className="mt-8 bg-white rounded-lg p-4 border">
            <div className="flex items-center justify-between">
              <div className="text-sm text-gray-600">
                Progress: {completedCount}/{pagination.total} tasks
              </div>
              <div className="text-sm font-medium text-blue-600">
                {Math.round((completedCount / pagination.total) * 100)}%
                Complete
              </div>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2 mt-2">
              <div
                className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                style={{
                  width: `${(completedCount / pagination.total) * 100}%`,
                }}
              ></div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
