from fastapi import FastAPI
from llm.llm_client import chat_with_gpt


app = FastAPI()


@app.get("/")
def read_root():
    return {"message": "Hello, this is your backend!"}

# checks that app is running
@app.get("/ping")
def ping():
    return {"status": "ok"}


@app.get("/gpt-test")
def test_gpt():
    instructions = "You are a helpful assistant."
    input = "Give me 10 cool names for a startup for an ai technical coding interviewer"
    
    reply = chat_with_gpt(instructions, input)
    return {"response": reply}
