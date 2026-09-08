import type { Dispatch, ReactNode, SetStateAction } from "react";
import {
  useSensors,
  useSensor,
  PointerSensor,
  DndContext,
  closestCenter,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import type { Todo } from "../bindings/todo/internal/model";
import { TodoService } from "../bindings/todo";

type Props = {
  noteID: string;
  todos: Todo[];
  setTodos: Dispatch<SetStateAction<Todo[]>>;
  sortingEnabled: boolean;
  renderTodo: (todo: Todo) => ReactNode;
};

type SortableTodoProps = {
  todo: Todo;
  children: ReactNode;
  disabled: boolean;
};

function SortableTodo({ todo, children, disabled }: SortableTodoProps) {
  const {
    attributes,
    listeners,
    setActivatorNodeRef,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: todo.id, disabled });

  return (
    <div
      ref={setNodeRef}
      className="sortable-todo-row"
      style={{
        position: "relative",
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.45 : 1,
        zIndex: isDragging ? 1 : "auto",
      }}
    >
      {children}
      <button
        ref={setActivatorNodeRef}
        type="button"
        className="drag-handle"
        aria-label={`拖动排序：${todo.content}`}
        title="拖动排序"
        disabled={disabled}
        style={{
          position: "absolute",
          top: "50%",
          left: 0,
          transform: "translateY(-50%)",
          zIndex: 2,
        }}
        {...attributes}
        {...listeners}
      >
        ⠿
      </button>
    </div>
  );
}

export function TodoList({
  noteID,
  todos,
  setTodos,
  sortingEnabled,
  renderTodo,
}: Props) {
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 6,
      },
    }),
  );

  async function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;

    if (!over || active.id === over.id) {
      return;
    }

    const oldIndex = todos.findIndex((todo) => todo.id === active.id);
    const newIndex = todos.findIndex((todo) => todo.id === over.id);

    if (oldIndex < 0 || newIndex < 0) {
      return;
    }

    const previousTodos = todos;
    const reorderedTodos = arrayMove(todos, oldIndex, newIndex);

    // 乐观更新：先立即更新界面，再持久化。
    setTodos(reorderedTodos);

    try {
      await TodoService.UpdateSortOrder(noteID, reorderedTodos.map((todo) => todo.id));
    } catch {
      // 保存失败则恢复原始顺序。
      setTodos(previousTodos);
      window.alert("排序保存失败，请重试。");
    }
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <SortableContext
        items={todos.map((todo) => todo.id)}
        strategy={verticalListSortingStrategy}
        disabled={!sortingEnabled}
      >
        <section className="todo-list">
          {todos.map((todo) => (
            <SortableTodo
              key={todo.id}
              todo={todo}
              disabled={!sortingEnabled}
            >
              {renderTodo(todo)}
            </SortableTodo>
          ))}
        </section>
      </SortableContext>
    </DndContext>
  );
}
