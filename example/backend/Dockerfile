FROM python:3.11-slim

WORKDIR /app

RUN pip install --upgrade pip

RUN pip install --no-cache-dir fastapi sqlmodel psycopg2-binary uvicorn

COPY . /app

EXPOSE 8000

ENV DATABASE_URL=postgresql://username:password@db:5432/mydatabase

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
