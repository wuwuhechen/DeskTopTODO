export type TodoItem = {
    id : string;
    noteId : string;
    content : string;
    completed : boolean;
    priority : string;
    sortOrder : number;
    createdAt : string;
    updatedAt : string;
};

export type Priority = "normal"|"low" | "medium" | "high";