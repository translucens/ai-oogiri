# AI-oogiri

Gemini Pro を使って謎かけをするデモアプリ

## 準備

* MySQL 8.0
    * `default_time_zone` は `Asia/Tokyo`
* Google Cloud プロジェクトの作成、権限の付与
    * [aiplatform.endpoints.predict
](https://cloud.google.com/vertex-ai/docs/reference/rest/v1/projects.locations.publishers.models/predict#iam-permissions) 権限が必要。Vertex AI User 
`roles/aiplatform.user` に含まれる

## 起動

1. `schema.sql` によりデータベースとテーブルを作成
1. `run.sh` の環境変数を書き換えた上で実行
1. http://localhost:8080 もしくはデプロイ先のURLをブラウザで開く