from pydantic import BaseModel
from typing import List

class Message(BaseModel):
    role: str
    content: str

class GPTRequest(BaseModel):
    instructions: str
    input: List[Message]

class GPTResponse(BaseModel):
    output_text: str
