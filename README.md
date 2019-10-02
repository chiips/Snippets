# Snippets

This repository contains snippets of a production-quality API and Single-Page Application (SPA). The code is not complete but provides examples of useful design patterns—such as dependency injection, responding appropriately to cancelled requests, and gracefully handling errors—which can be extended to further develop a complete API and SPA.

## API

The API folder contains snippets for developing a RESTful API in Go to communicate with any front-end client that accepts and sends JSON data and with a PostgreSQL database for long-term data storage.

The API folder contains an additional README.md file further discussing the code.

## SPA

The SPA folder contains snippets for developing a Single-Page Application in VueJS to communicate with any back-end server that accepts and sends JSON data.

The SPA folder contains an additional README.md file further discussing the code.

## nginx.conf

The nginx.conf file contains the nginx setup. Nginx serves as a reverse proxy so the SPA and API can easily communicate with each other and send cookies even when hosted on different servers. The nginx.conf file also includes code to limit the rate of requests to the server.

On initial nginx setup the file comes with additional code (e.g., to set up logs) which has been removed to highlight what's necessary for a rate-limited reverse proxy.

## Licence

This code is licensed under the MIT License.