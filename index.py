from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import subprocess
import traceback
from typing import Dict, List

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# In-memory storage for conversation history
conversation_history: Dict[str, List[str]] = {}

class Query(BaseModel):
    question: str
    conversation_id: str = None  # Optional conversation ID for context

def detect_language(text):
    if any("\u0600" <= char <= "\u06FF" for char in text):
        return "ar"
    else:
        return "en"

def build_prompt(question: str, conversation_id: str = None) -> str:
    """
    Builds a prompt for the AI model by including conversation history if a conversation_id is provided.
    """
    if not conversation_id:
        return question

    history = conversation_history.get(conversation_id, [])
    # Combine the conversation history with the new question
    full_prompt = "\n".join(history) + "\n" + question
    return full_prompt

@app.post("/chat/")
async def chat(query: Query):
    try:
        detected_language = detect_language(query.question)

        # Build the prompt with conversation history
        prompt = build_prompt(query.question, query.conversation_id)

        command = ["ollama", "run", "dolphin-llama3"]
        result = subprocess.run(
            command,
            input=prompt,
            capture_output=True,
            text=True,
            encoding='utf-8'
        )

        if result.returncode != 0:
            raise HTTPException(status_code=500, detail=f"Model error: {result.stderr.strip()}")

        response_text = result.stdout.strip()
        if not response_text:
            response_text = "No response from the model. Please try again."

        # Translate response if needed
        if detected_language == "ar" and response_text and not any("\u0600" <= char <= "\u06FF" for char in response_text):
            translation_command = ["ollama", "run", "translate-to-ar"]
            result = subprocess.run(
                translation_command,
                input=response_text,
                capture_output=True,
                text=True,
                encoding='utf-8'
            )
            response_text = result.stdout.strip() if result.stdout.strip() else response_text

        # Update conversation history if a conversation_id is provided
        if query.conversation_id:
            if query.conversation_id not in conversation_history:
                conversation_history[query.conversation_id] = []
            conversation_history[query.conversation_id].append(f"User: {query.question}")
            conversation_history[query.conversation_id].append(f"AI: {response_text}")

        return {"answer": response_text, "language": detected_language}

    except subprocess.CalledProcessError as e:
        raise HTTPException(status_code=500, detail=f"Subprocess error: {str(e)}")
    except UnicodeEncodeError as e:
        raise HTTPException(status_code=500, detail=f"Encoding error: {str(e)}")
    except Exception as e:
        error_details = traceback.format_exc()
        raise HTTPException(status_code=500, detail=f"Internal Server Error: {error_details}")