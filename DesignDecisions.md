
# Design Decisions for Challenge

This document outlines the architectural and design decisions made for the SpyCloud Challenge, detailing the backend data ingestion pipeline, API design, and frontend implementation.

---

## 1. The Architecture

### a) Data Ingestion Pipeline
- **S3 Bucket**:
  - Files are uploaded to an S3 bucket (`raw-data-bucket-7`).
  - The bucket is encrypted with AES256 and configured to block public access.
  - Versioning is enabled to maintain historical file versions.
- **Trigger**:
  - An AWS Lambda function (`DataIngestionFunction`) written in **Golang** is triggered whenever a new object is created in the S3 bucket.
- **Processing and Storage**:
  - The Lambda function parses the uploaded file and processes each row to extract credentials (username, domain, and password).
  - The processed data is stored in a DynamoDB table (`compromised-credentials`).

### b) API Layer
- **API Gateway**:
  - An API Gateway connects to a second Lambda function (`APIFetchFunction`) running a **FastAPI** Python application.
  - The endpoint `/credentials` supports filtering by `email`, `domain`, and `limit` parameters.
- **Lambda Integration**:
  - The API Gateway is integrated with the FastAPI Lambda using the AWS `aws_proxy` type, ensuring minimal overhead and direct invocation.

### c) Frontend
- **Locally Hosted**:
  - The frontend is hosted locally and communicates with the serverless API through a pre-configured proxy.
- **Technology Stack**:
  - Built using **React.js**, **TypeScript**, and **Material UI**, featuring a data grid for displaying and filtering credentials.

---

## 2. Data Ingestion

### Parsing Logic
- The file rows did not consistently follow the `username@domain:password` format.
- The parser was designed to handle variations such as `password:username@domain` and similar patterns.

### Database Table Structure
- DynamoDB table: `compromised-credentials`
  - **Attributes**:
    - `id` (Primary Key - Hash Key): Unique identifier for each credential.
    - `email` (String, nullable): Extracted email.
    - `domain` (String, nullable): Extracted domain.
    - `username` (String, nullable): Extracted username.
    - `password` (String, hashed): Hashed password.
  - **Indexes**:
    - `EmailIndex`: GSI for querying by email.
    - `DomainIndex`: GSI for querying by domain.

### Password Hashing
- The **bcrypt** algorithm was used to hash passwords for security.

### Performance Tuning
#### First Iteration:
- Sequentially processed each row, hashing the password and writing one record at a time to DynamoDB.
- Issues:
  - Timing out with ingestion times >90 seconds.
  - The bcrypt algorithm’s ~200-500 ms per hash and individual network calls caused delays.

#### Second Iteration:
- Leveraged **concurrent programming** to reduce ingestion time to 10 seconds:
  - **CPU Optimization**: AWS Lambda provides 1 CPU per 1768 MB of RAM, so 5304 MB RAM was allocated to utilize 3 CPUs.
  - **Batch Processing**:
    - Used **Golang channels** to process 25 lines at a time (DynamoDB’s bulk write limit).
    - Parsed, hashed, and bulk wrote items into DynamoDB to minimize network overhead.
  - Achieved an ingestion time of ~10 seconds.

---

## 3. API Design
- **Lambda + API Gateway**:
  - The FastAPI application, running in `APIFetchFunction`, leverages the **PynamoDB ORM** to interact with DynamoDB.
  - Supports filtering by `email`, `domain`, and `limit` at the `/credentials` path.
  - Enables efficient querying through GSIs (`EmailIndex`, `DomainIndex`).

---

## 4. UI Design
- **Frontend Implementation**:
  - Built with **React.js** and **TypeScript**.
  - Utilizes **Material UI DataGrid** for displaying, filtering, and managing credentials.
  - Communicates with the API Gateway for dynamic data fetching.
- **Features**:
  - Filters for email and domain.
  - Integration with the `/credentials` API endpoint to retrieve and display data.

---

These design decisions prioritize scalability, performance, and security while ensuring a user-friendly experience for interacting with the processed data.
