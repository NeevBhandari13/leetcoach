from collections import defaultdict
from fastapi import FastAPI
from llm.llm_client import chat_with_gpt
from uuid import uuid4
from interview.developer_prompts import STATE_DEVELOPER_PROMPTS


INSTRUCTIONS = """
You are an AI technical interviewer working with a candidate through a coding problem.

Your tone should be calm, supportive, and professional throughout the interview. You are not here to judge, but to help the candidate reason through the problem in a collaborative way.

You will progress through several stages of the interview, starting with a greeting, then presenting the problem, allowing for clarifying questions, discussing the initial solution, and guiding optimization. Avoid providing code or direct answers—focus on prompting the candidate to think critically and articulate their reasoning.
"""
# stores sessions in memory
# TODO: store in database
sessions = {}

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

@app.get("/start-interview")
def start_interview():
    session_id = str(uuid4()) # initialise a uuid for new session
    sessions[session_id] = [] # initialise a new session with empty list of messages
    
    state = "Greeting"
    developer_prompt = STATE_DEVELOPER_PROMPTS[state]
    instructions = INSTRUCTIONS
    input = [{"role": "developer", "content": developer_prompt}]

    response = chat_with_gpt(instructions, input)
    
    reply = response.output_text
    sessions[session_id].append({"role": "assistant", "content": reply})

    return {"session_id": session_id, "chat": sessions[session_id]}


 