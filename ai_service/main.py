from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import subprocess
import traceback

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

class Query(BaseModel):
    question: str
    conversation_id: str = None  # Optional conversation ID for context

def detect_language(text):
    if any("\u0600" <= char <= "\u06FF" for char in text):
        return "ar"
    else:
        return "en"

@app.post("/chat/")
async def chat(query: Query):
    try:
        detected_language = detect_language(query.question)

        command = ["ollama", "run", "dolphin-llama3"]
        result = subprocess.run(
            command,
            input=query.question,
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

        return {"answer": response_text, "language": detected_language}

    except subprocess.CalledProcessError as e:
        raise HTTPException(status_code=500, detail=f"Subprocess error: {str(e)}")
    except UnicodeEncodeError as e:
        raise HTTPException(status_code=500, detail=f"Encoding error: {str(e)}")
    except Exception as e:
        error_details = traceback.format_exc()
        raise HTTPException(status_code=500, detail=f"Internal Server Error: {error_details}")


