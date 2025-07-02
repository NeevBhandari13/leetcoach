from collections import defaultdict
from fastapi import FastAPI
from fastapi.responses import JSONResponse
from models.schemas import GPTRequest
from llm.llm_client import chat_with_gpt
from uuid import uuid4


# initialises FastAPI app
app = FastAPI()

# for testing
@app.get("/")
def read_root():
    return {"message": "Hello, this is your backend!"}

# checks that app is running
@app.get("/ping")
def ping():
    return {"status": "ok"}

# just a test
@app.get("/gpt-test")
def test_gpt():
    instructions = "You are a helpful assistant."
    input = [{"role": "user", "content": "Tell me a funny joke"},
             {"role": "assistant", "content": "Why did the chicken cross the road? To get to the other side!"},
             {"role": "user", "content": "Why is that funny?"},
             {"role": "assistant", "content": "Because it was a chicken joke!"},
             {"role": "user", "content": "What is a chicken joke?"},
             {"role": "developer", "content": "Talk in old English"}]
    
    response = chat_with_gpt(instructions, input)
    return {"response": response.output_text}

# @app.get("/start-interview")
# def start_interview():
#     session_id = str(uuid4()) # initialise a uuid for new session
#     sessions[session_id] = [] # initialise a new session with empty list of messages
    
#     state = "Greeting"
#     developer_prompt = STATE_DEVELOPER_PROMPTS[state]
#     instructions = INSTRUCTIONS
#     input = [{"role": "developer", "content": developer_prompt}]

#     response = chat_with_gpt(instructions, input)
    
#     reply = response.output_text
#     sessions[session_id].append({"role": "assistant", "content": reply})

#     return {"session_id": session_id, "chat": sessions[session_id]}

@app.post("/chat")
def chat(request: GPTRequest):
    instructions = request.instructions
    input = request.input
    response = chat_with_gpt(instructions, input)
    # JSONResponse creates JSON-encoded response with 200 status code and Content-Type: application/json header
    return JSONResponse(content={"response": response.output_text})

