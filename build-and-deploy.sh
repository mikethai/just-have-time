#! /bin/bash

export PROJECT_ID=kkchack22-just-have-time
export REGION=asia-east1
export CONNECTION_NAME=kkchack22-just-have-time:asia-east1:postgres-just-have-time-dev

gcloud builds submit \
    --tag gcr.io/$PROJECT_ID/just-have-time \
    --project $PROJECT_ID

gcloud run deploy just-have-time \
    --image gcr.io/$PROJECT_ID/just-have-time \
    --platform managed \
    --region $REGION \
    --allow-unauthenticated \
   --add-cloudsql-instances $CONNECTION_NAME \
    --project $PROJECT_ID