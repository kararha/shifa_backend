from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import subprocess
import traceback
from typing import Dict, List
import re
import asyncio
from cachetools import TTLCache
from concurrent.futures import ThreadPoolExecutor, TimeoutError
import threading
import sys

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

# Add cache with 1-hour TTL
response_cache = TTLCache(maxsize=100, ttl=3600)

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
    """Enhanced prompt building with medical Arabic support for DeepSeek"""
    language = detect_language(question)
    
    if language == "ar":
        question = preprocess_arabic(question)
        system_prompt = """<|system|>أنت مساعد طبي متخصص باللغة العربية. اتبع هذه القواعد:
1. قدم نصائح طبية أولية واضحة ودقيقة
2. استخدم لغة عربية فصحى مبسطة
3. كن مباشراً في إجاباتك
4. اذكر متى يجب استشارة الطبيب
5. تجنب التشخيص القطعي

<|user|>"""
    else:
        system_prompt = """<|system|>You are a specialized medical assistant. Follow these rules:
1. Provide clear initial medical advice
2. Use simple language
3. Be direct in your responses
4. Mention when to consult a doctor
5. Avoid definitive diagnosis

<|user|>"""

    if not conversation_id:
        return system_prompt + question

    history = conversation_history.get(conversation_id, [])
    context = "\nسياق المحادثة:\n" if language == "ar" else "\nConversation context:\n"
    full_prompt = system_prompt + context + "\n".join(history) + "\n\nالسؤال الحالي: " + question if language == "ar" else system_prompt + context + "\n".join(history) + "\n\nCurrent question: " + question
    
    return full_prompt

def run_model_sync(command: List[str], prompt: str) -> str:
    """Run the model synchronously with proper encoding"""
    try:
        # Force UTF-8 encoding for stdin/stdout
        process = subprocess.Popen(
            command,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=False  # Changed to handle bytes instead of text
        )
        
        # Encode prompt to UTF-8 bytes
        prompt_bytes = prompt.encode('utf-8')
        stdout, stderr = process.communicate(input=prompt_bytes)
        
        if process.returncode != 0:
            raise Exception(f"Model error: {stderr.decode('utf-8', errors='replace').strip()}")
            
        # Decode output with UTF-8
        return stdout.decode('utf-8', errors='replace').strip()
    except Exception as e:
        raise Exception(f"Model execution error: {str(e)}")

async def run_model_with_timeout(command: List[str], prompt: str, timeout: int = 15) -> str:
    """Run the model with timeout using thread pool"""
    loop = asyncio.get_event_loop()
    with ThreadPoolExecutor() as executor:
        try:
            return await loop.run_in_executor(
                executor,
                lambda: run_model_sync(command, prompt)
            )
        except TimeoutError:
            raise HTTPException(status_code=408, detail="Request timeout")
        except Exception as e:
            raise Exception(f"Model error: {str(e)}")

# Add model configuration
AVAILABLE_MODELS = {
    "medical": "dolphin-llama3",  # Primary model for medical queries
    "fallback": "mistral"       # Fallback model if primary fails         # For complex medical terminology
}

def select_model(query: str, language: str) -> str:
    """Select the appropriate model based on query complexity and language"""
    # Add your model selection logic here
    return AVAILABLE_MODELS["medical"]  # Default to medical model for now

async def try_alternate_model(command: List[str], prompt: str, current_model: str) -> tuple[str, str]:
    """Try an alternate model if the current one fails"""
    try:
        response = await run_model_with_timeout(command, prompt)
        return response, current_model
    except Exception:
        fallback_model = AVAILABLE_MODELS["fallback"]
        if current_model != fallback_model:
            command[2] = fallback_model
            try:
                response = await run_model_with_timeout(command, prompt)
                return response, fallback_model
            except Exception as e:
                raise Exception(f"Both primary and fallback models failed: {str(e)}")
        raise

@app.post("/chat/")
async def chat(query: Query):
    try:
        # Try to get cached response
        cache_key = f"{query.conversation_id}:{query.question}"
        if cache_key in response_cache:
            return response_cache[cache_key]

        detected_language = detect_language(query.question)
        
        if detected_language == "ar":
            query.question = preprocess_arabic(query.question)

        prompt = build_prompt(query.question, query.conversation_id)
        selected_model = select_model(query.question, detected_language)
        
        command = ["ollama", "run", selected_model]
        
        # Try primary model with fallback
        response_text, used_model = await try_alternate_model(command, prompt, selected_model)

        if not response_text:
            response_text = "لا يوجد رد من النموذج" if detected_language == "ar" else "No response from model"

        # Update conversation history
        if query.conversation_id:
            if query.conversation_id not in conversation_history:
                conversation_history[query.conversation_id] = []
            conversation_history[query.conversation_id].append(f"User: {query.question}")
            conversation_history[query.conversation_id].append(f"AI: {response_text}")

        result = {
            "answer": response_text,
            "language": detected_language,
            "preprocessed": detected_language == "ar",
            "model_used": used_model
        }

        # Cache the response
        response_cache[cache_key] = result
        
        return result

    except Exception as e:
        error_message = "حدث خطأ في النظام" if detected_language == "ar" else "System error occurred"
        error_details = traceback.format_exc()
        raise HTTPException(status_code=500, detail=f"{error_message}: {error_details}")