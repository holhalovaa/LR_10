from fastapi import FastAPI
import requests

app = FastAPI(title="Python Service")

@app.get("/call-go")
def call_go():
    try:
        response = requests.get("http://localhost:8081/")
        return {"status": "success", "go_response": response.text}
    except Exception as e:
        return {"status": "error", "message": str(e)}