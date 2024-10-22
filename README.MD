# Performance Test for Frontend and Backend Integration

## Objective

The goal of this test is to load and display a list of network assets in a performant and efficient manner. Each asset should be displayed in its own component, showcasing all fields in a manner that is both aesthetically pleasing and easy to navigate. A **search functionality** should be implemented to filter the results based on the "host" column.

Candidates are encouraged to improve both the **Go backend** and the **Next.js frontend** to achieve optimal performance and usability. **Efficiency and speed** are key, and we expect you to identify, fix, and document any performance bottlenecks you find.

### Requirements

1. **Performance Optimization**:
    - The app should handle large datasets efficiently. The backend must be optimized for fast API responses, and the frontend should render the data without unnecessary delays.

2. **Signature Calculation**:
    - The signature generation functionality **must be preserved**. Each asset, along with its associated ports and IPs, must have a signature calculated and included in the response.

3. **Aesthetic Frontend**:
    - The frontend must be user-friendly and visually appealing. Candidates should focus on presenting the data in a clear and structured way that enhances the user experience.

4. **Search Functionality**:
    - A search feature should be implemented in the frontend to filter the assets by the **"host"** column. The search should be efficient, responsive, and provide relevant results quickly.

5. **Problem Identification and Documentation**:
    - Candidates should **identify all performance problems** (both in the frontend and backend), **fix them**, and **document the changes**. A list of identified inefficiencies and their solutions must be provided.

## Notes

We strive for excellence in everything we do, and we are looking for top-tier talent to join our team. Code quality, efficiency, and adherence to best practices will be assessed.

## Backend Setup

1. **Generate the Database**  
   In the `Backend` folder:

   Run the following command to create and populate the SQLite database with 250,000 assets:

   ```bash
   go run cmd/db/main.go
   ```

2. **Start the API Server**  
   Run the following command to start the Go API server:

   ```bash
   go run cmd/api/main.go
   ```

   The API server will be running at `http://localhost:8080`.

## Frontend Setup

1. **Install Dependencies**  
   Navigate to the frontend directory and install the required dependencies:

   ```bash
   npm install
   ```

2. **Start the Next.js Application**  
   Run the following command to start the frontend application:

   ```bash
   npm run dev
   ```

   The frontend will be running at `http://localhost:3000`.

## API Endpoints

- **Get all assets or filter by host**:  
  This endpoint returns all assets :

  ```http
  GET http://localhost:8080/assets
  ```

---

### Key Points for Candidates

- Focus on improving **speed** and **efficiency** in both the backend API and the frontend rendering.
- Ensure that the **signature calculation** for each asset, port, and IP remains intact and is efficient.
- The **frontend design** should be visually appealing and ensure a smooth user experience.
- Implement a **search functionality** that allows filtering by the "host" column.
- Document all **identified performance problems**, your **solutions**, and any **changes** made to improve the application's performance.
