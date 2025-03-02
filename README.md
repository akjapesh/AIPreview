# AI-Driven UI Component Generator  

A powerful AI-driven UI component generator built with **Next.js (SSR + SSG)**, **Go + FastAPI**, and **Redis caching**, leveraging **OpenAI’s GPT-3** for intelligent code generation.  

## Features  

- **Efficient Backend Performance** – Built with **Go** and **FastAPI**, optimizing execution with **Gin**, **goroutines**, **connection pooling**, and **rate limiting**.  
- **Real-Time UX Enhancements** – Integrated **GraphQL event streaming** for seamless real-time updates.  
- **Enhanced AI Responses** – Implemented **RAG with LangChain**, improving retrieval accuracy by **20%** using **FAISS** for efficient vector search.  

## Installation & Running  

### Prerequisites  
- Node.js (for Next.js frontend)  
- Go (for backend)  
- FastAPI (Python)  
- Redis  

### Install Dependencies  
```sh
# Install frontend dependencies  
cd frontend  
yarn install  
# or  
npm install  

# Install backend dependencies  
cd backend  
go mod tidy  # For Go backend service

# Install LLM dependecyies
cd ml-services
pip install -r requirements.txt  # For python ml-service
```

### Run the Application
```
# Start the frontend  
cd frontend  
yarn dev  
# or  
npm run dev  

# Start the Go backend  
cd backend  
go run main.go  

# Start the FastAPI server  
cd ml-services  
python3 main.py 

```

### Tech Stack
Frontend – Next.js 
Backend – Go (Gin) & FastAPI
AI & Data Processing – OpenAI’s GPT-3, LangChain (RAG), FAISS
Database & Caching – Redis, Supabase
