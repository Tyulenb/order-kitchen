# Order-kitchen
Fast Food Restaurant simulator.
The main objective was to master gRPC and Redis while creating a functional application that mimics the operations of a fast food restaurant.
![demo](https://github.com/Tyulenb/order-kitchen/blob/main/docs/demo.png)
## Technology Stack
- gRPC: For efficient inter-service communication
- Redis: Used as the database for storing order data
## Poject Structure
Order-kitchen consists of three main components:
- **gRPC server**: Stores data about orders in Redis database, manages order processing.
- **kitchen**: 
	1. gRPC client which simulates work of kitchen; 
	2. Interacts with the server to retrieve and process orders
- **POS**:  
	1. A gRPC client that simulates the point of sales system.
	2. Displays the order results and interactions with customers.
	3. You can see a demonstration of this component in the image above.
## Download
```bash
git clone https://github.com/Tyulenb/order-kitchen.git
cd order-kitchen
docker compose up --build
```
