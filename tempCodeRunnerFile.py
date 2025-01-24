from ollama_api import OllamaClient

client = OllamaClient()
response = client.generate_completion(model="dolphin-llama3", prompt="Why is the sky blue?")
print(response)