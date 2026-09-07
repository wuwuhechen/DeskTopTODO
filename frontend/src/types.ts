export type TodoItem = {
    id : string;
    content : string;
    completed : boolean;
    priority : string;
    sortOrder : number;
    createdAt : string;
    updatedAt : string;
};

export type Priority = "normal"|"low" | "medium" | "high";