from typing import Optional, List
import os

from fastapi import FastAPI
from sqlmodel import Field, Session, SQLModel, create_engine, select


class Todo(SQLModel, table=True):
    id: Optional[int] = Field(default=None, primary_key=True)
    title: str
    completed: bool = Field(default=False)
    description: Optional[str]


# Replace with your PostgreSQL connection details
DATABASE_URL = os.getenv("DATABASE_URL", "sqlite:///./test.db")
engine = create_engine(DATABASE_URL, echo=True)


def create_db_and_tables():
    SQLModel.metadata.create_all(engine)


app = FastAPI()


@app.on_event("startup")
def on_startup() -> None:
    create_db_and_tables()


@app.post("/todos/", response_model=Todo)
def create_todo(todo: Todo) -> Todo:
    with Session(engine) as session:
        session.add(todo)
        session.commit()
        session.refresh(todo)
        return todo


@app.get("/todos/", response_model=List[Todo])
def read_todos() -> List[Todo]:
    with Session(engine) as session:
        todos = session.exec(select(Todo)).all()
        return todos


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}
