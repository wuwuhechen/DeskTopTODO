import { useState } from "react";
import type { TodoItem } from "./types";
import "./App.css";

export default function App(){
  const [todos, setTodos] = useState<TodoItem[]>([]);
  const [content, setContent] = useState<string>("");

  function addTodo(){
    const text = content.trim();
    if (!text) return;

    setTodos((current) => [
      ...current,
      {
        id: crypto.randomUUID(),
        content: text,
        completed: false,
      }
    ])
  }

  function toggleTodo(id: string){
    setTodos((current) => 
     current.map((todo)=>
      todo.id === id ? {...todo, completed: !todo.completed} : todo),
    );
  }

  function deleteTodo(id: string){
    setTodos((current) => current.filter((todo) => todo.id !== id));
  }

  return(
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
              onChange={() => toggleTodo(todo.id)}
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
  )
}