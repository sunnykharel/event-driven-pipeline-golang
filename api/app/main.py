from fastapi import FastAPI, Query, HTTPException
from typing import List, Optional
from mangum import Mangum
from app.models import CompromisedCredential
from app.schemas import CompromisedCredentialResponse

app = FastAPI()

@app.get("/credentials", response_model=List[CompromisedCredentialResponse])
def get_credentials(
    email: Optional[str] = Query(None),
    domain: Optional[str] = Query(None),
    limit: int = Query(10, gt=0)
):
    if email:
        results = list(CompromisedCredential.email_index.query(email))
    elif domain:
        results = list(CompromisedCredential.domain_index.query(domain))
    else:
        results = list(CompromisedCredential.scan(limit=limit))

    if not results:
        raise HTTPException(status_code=404, detail="No credentials found")

    return [
        CompromisedCredentialResponse(
            id=result.id,
            email=result.email,
            domain=result.domain,
            username=result.username,
            password=result.password
        )
        for result in results
    ]

handler = Mangum(app)
