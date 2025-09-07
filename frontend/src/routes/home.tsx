import React, { useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Trash2, Plus, Calendar, CheckCircle2, Circle } from "lucide-react";

export default function TodoListPage() {
  const todoAuthToken = localStorage.getItem("todo-auth-token");

  useEffect(() => {
    if (!todoAuthToken) {
      window.location.href = "/login";
    }
  }, [todoAuthToken]);

  const [todos, setTodos] = React.useState([
    {
      id: 1,
      text: "Complete project proposal",
      completed: false,
      priority: "high",
      dueDate: "2025-09-08",
    },
    {
      id: 2,
      text: "Review team feedback",
      completed: true,
      priority: "medium",
      dueDate: "2025-09-07",
    },
    {
      id: 3,
      text: "Schedule client meeting",
      completed: false,
      priority: "low",
      dueDate: "2025-09-10",
    },
    {
      id: 4,
      text: "Update documentation",
      completed: false,
      priority: "medium",
      dueDate: "2025-09-09",
    },
  ]);
  const [newTodo, setNewTodo] = React.useState("");
  const [filter, setFilter] = React.useState("all");

  const addTodo = () => {
    if (newTodo.trim()) {
      setTodos([
        ...todos,
        {
          id: Date.now(),
          text: newTodo,
          completed: false,
          priority: "medium",
          dueDate: new Date().toISOString().split("T")[0],
        },
      ]);
      setNewTodo("");
    }
  };

  const toggleTodo = (id: number) => {
    setTodos(
      todos.map((todo) =>
        todo.id === id ? { ...todo, completed: !todo.completed } : todo
      )
    );
  };

  const deleteTodo = (id: number) => {
    setTodos(todos.filter((todo) => todo.id !== id));
  };

  const filteredTodos = todos.filter((todo) => {
    if (filter === "completed") return todo.completed;
    if (filter === "pending") return !todo.completed;
    return true;
  });

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "high":
        return "bg-red-100 text-red-800 border-red-200";
      case "medium":
        return "bg-yellow-100 text-yellow-800 border-yellow-200";
      case "low":
        return "bg-green-100 text-green-800 border-green-200";
      default:
        return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const completedCount = todos.filter((todo) => todo.completed).length;
  const totalCount = todos.length;

  return (
    <div className="min-h-screen bg-gray-50 p-4">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">My Tasks</h1>
          <p className="text-gray-600">
            {completedCount} of {totalCount} tasks completed
          </p>
        </div>

        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="text-lg">Add New Task</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex gap-2">
              <Input
                placeholder="What needs to be done?"
                value={newTodo}
                onChange={(e) => setNewTodo(e.target.value)}
                onKeyPress={(e) => e.key === "Enter" && addTodo()}
                className="flex-1"
              />
              <Button
                onClick={addTodo}
                className="bg-blue-600 hover:bg-blue-700"
              >
                <Plus className="h-4 w-4 mr-2" />
                Add Task
              </Button>
            </div>
          </CardContent>
        </Card>

        <div className="flex gap-2 mb-6">
          <Button
            variant={filter === "all" ? "default" : "outline"}
            onClick={() => setFilter("all")}
            className={filter === "all" ? "bg-blue-600 hover:bg-blue-700" : ""}
          >
            All ({totalCount})
          </Button>
          <Button
            variant={filter === "pending" ? "default" : "outline"}
            onClick={() => setFilter("pending")}
            className={
              filter === "pending" ? "bg-blue-600 hover:bg-blue-700" : ""
            }
          >
            Pending ({totalCount - completedCount})
          </Button>
          <Button
            variant={filter === "completed" ? "default" : "outline"}
            onClick={() => setFilter("completed")}
            className={
              filter === "completed" ? "bg-blue-600 hover:bg-blue-700" : ""
            }
          >
            Completed ({completedCount})
          </Button>
        </div>

        <div className="space-y-3">
          {filteredTodos.length === 0 ? (
            <Card>
              <CardContent className="flex items-center justify-center py-8">
                <div className="text-center">
                  <CheckCircle2 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                  <p className="text-gray-500">
                    {filter === "completed"
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
                  <div className="flex items-center gap-3">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => toggleTodo(todo.id)}
                      className="p-1 h-6 w-6"
                    >
                      {todo.completed ? (
                        <CheckCircle2 className="h-5 w-5 text-green-600" />
                      ) : (
                        <Circle className="h-5 w-5 text-gray-400" />
                      )}
                    </Button>

                    <div className="flex-1">
                      <p
                        className={`${
                          todo.completed
                            ? "line-through text-gray-500"
                            : "text-gray-900"
                        }`}
                      >
                        {todo.text}
                      </p>
                      <div className="flex items-center gap-2 mt-2">
                        <Badge
                          className={`text-xs ${getPriorityColor(
                            todo.priority
                          )}`}
                        >
                          {todo.priority}
                        </Badge>
                        <div className="flex items-center gap-1 text-xs text-gray-500">
                          <Calendar className="h-3 w-3" />
                          {todo.dueDate}
                        </div>
                      </div>
                    </div>

                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => deleteTodo(todo.id)}
                      className="text-red-600 hover:text-red-700 hover:bg-red-50"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>

        {totalCount > 0 && (
          <div className="mt-8 bg-white rounded-lg p-4 border">
            <div className="flex items-center justify-between">
              <div className="text-sm text-gray-600">
                Progress: {completedCount}/{totalCount} tasks
              </div>
              <div className="text-sm font-medium text-blue-600">
                {Math.round((completedCount / totalCount) * 100)}% Complete
              </div>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2 mt-2">
              <div
                className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                style={{ width: `${(completedCount / totalCount) * 100}%` }}
              ></div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
