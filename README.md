
# Project

This repository contains the project. Below are detailed instructions on how to run the project locally, deploy the backend using AWS SAM, and connect the frontend to the serverless API. For more information about the design choices made during the project, refer to the [DesignDecisions.md](DesignDecisions.md) file.

---

## **Running the Frontend Locally**

The project is pre-configured with the URL of the serverless AWS API, making it easy to run the frontend locally without additional changes.

### Steps:
1. Navigate to the `ui` directory:
   ```bash
   cd ui
   ```

2. Install dependencies:
   ```bash
   yarn install
   ```

3. Build the application:
   ```bash
   yarn build
   ```

4. Start the development server:
   ```bash
   yarn start
   ```

The application will now run locally at `http://localhost:3000` and connect to the serverless API.

---

## **Deploying the Backend with AWS SAM**

This project uses AWS SAM (Serverless Application Model) as Infrastructure as Code (IaC) to define and deploy the backend. Follow these steps to deploy the backend to AWS and integrate it with the frontend.

### **Steps to Deploy the Backend and Run Frontend**

1. **Set Up Required Tools**  
   Ensure you have the following installed:
   - **AWS CLI**: Install from [here](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html).  
     Configure your AWS credentials:
     ```bash
     aws configure
     ```
   - **AWS SAM CLI**: Install from [here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html).  
   - **Node.js**: Download from [here](https://nodejs.org/).  
   - **Yarn**: Install globally using:
     ```bash
     npm install --global yarn
     ```

   Ensure your AWS IAM user has the following permissions:
   - Full access to **DynamoDB**, **Lambda**, **API Gateway**, **S3**.
   - Administrator privileges.

2. **Build the SAM Application**  
   From the root directory of the project, run:
   ```bash
   sam build
   ```

3. **Deploy the Application**  
   Deploy the backend using:
   ```bash
   sam deploy --guided
   ```
   Follow the on-screen prompts:
   - Stack Name: Provide a name for your stack.
   - Region: Specify your preferred AWS region.
   - Accept the default values unless customization is needed.

4. **Update the Frontend with API Gateway URL**  
   Once deployed, retrieve the API Gateway URL from the deployment output. Update the `proxy` field in `ui/package.json`:
   ```json
   {
     "proxy": "https://<api-gateway-url>"
   }
   ```

5. **Run the Frontend Locally**  
   Navigate to the `ui` directory and run:
   ```bash
   yarn build
   yarn start
   ```

---

## **Important Notes**

- The `design-decisions.md` file contains a detailed explanation of the architectural and design decisions made for the project.
- Always ensure your AWS credentials and IAM permissions are correctly configured before deploying to avoid errors.
- For production deployments, make sure to follow AWS security best practices and restrict IAM permissions.

Enjoy exploring and running the project! 🚀
