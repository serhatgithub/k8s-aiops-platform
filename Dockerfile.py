FROM python:3.9-slim
WORKDIR /app
RUN pip install fastapi uvicorn prometheus-api-client pandas requests
COPY ai_brain.py .
CMD ["python", "ai_brain.py"]