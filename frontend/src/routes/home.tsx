import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import {
  Trash2,
  Plus,
  Calendar,
  CheckCircle2,
  Circle,
  Search,
  Loader2,
  Filter,
  BarChart3,
  LogOut,
} from "lucide-react";
import { toast } from "sonner";
import { todoApi } from "@/services/todoApi";
import type { Todo, ApiError } from "@/types/api";
import useWebsocket from "@/hooks/useWebsocket";

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

  const onNewWsMessage = (event: MessageEvent<string>) => {
    try {
      const data = JSON.parse(event.data);
      console.log("WebSocket message received:", data);
      toast.success("New todo created");
      // Refresh todos when we receive updates
      fetchTodos();
    } catch (error) {
      console.error("Failed to parse WebSocket message:", error);
    }
  };

  const handleWebSocketOpen = () => {
    console.log("WebSocket connection opened");
  };

  const handleWebSocketClose = (event: CloseEvent) => {
    console.log("WebSocket connection closed:", event.code, event.reason);
  };

  useWebsocket(
    "ws://localhost:8085/ws",
    todoAuthToken as string,
    onNewWsMessage,
    handleWebSocketOpen,
    handleWebSocketClose
  );

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
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-slate-50">
      <div className="max-w-5xl mx-auto px-6 py-12 lg:px-8">
        {/* Header Section */}
        <div className="mb-16">
          <div className="flex items-center justify-between mb-8">
            <div>
              <h1 className="text-4xl font-light text-slate-900 mb-3 tracking-tight">
                Tasks
              </h1>
              <p className="text-slate-500 font-light">
                {completedCount} of {totalCount} completed
              </p>
            </div>
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2 text-slate-600">
                <BarChart3 className="h-5 w-5" />
                <span className="text-sm font-medium">
                  {totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0}%
                </span>
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  localStorage.removeItem("todo-auth-token");
                  window.location.href = "/login";
                }}
                className="p-2 h-auto w-auto text-slate-500 hover:text-slate-700 hover:bg-slate-100 transition-all duration-200"
              >
                <LogOut className="h-5 w-5" />
              </Button>
            </div>
          </div>
          
          {/* Progress Bar */}
          {totalCount > 0 && (
            <div className="w-full bg-slate-100 rounded-full h-1.5 mb-8">
              <div
                className="bg-gradient-to-r from-slate-600 to-slate-700 h-1.5 rounded-full transition-all duration-700 ease-out"
                style={{
                  width: `${(completedCount / totalCount) * 100}%`,
                }}
              />
            </div>
          )}
        </div>

        {/* Error Display */}
        {error && (
          <div className="mb-8 p-4 bg-red-50 border border-red-100 rounded-xl">
            <p className="text-red-700 text-sm font-medium">{error}</p>
          </div>
        )}

        {/* Add Task Section */}
        <div className="mb-12">
          <Card className="border-0 shadow-sm bg-white/70 backdrop-blur-sm">
            <CardContent className="p-8">
              <div className="space-y-6">
                <div className="flex items-center space-x-3 mb-6">
                  <Plus className="h-5 w-5 text-slate-600" />
                  <h2 className="text-lg font-medium text-slate-900">New Task</h2>
                </div>
                <div className="space-y-4">
                  <Input
                    placeholder="What needs to be done?"
                    value={newTodo}
                    onChange={(e) => setNewTodo(e.target.value)}
                    onKeyPress={(e) =>
                      e.key === "Enter" && !e.shiftKey && addTodo()
                    }
                    disabled={loading}
                    className="border-slate-200 focus:border-slate-400 focus:ring-slate-400 text-base py-3 px-4 rounded-lg bg-white/50"
                  />
                  <Input
                    placeholder="Add details (optional)"
                    value={newTodoDescription}
                    onChange={(e) => setNewTodoDescription(e.target.value)}
                    disabled={loading}
                    className="border-slate-200 focus:border-slate-400 focus:ring-slate-400 text-sm py-3 px-4 rounded-lg bg-white/50"
                  />
                  <Button
                    onClick={addTodo}
                    disabled={loading || !newTodo.trim()}
                    className="bg-slate-900 hover:bg-slate-800 text-white font-medium py-3 px-6 rounded-lg transition-all duration-200 shadow-sm hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {loading ? (
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    ) : (
                      <Plus className="h-4 w-4 mr-2" />
                    )}
                    Add Task
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Search and Filter Section */}
        <div className="mb-8 space-y-6">
          {/* Search */}
          <Card className="border-0 shadow-sm bg-white/70 backdrop-blur-sm">
            <CardContent className="p-6">
              <div className="relative">
                <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-slate-400 h-4 w-4" />
                <Input
                  placeholder="Search tasks..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-12 border-slate-200 focus:border-slate-400 focus:ring-slate-400 py-3 rounded-lg bg-white/50"
                />
              </div>
            </CardContent>
          </Card>

          {/* Filter Buttons */}
          <div className="flex items-center space-x-2">
            <Filter className="h-4 w-4 text-slate-500 mr-2" />
            <div className="flex space-x-1">
              <Button
                variant={filter === "all" ? "default" : "ghost"}
                onClick={() => setFilter("all")}
                className={`px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                  filter === "all"
                    ? "bg-slate-900 text-white shadow-sm"
                    : "text-slate-600 hover:text-slate-900 hover:bg-slate-100"
                }`}
                disabled={loading}
              >
                All
                <span className="ml-2 text-xs opacity-75">({pagination.total})</span>
              </Button>
              <Button
                variant={filter === "pending" ? "default" : "ghost"}
                onClick={() => setFilter("pending")}
                className={`px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                  filter === "pending"
                    ? "bg-slate-900 text-white shadow-sm"
                    : "text-slate-600 hover:text-slate-900 hover:bg-slate-100"
                }`}
                disabled={loading}
              >
                Active
                <span className="ml-2 text-xs opacity-75">({pagination.total - completedCount})</span>
              </Button>
              <Button
                variant={filter === "completed" ? "default" : "ghost"}
                onClick={() => setFilter("completed")}
                className={`px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                  filter === "completed"
                    ? "bg-slate-900 text-white shadow-sm"
                    : "text-slate-600 hover:text-slate-900 hover:bg-slate-100"
                }`}
                disabled={loading}
              >
                Done
                <span className="ml-2 text-xs opacity-75">({completedCount})</span>
              </Button>
            </div>
          </div>
        </div>

        {/* Tasks List */}
        <div className="space-y-3">
          {loading && todos.length === 0 ? (
            <Card className="border-0 shadow-sm bg-white/70 backdrop-blur-sm">
              <CardContent className="flex items-center justify-center py-16">
                <div className="text-center">
                  <Loader2 className="h-8 w-8 text-slate-400 mx-auto mb-4 animate-spin" />
                  <p className="text-slate-500 font-light">Loading tasks...</p>
                </div>
              </CardContent>
            </Card>
          ) : filteredTodos.length === 0 ? (
            <Card className="border-0 shadow-sm bg-white/70 backdrop-blur-sm">
              <CardContent className="flex items-center justify-center py-16">
                <div className="text-center">
                  <CheckCircle2 className="h-12 w-12 text-slate-300 mx-auto mb-4" />
                  <p className="text-slate-500 font-light">
                    {searchQuery.trim()
                      ? "No tasks match your search"
                      : filter === "completed"
                      ? "No completed tasks yet"
                      : filter === "pending"
                      ? "All tasks completed"
                      : "No tasks yet"}
                  </p>
                </div>
              </CardContent>
            </Card>
          ) : (
            filteredTodos.map((todo) => (
              <Card
                key={todo.id}
                className={`group border-0 shadow-sm transition-all duration-300 hover:shadow-md hover:-translate-y-0.5 ${
                  todo.completed 
                    ? "bg-white/40 backdrop-blur-sm" 
                    : "bg-white/70 backdrop-blur-sm hover:bg-white/90"
                }`}
              >
                <CardContent className="p-6">
                  <div className="flex items-start space-x-4">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => toggleTodo(todo.id)}
                      className="p-2 h-auto w-auto mt-0.5 hover:bg-transparent group-hover:scale-110 transition-transform duration-200"
                      disabled={loading}
                    >
                      {todo.completed ? (
                        <CheckCircle2 className="h-5 w-5 text-slate-600" />
                      ) : (
                        <Circle className="h-5 w-5 text-slate-300 hover:text-slate-500 transition-colors" />
                      )}
                    </Button>

                    <div className="flex-1 min-w-0">
                      <h3
                        className={`font-medium text-base leading-relaxed transition-all duration-200 ${
                          todo.completed
                            ? "line-through text-slate-400"
                            : "text-slate-900 group-hover:text-slate-700"
                        }`}
                      >
                        {todo.title}
                      </h3>
                      {todo.description && (
                        <p
                          className={`text-sm mt-2 leading-relaxed ${
                            todo.completed
                              ? "line-through text-slate-300"
                              : "text-slate-600"
                          }`}
                        >
                          {todo.description}
                        </p>
                      )}
                      <div className="flex items-center space-x-4 mt-4">
                        <div className="flex items-center space-x-1.5 text-xs text-slate-400">
                          <Calendar className="h-3.5 w-3.5" />
                          <span className="font-light">{formatDate(todo.created_at)}</span>
                        </div>
                        {todo.updated_at !== todo.created_at && (
                          <div className="flex items-center space-x-1.5 text-xs text-slate-400">
                            <span className="w-1 h-1 bg-slate-300 rounded-full" />
                            <span className="font-light">Updated {formatDate(todo.updated_at)}</span>
                          </div>
                        )}
                      </div>
                    </div>

                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => deleteTodo(todo.id)}
                      className="p-2 h-auto w-auto text-slate-400 hover:text-red-500 hover:bg-red-50 opacity-0 group-hover:opacity-100 transition-all duration-200"
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
      </div>
    </div>
  );
}
