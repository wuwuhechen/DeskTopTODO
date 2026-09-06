import { useState, useEffect } from "react";
import type { TodoItem } from "./types";
import { TodoService } from "../bindings/todo";
import "./App.css";

export default function App() {
  const [todos, setTodos] = useState<TodoItem[]>([]);
  const [content, setContent] = useState<string>("");

  useEffect(() => {
    TodoService.List()
      .then((items) => setTodos(items ?? []))
      .catch(console.error);
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

  async function toggleTodo(todo:  TodoItem) {
    try {
      await TodoService.Toggle(todo.id, !todo.completed);

      setTodos((current) =>
        current.map((item) =>
          item.id === todo.id
            ? { ...item, completed: !item.completed }
            : item,
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

  return (
    <main className="sticky-note">
      <header className="note-header">
        <h1>今天要做什么</h1>
      </header>

      <section className="todo-list">
        {todos.map((todo) => (
          <div key={todo.id} className="todo-item">
            <input
              type="checkbox"
              checked={todo.completed}
              onChange={() => toggleTodo(todo)}
            />
            <span className={todo.completed ? "completed" : ""}>
              {todo.content}
            </span>
            <button onClick={() => deleteTodo(todo.id)}>删除</button>
          </div>
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
