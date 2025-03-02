from fastapi import FastAPI, HTTPException
from fastapi.responses import JSONResponse,StreamingResponse
from services.embedding_services import ComponentRetriever
from services.retrieval_services import ComponentGenerator
from pydantic import BaseModel
import os
from dotenv import load_dotenv
from huggingface_hub import login
import logging

load_dotenv()

app = FastAPI()
hf_token = os.getenv("HUGGINGFACE_ACCESS_TOKEN")
retriever = ComponentRetriever()
generator = ComponentGenerator(hf_token)



class GenerationRequest(BaseModel):
    prompt: str


@app.post("/generate")
async def generate_component(request: GenerationRequest):
    try:
        # Login to Hugging Face
        login(hf_token)

        # Retrieve relevant components
        context = retriever.retrieve_components(request.prompt)

        # Generate new component
        generated_code = generator.generate_component(request.prompt, context)

        print(generated_code)

         # Ensure the response is JSON-serializable
        if isinstance(generated_code, StreamingResponse):
            # async for chunk in generated_code.body_iterator:
                
            #     print("Chunk received:", chunk)  # ✅ Print each chunk

            return generated_code

        if not generated_code:
            logging.error("Failed to generate component for prompt: %s", request.prompt)
            return JSONResponse(status_code=500, content={"detail": "Component generation failed"})
       

        return JSONResponse(
            status_code=200,
            content={
                "code": generated_code,
                "context": [doc.metadata["name"] for doc in context]
            }
        )

    except Exception as e:
        logging.exception("Unexpected error during component generation")
        return JSONResponse(status_code=500, content={"detail": f"Internal Server Error: {str(e)}"})

    

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)