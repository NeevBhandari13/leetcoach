import os
from openai import OpenAI
from dotenv import load_dotenv

load_dotenv()

client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))

def chat_with_gpt(instructions: str = None, input: str = None, model="gpt-4o-mini"):
    response = client.responses.create(model=model,
    instructions=instructions,
    input=input)
    return response


