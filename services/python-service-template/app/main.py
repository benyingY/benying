from fastapi import FastAPI

app = FastAPI(title="service-template")


@app.get("/healthz")
def healthz():
    return {"status": "ok"}
