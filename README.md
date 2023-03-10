# SmolAge Image rendering service

## Overview
This repository contains a simple image manipulation API that can be used to modify and return images via a RESTful interface. The API allows the caller to specify one or more addons to be added to a base image. The resulting image is returned in PNG format.

## Requirements

The following packages are required to run the API:

- net/http
- image
- github.com/disintegration/gift
- github.com/nfnt/resize

## Usage

The API is invoked by sending a GET request to the `/api/smol` endpoint with the desired image id and addons specified in the URL. The `id` represents the base image, and the `addons` represent the images to be added to the base image. The response will be a PNG image containing the manipulated image.

`GET /api/smol?id=1&addons=2,3`

This request will return an image with the base image 1 and addons 2 and 3 applied.

## Running the API

To run the API locally, clone the repository and navigate to the root directory. Then, run the following command:

`go run main.go`

The API will start on localhost:8080 by default.


## Deploying to a Cloud Service


We are using Vercel's GO feature. It requires having a "api" subfolder with an index.go file inside. 

Please check https://vercel.com/docs/concepts/functions/serverless-functions/runtimes/go#advanced-go-usage
