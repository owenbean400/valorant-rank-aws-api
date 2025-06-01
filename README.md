# Valorant Rank Data

This is for those who want to own their Valorant rank data summary information in AWS DynamoDB and have an API endpoint for fetching history.

## Requirements

Below are requirements cloud infrastructure to setup this project for personal use.

- AWS
  - DynamoDB
  - Lambda Function

## Setup

This project was built for users to setup their own application on AWS.

### API

[Valorant API Doc Authentication](https://docs.henrikdev.xyz/authentication-and-authorization)

### AWS

This project was implemented to run on AWS cloud infrastructure. Below is some help with setting up AWS. Note that it is beneficial to read through AWS documentation on how lambda and dynamoDB works.

#### Lambda Function

#### DynamoDB

#### Environment Variable

| Key | Description | Default Value |
| -- | -- | -- |
| PLAYER_PUUID | The Valorant Player UUID from Riot Games | `""` |
| VALORANK_API_KEY | The Valorant API key from henrikdev.xyz restful API | `""` |
| AWS_DYNAMODB_ROLE_ARN | The AWS role ARN consume from STS setup | `""` |
| AWS_DYNAMODB_TABLE_ARN | The AWS table ARN to connect to cloud database | `""` |
| AWS_DYNAMODB_SESSION_NAME | The AWS Session name of STS token | `"valorant-dynamodb-session"` |

## Use Case

I am using this project to display history of my last 10 rank games on gaming website I've made to mess with who look up my website while gaming.

[Valorant BeanBaller Website](www.beanballer.com)