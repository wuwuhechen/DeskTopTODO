import { useState, useEffect } from "react";
import { TodoItem as TodoRow } from "./TodoItem";
import type { Priority, TodoItem } from "./types";
import { TodoService } from "../bindings/todo";
import "./App.css";

export default function App() {
  const [todos, setTodos] = useState<TodoItem[]>([]);
  const [content, setContent] = useState<string>("");

  const [editingID, setEditingID] = useState<string | null>(null);
  const [draft, setDraft] = useState<string>("");

  const [showCompleted, setShowCompleted] = useState("false");

  const visibleTodos =
    showCompleted === "true" ? todos : todos.filter((todo) => !todo.completed);

  useEffect(() => {
    TodoService.List()
      .then((items) => setTodos(items ?? []))
      .catch(console.error);

    fetchShowCompleted().catch(console.error);
  }, []);

  async function addTodo() {
    const text = content.trim();
    if (!text) return;

    try {
      const todo = await TodoService.Create(text);
      if (todo === null) {
        throw new Error("创建待办未返回数据");
      }
      setTodos((current) => [...current, todo]);
      setContent("");
    } catch (error) {
      console.error(error);
    }
  }

  async function toggleTodo(todo: TodoItem) {
    try {
      await TodoService.Toggle(todo.id, !todo.completed);

      setTodos((current) =>
        current.map((item) =>
          item.id === todo.id ? { ...item, completed: !item.completed } : item,
        ),
      );
    } catch (error) {
      console.error(error);
    }
  }

  async function deleteTodo(id: string) {
    try {
      await TodoService.Delete(id);
      setTodos((current) => current.filter((todo) => todo.id !== id));
    } catch (error) {
      console.error(error);
    }
  }

  function startEditing(todo: TodoItem) {
    setEditingID(todo.id);
    setDraft(todo.content);
  }

  function cancelEditing() {
    setEditingID(null);
    setDraft("");
  }

  async function saveEditing(todo: TodoItem) {
    const text = draft.trim();
    if (!text) return;

    await TodoService.UpdateContent(todo.id, text);

    setTodos((current) =>
      current.map((item) =>
        item.id === todo.id ? { ...item, content: text } : item,
      ),
    );

    cancelEditing();
  }

  async function fetchShowCompleted() {
    try {
      const value = await TodoService.GetShowCompleted();
      setShowCompleted(value);
    } catch (error) {
      console.error(error);
    }
  }

  async function toggleShowCompleted() {
    const newValue = showCompleted === "true" ? "false" : "true";
    setShowCompleted(newValue);
    await TodoService.SetShowCompleted(newValue);
  }

  async function updatePriority(id: string, priority: Priority) {
    const previousTodos = todos;

    setTodos((current) =>
      current.map((todo) => (todo.id === id ? { ...todo, priority } : todo)),
    );
    
    try {
      await TodoService.UpdatePriority(id, priority);
    } catch (error) {
      setTodos(previousTodos);
      alert("更新优先级失败，请重试");
    }

  }

  return (
    <main className="sticky-note">
      <header className="note-header">
        <div>
          <h1>今天要做什么</h1>
          <span>
            已完成 {todos.filter((todo) => todo.completed).length}/{todos.length}
          </span>
        </div>

        <button
          className="toggle-completed-button"
          type="button"
          onClick={toggleShowCompleted}
        >
          {showCompleted === "true" ? "隐藏已完成" : "显示已完成"}
        </button>
      </header>

      <section className="todo-list">
        {visibleTodos.map((todo) => (
          <TodoRow
            key={todo.id}
            todo={todo}
            isEditing={editingID === todo.id}
            draft={draft}
            onToggle={toggleTodo}
            onStartEditing={startEditing}
            onDraftChange={setDraft}
            onSave={saveEditing}
            onCancel={cancelEditing}
            onDelete={deleteTodo}
            onPriorityChange={updatePriority}
          />
        ))}
      </section>

      <footer>
        <input
          value={content}
          placeholder="输入待办事项"
          onChange={(e) => setContent(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              addTodo();
              setContent("");
            }
          }}
        />
      </footer>
    </main>
  );
}
