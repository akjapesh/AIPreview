from langchain_community.vectorstores import Chroma
from sentence_transformers import SentenceTransformer
from langchain_core.documents import Document 
import json
import os

class SentenceTransformerEmbedder:
    def __init__(self, model_name='all-MiniLM-L6-v2'):
        self.model = SentenceTransformer(model_name)

    def embed_query(self, text):
        # Convert text to embedding
        return self.model.encode(text).tolist()
class ComponentRetriever:
    def __init__(self):
        self.embedding = SentenceTransformerEmbedder()
        self.db = self._initialize_db()
        
    def _initialize_db(self):
        if not os.path.exists("vector_db"):
            os.makedirs("vector_db")
            
        if os.listdir("vector_db"):
            return Chroma(
                persist_directory="vector_db",
                embedding_function=self.embedding
            )
            
        return self._create_vector_db()
    
    def _create_vector_db(self):
        with open("data/react_components.json") as f:
            components = json.load(f)
        
        docs = [
            Document(
                page_content=f"{comp['description']}\n\n{comp['code']}",
                metadata={
                    "name": comp["name"],
                    "tags": comp["tags"]
                }
            )
            for comp in components
        ]
        
        return Chroma.from_documents(
            documents=docs,
            embedding=self.embedding,
            persist_directory="vector_db"
        )
        
    def retrieve_components(self, query, k=3):
        return self.db.similarity_search(query, k=k)