import os
import uuid
from flask import Flask, request, jsonify, render_template
import requests
from flask_cors import CORS

app = Flask(__name__)
CORS(app)

# API Configuration
# API_KEY = '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'
API_KEY = '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'
API_URL = 'https://openrouter.ai/api/v1/chat/completions'

headers = {
    'Authorization': f'Bearer {API_KEY}',
    'Content-Type': 'application/json'
}

SYSTEM_PROMPTS = {
    'ar': """أنت مساعد طبي ذكي، مهمتك:
1. اسأل عن الأعراض بشكل متتابع
2. قدم نصائح محددة بناء على الإجابات
3. لا تكرر التحيات أو التعريفات
4. ذكر الإخلاء القانوني مرة واحدة فقط
5. استخدم نقاطًا محددة عند تقديم النصائح
6. حافظ على إجابات موجزة ومركزة""",
    
    'en': """You are a medical AI assistant. Your role:
1. Ask follow-up questions sequentially 
2. Provide specific advice based on responses
3. Never repeat greetings/introductions
4. Mention disclaimer once initially
5. Use bullet points for advice
6. Keep responses concise and focused"""
}

conversation_history = {}

def detect_language(text):
    """Enhanced language detection with context awareness"""
    if any('\u0600' <= char <= '\u06FF' for char in text):
        return 'ar'
    # Check conversation history if empty input
    return 'en'

def get_system_prompt(language):
    """Return appropriate system prompt"""
    base_prompt = SYSTEM_PROMPTS[language]
    if language == 'ar':
        return f"{base_prompt}\n\nملاحظة: هذه استشارة أولية - يجب مراجعة طبيب للتشخيص الدقيق"
    return f"{base_prompt}\n\nNote: This is preliminary advice - consult a doctor for proper diagnosis"

def manage_conversation(session_id, question):
    """Manage conversation flow and history"""
    language = detect_language(question)
    
    if session_id not in conversation_history:
        conversation_history[session_id] = {
            'history': [{"role": "system", "content": get_system_prompt(language)}],
            'state': 'initial',
            'language': language
        }
        # Add initial greeting
        greeting = initial_message(language)
        conversation_history[session_id]['history'].append(
            {"role": "assistant", "content": greeting}
        )
    
    # Add user message
    conv = conversation_history[session_id]
    conv['history'].append({"role": "user", "content": question})
    
    # Generate response
    response = generate_ai_response(conv['history'], conv['language'])
    conv['history'].append({"role": "assistant", "content": response})
    
    return response

def initial_message(language):
    """Language-specific initial message"""
    if language == 'ar':
        return "مرحبًا، كيف يمكنني مساعدتك اليوم؟ يرجى وصف الأعراض الرئيسية."
    return "Hello, how can I assist you today? Please describe your main symptoms."

def generate_ai_response(history, lang):
    """Generate AI response with conversation context"""
    data = {
        "model": "meta-llama/llama-3-70b-instruct",
        "messages": history,
        "temperature": 0.3,
        "max_tokens": 500
    }

    try:
        response = requests.post(API_URL, json=data, headers=headers)
        if response.status_code == 200:
            return postprocess_response(
                response.json()['choices'][0]['message']['content'],
                lang
            )
        return error_message(lang)
    except Exception as e:
        print(f"API Error: {str(e)}")
        return error_message(lang)

def postprocess_response(text, lang):
    """Clean up the AI response"""
    # Remove redundant greetings
    remove_phrases = {
        'ar': ["مرحبًا", "أهلاً", "أنت دكتور ذكاء اصطناعي"],
        'en': ["Hello", "Hi there", "As an AI assistant"]
    }
    for phrase in remove_phrases[lang]:
        text = text.replace(phrase, "")
    return text.strip()

def error_message(lang):
    """Return appropriate error message"""
    if lang == 'ar':
        return "حدث خطأ في النظام. يرجى المحاولة مرة أخرى لاحقًا."
    return "System error. Please try again later."

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/ask', methods=['POST'])
def handle_query():
    data = request.json
    question = data.get('question', '').strip()
    session_id = data.get('session_id', str(uuid.uuid4()))
    
    if not question:
        return jsonify({"error": "No question provided"}), 400
    
    response = manage_conversation(session_id, question)
    return jsonify({
        "response": response,
        "session_id": session_id,
        "symptoms": []  # Added empty list for front-end compatibility
    })

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000, debug=True)  # Run the app on port 8000