from pynamodb.models import Model
from pynamodb.attributes import UnicodeAttribute
from pynamodb.indexes import GlobalSecondaryIndex, AllProjection
import os

# Define the Domain Index
class DomainIndex(GlobalSecondaryIndex):
    class Meta:
        index_name = "DomainIndex"
        projection = AllProjection()

    domain = UnicodeAttribute(hash_key=True)

# Define the Email Index
class EmailIndex(GlobalSecondaryIndex):
    class Meta:
        index_name = "EmailIndex"
        projection = AllProjection()

    email = UnicodeAttribute(hash_key=True)

# Define the main table model
class CompromisedCredential(Model):
    class Meta:
        table_name = os.getenv("COMPROMISEDCREDENTIALS_TABLE_NAME", "compromised-credentials")
        region = "us-east-1"

    id = UnicodeAttribute(hash_key=True) 
    email = UnicodeAttribute(null=True)
    domain = UnicodeAttribute(null=True)
    username = UnicodeAttribute(null=True)
    password = UnicodeAttribute(null=True)

    # Attach the indexes
    domain_index = DomainIndex()
    email_index = EmailIndex()
