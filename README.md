# Micro Pizzas

This is a monorepo for microservices I built to learn about backend a little more

## The strucure

- 🤑orders: Orders handle shopping cart and the actual orders
- 🍕products: Products manage the products we have in our menu and ones we have ingredients for
- 😎users: Users handle user authorizations and all that jazz
- 🧑‍⚕️sentry: Sentry monitors our microservices and aggregates their logs

Every microservice above has its own database. They communicate with eachother through GRPC.
