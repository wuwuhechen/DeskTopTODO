import type { Priority, TodoItem } from "./types";

type Props = {
  todo: TodoItem;
  isEditing: boolean;
  draft: string;
  onToggle: (todo: TodoItem) => Promise<void>;
  onStartEditing: (todo: TodoItem) => void;
  onDraftChange: (value: string) => void;
  onSave: (todo: TodoItem) => Promise<void>;
  onCancel: () => void;
  onDelete: (id: string) => Promise<void>;
  onPriorityChange: (id: string, priority: Priority) => Promise<void>;
};

const priorityLabels: Record<Priority, string> = {
  normal: "无优先级",
  low: "低优先级",
  medium: "中优先级",
  high: "高优先级",
};

export function TodoItem({
  todo,
  isEditing,
  draft,
  onToggle,
  onStartEditing,
  onDraftChange,
  onSave,
  onCancel,
  onDelete,
  onPriorityChange,
}: Props) {
  return (
    <article className={`todo-item priority-${todo.priority}`}>
      <span
        className="priority-dot"
        title={`优先级：${priorityLabels[todo.priority]}`}
      />

      <input
        type="checkbox"
        checked={todo.completed}
        onChange={() => void onToggle(todo)}
      />

      {isEditing ? (
        <input
          className="todo-edit-input"
          value={draft}
          autoFocus
          onChange={(event) => onDraftChange(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              void onSave(todo);
            } else if (event.key === "Escape") {
              onCancel();
            }
          }}
          onBlur={() => void onSave(todo)}
        />
      ) : (
        <span
          className={`todo-content ${todo.completed ? "completed" : ""}`}
          onDoubleClick={() => onStartEditing(todo)}
          title="双击编辑"
        >
          {todo.content}
        </span>
      )}

      <select
        className="priority-select"
        aria-label="选择待办优先级"
        value={todo.priority}
        onChange={(event) =>
          void onPriorityChange(todo.id, event.target.value as Priority)
        }
      >
        <option value="normal">无优先级</option>
        <option value="low">低优先级</option>
        <option value="medium">中优先级</option>
        <option value="high">高优先级</option>
      </select>

      <button
        className="delete-button"
        type="button"
        onClick={() => void onDelete(todo.id)}
      >
        删除
      </button>
    </article>
  );
}
