import json
import requests
import time
from langchain_core.prompts import ChatPromptTemplate
from langchain_community.llms import HuggingFaceHub
# from langchain_core.output_parsers import StrOutputParser
from langchain.chains import LLMChain
from transformers import AutoTokenizer, AutoModelForCausalLM
from fastapi.responses import StreamingResponse
from fastapi import  HTTPException


class ComponentGenerator:
    def __init__(self, hf_token):
        self.hf_token = hf_token
        self.api_url = "https://api-inference.huggingface.co/models/mistralai/Mistral-7B-Instruct-v0.3"
        self.headers = {
            "Authorization": f"Bearer {self.hf_token}",
            "Content-Type": "application/json"
        }
        self.prompt_template = ChatPromptTemplate.from_messages([
            ("system", """You are an expert React developer. Generate components using:
            - BaseWeb UI components
            - Tailwind CSS for styling
            - TypeScript
            
            Follow these rules:
            1. Use functional components
            2. Export as default
            3. Include all necessary imports
            4. Use proper TypeScript types
            
            Context components:
            {context}"""),
            ("human", "User request: {prompt}")
        ])
        self.retries=5
        self.initial_wait=10

    def generate_component(self, prompt, context):
        formatted_context = "\n\n".join([
            f"Component {i+1}:\n{doc.page_content}" 
            for i, doc in enumerate(context)
        ])
        
        full_prompt = self.prompt_template.format(
            context=formatted_context,
            prompt=prompt
        )

        def generate():
            attempt = 0
            wait_time = self.initial_wait
            buffer = ""

            while attempt < self.retries:
                try:
                    response = requests.post(
                        self.api_url,
                        headers=self.headers,
                        json={
                            "inputs": full_prompt,
                            "parameters": {
                                "stream": True,
                                "return_full_text": False,
                                "max_new_tokens": 1024
                            }
                        },
                        stream=True
                    )

                    if response.status_code == 503:
                        error = response.json().get("error", {})
                        if "estimated_time" in error:
                            time.sleep(error["estimated_time"] + 5)
                            continue

                    response.raise_for_status()

                    for byte_chunk in response.iter_content(chunk_size=512):
                        if byte_chunk:
                            try:
                                chunk = byte_chunk.decode("utf-8")
                                print('Raw chunk----->',chunk)
                                
                                
                                    
                                if chunk:
                                   
                                    buffer += chunk
                                    print("-----buff-----",buffer)
                                    yield f"data: {json.dumps({'content': chunk})}\n\n"
                            except json.JSONDecodeError:
                                continue
                    # yield f"data: {json.dumps({'content': buffer})}\n\n"
                    yield "data: [DONE]\n\n"
                    return

                except Exception as e:
                    attempt += 1
                    print(f"Attempt {attempt} failed: {str(e)}")
                    if attempt < self.retries:
                        print(f"Retrying in {wait_time} seconds...")
                        time.sleep(wait_time)
                        wait_time *= 2
                    else:
                        yield f"data: {json.dumps({'error': str(e)})}\n\n"
                        yield "data: [DONE]\n\n"

        return StreamingResponse(
            generate(),
            media_type="text/event-stream",
            headers={"Cache-Control": "no-cache"}
        )