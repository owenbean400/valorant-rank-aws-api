# Valorant Rank Data

This project provides a way to **store and manage Valorant rank data** using **AWS DynamoDB**, with an **API endpoint** for accessing rank history and clip data.

## 🚀 Overview

The project runs on **AWS** using a **Golang monolithic server architecture** deployed via **AWS Lambda**.
It stores and retrieves Valorant player rank history and gameplay clips within DynamoDB.

## 🧰 Requirements

### Cloud Infrastructure

* **AWS Services**

  * DynamoDB
  * Lambda

### Local Setup

* **Golang** installed on the build machine

## ☁️ AWS Setup

### DynamoDB

Two DynamoDB tables are required: one for rank history and one for clips.

#### **Rank Table**

Stores Valorant rank history for a specific player.

| Attribute      | Type   | Description                                    |
| -------------- | ------ | ---------------------------------------------- |
| `puuid_match`  | String | Partition key — player match unique identifier |
| `raw_date_int` | Number | Sort key — timestamp as integer                |

#### **Clips Table**

Stores Valorant gameplay clips.
Each clip can be a file URL or an external link from **YouTube** or **Twitch**.

| Attribute | Type   | Description                    |
| --------- | ------ | ------------------------------ |
| `uuid`    | String | Partition key — unique clip ID |

### Lambda Function

Use the provided build scripts to compile and package the Lambda function:

* **Linux/MacOS:** `build.sh`
* **Windows:** `build.ps1`

These scripts generate the `function.zip` file for deployment to AWS Lambda.

**Recommended Lambda Configuration**

* **Runtime:** Amazon Linux 2023
* **Architecture:** x86_64

### Environment Variables

The following environment variables are required for the Lambda configuration:

| Key                           | Description                                                                                                        | Required |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------ | -------- |
| `PLAYER_PUUID`                | Valorant player UUID from Riot Games.                                                                              | ✅        |
| `VALORANT_API_KEY`            | API key from [henrikdev.xyz](https://docs.henrikdev.xyz/authentication-and-authorization) used to fetch rank data. | ✅        |
| `AWS_DYNAMODB_RANK_TABLE_ARN` | ARN of the DynamoDB table for rank history.                                                                        | ✅        |
| `AWS_DYNAMODB_CLIP_TABLE_ARN` | ARN of the DynamoDB table for clips.                                                                               | ✅        |
| `API_POST_PASSWORD`           | Optional password for API POST request authentication.                                                             | ❌        |

### IAM Role Permissions

The Lambda function requires access permissions to the DynamoDB tables as shown below:

```json
"Statement": [
  {
    "Sid": "AllowClipTableAccess",
    "Effect": "Allow",
    "Action": [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:Scan"
    ],
    "Resource": "AWS_DYNAMODB_TABLE_VALORANT_CLIPS_ARN"
  },
  {
    "Sid": "AllowRankTableAccess",
    "Effect": "Allow",
    "Action": [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:Query"
    ],
    "Resource": "AWS_DYNAMODB_TABLE_VALORANT_RANK_ARN"
  }
]
```

## 💡 Example Use Case

This project is used to display **BeanBaller's 10 most recent ranked matches** and related **gameplay clips** on a gaming website.

* 🎮 [Valorant BeanBaller Website](https://www.beanballer.com)
* 📘 [Valorant BeanBaller Swagger Docs](https://www.beanballer.com/doc/api)