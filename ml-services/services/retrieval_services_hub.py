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
    def __init__(self,hf_token):
        self.tokenizer = AutoTokenizer.from_pretrained("mistralai/Mistral-7B-v0.1")
        self.model = AutoModelForCausalLM.from_pretrained("mistralai/Mistral-7B-v0.1")
        self.hf_token = hf_token
        self.api_url = "https://api-inference.huggingface.co/models/mistralai/Mistral-7B-Instruct-v0.3/v1/chat/completions"
        self.headers = {"Authorization": f"Bearer {self.hf_token}"}
        
        self.llm = HuggingFaceHub(
            repo_id="mistralai/Mistral-7B-Instruct-v0.3",
            model_kwargs={
                "temperature": 0.7,
                "max_new_tokens": 1024,
                "streaming": True,
                "max_retries": 5,
                "wait_time": 2
            },
            huggingfacehub_api_token=self.hf_token
        )
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
        
        # self.chain = self.prompt_template | self.model | StrOutputParser()
    # def query_huggingface_api(self, prompt, retries=3, wait_time=5):
    #     for i in range(retries):
    #         try:
    #             response = requests.post(
    #             self.api_url,
    #             headers=self.headers,
    #             json={
    #                 "messages": prompt,
    #                 "parameters": {
    #                     "stream": True,  # Enable streaming
    #                     "return_full_text": False
    #                 }
    #             },
    #             stream=True
    #         )
    #             if response.status_code == 503:
    #                 try:
    #                     error_data = response.json()
    #                     if "estimated_time" in error_data:
    #                         wait = error_data["estimated_time"] + 5
    #                         print(f"Model loading, waiting {wait} seconds...")
    #                         time.sleep(wait)
    #                         continue
    #                 except json.JSONDecodeError:
    #                     print("Received non-JSON 503 response, retrying...")
    #                     time.sleep(wait_time)
    #                     wait_time *= 2
    #                     continue
    #             response.raise_for_status() #if any errors raise immediately
    #             return response
    #         except requests.exceptions.HTTPError as e:
    #             if response.status_code == 429:
    #                 print(f"Rate limit hit. Retrying in {wait_time} seconds...")
    #                 time.sleep(wait_time)
    #                 wait_time *= 2
    #             else:
    #                 # Return a clean error message
    #                 return f"HTTP Error {response.status_code}: {response.text}"
    #         except Exception as e:
    #             return f"Error: {str(e)}"
    #     return f"Failed after {retries} retries"

    def generate_component(self, prompt, context):
        
        formatted_context="\n\n".join([
            f"Component {i+1}:\n{doc.page_content}" 
            for i, doc in enumerate(context)
        ])


        #### to be used if huggingface interface api called
        chain = LLMChain(llm=self.llm, prompt=self.prompt_template)
        
        def generate():
            attempt = 0
            wait_time = self.initial_wait
            
            while attempt < self.retries:
                try:
                    for chunk in chain.stream({
                        "context": formatted_context,
                        "prompt": prompt
                    }):
                        print("RAW CHUNK------:", chunk)
                        # Filter out system message echoes
                        if "System: You are an expert" not in chunk['text']:
                            yield f"data: {json.dumps({'content': chunk['text']})}\n\n"
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


        #### to be used if local model used
        # inputs = self.tokenizer(formatted_prompt, return_tensors="pt", truncation=True, max_length=4096)

        # # Generate output from Mistral model
        # output_ids = self.model.generate(**inputs, max_new_tokens=500)

        # # Decode output
        # response = self.tokenizer.decode(output_ids[0], skip_special_tokens=True)

        # return response