from pydantic import BaseModel
from typing import Optional

class CompromisedCredentialResponse(BaseModel):
    id: str
    email: Optional[str]
    domain: Optional[str]
    username: Optional[str]
    password: Optional[str]
