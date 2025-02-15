from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import subprocess
import traceback
from typing import Dict, List
import re

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
    # Improved Arabic detection including common Arabic characters and diacritics
    arabic_pattern = re.compile(r'[\u0600-\u06FF\u0750-\u077F\u08A0-\u08FF\uFB50-\uFDFF\uFE70-\uFEFF]+')
    arabic_match = arabic_pattern.search(text)
    return "ar" if arabic_match else "en"

def preprocess_arabic(text: str) -> str:
    """Preprocess Arabic text for better model handling"""
    # Remove excessive whitespace while preserving Arabic text
    text = ' '.join(text.split())
    # Normalize Arabic characters
    text = text.replace('ی', 'ي').replace('ک', 'ك')
    return text

def build_prompt(question: str, conversation_id: str = None) -> str:
    """Enhanced prompt building with Arabic support"""
    language = detect_language(question)
    base_prompt = ""
    
    if language == "ar":
        question = preprocess_arabic(question)
        base_prompt = "أجب على السؤال التالي باللغة العربية:\n"
    
    if not conversation_id:
        return base_prompt + question

    history = conversation_history.get(conversation_id, [])
    full_prompt = base_prompt + "\n".join(history) + "\n" + question
    return full_prompt

@app.post("/chat/")
async def chat(query: Query):
    try:
        detected_language = detect_language(query.question)
        
        if detected_language == "ar":
            query.question = preprocess_arabic(query.question)

        prompt = build_prompt(query.question, query.conversation_id)

        # Add Arabic system prompt as part of the main prompt
        if detected_language == "ar":
            system_prompt = "أنت مساعد ذكي متخصص في التواصل باللغة العربية. يرجى الرد بالعربية الفصحى.\n\n"
            prompt = system_prompt + prompt

        # Use mistral model without --system flag
        command = ["ollama", "run", "mistral"]
        
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
            response_text = "لا يوجد رد من النموذج. يرجى المحاولة مرة أخرى." if detected_language == "ar" else "No response from the model. Please try again."

        # No need for translation if the model responds in Arabic for Arabic queries
        
        # Update conversation history if a conversation_id is provided
        if query.conversation_id:
            if query.conversation_id not in conversation_history:
                conversation_history[query.conversation_id] = []
            conversation_history[query.conversation_id].append(f"User: {query.question}")
            conversation_history[query.conversation_id].append(f"AI: {response_text}")

        return {
            "answer": response_text, 
            "language": detected_language,
            "preprocessed": detected_language == "ar"
        }

    except Exception as e:
        error_message = "حدث خطأ في النظام" if detected_language == "ar" else "System error occurred"
        error_details = traceback.format_exc()
        raise HTTPException(status_code=500, detail=f"{error_message}: {error_details}")